import {
  RequestsSenpost,
  ResponsesSendpost,
  SendpostApi,
  RequestsParameters,
} from "@/api";
import { apiClient } from "./apiClient";

const sendpostApi = new SendpostApi(apiClient);

export async function getSendposts(): Promise<ResponsesSendpost[]> {
  const response = await sendpostApi.getSendposts();
  return response.data;
}

export async function createSendpost(
  data: RequestsSenpost
): Promise<ResponsesSendpost> {
  const response = await sendpostApi.createSendpost(data);
  return response.data;
}

export async function getSendpost(id: number): Promise<ResponsesSendpost> {
  const response = await sendpostApi.getSendpost(id);
  return response.data;
}

export async function deleteSendpost(id: number): Promise<void> {
  await sendpostApi.deleteSendpost(id);
}

export async function copySendpost(
  id: number,
  data: RequestsSenpost
): Promise<ResponsesSendpost> {
  const response = await sendpostApi.copySendpost(id, data);
  return response.data;
}

export async function addUpdateSendpostParameters(
  id: number,
  parameters: RequestsParameters
) {
  await sendpostApi.addUpdateSendpostParameters(id, parameters);
}

export async function deleteSendpostParameter(id: number, key: string) {
  await sendpostApi.deleteSendpostParameter(id, key);
}
