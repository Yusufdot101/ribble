"use client";
import { useAuthStore } from "@/store/useAuthStore";
import { deleteMessage, editMessage, MessageType } from "@/utils/messages";
import { flip, offset, shift, useFloating } from "@floating-ui/react";
import { RefObject, useEffect, useState } from "react";

interface Props {
    message: MessageType;
    menuIsOpen: boolean;
    selectedMessageID: number;
    handleRightClick: (messageID: number) => void;
    containerRef: RefObject<HTMLDivElement | null>;
    isEditing: boolean;
    editingMessageID: number | undefined;
    handleClickEdit: (messageID: number) => void;
    handleCancelMessageEdit: () => void;
}

const Message = ({
    message,
    menuIsOpen,
    selectedMessageID,
    handleRightClick,
    containerRef,
    isEditing,
    handleClickEdit,
    editingMessageID,
    handleCancelMessageEdit,
}: Props) => {
    const userID = useAuthStore((state) => state.userID);
    const messageIsEdited = message.CreatedAt !== message.UpdatedAt;
    const date = new Date(
        messageIsEdited ? message.UpdatedAt : message.CreatedAt,
    );
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

    const [newContent, setNewContent] = useState(message.Content);
    useEffect(() => {
        if (isEditing) {
            (() => setNewContent(message.Content))();
        }
    }, [isEditing, message.Content]);
    const isEditingCurrentMessage =
        isEditing && editingMessageID === message.ID;

    const createdDate = new Date(message.CreatedAt);
    const currentDate = new Date();
    const diff = currentDate.getTime() - createdDate.getTime();
    const ONE_HOUR = 60 * 60 * 1000;
    const createdlessThanHourAgo = diff < ONE_HOUR;
    return (
        <div
            className={`${isEditingCurrentMessage ? "w-full" : ""} ${message.SenderID === userID ? "ml-auto" : "mr-auto"} flex flex-col rounded-[4px]`}
            onKeyDown={(e) => {
                if (e.key === "Escape") {
                    handleCancelMessageEdit();
                    return;
                }
                if (e.key !== "Delete") return;
                e.preventDefault();

                if (message.SenderID !== userID) return;
                handleRightClick(message.ID);

                const rect = containerRef.current?.getBoundingClientRect();
                if (!rect) return;

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
            <div
                tabIndex={0}
                hidden={isEditingCurrentMessage}
                className={`${message.SenderID === userID ? "bg-accent/80" : "bg-foreground/20"} flex flex-col w-fit py-[4px] px-[8px] rounded-[4px]`}
            >
                {message.Content}
                <span className="text-[12px] opacity-75 text-right">
                    {messageIsEdited ? "Edited " : ""}
                    {formattedDate}
                </span>

                {menuIsOpen && selectedMessageID === message.ID && (
                    <div
                        ref={refs.setFloating ?? undefined}
                        style={{
                            ...floatingStyles,
                        }}
                        className="absolute bg-foreground p-[4px] rounded-[4px] text-background shadow-lg z-1 flex flex-col"
                    >
                        <button
                            aria-label="delete message"
                            className="hover:bg-red-500 hover:text-foreground cursor-pointer duration-300 rounded-[2px] p-[4px]"
                            onClick={handleDelete}
                        >
                            Delete
                        </button>

                        {new Date(message.CreatedAt) -
                            new Date(message.CreatedAt) <
                        60 * 60 * 60
                            ? ""
                            : ""}
                        {createdlessThanHourAgo ? (
                            <button
                                aria-label="edit message"
                                className="hover:bg-blue-500 hover:text-foreground cursor-pointer duration-300 rounded-[2px] p-[4px]"
                                onClick={() => {
                                    handleClickEdit(message.ID);
                                }}
                            >
                                EDIT
                            </button>
                        ) : undefined}
                    </div>
                )}
            </div>

            {isEditingCurrentMessage && (
                <div className="w-full flex flex-col gap-y-[4px]">
                    <textarea
                        className="bg-foreground/90 w-full text-background outline-none border-none p-2 rounded-[4px] min-h-[32px] max-h-[128px] resize-none overflow-hidden leading-tight"
                        value={newContent}
                        onInput={(e) => {
                            const el = e.currentTarget;

                            el.style.height = "auto";
                            el.style.height =
                                Math.min(el.scrollHeight, 128) + "px";
                        }}
                        onChange={(e) => {
                            setNewContent(e.target.value);
                        }}
                        onKeyDown={(e) => {
                            if (e.key === "Enter" && !e.shiftKey) {
                                e.preventDefault(); // prevents newline
                                handleCancelMessageEdit();
                                if (!createdlessThanHourAgo) return;
                                editMessage(message.ID, newContent);
                            }
                        }}
                    />

                    <div className="flex justify-end gap-x-[4px]">
                        <button
                            aria-label="cancel edit"
                            className="hover:bg-red-500 hover:text-foreground cursor-pointer duration-300 rounded-[2px] p-[4px]"
                            onClick={handleCancelMessageEdit}
                        >
                            Cancel
                        </button>
                        <button
                            aria-label="send edited message"
                            className="hover:bg-accent hover:text-foreground cursor-pointer duration-300 rounded-[2px] p-[4px]"
                            onClick={() => {
                                handleCancelMessageEdit();
                                if (!createdlessThanHourAgo) return;
                                editMessage(message.ID, newContent);
                            }}
                        >
                            Send
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Message;
