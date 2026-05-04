"use client";
import Contacts from "@/components/Contacts";
import { getAddableChatUsers, UserType } from "@/utils/users";
import { useCallback, useEffect, useState } from "react";
import BackArrowButton from "./BackArrowButton";
import XButton from "./XButton";
import SearchBar from "./SearchBar";
import { addUsersToGroup } from "@/utils/groups";

interface Props {
    handleClose: () => void;
    addToGroupIsOpen: boolean;
    chatID: number;
}

const AddUsersToGroup = ({ handleClose, addToGroupIsOpen, chatID }: Props) => {
    const [selectedUsers, setSelectedUsers] = useState<UserType[]>([]);
    const handleClick = async (clickedUser: UserType) => {
        setSelectedUsers((prev) => {
            return prev.includes(clickedUser)
                ? prev.filter((user) => user !== clickedUser)
                : [...prev, clickedUser];
        });
    };
    const removeUser = (userID: number) => {
        setSelectedUsers((prev) => {
            return prev.filter((user) => user.id !== userID);
        });
    };

    const [isLoading, setIsLoading] = useState(false);

    const [users, setUsers] = useState<UserType[]>([]);
    const searchUsers = useCallback(
        async (email: string = "") => {
            setIsLoading(true);
            try {
                const users = await getAddableChatUsers(chatID, email);
                setUsers(users ?? []);
            } finally {
                setIsLoading(false);
            }
        },
        [chatID],
    );

    useEffect(() => {
        (() => searchUsers())();
    }, [searchUsers]);

    const addToGroup = async () => {
        if (selectedUsers.length === 0) return;
        await addUsersToGroup(
            chatID,
            selectedUsers.map((user) => user.id),
        );
        handleClose();
    };

    return (
        <div
            className={`${addToGroupIsOpen ? "translate-x-0" : "translate-x-full"} transition-transform absolute w-full bg-background z-10 top-1/2 translate-y-1/2 duration-300 flex-1 flex overflow-x-hidden`}
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
