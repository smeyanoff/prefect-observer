import { StageInfoApi } from "@/api";
import { apiClient } from "./apiClient";

const sendpostApi = new StageInfoApi(apiClient);

export async function getStageParameters(deploymentID: string) {
  return (await sendpostApi.getStageParameters(deploymentID)).data;
}
