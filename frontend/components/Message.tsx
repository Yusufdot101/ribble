"use client";
import { useAuthStore } from "@/store/useAuthStore";
import { MessageType } from "@/utils/messages";

interface Props {
    message: MessageType;
}

const Message = ({ message }: Props) => {
    const userID = useAuthStore((state) => state.userID);
    const date = new Date(message.CreatedAt);
    const formattedDate = new Intl.DateTimeFormat("en-GB", {
        year: "numeric",
        month: "long",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
        hour12: false,
    }).format(date);
    return (
        <div
            className={`${message.SenderID === userID ? "bg-accent/80 ml-auto" : "bg-foreground/20 mr-auto"} py-[4px] px-[8px] rounded-[4px] w-fit flex flex-col`}
        >
            {message.Content}
            <span className="text-[12px] opacity-75 text-right">
                {formattedDate}
            </span>
        </div>
    );
};

export default Message;
