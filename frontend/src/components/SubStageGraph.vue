<template>
  <div class="stage-graph">
    <div class="graph-container">
      <div v-for="stage in stages" :key="stage.id" class="stage-node-wrapper">
        <div
          class="stage-node"
          :class="[getStageClass(stage)]"
          @click="openStageDetail(stage)"
        >
          <span class="stage-type-badge">{{ stage.type }}</span>
          <span class="stage-state-icon">{{
            getStageStateSymbol(stage.state)
          }}</span>
        </div>
      </div>
    </div>
  </div>

  <BaseModal
    v-if="showCreateStageModal"
    @close="closeCreateStageModal"
    :sendpostId="sendpostId"
    :parentStageId="parentStageId"
  >
    <template #header>
      <h3>Create stage</h3>
    </template>
    <template #body>
      <CreateStageForm
        :sendpostId="sendpostId"
        :parentStageId="props.stageId"
        @stageAdded="handleStage"
      />
    </template>
  </BaseModal>

  <BaseModal v-if="showDetailedStageModal" @close="closeStageDetailModal">
    <template #header>
      <h3>Details</h3>
    </template>
    <template #body>
      <StageModal
        :sendpostId="sendpostId"
        :stageId="selectedStage!.id"
        :isBlocked="selectedStage!.is_blocked"
        @handleStage="handleStage"
      />
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";

import BaseModal from "./BaseModal.vue";

import StageModal from "./StageModal.vue";

import CreateStageForm from "./CreateStageForm.vue";

import type { Stage } from "@/types";

import { ValueStateType } from "@/api";

import { getSubStages } from "@/services";

const showCreateStageModal = ref(false);

const showDetailedStageModal = ref(false);

const props = defineProps(["stageId", "sendpostId"]);

const parentStageId = ref<number | null>(null);

const selectedStage = ref<Stage | null>(null);

const stages = ref<Stage[]>([]);

const fetchStages = async (): Promise<void> => {
  try {
    const data = await getSubStages(props.sendpostId, props.stageId);

    stages.value = Array.isArray(data) ? data : [];
  } catch (error) {
    console.error("Error get stages", error);

    stages.value = [];
  }
};

onMounted(fetchStages);

const handleStage = async (): Promise<void> => {
  console.log("[SubStageGraph] handleStage", selectedStage.value);

  closeStageDetailModal();
  closeCreateStageModal();

  await fetchStages();
};

const openAddStage = (parentId: number): void => {
  if (stages.value.length === 5) {
    console.error("Maximum number of stages reached");
    return;
  }
  showCreateStageModal.value = true;

  parentStageId.value = parentId;
};

defineExpose({ openAddStage });

const closeCreateStageModal = (): void => {
  showCreateStageModal.value = false;

  parentStageId.value = null;
};

const openStageDetail = (stage: Stage): void => {
  showDetailedStageModal.value = true;

  selectedStage.value = stage;
};

const closeStageDetailModal = (): void => {
  showDetailedStageModal.value = false;

  selectedStage.value = null;
};

const getStageClass = (stage: Stage): string => {
  switch (stage.state) {
    case ValueStateType.Failed:
      return "error-stage animate-blink";

    case ValueStateType.NeverRunning:
      return "never-running";

    case ValueStateType.Running:
      return "running-stage animate-blink";

    case ValueStateType.Completed:
      return "completed-stage";

    default:
      return "never-running";
  }
};

const getStageStateSymbol = (state: string): string => {
  switch (state) {
    case "RUNNING":
      return "▶";

    case "COMPLETED":
      return "✓";

    case "NEVERRUNNING":
      return "•";

    default:
      return "?";
  }
};
</script>
