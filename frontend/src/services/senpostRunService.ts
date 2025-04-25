import { SendpostRunnerApi } from "@/api";
import { apiClient } from "./apiClient";

const sendpostApi = new SendpostRunnerApi(apiClient);

export const RunSendpost = async (sendpostId: number): Promise<string> => {
  const response = await sendpostApi.sendpostsSendpostIdRunPost(sendpostId);
  return response.data;
};
