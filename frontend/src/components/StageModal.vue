<template>
  <div v-if="loading">Загрузка...</div>

  <div v-else-if="stage">
    <div class="modal-group">
      <p><strong>Type:</strong> {{ stage.type }}</p>
      <p><strong>Deployment ID:</strong> {{ stage.deployment_id }}</p>
      <p><strong>State:</strong> {{ stage.state }}</p>
    </div>
    <div v-if="stage.type === ValueStageType.ParallelStage">
      <button class="modal-group-button" @click="callCreateSubStage(stage.id)">
        <span class="material-symbols-outlined"> playlist_add </span>
      </button>
      <SubStageGraph
        ref="subStageReg"
        :stageId="stage.id"
        :sendpostId="props.sendpostId"
      />
    </div>

    <div class="modal-group-button">
      <button @click="blockStageItem" v-if="isBlocked">
        <span class="material-symbols-outlined"> lock </span>
      </button>
      <button @click="blockStageItem" v-if="!isBlocked">
        <span class="material-symbols-outlined"> lock_open </span>
      </button>
      <button @click="deleteStageItem" v-if="!isBlocked">
        <span class="material-symbols-outlined"> delete </span>
      </button>
    </div>
  </div>

  <div v-else>
    <p>Error download data</p>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";

import SubStageGraph from "./SubStageGraph.vue";

import {
  deleteStage,
  getStageDetailedInfo,
  blockUnblockStage,
} from "@/services/stageService";

import { StageDetailed } from "@/types";
import { ValueStageType } from "@/api";

const props = defineProps(["sendpostId", "stageId", "isBlocked"]);

const emit = defineEmits(["handleStage"]);

const stage = ref<StageDetailed | null>(null);

const loading = ref(false);

const subStageReg = ref<InstanceType<typeof SubStageGraph> | null>(null);

const callCreateSubStage = (parentId: number) => {
  if (subStageReg.value) {
    subStageReg.value.openAddStage(parentId);
  }
};

// Загрузка данных при изменении stageId

watch(
  () => props.stageId,

  async (newStageId) => {
    if (!newStageId) return;

    loading.value = true;

    try {
      stage.value = await getStageDetailedInfo(props.sendpostId, newStageId);
    } catch (error) {
      console.error("Error download stage", error);

      stage.value = null;
    } finally {
      loading.value = false;
    }
  },

  { immediate: true }
);

const deleteStageItem = async () => {
  if (!stage.value || !confirm("Delete stage?")) return;

  try {
    await deleteStage(props.sendpostId, stage.value.id);

    emit("handleStage");
  } catch (error) {
    console.error("Error delete stage", error);
  }
};

const blockStageItem = async () => {
  if (!stage.value || !confirm("Block stage?")) return;

  try {
    await blockUnblockStage(props.sendpostId, stage.value.id);

    emit("handleStage");
  } catch (error) {
    console.error("Error block stage", error);
  }
};
</script>

<style scoped>
.modal-group-button {
  display: flex;
  padding: 1px;
}
</style>
