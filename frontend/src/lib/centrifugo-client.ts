/**
 * Centrifugo WebSocket client wrapper for real-time messaging
 */

export class CentrifugoClient {
  private ws: WebSocket | null = null;
  private url: string;
  private token: string = '';
  private userId: string = '';
  private subscriptions: Map<string, (data: any) => void> = new Map();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;

  constructor(wsUrl: string, userId: string) {
    this.url = wsUrl;
    this.userId = userId;
  }

  async connect(token: string): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        this.token = token;
        const wsUrl = `${this.url}?token=${token}`;
        this.ws = new WebSocket(wsUrl);

        this.ws.onopen = () => {
          console.log('Centrifugo connected');
          this.reconnectAttempts = 0;
          resolve();
        };

        this.ws.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data);
            this.handleMessage(data);
          } catch (e) {
            console.error('Failed to parse message:', e);
          }
        };

        this.ws.onerror = (error) => {
          console.error('Centrifugo error:', error);
          reject(error);
        };

        this.ws.onclose = () => {
          console.log('Centrifugo disconnected');
          this.attemptReconnect();
        };
      } catch (error) {
        reject(error);
      }
    });
  }

  private handleMessage(data: any): void {
    // Handle different message types from Centrifugo
    if (data.type === 'publication') {
      const channel = data.channel;
      const callback = this.subscriptions.get(channel);
      if (callback) {
        callback(data.data);
      }
    }
  }

  subscribe(
    channel: string,
    callback: (data: any) => void
  ): void {
    this.subscriptions.set(channel, callback);
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(
        JSON.stringify({
          type: 'subscribe',
          channel,
        })
      );
    }
  }

  unsubscribe(channel: string): void {
    this.subscriptions.delete(channel);
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(
        JSON.stringify({
          type: 'unsubscribe',
          channel,
        })
      );
    }
  }

  private attemptReconnect(): void {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      const delay = this.reconnectDelay * this.reconnectAttempts;
      console.log(`Attempting to reconnect in ${delay}ms...`);
      setTimeout(() => {
        this.connect(this.token).catch((error) => {
          console.error('Reconnection failed:', error);
        });
      }, delay);
    } else {
      console.error('Max reconnection attempts reached');
    }
  }

  disconnect(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.subscriptions.clear();
  }

  isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
  }
}
