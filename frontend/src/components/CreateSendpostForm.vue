<template>
  <form @submit.prevent="handleSubmit">
    <div class="modal-group">
      <label for="sendpostName">Sendpost name</label>
      <input
        v-model="sendpostName"
        type="text"
        placeholder="Enter sendpost name"
        required
      />
    </div>
    <div class="modal-group">
      <label for="sendpostDescription">Description</label>
      <textarea
        v-model="sendpostDescription"
        placeholder="Enter sendpost description"
      ></textarea>
    </div>
    <div class="modal-group">
      <label for="date" type="date">Sendpost date</label>
      <input id="date" type="date" v-model="selectedDate" />
    </div>
    <div class="modal-group-button">
      <button type="submit">Create</button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref } from "vue";

import { createSendpost, copySendpost } from "@/services/sendpostService";
import type { RequestsSenpost } from "@/api";

const formatDate = (date: Date): string => {
  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, "0");
  const day = date.getDate().toString().padStart(2, "0");

  return `${year}-${month}-${day}`;
};

const emit = defineEmits(["sendpostCreated"]);

const props = defineProps(["copySendpost"]);

const sendpostName = ref("");
const sendpostDescription = ref("");
const selectedDate = ref<string>(formatDate(new Date()));

const handleSubmit = async () => {
  console.log("[CreateSendpostForm] handleSubmit");
  try {
    const data: RequestsSenpost = {
      sendpost_name: sendpostName.value,
      description: sendpostDescription.value,
      global_parameters: {
        mailing_date: selectedDate.value,
      },
    };
    if (props.copySendpost) {
      await copySendpost(props.copySendpost, data);
    } else {
      await createSendpost(data);
    }
    sendpostName.value = "";

    emit("sendpostCreated");
  } catch (error) {
    console.error("Error create sendpost", error);
  }
};
</script>
