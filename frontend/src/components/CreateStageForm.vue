<template>
  <form @submit.prevent="handleSubmit">
    <div class="modal-group" v-if="!parentStageId">
      <label for="stageType">Stage type</label>
      <select v-model="stageType" required>
        <option v-for="value in stageTypes" :key="value" :value="value">
          {{ value }}
        </option>
      </select>
    </div>

    <div class="modal-group" v-if="stageType !== ValueStageType.ParallelStage">
      <label for="deploymentId">Deployment ID</label>
      <input v-model="deploymentId" type="text" required />
    </div>

    <div
      class="modal-group"
      v-if="hasParams && stageType === ValueStageType.SequentialStage"
    >
      <label for="parameters">Params</label>
      <div
        v-for="(value, key) in filteredParams"
        :key="key"
        class="modal-group"
      >
        <label for="param">{{ key }}</label>
        <input v-model="parameters[key]" type="text" />
      </div>
    </div>

    <div v-if="hasParams">
      <label for="showAll">Show all params</label>
      <input type="checkbox" v-model="showAll" />
    </div>

    <div class="modal-group-button">
      <button type="submit">Create</button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref, watch, computed } from "vue";

import { addStage, addSubStage, getStageParameters } from "@/services";
import { RequestsStage, ValueStageType } from "@/api";

const props = defineProps<{
  sendpostId: number;
  previousStageId?: number;
  parentStageId?: number;
}>();
const emit = defineEmits(["stageAdded"]);

const stageTypes = Object.values(ValueStageType);

const stageType = ref(stageTypes[1]);

const deploymentId = ref("");

const parameters = ref<{ [key: string]: string }>({});

const hasParams = computed(() => Object.keys(parameters.value).length > 0);

const showAll = ref(false);

watch(deploymentId, async (newValue) => {
  if (newValue && newValue.trim() !== "") {
    try {
      const params = await getStageParameters(newValue);
      parameters.value = params;
    } catch (error) {
      console.error(error);
      alert("Error getting stage parameters");
    }
  } else {
    parameters.value = {};
  }
});

const filteredParams = computed(() => {
  return showAll.value
    ? parameters.value
    : Object.fromEntries(
        Object.entries(parameters.value).filter(([_, value]) => !value)
      );
});

const handleSubmit = async () => {
  console.log(
    "[CreateStageForm] handleSubmit",
    "stageType:",
    stageType.value,
    deploymentId.value,
    props.previousStageId,
    props.parentStageId
  );
  try {
    let id_depl = deploymentId.value;
    if (stageType.value === ValueStageType.ParallelStage) {
      id_depl = "-";
    }

    const newStage: RequestsStage = {
      type: stageType.value,
      deployment_id: id_depl,
      stage_parameters: parameters.value,
    };

    // if add-sub stage
    if (props.parentStageId) {
      await addSubStage(props.sendpostId, props.parentStageId, newStage);

      return;
    }

    newStage.previous_stage_id = props.previousStageId;

    await addStage(props.sendpostId, newStage);
  } catch (error) {
    console.error("Error add stage", error);
  } finally {
    emit("stageAdded");
  }
};
</script>
