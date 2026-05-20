"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import {
    ChatType,
    ConversationDataType,
    getChatByUserIds,
    getConversations,
} from "@/utils/chats";
import CreateGroupButton from "./CreateGroupButton";
import CreateGroup from "./CreateGroup";
import SearchBar from "./SearchBar";
import ContactsSection from "./ContactsSection";
import ChatsSection from "./ChatsSection";
import GroupsSection from "./GroupsSection";
import { UserType } from "@/utils/users";
import { useSocket } from "@/providers/socket-provider";
import { useAuthStore } from "@/store/useAuthStore";

const Sidebar = () => {
    const [activeUser, setActiveUser] = useState<number>();
    const [activeChat, setActiveChat] = useState<number>();
    const router = useRouter();
    const [isLoading, setIsLoading] = useState(false);
    const loggedInUserId = useAuthStore((state) => state.userId);
    const [coversationData, setConverastionData] =
        useState<ConversationDataType>();

    const handleChatClick = async (chatId: number) => {
        setActiveChat(chatId);
        setActiveUser(-1);
        router.push(`/chats/${chatId}`);
    };

    const updateGroups = (chat: ChatType) => {
        setConverastionData((prev) => {
            return prev
                ? {
                      ...prev,
                      groups: [...(prev.groups ?? []), chat],
                  }
                : prev;
        });
    };

    const updateChats = (chat: ChatType, userIds: number[]) => {
        setConverastionData((prev) => {
            return prev
                ? {
                      ...prev,
                      chats: [...(prev.chats ?? []), chat],
                      contacts: prev.contacts && [
                          ...prev.contacts.filter(
                              (contact) => !userIds.includes(contact.id),
                          ),
                      ],
                  }
                : prev;
        });
    };

    const socket = useSocket();

    const handleMessage = useCallback(
        (event: MessageEvent) => {
            const data = JSON.parse(event.data);

            if (data.type === "error") {
                console.error(data.message);
                return;
            }

            const payload = data.payload;
            if (data.subType === "chatCreated") {
                const chat = payload.chat as ChatType;
                const userIds = payload.userIds as number[];
                if (chat.isGroup) {
                    updateGroups(chat);
                } else {
                    updateChats(chat, userIds);
                }
            }

            if (["userRemoved", "userBanned"].includes(data.subType)) {
                const removedUser = payload.target as UserType;
                if (removedUser.id === loggedInUserId) {
                    setConverastionData((prev) => {
                        return {
                            ...prev,
                            groups: prev
                                ? prev.groups.filter(
                                      (group) => group.id !== payload.chatId,
                                  )
                                : [],

                            chats: prev?.chats ?? [],
                            contacts: prev?.contacts ?? [],
                        };
                    });
                    return;
                }
            }

            if (data.subType === "usersAdded") {
                const addedUsers = payload.addedUsers as UserType[];
                if (addedUsers.some((user) => user.id === loggedInUserId)) {
                    setConverastionData((prev) => {
                        return {
                            ...prev,
                            groups: [...(prev?.groups ?? []), payload.chat],
                            chats: prev?.chats ?? [],
                            contacts: prev?.contacts ?? [],
                        };
                    });
                    return;
                }
            }
        },
        [loggedInUserId],
    );

    useEffect(() => {
        socket?.addEventListener("message", handleMessage);
        return () => {
            socket?.removeEventListener("message", handleMessage);
        };
    }, [socket, handleMessage]);

    const handleUserClick = async (user: UserType) => {
        setActiveUser(user.id);
        setActiveChat(-1);
        const chat = await getChatByUserIds([user.id]);
        if (!chat) return;
        setActiveChat(chat.id);
        router.push(`/chats/${chat.id}`);
    };

    const fetchData = async (query: string = "") => {
        setIsLoading(true);
        const data = await getConversations(query);
        setIsLoading(false);
        if (!data) return;
        setConverastionData(data);
    };

    useEffect(() => {
        const activeChat = window.location.pathname.split("/").at(-1);
        (() => {
            setActiveChat(activeChat ? +activeChat : -1);
            fetchData();
        })();
    }, []);

    const [isCreatingGroup, setIsCreatingGroup] = useState(false);
    return (
        <div className="flex-1 flex flex-col gap-y-[8px] relative overflow-hidden">
            <div
                className={`${isCreatingGroup ? "-translate-x-full invisible" : "translate-x-0"} h-full transition-transform duration-300 ease-in-out flex flex-col gap-y-[8px]`}
            >
                <CreateGroupButton
                    handleClick={() => setIsCreatingGroup(true)}
                />

                <SearchBar handleEnter={(query: string) => fetchData(query)} />

                <div className="flex-1 min-[900px]:border-r-1 border-foreground flex flex-col gap-y-[8px] overflow-auto">
                    <GroupsSection
                        isLoading={isLoading}
                        selectedChats={activeChat ? [activeChat] : []}
                        chats={coversationData?.groups ?? []}
                        handleChatClick={handleChatClick}
                    />

                    <ChatsSection
                        isLoading={isLoading}
                        selectedChats={activeChat ? [activeChat] : []}
                        chats={coversationData?.chats ?? []}
                        handleChatClick={handleChatClick}
                    />

                    <ContactsSection
                        selectedUsers={activeUser ? [activeUser] : []}
                        handleUserClick={handleUserClick}
                        isLoading={isLoading}
                        users={coversationData?.contacts ?? []}
                    />
                </div>
            </div>

            <CreateGroup
                handleClose={() => setIsCreatingGroup(false)}
                handleUpdateGroups={(chat: ChatType) => {
                    setActiveChat(chat.id);
                }}
                createGroupOpen={isCreatingGroup}
            />
        </div>
    );
};

export default Sidebar;
