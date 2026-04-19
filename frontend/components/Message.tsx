"use client";
import { useAuthStore } from "@/store/useAuthStore";
import { deleteMessage, MessageType } from "@/utils/messages";
import { flip, offset, shift, useFloating } from "@floating-ui/react";
import { RefObject } from "react";

interface Props {
    message: MessageType;
    menuIsOpen: boolean;
    selectedMessageID: number;
    handleRightClick: (messageID: number) => void;
    containerRef: RefObject<HTMLDivElement | null>;
}

const Message = ({
    message,
    menuIsOpen,
    selectedMessageID,
    handleRightClick,
    containerRef,
}: Props) => {
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

    const { refs, floatingStyles } = useFloating({
        open: menuIsOpen,
        middleware: [offset(), flip(), shift()],
        placement: "right-start",
    });

    const handleDelete = () => {
        if (!confirm("Do you want delete this message? ")) return;
        deleteMessage(message.ID);
    };

    const rect = containerRef.current?.getBoundingClientRect();
    if (!rect) return;

    return (
        <div
            tabIndex={0}
            className={`${message.SenderID === userID ? "bg-accent/80 ml-auto" : "bg-foreground/20 mr-auto"} py-[4px] px-[8px] rounded-[4px] w-fit flex flex-col`}
            onKeyDown={(e) => {
                if (e.key !== "Delete") return;
                e.preventDefault();

                if (message.SenderID !== userID) return;
                handleRightClick(message.ID);
                e.preventDefault();

                refs.setPositionReference({
                    getBoundingClientRect() {
                        return {
                            width: 0,
                            height: 0,

                            x: rect.left + rect.width / 2,
                            y: rect.top + rect.height / 2,

                            top: rect.top + rect.height / 2,
                            left: rect.left + rect.width / 2,
                            right: rect.left + rect.width / 2,
                            bottom: rect.top + rect.height / 2,
                        };
                    },
                });
            }}
            onContextMenu={(e) => {
                if (message.SenderID !== userID) return;
                handleRightClick(message.ID);
                e.preventDefault();

                refs.setPositionReference({
                    getBoundingClientRect() {
                        return {
                            width: 0,
                            height: 0,
                            x: e.clientX,
                            y: e.clientY,
                            top: e.clientY,
                            left: e.clientX,
                            right: e.clientX,
                            bottom: e.clientY,
                        };
                    },
                });
            }}
        >
            {message.Content}
            <span className="text-[12px] opacity-75 text-right">
                {formattedDate}
            </span>

            {menuIsOpen && selectedMessageID === message.ID && (
                <div
                    ref={refs.setFloating ?? undefined}
                    style={{
                        ...floatingStyles,
                    }}
                    className="absolute bg-foreground p-[4px] rounded-[4px] text-background shadow-lg z-1"
                >
                    <button
                        aria-label="delete message"
                        className="hover:bg-red-500 hover:text-foreground cursor-pointer duration-300 rounded-[2px] p-[4px]"
                        onClick={handleDelete}
                    >
                        Delete
                    </button>
                </div>
            )}
        </div>
    );
};

export default Message;
