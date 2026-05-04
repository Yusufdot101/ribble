"use client";
import { useState } from "react";
import AddUsersToGroup from "./AddUsersToGroup";
import GroupMembers from "./GroupMembers";

type Props = {
    chatID: number;
    currentGroupUsers: number[];
    hasPermission: (permissionName: string) => boolean;
};

const Menu = ({ chatID, hasPermission }: Props) => {
    const [menuIsOpen, setMenuIsOpen] = useState(false);

    const [addToGroupIsOpen, setAddToGroupIsOpen] = useState(false);
    const handleClose = () => {
        setAddToGroupIsOpen(false);
    };

    const [groupMembersIsOpen, setGroupMembersIsOpen] = useState(false);
    const handleCloseGroupMembers = () => {
        setGroupMembersIsOpen(false);
    };

    return (
        <div className="relative flex flex-col max-h-screen flex-1">
            <button
                onClick={() => setMenuIsOpen((prev) => !prev)}
                className="cursor-pointer ml-auto hover:text-accent duration-300 active:text-foreground"
            >
                Menu
            </button>

            <div
                className={`${menuIsOpen ? "max-h-96 p-[4px]" : "max-h-0 invisible p-0"} z-10 duration-300 w-50 border-1 border-foreground rounded-[4px] flex flex-col bg-background gap-y-[4px] absolute right-0 overflow-hidden top-[28px]`}
            >
                {hasPermission("add users to group") && (
                    <button
                        onClick={(e) => {
                            e.stopPropagation();
                            setGroupMembersIsOpen(false);
                            setAddToGroupIsOpen((prev) => !prev);
                            setMenuIsOpen(false);
                        }}
                        className="cursor-pointer hover:bg-foreground/20 active:bg-background duration-300"
                    >
                        Add member
                    </button>
                )}
                <button
                    onClick={(e) => {
                        e.stopPropagation();
                        setAddToGroupIsOpen(false);
                        setGroupMembersIsOpen((prev) => !prev);
                        setMenuIsOpen(false);
                    }}
                    className="cursor-pointer hover:bg-foreground/20 active:bg-background duration-300"
                >
                    Group members
                </button>
                <button
                    onClick={(e) => {
                        e.stopPropagation();
                    }}
                    className="cursor-pointer hover:bg-foreground/20 active:bg-background duration-300"
                >
                    {/*TODO: implement delete chat functionality*/}
                    Delete Chat
                </button>
            </div>

            <AddUsersToGroup
                addToGroupIsOpen={addToGroupIsOpen}
                handleClose={handleClose}
                chatID={chatID}
            />

            <GroupMembers
                groupMembersIsOpen={groupMembersIsOpen}
                handleClose={handleCloseGroupMembers}
                chatID={chatID}
                hasPermission={hasPermission}
            />
        </div>
    );
};

export default Menu;
