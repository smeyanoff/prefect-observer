<template>
  <div class="stage-graph">
    <div class="graph-container">
      <div
        v-if="stages.length === 0"
        class="add-stage-edge"
        @click="openAddStage(undefined, false)"
      >
        +
      </div>
      <div v-else class="stages-wrapper">
        <div
          v-for="(stage, index) in stages"
          :key="stage.id"
          class="stage-node-wrapper"
        >
          <div
            :class="['add-stage-edge', { blocked: stage.is_blocked }]"
            @click="openAddStage(getPreviousStageId(index), stage.is_blocked)"
          >
            +
          </div>
          <div
            class="stage-node"
            :class="[getStageClass(stage), { blocked: stage.is_blocked }]"
            @click="openStageDetail(stage)"
          >
            <span class="stage-type-badge">{{ stage.type }}</span>
            <span class="stage-state-icon">{{
              getStageStateSymbol(stage.state)
            }}</span>
          </div>
        </div>
        <div
          :class="[
            'add-stage-edge',
            { blocked: stages[stages.length - 1].is_blocked },
          ]"
          @click="
            openAddStage(
              getPreviousStageId(stages.length),
              stages[stages.length - 1].is_blocked
            )
          "
        >
          +
        </div>
      </div>
    </div>

    <BaseModal v-if="showCreateStageModal" @close="closeCreateStageModal">
      <template #header>
        <h3>Create stage</h3>
      </template>
      <template #body>
        <CreateStageForm
          :sendpostId="sendpostId"
          :previousStageId="previousStageId"
          @stageAdded="handleStage"
        />
      </template>
    </BaseModal>

    <BaseModal v-if="selectedStage" @close="closeStageDetailModal">
      <template #header>
        <h3>Details</h3>
      </template>
      <template #body>
        <StageModal
          :sendpostId="sendpostId"
          :stageId="selectedStage.id"
          :isBlocked="selectedStage.is_blocked"
          @handleStage="handleStage"
        />
      </template>
    </BaseModal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";

import { getSendpostStages } from "@/services";

import CreateStageForm from "./CreateStageForm.vue";

import StageModal from "./StageModal.vue";

import BaseModal from "./BaseModal.vue";

import type { Stage } from "@/types";

import { ValueStateType } from "@/api";

// Определяем типы для props и emits

const props = defineProps(["sendpostId"]);

// Реактивное состояние

const stages = ref<Stage[]>([]);

const showCreateStageModal = ref(false);

const previousStageId = ref<number>();

const selectedStage = ref<Stage | null>(null);

const fetchStages = async (): Promise<void> => {
  try {
    console.info("[StageGraph] fetchStages", props.sendpostId);

    const data = await getSendpostStages(props.sendpostId);

    stages.value = Array.isArray(data) ? data : [];
  } catch (error) {
    console.error("Error get stages", error);

    stages.value = [];
  }
};

defineExpose({ fetchStages });

onMounted(fetchStages);

watch(
  () => props.sendpostId,

  async (newValue, oldValue) => {
    if (newValue !== oldValue) {
      await fetchStages();
    }
  }
);

// Функция для получения id предыдущей стадии по индексу

const getPreviousStageId = (index: number): number | undefined => {
  if (index === 0) {
    console.info("[StageGraph] getPreviousStageId", null);

    return;
  }

  const prevId = stages.value[index - 1].id;

  console.info("[StageGraph] getPreviousStageId", prevId);

  return prevId;
};

const openAddStage = (
  prevId: number | undefined,
  is_blocked: boolean
): void => {
  if (is_blocked) return;
  if (stages.value.length === 12) {
    console.error("max length reached");
    return;
  }
  previousStageId.value = prevId;

  showCreateStageModal.value = true;
};

const closeCreateStageModal = (): void => {
  showCreateStageModal.value = false;
};

const openStageDetail = (stage: Stage): void => {
  selectedStage.value = stage;
};

const closeStageDetailModal = (): void => {
  selectedStage.value = null;
};

const handleStage = async (): Promise<void> => {
  console.info("[StageGraph] handleStage", selectedStage.value);

  closeStageDetailModal();
  closeCreateStageModal();

  await fetchStages();
};

// Возвращает символ, соответствующий состоянию стадии

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

// Возвращает CSS-класс для стадии в зависимости от её состояния

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
</script>

<style>
.stage-graph {
  margin-top: 5px;
  background-color: #f9f9f9;
  padding: 20px;
  border-radius: 8px;
  overflow-y: hidden;
  overflow-x: auto;
}

.graph-container {
  display: flex;
  align-items: center;
}

.stages-wrapper {
  display: flex;
  align-items: center;
}

.stage-node-wrapper {
  display: flex;
  align-items: center;
  margin-right: 5px;
}

.stage-node {
  position: relative;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 2px solid #ccc;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  transition: transform 0.2s;
  margin: 0 10px;
}

.stage-node:hover {
  transform: scale(1.05);
}

.stage-type-badge {
  position: absolute;

  top: -15px;

  background-color: #6c757d;

  color: #fff;

  font-size: 10px;

  padding: 2px 5px;

  border-radius: 3px;

  pointer-events: none;
}

.stage-state-icon {
  font-size: 14px;

  color: #000;
}

.add-stage-edge {
  width: 20px;

  height: 20px;

  border-radius: 50%;

  background-color: #a2b4c5;

  display: flex;

  justify-content: center;

  align-items: center;

  cursor: pointer;

  font-weight: bold;

  transition: background-color 0.3s, transform 0.2s;
}

.add-stage-edge:hover {
  background-color: #ced4da;

  transform: scale(1.1);
}

.never-running {
  background-color: #b0b0b0;
}

.running-stage {
  background-color: #28a745;
}

.completed-stage {
  background-color: #007bff;
}

.error-stage {
  background-color: #dc3545;
}

/* Анимация мигания */
@keyframes blink {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
  100% {
    opacity: 1;
  }
}

.animate-blink {
  animation: blink 1s infinite;
}

.no-stages {
  margin-top: 10px;
}
</style>
