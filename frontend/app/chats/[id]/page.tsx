"use client";
import BackArrowButton from "@/components/BackArrowButton";
import Message from "@/components/Message";
import MessageInput, { WebsocketMsg } from "@/components/MessageInput";
import { useAuthStore } from "@/store/useAuthStore";
import { BASE_CHAT_SERVICE_API_URL } from "@/utils/api";
import { ChatType, getChatByID, getChatUsers } from "@/utils/chats";
import { messageStore } from "@/utils/messagesStore";
import {
    getChatMessages,
    MessageType,
    syncChatMessages,
} from "@/utils/messages";
import { UserType } from "@/utils/users";
import { useParams, useRouter } from "next/navigation";
import { useCallback, useEffect, useRef, useState } from "react";
import { useOnlineStatus } from "@/hooks/useOnlineStatus";

const ChatPage = () => {
    const isOnline = useOnlineStatus();
    const newMessageID = (): string => {
        return crypto.randomUUID();
    };

    const newNumberID = (): number => {
        return Math.floor(Math.random() * 100_000_000) + 1;
    };

    const params = useParams();
    const chatID = params.id;
    const [messages, setMessages] = useState<MessageType[]>([]);
    const messagesRef = useRef(messages);
    const pendingMessages = useRef(new Map<string, WebsocketMsg>());
    const [chat, setChat] = useState<ChatType>();
    const [chatUsers, setChatUsers] = useState<UserType[]>([]);
    const socketRef = useRef<WebSocket | null>(null);
    const messagesContainer = useRef<HTMLDivElement | null>(null);
    const loggedInUserID = useAuthStore((state) => state.userID);

    useEffect(() => {
        if (!chatID) return;
        (async () => {
            const messages = await getChatMessages(+chatID);
            setMessages(messages);

            // add local messages
            const savedMessages = await messageStore.getByChat(+chatID);
            const queuedMessages = savedMessages.map((message, i) => ({
                Content: message.content,
                ChatID: message.chatID,
                Status: message.status ?? "pending",
                ClientID: message.clientID,
                SenderID: message.senderID,
                CreatedAt: message.CreatedAt ?? "",
                Deleted: false,
                ID: -(i + 1),
                DeletedAt: null,
                UpdatedAt: message.CreatedAt ?? "",
            }));

            for (const msg of savedMessages) {
                if (msg.status === "pending") {
                    pendingMessages.current.set(msg.clientID, msg);
                }
            }

            setMessages((prev) => [...prev, ...queuedMessages]);

            const chat = await getChatByID(+chatID);
            setChat(chat);

            const chatUsers = await getChatUsers(+chatID);
            if (!chatUsers) return;
            setChatUsers(chatUsers);
        })();
    }, [chatID]);

    const accessToken = useAuthStore((state) => state.accessToken);

    const manualClose = useRef(false);
    const retry = useRef(1);
    const reconnectTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
    const closeSocket = useCallback(() => {
        console.log("closed");
        manualClose.current = true;
        if (reconnectTimer.current) {
            clearTimeout(reconnectTimer.current);
            reconnectTimer.current = null;
        }
        socketRef.current?.close();
    }, []);

    const connect = useCallback(
        function run() {
            console.log("running");
            manualClose.current = false;
            if (!accessToken || !chatID || !BASE_CHAT_SERVICE_API_URL) return;
            const wsUrl = new URL(BASE_CHAT_SERVICE_API_URL);
            wsUrl.protocol = wsUrl.protocol === "https:" ? "wss:" : "ws:";
            wsUrl.pathname = `${wsUrl.pathname.replace(/\/$/, "")}/ws`;

            const socket = new WebSocket(wsUrl.toString());
            socketRef.current = socket;

            socket.onopen = async () => {
                console.log("connected");
                retry.current = 1;
                socket.send(
                    JSON.stringify({
                        token: accessToken,
                    }),
                );
                const lastMessage = messagesRef.current.at(-1);
                if (!chatID || !lastMessage) return;
                const missedMessages = await syncChatMessages(
                    +chatID,
                    lastMessage.ID,
                );
                setMessages((prev) => {
                    const seen = new Set(prev.map((m) => m.ID));
                    const uniqueMissed = missedMessages.filter(
                        (m) => !seen.has(m.ID),
                    );
                    return [...prev, ...uniqueMissed];
                });

                for (const [, pendingMessage] of pendingMessages.current) {
                    socketRef.current?.send(JSON.stringify(pendingMessage));
                }
            };

            socket.onmessage = (event) => {
                const data = JSON.parse(event.data);
                console.log("receieved: ", data);
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
                    return;
                }

                if (data.type === "nack") {
                    setMessages((prev) =>
                        prev.map((message) =>
                            message.ClientID === data.clientID
                                ? { ...message, Status: "failed" }
                                : message,
                        ),
                    );
                    const msg = pendingMessages.current.get(data.clientID);
                    if (msg) {
                        messageStore.update({ ...msg, status: "failed" });
                    }
                    pendingMessages.current.delete(data.clientID);
                    return;
                }

                if (data.type === "ack") {
                    setMessages((prev) =>
                        prev.map((message) =>
                            message.ClientID === data.clientID
                                ? { ...message, Status: "delivered" }
                                : message,
                        ),
                    );
                    pendingMessages.current.delete(data.clientID);
                    messageStore.delete(data.clientID);
                    return;
                }
                setMessages((prev) => {
                    if (!chatID) return prev;
                    const incoming = data as MessageType;
                    if (incoming.ChatID !== +chatID) return prev;
                    if (prev.some((m) => m.ID === incoming.ID)) return prev;
                    return [...prev, incoming];
                });
            };

            socket.onclose = () => {
                if (manualClose.current === true) return;
                reconnectTimer.current = setTimeout(() => {
                    run();
                }, retry.current * 1000);

                retry.current = Math.min(retry.current * 2, 60);
            };

            socket.onerror = (err) => {
                console.error("socket error", err);
            };
        },
        [accessToken, chatID],
    );

    const sendMessage = (message: string) => {
        if (!chatID || message.trim() === "" || !loggedInUserID) return;

        const socket = socketRef.current;
        if (!socket || socket.readyState !== WebSocket.OPEN) return;

        const clientID = newMessageID();

        const creationDate = new Date().toISOString();
        const msg: WebsocketMsg = {
            status: "pending",
            chatID: +chatID,
            senderID: loggedInUserID,
            clientID: clientID,
            content: message,
            type: "message",
            CreatedAt: creationDate,
        };

        setMessages((prev) => {
            if (!chatID) return prev;
            const newMessage: MessageType = {
                ClientID: clientID,
                Status: "pending",
                ChatID: +chatID,
                ID: -newNumberID(),
                Content: msg.content,
                CreatedAt: msg.CreatedAt ?? creationDate,
                Deleted: false,
                DeletedAt: null,
                SenderID: loggedInUserID,
                UpdatedAt: creationDate,
            };
            return [...prev, newMessage];
        });
        pendingMessages.current.set(clientID, msg);
        messageStore.add(msg);
        socket.send(JSON.stringify(msg));
    };

    useEffect(() => {
        connect();
        return () => {
            closeSocket();
        };
    }, [accessToken, chatID, closeSocket, connect]);

    useEffect(() => {
        messagesContainer.current?.scrollTo({
            top: messagesContainer.current.scrollHeight,
            behavior: "smooth",
        });
        messagesRef.current = messages;
    }, [messages]);

    const [menuIsOpen, setMenuIsOpen] = useState(false);
    const [selectedMessageID, setSelectedMessageID] = useState<number>();

    const containerRef = useRef<HTMLDivElement>(null);

    const [isEditingMessage, setIsEditingMessage] = useState(false);
    const [editingMessageID, setEditingMessageID] = useState<number>();

    const router = useRouter();

    return (
        <div
            ref={containerRef}
            className="flex-1 min-h-0 flex flex-col gap-y-[8px] relative"
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
            <div
                className={`${!isOnline ? "" : "opacity-0 invisible"} duration-300 right-1/2 translate-x-1/2 bg-red-500 p-[4px] rounded-[4px] absolute`}
            >
                <span className="text-[16px]">Currently offline</span>
            </div>
            <div className="flex w-full h-[32px] gap-x-[8px] items-center">
                <BackArrowButton
                    handleClick={() => router.back()}
                    text="Chat"
                />
            </div>
            <div className="flex justify-center shrink-0">
                <div className="flex gap-x-[4px]">
                    {chat?.name !== ""
                        ? chat?.name
                        : chatUsers
                              .filter((user) => user.id !== loggedInUserID)
                              .map((user) => (
                                  <div key={user.id}>{user.name}</div>
                              ))}
                </div>
            </div>

            <div
                ref={messagesContainer}
                className="flex-1 min-h-0 overflow-y-auto flex flex-col gap-y-[8px] p-[4px]"
            >
                {messages
                    .sort((a, b) => a.CreatedAt.localeCompare(b.CreatedAt))
                    .map((message) => (
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
                            username={
                                chat?.isGroup
                                    ? chatUsers.filter(
                                          (user) =>
                                              user.id === message.SenderID,
                                      )[0]?.name
                                    : undefined
                            }
                        />
                    ))}
            </div>

            <MessageInput
                handleSend={(message: string) => {
                    sendMessage(message);
                }}
            />
        </div>
    );
};

export default ChatPage;
