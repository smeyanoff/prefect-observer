import { createRouter, createWebHistory } from "vue-router";
import SendpostRunsView from "@/views/SendpostRunsView.vue";

const routes = [
  { path: "/runs", name: "sendpostRuns", component: SendpostRunsView },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
