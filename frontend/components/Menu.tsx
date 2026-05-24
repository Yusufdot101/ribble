"use client";
import { useState } from "react";
import AddUsersToGroup from "./AddUsersToGroup";
import GroupMembers from "./GroupMembers";
import GroupRolesPermissions from "./GroupRolesPermissions";

type Props = {
    chatId: number;
    isGroup: boolean;
    currentGroupUsers: number[];
    hasPermission: (permissionName: string) => boolean;
    reloadPermissions: () => void;
};

const Menu = ({ chatId, hasPermission, reloadPermissions, isGroup }: Props) => {
    const [menuIsOpen, setMenuIsOpen] = useState(false);

    const [addToGroupIsOpen, setAddToGroupIsOpen] = useState(false);
    const handleAddToGroupClose = () => {
        setAddToGroupIsOpen(false);
    };
    const handleAddToGroupOpen = () => {
        setMenuIsOpen(false);
        setGroupMembersIsOpen(false);
        setGroupPermissionsIsOpen(false);
        setAddToGroupIsOpen((prev) => !prev);
    };

    const [groupMembersIsOpen, setGroupMembersIsOpen] = useState(false);
    const handleCloseGroupMembers = () => {
        setGroupMembersIsOpen(false);
    };
    const handleGroupMembersOpen = () => {
        setMenuIsOpen(false);
        setGroupPermissionsIsOpen(false);
        setAddToGroupIsOpen(false);
        setGroupMembersIsOpen((prev) => !prev);
    };

    const [groupPermissionsIsOpen, setGroupPermissionsIsOpen] = useState(false);
    const handleCloseGroupPermissions = () => {
        setGroupPermissionsIsOpen(false);
    };
    const handleGroupPermissionsOpen = () => {
        setMenuIsOpen(false);
        setGroupMembersIsOpen(false);
        setAddToGroupIsOpen(false);
        setGroupPermissionsIsOpen((prev) => !prev);
    };

    return (
        <div className="relative flex flex-col max-h-screen flex-1">
            <button
                onClick={() => setMenuIsOpen((prev) => !prev)}
                className="cursor-pointer ml-auto hover:text-accent duration-300 active:text-foreground z-100"
            >
                Menu
            </button>

            <div
                className={`${menuIsOpen ? "max-h-96 p-[4px]" : "max-h-0 invisible p-0"} z-100 duration-300 w-50 border-1 border-foreground rounded-[4px] flex flex-col bg-background gap-y-[4px] absolute right-0 overflow-hidden top-[28px]`}
            >
                {hasPermission("add users to group") && isGroup && (
                    <button
                        onClick={(e) => {
                            e.stopPropagation();
                            handleAddToGroupOpen();
                            setMenuIsOpen(false);
                        }}
                        className="cursor-pointer hover:bg-foreground/20 active:bg-background duration-300"
                    >
                        Add member
                    </button>
                )}
                {isGroup && (
                    <button
                        onClick={(e) => {
                            e.stopPropagation();
                            handleGroupMembersOpen();
                        }}
                        className="cursor-pointer hover:bg-foreground/20 active:bg-background duration-300"
                    >
                        Group members
                    </button>
                )}

                {isGroup && (
                    <button
                        onClick={(e) => {
                            e.stopPropagation();
                            handleGroupPermissionsOpen();
                        }}
                        className="cursor-pointer hover:bg-foreground/20 active:bg-background duration-300"
                    >
                        Group Permissions
                    </button>
                )}
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
                handleClose={handleAddToGroupClose}
                chatId={chatId}
            />

            <GroupMembers
                groupMembersIsOpen={groupMembersIsOpen}
                handleClose={handleCloseGroupMembers}
                chatId={chatId}
                hasPermission={hasPermission}
                reloadPermissions={reloadPermissions}
            />

            <GroupRolesPermissions
                groupPermissionsIsOpen={groupPermissionsIsOpen}
                handleClose={handleCloseGroupPermissions}
                chatId={chatId}
                hasPermission={hasPermission}
                reloadPermissions={reloadPermissions}
            />
        </div>
    );
};

export default Menu;
