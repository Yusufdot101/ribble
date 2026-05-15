"use client";

import { useAuthStore } from "@/store/useAuthStore";
import { BASE_CHAT_SERVICE_API_URL } from "@/utils/api";
import { socketManager } from "@/utils/socket-manager";
import { createContext, useContext, useEffect, useRef, useState } from "react";

const SocketContext = createContext<WebSocket | null>(null);

export const SocketProvider = ({ children }: { children: React.ReactNode }) => {
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const accessToken = useAuthStore((state) => state.accessToken);

    const reconnectTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
    const disposedRef = useRef(false);
    const retryCounter = useRef(1);

    useEffect(() => {
        if (!accessToken || !BASE_CHAT_SERVICE_API_URL) return;

        disposedRef.current = false;

        const clearReconnectTimer = () => {
            if (reconnectTimer.current) {
                clearTimeout(reconnectTimer.current);
                reconnectTimer.current = null;
            }
        };

        const connect = () => {
            console.log("running");
            if (disposedRef.current) return;
            clearReconnectTimer();

            const ws = socketManager.connect();
            if (!ws) return;
            setSocket(ws);

            const sendToken = () => {
                if (disposedRef.current) return;
                if (ws.readyState !== WebSocket.OPEN) return;
                ws.send(JSON.stringify({ token: accessToken }));
            };

            ws.onopen = async () => {
                console.log("connected");
                retryCounter.current = 1;
                sendToken();
            };

            ws.onmessage = (event) => {
                const data = JSON.parse(event.data);
                console.log("received:", data);

                if (data.type === "error") {
                    console.error(data.message);
                    return;
                }
            };

            ws.onerror = (err) => {
                console.error("socket error", err);
            };

            ws.onclose = () => {
                console.log("closed");
                if (disposedRef.current) return;

                clearReconnectTimer();

                reconnectTimer.current = setTimeout(() => {
                    connect();
                }, retryCounter.current * 1000);

                retryCounter.current = Math.min(retryCounter.current * 2, 30);
            };
        };

        connect();

        return () => {
            disposedRef.current = true;
            clearReconnectTimer();
            socketManager.close();
            setSocket(null);
        };
    }, [accessToken]);

    return (
        <SocketContext.Provider value={socket}>
            {children}
        </SocketContext.Provider>
    );
};

export const useSocket = () => {
    return useContext(SocketContext);
};
