"use client";

import BackArrowButton from "@/components/BackArrowButton";
import Message from "@/components/Message";
import MessageInput, { WebsocketMsg } from "@/components/MessageInput";
import Menu from "@/components/Menu";
import { useAuthStore } from "@/store/useAuthStore";
import { ChatType, getChatByID, getChatUsers } from "@/utils/chats";
import { messageStore } from "@/utils/messagesStore";
import {
    getChatMessages,
    MessageType,
    syncChatMessages,
} from "@/utils/messages";
import { UserType } from "@/utils/users";
import { getUserPermissions, PermissionType } from "@/utils/permissions";
import { useOnlineStatus } from "@/hooks/useOnlineStatus";
import { useParams, useRouter } from "next/navigation";
import { useCallback, useEffect, useRef, useState } from "react";
import { useSocket } from "@/providers/socket-provider";

const ChatPage = () => {
    const isOnline = useOnlineStatus();
    const params = useParams();
    const chatID = params.id;
    const loggedInUserID = useAuthStore((state) => state.userID);
    const router = useRouter();

    const [messages, setMessages] = useState<MessageType[]>([]);
    const messagesRef = useRef(messages);
    const pendingMessages = useRef(new Map<string, WebsocketMsg>());

    const [chat, setChat] = useState<ChatType>();
    const [chatUsers, setChatUsers] = useState<UserType[]>([]);
    const [permissions, setPermissions] = useState<PermissionType[]>([]);

    const messagesContainer = useRef<HTMLDivElement | null>(null);
    const containerRef = useRef<HTMLDivElement>(null);

    const [menuIsOpen, setMenuIsOpen] = useState(false);
    const [selectedMessageID, setSelectedMessageID] = useState<number>();
    const [isEditingMessage, setIsEditingMessage] = useState(false);
    const [editingMessageID, setEditingMessageID] = useState<number>();

    const socket = useSocket();

    const newMessageID = (): string => crypto.randomUUID();
    const newNumberID = (): number =>
        Math.floor(Math.random() * 100_000_000) + 1;

    useEffect(() => {
        messagesRef.current = messages;
    }, [messages]);

    useEffect(() => {
        if (!chatID) return;

        let cancelled = false;

        (async () => {
            const chatNum = +chatID;
            if (chatNum <= 0) return;

            const { messages } = await getChatMessages(chatNum);
            if (cancelled) return;
            setMessages(messages);

            pendingMessages.current.clear();

            const savedMessages = await messageStore.getByChat(chatNum);
            if (cancelled) return;

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

            const chat = await getChatByID(chatNum);
            if (cancelled) return;
            setChat(chat);

            const chatUsers = await getChatUsers(chatNum);
            if (cancelled) return;
            setChatUsers(chatUsers ?? []);
        })();

        return () => {
            cancelled = true;
        };
    }, [chatID]);

    useEffect(() => {
        (() => setPermissions([]))();

        if (!chatID) return;

        const chatNum = +chatID;
        if (chatNum <= 0) return;

        let cancelled = false;

        (async () => {
            const { permissions } = await getUserPermissions(chatNum);
            if (cancelled) return;
            setPermissions(permissions ?? []);
        })();

        return () => {
            cancelled = true;
        };
    }, [chatID]);

    const hasPermission = (permissionName: string): boolean => {
        return permissions.some(
            (permission) => permission.name === permissionName,
        );
    };

    const handleOpen = useCallback(async () => {
        if (!chatID) return;
        const chatNum = +chatID;
        if (chatNum <= 0) return;

        const lastAckedID = messagesRef.current.reduce(
            (max, message) => (message.ID > max ? message.ID : max),
            0,
        );
        if (lastAckedID > 0) {
            const { messages: missedMessages } = await syncChatMessages(
                chatNum,
                lastAckedID,
            );

            setMessages((prev) => {
                const seen = new Set(prev.map((m) => m.ID));
                const uniqueMissed = missedMessages.filter(
                    (m) => !seen.has(m.ID),
                );
                return [...prev, ...uniqueMissed];
            });
        }

        for (const [, pendingMessage] of pendingMessages.current) {
            if (pendingMessage.status !== "pending") {
                continue;
            }
            socket?.send(JSON.stringify(pendingMessage));
        }
    }, [chatID, socket]);

    const handleMessage = useCallback(
        (event: MessageEvent) => {
            const data = JSON.parse(event.data);

            if (data.type === "error") {
                console.error(data.message);
                return;
            }

            if (data.type === "messageDeleted") {
                setMessages((prev) =>
                    prev.map((msg) =>
                        msg.ID === data.messageID
                            ? {
                                  ...msg,
                                  Deleted: true,
                                  Content: data.content,
                              }
                            : msg,
                    ),
                );
                return;
            }

            if (data.type === "messageEdited") {
                setMessages((prev) =>
                    prev.map((msg) =>
                        msg.ID === data.messageID
                            ? {
                                  ...msg,
                                  Content: data.newContent,
                                  UpdatedAt: data.updatedAt,
                              }
                            : msg,
                    ),
                );
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
                            ? {
                                  ...message,
                                  Status: "delivered",
                                  ID: data.id,
                              }
                            : message,
                    ),
                );

                const pending = pendingMessages.current.get(data.clientID);
                if (pending) {
                    pendingMessages.current.set(data.clientID, {
                        ...pending,
                        status: "delivered",
                    });

                    messageStore.update({
                        ...pending,
                        status: "delivered",
                    });
                }

                return;
            }

            const incoming = data as MessageType;

            if (pendingMessages.current.has(incoming.ClientID)) {
                pendingMessages.current.delete(incoming.ClientID);
                messageStore.delete(incoming.ClientID);
                return;
            }

            setMessages((prev) => {
                if (!chatID) return prev;
                if (incoming.ChatID !== +chatID) return prev;
                return [...prev, incoming];
            });
        },
        [chatID],
    );

    useEffect(() => {
        if (!socket) return;
        if (socket.readyState === WebSocket.OPEN) {
            handleOpen();
        }

        socket?.addEventListener("open", handleOpen);
        socket?.addEventListener("message", handleMessage);
        return () => {
            socket?.removeEventListener("open", handleOpen);
            socket?.removeEventListener("message", handleMessage);
        };
    }, [socket, handleMessage, handleOpen]);

    const sendMessage = (message: string) => {
        if (!chatID || message.trim() === "" || !loggedInUserID) return;

        if (!socket || socket.readyState !== WebSocket.OPEN) return;

        const clientID = newMessageID();
        const creationDate = new Date().toISOString();

        const msg: WebsocketMsg = {
            status: "pending",
            chatID: +chatID,
            senderID: loggedInUserID,
            clientID,
            content: message,
            type: "message",
            CreatedAt: creationDate,
        };

        setMessages((prev) => {
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
        messagesContainer.current?.scrollTo({
            top: messagesContainer.current.scrollHeight,
            behavior: "smooth",
        });
    }, [messages]);

    return (
        <div
            ref={containerRef}
            className="flex-1 min-h-0 flex flex-col gap-y-[8px] overflow-x-clip"
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
                className={`${isOnline ? "opacity-0 invisible" : "opacity-100"} duration-300 right-1/2 translate-x-1/2 bg-red-500 p-[4px] rounded-[4px] absolute`}
            >
                <span className="text-[16px]">Currently offline</span>
            </div>

            <div className="flex">
                <div className="flex h-[32px] gap-x-[8px] items-center min-[900px]:hidden">
                    <BackArrowButton
                        handleClick={() => router.back()}
                        text="Chat"
                    />
                </div>

                <div className="flex-1 flex">
                    {chatID && (
                        <Menu
                            chatID={+chatID}
                            currentGroupUsers={chatUsers.map((user) => user.id)}
                            hasPermission={hasPermission}
                        />
                    )}
                </div>
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
                            hasPermission={hasPermission}
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
