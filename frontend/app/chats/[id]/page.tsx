"use client";

import BackArrowButton from "@/components/BackArrowButton";
import Message from "@/components/Message";
import MessageInput, { outgoingMsg } from "@/components/MessageInput";
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
    const chatId = params.id;
    const loggedInUserId = useAuthStore((state) => state.userId);
    const router = useRouter();

    const [messages, setMessages] = useState<MessageType[]>([]);
    const messagesRef = useRef(messages);
    const pendingMessages = useRef(new Map<string, MessageType>());

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
        Math.floor(Math.random() * 1_000_000_000) + 1;

    useEffect(() => {
        messagesRef.current = messages;
    }, [messages]);

    useEffect(() => {
        if (!chatId) return;

        let cancelled = false;

        (async () => {
            const chatNum = +chatId;
            if (chatNum <= 0) return;

            const { messages } = await getChatMessages(chatNum);
            if (cancelled) return;
            setMessages(messages);

            pendingMessages.current.clear();

            const savedMessages = await messageStore.getByChat(chatNum);
            if (cancelled) return;
            for (const msg of savedMessages) {
                if (msg.status === "pending") {
                    pendingMessages.current.set(msg.clientId, msg);
                }
            }

            setMessages((prev) => {
                const existingIds = new Set(prev.map((m) => m.clientId));
                const uniqueSaved = savedMessages.filter(
                    (m) => !existingIds.has(m.clientId),
                );
                return [...prev, ...uniqueSaved];
            });

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
    }, [chatId]);

    const fetchPermissions = useCallback(
        async (cancelled: boolean = false) => {
            (() => setPermissions([]))();
            if (!chatId) return;
            const chatNum = +chatId;
            if (chatNum <= 0) return;
            if (cancelled) return;
            const { permissions } = await getUserPermissions(chatNum);
            (() => setPermissions(permissions ?? []))();
        },
        [chatId],
    );
    useEffect(() => {
        let cancelled = false;
        fetchPermissions(cancelled);
        return () => {
            cancelled = true;
        };
    }, [fetchPermissions]);

    const hasPermission = (permissionName: string): boolean => {
        return permissions.some(
            (permission) => permission.name === permissionName,
        );
    };

    const handleOpen = useCallback(async () => {
        if (!chatId) return;
        const chatNum = +chatId;
        if (chatNum <= 0) return;

        const lastAckedID = messagesRef.current.reduce(
            (max, message) => (message.id > max ? message.id : max),
            0,
        );
        if (lastAckedID > 0) {
            const { messages: missedMessages } = await syncChatMessages(
                chatNum,
                lastAckedID,
            );

            setMessages((prev) => {
                const seen = new Set(prev.map((m) => m.id));
                const uniqueMissed = missedMessages.filter(
                    (m) => !seen.has(m.id),
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
    }, [chatId, socket]);

    const handleMessage = useCallback(
        (event: MessageEvent) => {
            const data = JSON.parse(event.data);
            if (
                ![
                    "message",
                    "ack",
                    "nack",
                    "usersAdded",
                    "userRemoved",
                    "userBanned",
                ].includes(data.type)
            )
                return;

            const payload = data.payload;
            if (data.subType === "usersAdded") {
                const users = payload.addedUsers as UserType[];
                setChatUsers((prev) => [
                    ...prev,
                    ...users.filter((user) => user.id !== loggedInUserId),
                ]);
            }

            if (["userBanned", "userRemoved"].includes(data.subType)) {
                const removedUser = payload.target as UserType;
                setChatUsers((prev) =>
                    prev.filter((user) => user.id !== removedUser.id),
                );
            }

            if (data.subType === "messageDeleted") {
                setMessages((prev) =>
                    prev.map((msg) =>
                        msg.id === payload.id
                            ? {
                                  ...msg,
                                  deleted: true,
                                  content: payload.content,
                              }
                            : msg,
                    ),
                );
                return;
            }

            if (data.subType === "messageEdited") {
                setMessages((prev) =>
                    prev.map((msg) =>
                        msg.id === payload.id
                            ? {
                                  ...msg,
                                  content: payload.content,
                                  updatedAt: payload.updatedAt,
                              }
                            : msg,
                    ),
                );
                return;
            }

            if (data.type === "nack") {
                const clientId = payload.clientId;
                setMessages((prev) =>
                    prev.map((message) =>
                        message.clientId === clientId
                            ? { ...message, status: "failed" }
                            : message,
                    ),
                );

                const msg = pendingMessages.current.get(clientId);
                if (msg) {
                    messageStore.update({ ...msg, status: "failed" });
                }

                pendingMessages.current.delete(clientId);
                return;
            }

            if (data.type === "ack") {
                const clientId = payload.clientId;
                setMessages((prev) => {
                    return prev.map((message) =>
                        message.clientId === clientId
                            ? {
                                  ...message,
                                  status: "delivered",
                                  id: data.payload.id,
                              }
                            : message,
                    );
                });

                const pending = pendingMessages.current.get(clientId);
                if (pending) {
                    pendingMessages.current.set(clientId, {
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

            if (!payload.message) return;
            const incoming = payload.message as MessageType;

            if (pendingMessages.current.has(payload.clientId)) {
                pendingMessages.current.delete(payload.clientId);
                messageStore.delete(payload.clientId);
                return;
            }

            setMessages((prev) => {
                if (!chatId) return prev;
                if (incoming.chatId !== +chatId) return prev;
                return [...prev, incoming];
            });
        },
        [chatId, loggedInUserId],
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
        if (!chatId || message.trim() === "" || !loggedInUserId) return;

        if (!socket || socket.readyState !== WebSocket.OPEN) return;

        const clientId = newMessageID();
        const creationDate = new Date().toISOString();

        const msg: MessageType = {
            status: "pending",
            chatId: +chatId,
            senderId: loggedInUserId,
            clientId: clientId,
            content: message,
            createdAt: creationDate,
            deleted: false,
            id: -newNumberID(),
            updatedAt: creationDate,
            deletedAt: null,
        };

        setMessages((prev) => {
            return [...prev, msg];
        });

        pendingMessages.current.set(clientId, msg);
        messageStore.add(msg);

        const websocketMsg: outgoingMsg = {
            type: "newMessage",
            payload: msg,
        };
        socket.send(JSON.stringify(websocketMsg));
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
            className="flex-1 min-h-0 flex flex-col gap-y-[8px] overflow-x-clip relative"
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
                <div className="flex h-[32px] gap-x-[8px] items-center min-[900px]:hidden z-10">
                    <BackArrowButton
                        handleClick={() => router.back()}
                        text="Chat"
                    />
                </div>

                <div className="absolute w-full flex-1 flex">
                    {chatId && (
                        <Menu
                            chatId={+chatId}
                            currentGroupUsers={chatUsers.map((user) => user.id)}
                            hasPermission={hasPermission}
                            reloadPermissions={() => fetchPermissions()}
                        />
                    )}
                </div>
            </div>

            <div className="flex justify-center shrink-0">
                <div className="flex gap-x-[4px]">
                    {chat?.name !== ""
                        ? chat?.name
                        : chatUsers
                              .filter((user) => user.id !== loggedInUserId)
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
                    .sort((a, b) => a.createdAt.localeCompare(b.createdAt))
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
                            key={message.id}
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
                                              user.id === message.senderId,
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
