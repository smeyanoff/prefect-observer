import { apiClient } from "./apiClient";
import {
  RequestsStage,
  StageApi,
  ResponsesStage,
  ResponsesStageDetailed,
  RequestsParameters,
} from "@/api";

const stageApi = new StageApi(apiClient);

export async function getSendpostStages(
  sendpostId: number
): Promise<ResponsesStage[]> {
  const response = await stageApi.getSendpostStages(sendpostId);
  return response.data;
}

export async function addStage(
  sendpostId: number,
  data: RequestsStage
): Promise<ResponsesStage> {
  const response = await stageApi.addStageToSendpost(sendpostId, data);
  return response.data;
}

export async function getStageDetailedInfo(
  sendpostId: number,
  stageId: number
): Promise<ResponsesStageDetailed> {
  const response = await stageApi.getStageDetailedInfo(sendpostId, stageId);
  return response.data;
}

export async function deleteStage(
  sendpostId: number,
  stageId: number
): Promise<void> {
  await stageApi.deleteStage(sendpostId, stageId);
}

export async function addSubStage(
  sendpostId: number,
  stageId: number,
  data: RequestsStage
): Promise<ResponsesStage> {
  const response = await stageApi.addSubStage(sendpostId, stageId, data);
  return response.data;
}

export async function blockUnblockStage(
  sendpostId: number,
  stageId: number
): Promise<void> {
  await stageApi.blockUnblockStage(sendpostId, stageId);
}

export async function getSubStages(
  sendpostId: number,
  stageId: number
): Promise<ResponsesStage[]> {
  const response = await stageApi.getSubStages(sendpostId, stageId);
  return response.data;
}

export async function updateParameters(
  sendpostId: number,
  stageId: number,
  data: RequestsParameters
): Promise<void> {
  await stageApi.updateStageParameters(sendpostId, stageId, data);
}
