"use client";
import { useAuthStore } from "@/store/useAuthStore";
import { MessageType } from "@/utils/messages";

interface Props {
    message: MessageType;
}

const Message = ({ message }: Props) => {
    const userID = useAuthStore((state) => state.userID);
    return (
        <div
            className={`${message.SenderID === userID ? "bg-accent/80 ml-auto" : "bg-foreground/20 mr-auto"} py-[4px] px-[8px] rounded-[4px] w-fit`}
        >
            {message.Content}
        </div>
    );
};

export default Message;
