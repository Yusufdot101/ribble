"use client";
import { useAuthStore } from "@/store/useAuthStore";
import { ChatType, getChatUsers } from "@/utils/chats";
import { UserType } from "@/utils/users";
import { useEffect, useState } from "react";

interface Props {
    activeChats: number[];
    chat: ChatType;
    handleClick: (chatId: number) => void;
}

const ChatCard = ({ activeChats, chat, handleClick }: Props) => {
    const [chatUser, setChatUser] = useState<UserType>();
    const loggedInUserId = useAuthStore((state) => state.userId);

    useEffect(() => {
        let cancelled = false;
        (async () => {
            const chatUsers = await getChatUsers(+chat.id);
            if (!chatUsers) return;
            if (cancelled || !chatUsers || chatUsers.length === 0) return;
            setChatUser(
                chatUsers[0].id === loggedInUserId
                    ? (chatUsers[1] ?? chatUsers[0])
                    : chatUsers[0],
            );
        })();
        return () => {
            cancelled = true;
        };
    }, [chat.id, loggedInUserId]);

    return (
        <div
            tabIndex={0}
            onClick={() => handleClick(chat.id)}
            className={`${activeChats?.includes(chat.id) ? "bg-foreground/20" : "hover:bg-foreground/10 focus:bg-foreground/10"} border-foreground p-[4px] cursor-pointer duration-300 h-[64px]`}
            onKeyDown={(e) => {
                if (e.key === "Enter") {
                    handleClick(chat.id);
                }
            }}
        >
            {chat.name ? (
                <p className="min-[620px]:text-[20px]">{chat.name}</p>
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
