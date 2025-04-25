export class WsConnector {
  private ws: WebSocket | null = null;

  private baseUrl: string;

  /**
  
     * @param baseUrl - базовый URL для WebSocket, например: "ws://localhost:8080"
  
     */

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  /**
  
     * Establishes a WebSocket connection for sendpost notifications.
  
     * Automatically constructs the connection URL using the sendpostID.
  
     * @param sendpostID The ID of the sendpost to subscribe for notifications.
  
     * @param onMessage Handler for incoming messages.
  
     * @param onError (Optional) Handler for error events.
  
     * @param onOpen (Optional) Handler for the open event.
  
     * @param onClose (Optional) Handler for the close event.
  
     */

  public getNotifications(
    sendpostID: number,

    onMessage: (event: MessageEvent) => void,

    onError?: (error: Event) => void,

    onOpen?: (event: Event) => void,

    onClose?: (event: CloseEvent) => void
  ): void {
    console.log(
      `[WsConnector] Connecting to ${this.baseUrl}/sendposts/${sendpostID}/run/ws`
    );

    const url = `${this.baseUrl}/sendposts/${sendpostID}/run/ws`;

    this.ws = new WebSocket(url);

    this.ws.onmessage = onMessage;

    if (onError) this.ws.onerror = onError;

    if (onOpen) this.ws.onopen = onOpen;

    if (onClose) this.ws.onclose = onClose;
  }

  public send(message: string): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(message);
    } else {
      console.error("WebSocket is not open.");
    }
  }

  public close(): void {
    console.log("[WsConnector] Closing WebSocket connection.");
    if (this.ws) {
      this.ws.close();

      this.ws = null;
    }
  }
}
