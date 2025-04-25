<template>
  <div class="view">
    <div class="view-header">
      <span>
        <h1>Runs</h1>
      </span>
      <span>
        <button @click="handleCreateSendpost(null)">
          <span class="material-symbols-outlined"> playlist_add </span>
        </button>
      </span>
    </div>
    <div class="sendpost-list">
      <SendpostRunItem
        v-for="post in sendposts"
        :key="post.id"
        :sendpost="post"
        :isExpanded="expandSendpost[post.id]"
        @refresh="fetchSendposts"
        @toggleExpand="toggleExpand(post.id)"
        @handleCopySendpost="handleCreateSendpost(post.id)"
        @handleSendpostRun="handleSendpostRun(post)"
        @handleSendpostFailed="handleSendpostFailed(post)"
        @handleSendpostCompleted="handleSendpostCompleted(post)"
      />
    </div>

    <BaseModal v-if="showCreateSendpost" @close="showCreateSendpost = false">
      <template #header>
        <h2>Create</h2>
      </template>
      <template #body>
        <CreateSendpostForm
          :copySendpost="copySendpost"
          @sendpostCreated="handleSendpostCreated"
        />
      </template>
    </BaseModal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";

import SendpostRunItem from "@/components/SendpostRunItem.vue";

import CreateSendpostForm from "@/components/CreateSendpostForm.vue";

import BaseModal from "@/components/BaseModal.vue";

import { getSendposts } from "@/services/sendpostService";

import type { ExtendedSendpost } from "@/types";

import { ValueStateType } from "@/api";

const sendposts = ref<ExtendedSendpost[]>([]);

const expandSendpost = ref<Record<number, boolean>>({});

const showCreateSendpost = ref(false);

const copySendpost = ref<number | null>(null);

const fetchSendposts = async () => {
  console.log("[SendpostView] fetchSendposts");
  try {
    const data = await getSendposts();

    sendposts.value = data;
  } catch (error) {
    console.error("[SendpostView] Error get sendposts", error);
  }
};

const handleCreateSendpost = (sendpostID: number | null) => {
  console.log("[SendpostView] handleCreateSendpost: sendpostID: ", sendpostID);
  if (sendpostID) {
    copySendpost.value = sendpostID;
  }

  showCreateSendpost.value = true;
};

const handleSendpostRun = (post: ExtendedSendpost) => {
  console.log("[SendpostView] handleSendpostRun");
  post.state = ValueStateType.Running;
};

const handleSendpostCompleted = (post: ExtendedSendpost) => {
  console.log("[SendpostView] handleSendpostCompleted");
  post.state = ValueStateType.Completed;
};

const handleSendpostFailed = (post: ExtendedSendpost) => {
  console.log("[SendpostView] handleSendpostFailed");
  post.state = ValueStateType.Failed;
};

const toggleExpand = (postId: number) => {
  expandSendpost.value = {
    ...expandSendpost.value,

    [postId]: !expandSendpost.value[postId],
  };
};

const handleSendpostCreated = async () => {
  console.log("[SendpostView] handleSendpostCreated");
  showCreateSendpost.value = false;

  await fetchSendposts();
};

onMounted(fetchSendposts);
</script>
