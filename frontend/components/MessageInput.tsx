"use client";

import { useState } from "react";

interface Props {
    handleSend: (message: string) => void;
}

export type WebsocketMsg = {
    status?: "pending" | "delivered" | "failed";
    clientID: string;
    senderID: number;
    chatID: number;
    type: string;
    content: string;
    CreatedAt?: string;
};

export interface MessageType {
    ClientID?: string;
    ID: number;
    ChatID: number;
    SenderID: number;
    Content: string;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt: string | null;
    Deleted: boolean;
    Status: "pending" | "delivered" | "failed";
}

const MessageInput = ({ handleSend }: Props) => {
    const [message, setMessage] = useState("");

    return (
        <div className="flex w-full mt-auto relative">
            <textarea
                placeholder="Type a message here..."
                className="bg-foreground w-full text-background outline-none border-none p-2 rounded-[4px] min-h-[32px] max-h-[128px] resize-none overflow-hidden leading-tight"
                value={message}
                onInput={(e) => {
                    const el = e.currentTarget;

                    el.style.height = "auto";
                    el.style.height = Math.min(el.scrollHeight, 128) + "px";
                }}
                onChange={(e) => {
                    setMessage(e.target.value);
                }}
                onKeyDown={(e) => {
                    if (e.key === "Enter" && !e.shiftKey) {
                        e.preventDefault(); // prevents newline
                        handleSend(message);
                        setMessage("");
                    }
                }}
            />
            <svg
                viewBox="0 0 24 24"
                fill="currentColor"
                xmlns="http://www.w3.org/2000/svg"
                className={`${message.trim() === "" ? "opacity-0 z-[-1]" : ""} transition-all duration-300 text-accent absolute right-0 top-0 h-full max-h-[32px] cursor-pointer`}
                role="button"
                aria-disabled={message === ""}
                onClick={() => {
                    handleSend(message);
                    setMessage("");
                }}
            >
                <g id="SVGRepo_bgCarrier" strokeWidth="0"></g>
                <g
                    id="SVGRepo_tracerCarrier"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                ></g>
                <g id="SVGRepo_iconCarrier">
                    <path
                        d="M11.5003 12H5.41872M5.24634 12.7972L4.24158 15.7986C3.69128 17.4424 3.41613 18.2643 3.61359 18.7704C3.78506 19.21 4.15335 19.5432 4.6078 19.6701C5.13111 19.8161 5.92151 19.4604 7.50231 18.7491L17.6367 14.1886C19.1797 13.4942 19.9512 13.1471 20.1896 12.6648C20.3968 12.2458 20.3968 11.7541 20.1896 11.3351C19.9512 10.8529 19.1797 10.5057 17.6367 9.81135L7.48483 5.24303C5.90879 4.53382 5.12078 4.17921 4.59799 4.32468C4.14397 4.45101 3.77572 4.78336 3.60365 5.22209C3.40551 5.72728 3.67772 6.54741 4.22215 8.18767L5.24829 11.2793C5.34179 11.561 5.38855 11.7019 5.407 11.8459C5.42338 11.9738 5.42321 12.1032 5.40651 12.231C5.38768 12.375 5.34057 12.5157 5.24634 12.7972Z"
                        stroke="#000000"
                        strokeWidth="1"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                    ></path>
                </g>
            </svg>
        </div>
    );
};

export default MessageInput;
