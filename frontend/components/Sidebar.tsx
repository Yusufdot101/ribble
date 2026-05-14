"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import {
    ChatType,
    ConversationDataType,
    getChatByUserIDs,
    getConversations,
} from "@/utils/chats";
import CreateGroupButton from "./CreateGroupButton";
import CreateGroup from "./CreateGroup";
import SearchBar from "./SearchBar";
import ContactsSection from "./ContactsSection";
import ChatsSection from "./ChatsSection";
import GroupsSection from "./GroupsSection";
import { UserType } from "@/utils/users";

const Sidebar = () => {
    const [activeUser, setActiveUser] = useState<number>();
    const [activeChat, setActiveChat] = useState<number>();
    const router = useRouter();

    const handleChatClick = async (chatID: number) => {
        setActiveChat(chatID);
        setActiveUser(-1);
        router.push(`/chats/${chatID}`);
    };

    const [isLoading, setIsLoading] = useState(false);

    const [coversationData, setConverastionData] =
        useState<ConversationDataType>();

    const updateGroups = (chat: ChatType) => {
        setConverastionData((prev) => {
            return prev
                ? {
                      ...prev,
                      groups: [...(prev.groups ?? []), chat],
                  }
                : prev;
        });
        setActiveChat(chat.id);
    };

    const updateChats = (chat: ChatType, user: UserType) => {
        setConverastionData((prev) => {
            return prev
                ? {
                      ...prev,
                      chats: [...(prev.chats ?? []), chat],
                      contacts: [
                          ...prev.contacts.filter(
                              (contact) => contact.id !== user.id,
                          ),
                      ],
                  }
                : prev;
        });
        setActiveChat(chat.id);
    };

    const handleUserClick = async (user: UserType) => {
        setActiveUser(user.id);
        setActiveChat(-1);
        const chat = await getChatByUserIDs([user.id]);
        if (!chat) return;
        setActiveChat(chat.id);
        updateChats(chat, user);
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
                    updateGroups(chat);
                }}
                createGroupOpen={isCreatingGroup}
            />
        </div>
    );
};

export default Sidebar;
