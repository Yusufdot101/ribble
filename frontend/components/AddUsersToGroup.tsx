"use client";
import Contacts from "@/components/Contacts";
import { getAddableChatUsers, UserType } from "@/utils/users";
import { useCallback, useEffect, useState } from "react";
import BackArrowButton from "./BackArrowButton";
import XButton from "./XButton";
import SearchBar from "./SearchBar";
import { addUsersToGroup } from "@/utils/groups";
import { useSocket } from "@/providers/socket-provider";
import { useAuthStore } from "@/store/useAuthStore";

interface Props {
    handleClose: () => void;
    addToGroupIsOpen: boolean;
    chatId: number;
}

const AddUsersToGroup = ({ handleClose, addToGroupIsOpen, chatId }: Props) => {
    const [selectedUsers, setSelectedUsers] = useState<UserType[]>([]);
    const handleClick = async (clickedUser: UserType) => {
        setSelectedUsers((prev) => {
            return prev.includes(clickedUser)
                ? prev.filter((user) => user !== clickedUser)
                : [...prev, clickedUser];
        });
    };
    const removeUser = (userId: number) => {
        setSelectedUsers((prev) => {
            return prev.filter((user) => user.id !== userId);
        });
    };

    const [isLoading, setIsLoading] = useState(false);

    const [users, setUsers] = useState<UserType[]>([]);
    const searchUsers = useCallback(
        async (email: string = "") => {
            setIsLoading(true);
            try {
                const { users } = await getAddableChatUsers(chatId, email);
                setUsers(users ?? []);
            } finally {
                setIsLoading(false);
            }
        },
        [chatId],
    );

    useEffect(() => {
        (() => searchUsers())();
    }, [searchUsers]);

    const addToGroup = async () => {
        if (selectedUsers.length === 0) return;
        await addUsersToGroup(
            chatId,
            selectedUsers.map((user) => user.id),
        );
        handleClose();
    };

    const socket = useSocket();

    const loggedInUserId = useAuthStore((state) => state.userId);

    const handleMessage = useCallback(
        (event: MessageEvent) => {
            const data = JSON.parse(event.data);

            if (data.type === "error") {
                console.error(data.message);
                return;
            }

            const payload = data.payload;
            if (data.subType === "usersAdded") {
                const users = payload.addedUsers as UserType[];
                setUsers((prev) => {
                    return prev.filter(
                        (user) => !users.some((u) => u.id === user.id),
                    );
                });
                setSelectedUsers((prev) => {
                    return prev.filter(
                        (user) => !users.some((u) => u.id === user.id),
                    );
                });
            }

            if (["userBanned"].includes(data.subType)) {
                const bannedUser = payload.target as UserType;
                setUsers((prev) => {
                    return prev.filter((user) => user.id !== bannedUser.id);
                });
            }

            if (["userRemoved"].includes(data.subType)) {
                const removedUser = payload.target as UserType;
                if (removedUser.id === loggedInUserId) return;
                setUsers((prev) => {
                    return [...prev, removedUser];
                });
            }

            if (["userUnbanned"].includes(data.subType)) {
                const unbannedUser = payload.target as UserType;
                if (unbannedUser.id === loggedInUserId) return;
                setUsers((prev) => {
                    return [...prev, unbannedUser];
                });
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

    return (
        <div
            className={`${addToGroupIsOpen ? "translate-x-0" : "translate-x-full"} transition-transform absolute w-full overflow-y-auto bg-background z-10 duration-300 flex-1 flex overflow-x-hidden`}
        >
            <div className="h-full transition-transform duration-300 ease-in-out flex flex-1 flex-col gap-y-[8px]">
                <div className="flex w-full h-[32px] gap-x-[8px] items-center">
                    <BackArrowButton
                        handleClick={handleClose}
                        text="Add group members"
                    />
                </div>

                <div className="flex flex-col gap-y-[8px] relative flex-1">
                    {selectedUsers?.length !== 0 && (
                        <div className="flex flex-col space-y-[8px] max-h-[180px] overflow-y-scroll border-b-1 border-foreground transition-all duration-300 ease-in-out">
                            {selectedUsers.map((user) => (
                                <div
                                    key={user.id}
                                    className="w-full flex justify-between items-center transition-all duration-300 ease-in-out"
                                >
                                    <div className="flex flex-col">
                                        <p className="min-[620px]:text-[20px]">
                                            {user.name} :
                                        </p>
                                        <p className="min-[620px]:text-[16px]">
                                            {user.email}
                                        </p>
                                    </div>
                                    <XButton
                                        handleClick={() => removeUser(user.id)}
                                    />
                                </div>
                            ))}
                        </div>
                    )}

                    <SearchBar
                        placeholder="Search group members"
                        handleEnter={searchUsers}
                    />

                    <div className="h-full max-h-[180px] overflow-y-scroll">
                        <Contacts
                            isLoading={isLoading}
                            users={users}
                            handleUserClick={handleClick}
                            selectedUsers={selectedUsers.map((user) => user.id)}
                            excludeUsers={selectedUsers.map((user) => user.id)}
                        />
                    </div>

                    <button
                        aria-label="create group"
                        className="p-[4px] bg-accent text-white rounded-[4px] hover:bg-accent/80 active:bg-accent duration-300 cursor-pointer w-full"
                        onClick={addToGroup}
                    >
                        Add to group
                    </button>
                </div>
            </div>
        </div>
    );
};

export default AddUsersToGroup;
