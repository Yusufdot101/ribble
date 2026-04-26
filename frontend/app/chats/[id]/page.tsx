"use client";
import BackArrowButton from "@/components/BackArrowButton";
import Message from "@/components/Message";
import MessageInput from "@/components/MessageInput";
import { useAuthStore } from "@/store/useAuthStore";
import { BASE_CHAT_SERVICE_API_URL } from "@/utils/api";
import { getChatMessages, MessageType } from "@/utils/messages";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";

const ChatPage = () => {
    const params = useParams();
    const chatID = params.id;
    const [messages, setMessages] = useState<MessageType[]>([]);
    const socketRef = useRef<WebSocket | null>(null);

    useEffect(() => {
        if (!chatID) return;
        (async () => {
            const messages = await getChatMessages(+chatID);
            setMessages(messages);
        })();
    }, [chatID]);

    const accessToken = useAuthStore((state) => state.accessToken);
    useEffect(() => {
        if (!accessToken || !chatID || !BASE_CHAT_SERVICE_API_URL) return;

        const wsUrl = new URL(BASE_CHAT_SERVICE_API_URL);
        wsUrl.protocol = wsUrl.protocol === "https:" ? "wss:" : "ws:";
        wsUrl.pathname = `${wsUrl.pathname.replace(/\/$/, "")}/messages/new`;

        const socket = new WebSocket(
            // `ws://${BASE_CHAT_SERVICE_API_URL}/messages/new`,
            wsUrl.toString(),
        );
        socketRef.current = socket;

        socket.onopen = () => {
            console.log("connected");
            socket.send(
                JSON.stringify({
                    token: accessToken,
                }),
            );
        };

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            console.log("received: ", data);
            if (data.type === "error") {
                console.error(data.message);
                return;
            }

            if (data.type === "messageDeleted") {
                setMessages((prev) => {
                    return prev.map((msg) =>
                        msg.ID === data.messageID
                            ? { ...msg, Deleted: true }
                            : msg,
                    );
                });
                return;
            }

            if (data.type === "messageEdited") {
                setMessages((prev) => {
                    return prev.map((msg) =>
                        msg.ID === data.messageID
                            ? {
                                  ...msg,
                                  Content: data.newContent,
                                  UpdatedAt: data.updatedAt,
                              }
                            : msg,
                    );
                });
            }

            setMessages((prev) => {
                if (!chatID) return prev;
                if ((data as MessageType).ChatID !== +chatID) return prev;
                return [...prev, data];
            });
        };

        socket.onclose = () => {
            console.log("socket closed");
        };

        socket.onerror = (err) => {
            console.error("socket error", err);
        };

        return () => {
            socket.close();
        };
    }, [accessToken, chatID]);

    const [menuIsOpen, setMenuIsOpen] = useState(false);
    const [selectedMessageID, setSelectedMessageID] = useState<number>();

    const containerRef = useRef<HTMLDivElement>(null);

    const [isEditingMessage, setIsEditingMessage] = useState(false);
    const [editingMessageID, setEditingMessageID] = useState<number>();

    const router = useRouter();
    return (
        <div
            ref={containerRef}
            className="flex-1 min-h-0 flex flex-col gap-y-[8px]"
            onClick={() => {
                setMenuIsOpen(false);
            }}
            onKeyDown={(e) => {
                if (e.key !== "Escape") return;
                e.preventDefault();
                setMenuIsOpen(false);
                setIsEditingMessage(false);
            }}
        >
            <div className="flex w-full h-[32px] gap-x-[8px] items-center">
                <BackArrowButton
                    handleClick={() => router.back()}
                    text="Chat"
                />
            </div>
            <div className="flex justify-center shrink-0">Username/email</div>

            <div className="flex-1 min-h-0 overflow-y-auto flex flex-col gap-y-[8px] p-[4px]">
                {messages?.map((message) => (
                    <Message
                        containerRef={containerRef}
                        menuIsOpen={menuIsOpen}
                        selectedMessageID={selectedMessageID ?? -1}
                        handleRightClick={(messageID: number) => {
                            setMenuIsOpen(true);
                            setSelectedMessageID(messageID);
                        }}
                        key={message.ID}
                        message={message}
                        editingMessageID={editingMessageID}
                        isEditing={isEditingMessage}
                        handleClickEdit={(messageID: number) => {
                            setIsEditingMessage(true);
                            setEditingMessageID(messageID);
                        }}
                        handleCancelMessageEdit={() =>
                            setIsEditingMessage(false)
                        }
                    />
                ))}
            </div>

            <MessageInput socketRef={socketRef} chatID={chatID} />
        </div>
    );
};

export default ChatPage;
