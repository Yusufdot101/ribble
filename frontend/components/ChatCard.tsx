"use client";
import { useAuthStore } from "@/store/useAuthStore";
import { ChatType, getChatUsers } from "@/utils/chats";
import { UserType } from "@/utils/users";
import { useEffect, useState } from "react";

interface Props {
    activeChats: number[];
    chat: ChatType;
    handleClick: (chatID: number) => void;
}

const ChatCard = ({ activeChats, chat, handleClick }: Props) => {
    const [chatUser, setChatUser] = useState<UserType>();
    const loggedInUserID = useAuthStore((state) => state.userID);
    useEffect(() => {
        (async () => {
            const chatUsers = await getChatUsers(+chat.ID);
            if (!chatUsers) return;
            setChatUser(
                chatUsers[0].id === loggedInUserID
                    ? chatUsers[1]
                    : chatUsers[0],
            );
        })();
    }, [chat, loggedInUserID]);

    return (
        <div
            tabIndex={0}
            onClick={() => handleClick(chat.ID)}
            className={`${activeChats?.includes(chat.ID) ? "bg-foreground/20" : ""} border-foreground p-[4px] cursor-pointer duration-300 h-[64px]`}
            onKeyDown={(e) => {
                if (e.key === "Enter") {
                    handleClick(chat.ID);
                }
            }}
        >
            {chat.Name ? (
                <p className="min-[620px]:text-[20px]">{chat.Name}</p>
            ) : (
                <>
                    <p className="min-[620px]:text-[20px]">{chatUser?.name}</p>
                    <p className="min-[620px]:text-[16px]">{chatUser?.email}</p>
                </>
            )}
        </div>
    );
};

export default ChatCard;
