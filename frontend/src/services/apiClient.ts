import { Configuration } from "@/api";
import { BASE_PATH } from "./base";

export const apiClient = new Configuration({
  basePath: BASE_PATH,
});
