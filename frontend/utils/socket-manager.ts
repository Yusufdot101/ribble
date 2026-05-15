import { BASE_CHAT_SERVICE_API_URL } from "./api";

class SocketManager {
    socket: WebSocket | null = null;

    connect(): WebSocket | null {
        if (
            this.socket &&
            (this.socket.readyState === WebSocket.OPEN ||
                this.socket.readyState === WebSocket.CONNECTING)
        ) {
            return this.socket;
        }

        const wsUrl = new URL(BASE_CHAT_SERVICE_API_URL);
        wsUrl.protocol = wsUrl.protocol === "https:" ? "wss:" : "ws:";
        wsUrl.pathname = `${wsUrl.pathname.replace(/\/$/, "")}/ws`;

        const socket = new WebSocket(wsUrl.toString());
        this.socket = socket;
        return socket;
    }

    getSocket(): WebSocket | null {
        return this.socket;
    }

    close() {
        if (!this.socket) return;

        this.socket.close();
        this.socket = null;
    }
}

export const socketManager = new SocketManager();
