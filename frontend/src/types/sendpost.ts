import { ResponsesSendpost } from "@/api";

export type Sendpost = ResponsesSendpost;

export interface ExtendedSendpost extends Sendpost {
  isExpanded?: boolean;
}

export type Sendposts = Array<ExtendedSendpost>;
