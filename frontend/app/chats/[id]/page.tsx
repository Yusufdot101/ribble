"use client";
import Message from "@/components/Message";
import MessageInput from "@/components/MessageInput";
import { getChatMessages, MessageType } from "@/utils/messages";
import { useParams } from "next/navigation";
import { useEffect, useState } from "react";

const ChatPage = () => {
    const params = useParams();
    const chatID = params.id;
    const [messages, setMessages] = useState<MessageType[]>([]);
    useEffect(() => {
        if (!chatID) return;
        (async () => {
            const messages = await getChatMessages(+chatID);
            setMessages(messages);
        })();
    }, [chatID]);

    return (
        <div className="flex flex-col h-full gap-y-[8px]">
            <div className="flex justify-center">Username/email</div>

            <div className="flex flex-col gap-y-[8px]">
                {messages?.map((message) => (
                    <Message key={message.ID} message={message} />
                ))}
            </div>

            <MessageInput />
        </div>
    );
};

export default ChatPage;
