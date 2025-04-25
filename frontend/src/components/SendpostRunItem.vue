<template>
  <div class="sendpost-item">
    <div class="sendpost-header">
      <button @click="emit('toggleExpand')">
        <span v-if="!isExpanded" class="material-symbols-outlined">
          expand_content
        </span>
        <span v-if="isExpanded" class="material-symbols-outlined"> hide </span>
      </button>
      <span>{{ sendpost.name }}</span>
      <span
        ><strong>{{
          props.sendpost.global_parameters
            ? props.sendpost.global_parameters["mailing_date"]
            : ""
        }}</strong></span
      >
      <div class="actions">
        <button
          @click="runSendpost"
          :class="[{ blocked: sendpost.state === ValueStateType.Running }]"
        >
          <span>Run</span>
        </button>
        <button @click="$emit('handleCopySendpost')">
          <span class="material-symbols-outlined"> content_copy </span>
        </button>
        <button @click="deleteSendpostItem">
          <span class="material-symbols-outlined"> delete </span>
        </button>
      </div>
    </div>
    <div class="sendpost-info">
      <div class="sendpost-state">
        <span>{{ sendpost.state }}</span>
      </div>
    </div>
    <div v-if="isExpanded" class="sendpost-details">
      <label for="description"><strong>Description:</strong></label>
      <p id="description">{{ sendpost.description }}</p>
      <StageGraph ref="stageGraphReg" :sendpostId="sendpost.id" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from "vue";

import { ValueStateType } from "@/api";

import { RunSendpost, wsService } from "@/services";

import StageGraph from "./StageGraph.vue";

import { deleteSendpost } from "@/services/sendpostService";

import type { ExtendedSendpost } from "@/types";

const emit = defineEmits([
  "deleteSendpost",
  "refresh",
  "toggleExpand",
  "handleSendpostRun",
  "handleCopySendpost",
  "handleSendpostFailed",
  "handleSendpostCompleted",
]);

const props = defineProps<{
  sendpost: ExtendedSendpost;

  isExpanded: boolean;
}>();

const isExpanded = ref(props.isExpanded);

const stageGraphReg = ref<InstanceType<typeof StageGraph> | null>(null);

const handlers = {
  onRun: () => emit("handleSendpostRun"),
  onFailed: () => {
    callFetchStages();
    emit("handleSendpostFailed");
    wsService.close(props.sendpost.id);
  },
  onCompleted: () => {
    callFetchStages();
    emit("handleSendpostCompleted");
    wsService.close(props.sendpost.id);
  },
  onUpdated: () => callFetchStages(),
  onError: () => {
    callFetchStages();
  },
};

watch(
  () => props.isExpanded,

  (newValue) => {
    isExpanded.value = newValue;
  }
);

const callFetchStages = () => {
  console.log("[SendpostItem] callFetchStages");
  stageGraphReg.value?.fetchStages();
};

onMounted(() => {
  wsService.connect(props.sendpost.id, handlers);
});

onUnmounted(() => {
  wsService.close(props.sendpost.id);
});

const runSendpost = () => {
  console.log("[SendpostItem] runSendpost");
  if (props.sendpost.state === ValueStateType.Running) {
    alert("Sendpost is already running");
    return;
  }
  try {
    RunSendpost(props.sendpost.id);
    wsService.connect(props.sendpost.id, handlers);
    emit("handleSendpostRun");
  } catch (error) {
    console.error("[SendpostItem] runSendpost: Error start sendpost", error);
  }
};

const deleteSendpostItem = async () => {
  console.log("[SendpostItem] deleteSendpostItem");
  if (confirm("Are you shure?")) {
    try {
      await deleteSendpost(props.sendpost.id);
      emit("deleteSendpost", props.sendpost.id);
      emit("refresh");
    } catch (error) {
      console.error(
        "[SendpostItem] deleteSendpostItem: Error delete sendpost:",
        error
      );
    }
  }
};
</script>

<style>
.blocked {
  opacity: 0.5;
}
.sendpost-item {
  border: 1px solid #ccc;

  padding: 10px;

  margin-bottom: 10px;
}
.sendpost-details {
  margin-top: 10px;

  padding-left: 20px;

  border-left: 2px solid #ccc;
}
.sendpost-details p {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: pre-wrap;
  word-wrap: break-word;
}
.sendpost-state {
  display: flex;
  align-items: start;
}
</style>
