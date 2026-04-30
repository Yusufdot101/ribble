"use client";
import { getUserPermissions, PermissionType } from "@/utils/permissions";
import { useEffect, useState } from "react";
import AddUsersToGroup from "./AddUsersToGroup";

type Props = {
    chatID: number;
    currentGroupUsers: number[];
};

const Menu = ({ chatID, currentGroupUsers }: Props) => {
    const [menuIsOpen, setMenuIsOpen] = useState(false);
    const [permissions, setPermissions] = useState<PermissionType[]>([]);
    useEffect(() => {
        (() => setPermissions([]))();
        if (chatID <= 0) return;
        (async () => {
            const permissions = await getUserPermissions(chatID);
            setPermissions(permissions ?? []);
        })();
    }, [chatID]);

    const hasPermission = (permissionName: string): boolean => {
        return (
            permissions.filter(
                (permission) => permission.name === permissionName,
            ).length !== 0
        );
    };

    const [addToGroupIsOpen, setAddToGroupIsOpen] = useState(false);
    const handleClose = () => {
        setAddToGroupIsOpen(false);
    };

    return (
        <div className="relative flex flex-col max-h-screen overflow-x-clip">
            <button
                onClick={() => setMenuIsOpen((prev) => !prev)}
                className="cursor-pointer ml-auto"
            >
                Menu
            </button>

            <div
                className={`${menuIsOpen ? "max-h-96 p-[4px]" : "max-h-0 invisible p-0"} duration-300 w-50 border-1 border-foreground rounded-[4px] flex flex-col bg-background gap-y-[4px] absolute right-0 overflow-hidden top-[28px]`}
            >
                {hasPermission("add users to group") && (
                    <button
                        onClick={(e) => {
                            e.stopPropagation();
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
                currentGroupUsers={currentGroupUsers}
                chatID={chatID}
            />
        </div>
    );
};

export default Menu;
