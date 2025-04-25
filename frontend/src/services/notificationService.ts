import { WsConnector } from "@/ws";
import { ValueStateType } from "@/api";
import { BASE_PATH } from "./base";

type Handlers = {
  onRun?: () => void;
  onFailed?: () => void;
  onCompleted?: () => void;
  onError?: () => void;
  onUpdated?: () => void;
};

export class WsService {
  private readonly baseUrl: string;
  private readonly connections: Map<
    number,
    { connector: WsConnector; handlers: Handlers }
  > = new Map();

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  public connect(sendpostId: number, handlers: Handlers): void {
    if (this.connections.has(sendpostId)) {
      console.warn(`[WsService] Already connected to sendpost ${sendpostId}`);
      return;
    }

    const connector = new WsConnector(this.baseUrl);

    connector.getNotifications(
      sendpostId,
      (event) => this.handleMessage(sendpostId, event),
      (errorEvent) => this.handleError(sendpostId, errorEvent),
      undefined,
      (closeEvent) => {
        console.warn(
          `[WsService] Closed WS for sendpost ${sendpostId}`,
          closeEvent
        );

        // code 1006 = abnormal closure
        if (closeEvent.code === 1006) {
          this.handleError(sendpostId, closeEvent);
        }

        this.connections.delete(sendpostId);
      }
    );

    this.connections.set(sendpostId, { connector, handlers });
    console.log(`[WsService] Connected to sendpost ${sendpostId}`);
  }

  private handleMessage(sendpostId: number, event: MessageEvent) {
    const entry = this.connections.get(sendpostId);
    if (!entry) return;

    const { handlers } = entry;

    console.log(
      `[WsService] Received WS message for sendpost ${sendpostId}:`,
      event.data
    );

    switch (event.data) {
      case ValueStateType.Running:
        handlers.onRun?.();
        break;
      case ValueStateType.Failed:
        handlers.onFailed?.();
        this.close(sendpostId);
        break;
      case ValueStateType.Completed:
        handlers.onCompleted?.();
        this.close(sendpostId);
        break;
      case ValueStateType.Updated:
        handlers.onUpdated?.();
        break;
      default:
        handlers.onError?.();
        break;
    }
  }

  private handleError(sendpostId: number, event: Event | CloseEvent) {
    console.warn(
      `[WsService] WebSocket error for sendpost ${sendpostId}`,
      event
    );

    const entry = this.connections.get(sendpostId);
    entry?.handlers.onError?.();
    this.close(sendpostId);
  }

  public close(sendpostId: number): void {
    const entry = this.connections.get(sendpostId);
    if (entry) {
      entry.connector.close();
      this.connections.delete(sendpostId);
    }
  }

  public closeAll(): void {
    for (const [_, entry] of this.connections.entries()) {
      entry.connector.close();
    }
    this.connections.clear();
  }
}

const wsService = new WsService(BASE_PATH);

export default wsService;
