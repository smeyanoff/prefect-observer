<template>
  <div>
    <div
      class="sendpost-item"
      v-for="post in props.sendposts"
      :key="post.id"
      @click="chooseTemplate(post)"
    >
      <span>{{ post.name }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Sendposts, ExtendedSendpost } from "@/types";

const props = defineProps<{
  sendposts: Sendposts;
}>();

const emit = defineEmits(["chooseTemplate"]);

const chooseTemplate = (post: ExtendedSendpost) => {
  emit("chooseTemplate", post);
};
</script>

<style lang="css" scoped>
.sendpost-item {
  border-radius: 8px;
}
.sendpost-item:hover {
  background-color: #a6523f;
  color: #fff;
  cursor: pointer;
  transition: all 0.3s ease-in-out;
}
</style>
