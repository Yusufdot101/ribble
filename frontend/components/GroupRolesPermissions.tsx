"use client";

import { useCallback, useEffect, useState } from "react";
import BackArrowButton from "./BackArrowButton";
import { useSocket } from "@/providers/socket-provider";
import {
    getRolesPermissions,
    updateRolePermissions,
} from "@/utils/permissions";

interface Props {
    handleClose: () => void;
    groupPermissionsIsOpen: boolean;
    chatId: number;
    hasPermission: (permissionName: string) => boolean;
    reloadPermissions: () => void;
}

const GroupRolesPermissions = ({
    handleClose,
    groupPermissionsIsOpen,
    chatId,
    hasPermission,
    reloadPermissions,
}: Props) => {
    const [memberPermissions, setMemberPermissions] = useState<
        Map<string, boolean>
    >(
        new Map<string, boolean>([
            ["send message", false],
            ["add users to group", false],
        ]),
    );

    const [adminPermissions, setAdminPermissions] = useState<
        Map<string, boolean>
    >(
        new Map<string, boolean>([
            ["send message", false],
            ["add users to group", false],
            ["remove users from group", false],
            ["delete messages", false],
            ["ban users", false],
            ["promote members", false],
            ["demote admins", false],
        ]),
    );

    const socket = useSocket();

    const handleMessage = useCallback(
        (event: MessageEvent) => {
            let data: any;
            try {
                data =
                    typeof event.data === "string"
                        ? JSON.parse(event.data)
                        : null;
            } catch (error) {
                console.error(error);
                return;
            }
            if (!data) return;

            if (data.type === "error") {
                console.error(data.message);
                return;
            }

            const payload = data.payload;
            if (data.subType === "updatedRolePermissions" && payload) {
                const role = payload.role;
                const action = payload.action;
                const permission = payload.permission;
                if (role === "admin") {
                    setAdminPermissions((prev) => {
                        const newMap = new Map(prev);
                        newMap.set(
                            permission,
                            action === "grant" ? true : false,
                        );
                        return newMap;
                    });
                } else {
                    setMemberPermissions((prev) => {
                        const newMap = new Map(prev);
                        newMap.set(
                            permission,
                            action === "grant" ? true : false,
                        );
                        return newMap;
                    });
                }
                reloadPermissions();
            }
        },
        [reloadPermissions],
    );

    useEffect(() => {
        socket?.addEventListener("message", handleMessage);
        return () => {
            socket?.removeEventListener("message", handleMessage);
        };
    }, [socket, handleMessage]);

    const toggleMemberPermissions = (key: string, checked: boolean) => {
        if (!hasPermission("update permissions")) {
            alert("not permitted");
            return;
        }
        setMemberPermissions((prev) => {
            const newMap = new Map(prev);
            newMap.set(key, checked);
            return newMap;
        });

        let action: "grant" | "revoke";
        if (checked) {
            action = "grant";
        } else {
            action = "revoke";
        }
        updatePermission("member", key, action);
    };
    const toggleAdminPermissions = (key: string, checked: boolean) => {
        if (!hasPermission("update permissions")) {
            alert("not permitted");
            return;
        }
        setAdminPermissions((prev) => {
            const newMap = new Map(prev);
            newMap.set(key, checked);
            return newMap;
        });

        let action: "grant" | "revoke";
        if (checked) {
            action = "grant";
        } else {
            action = "revoke";
        }
        updatePermission("admin", key, action);
    };

    const updatePermission = async (
        role: string,
        permission: string,
        action: "revoke" | "grant",
    ) => {
        const success = await updateRolePermissions(
            chatId,
            role,
            permission,
            action,
        );
        if (success) {
            return true;
        }
        // revert the ui on failure
        if (role === "admin") {
            setAdminPermissions((prev) => {
                const newMap = new Map(prev);
                newMap.set(permission, action === "revoke" ? true : false);
                return newMap;
            });
            return;
        }
        setMemberPermissions((prev) => {
            const newMap = new Map(prev);
            newMap.set(permission, action === "revoke" ? true : false);
            return newMap;
        });
    };

    useEffect(() => {
        (async () => {
            const { rolePermissions } = await getRolesPermissions(chatId, [
                "admin",
                "member",
            ]);
            if (!rolePermissions) return;
            const adminPermissions = rolePermissions.admin;
            const memberPermissions = rolePermissions.member;

            const adminSet = new Set(
                (rolePermissions.admin ?? []).map(
                    (permission) => permission.name,
                ),
            );
            const memberSet = new Set(
                (rolePermissions.member ?? []).map(
                    (permission) => permission.name,
                ),
            );

            setAdminPermissions((prev) => {
                const newMap = new Map<string, boolean>();
                for (const key of prev.keys())
                    newMap.set(key, adminSet.has(key));
                return newMap;
            });

            setMemberPermissions((prev) => {
                const newMap = new Map<string, boolean>();
                for (const key of prev.keys())
                    newMap.set(key, memberSet.has(key));
                return newMap;
            });
        })();
    }, [chatId]);

    return (
        <div
            className={`${groupPermissionsIsOpen ? "translate-x-0" : "translate-x-full"} transition-transform absolute w-full h-screen overflow-y-hidden bg-background duration-300 flex-1 flex overflow-x-hidden z-10 flex-col gap-y-[1px]`}
        >
            <div className="flex w-full h-[32px] gap-x-[8px] items-center">
                <BackArrowButton
                    handleClick={handleClose}
                    text="Group Permissions"
                />
            </div>
            <div className="flex flex-col gap-y-[2px]">
                <div className="flex flex-col p-[8px] w-full">
                    <span className="text-foreground/80">Members can: </span>
                    {[...memberPermissions].map(([key, value]) => (
                        <div
                            key={key}
                            className="flex px-[8px] items-center justify-between"
                        >
                            <div className="flex items-center w-fit gap-x-[8px] h-[32px]">
                                <span>{key}</span>
                            </div>

                            <div className="flex">
                                <label className="relative inline-flex cursor-pointer items-center">
                                    <input
                                        aria-label={`Member permission: ${key}`}
                                        type="checkbox"
                                        className="peer sr-only"
                                        checked={value}
                                        onChange={(e) => {
                                            toggleMemberPermissions(
                                                key,
                                                e.target.checked,
                                            );
                                        }}
                                    />

                                    <div
                                        className="relative h-6 w-11 rounded-full bg-foreground 
                                    transition-colors duration-300 peer-checked:bg-accent 
                                    after:absolute after:left-0.5 after:top-0.5 after:h-5 after:w-5
                                    after:rounded-full after:bg-white after:transition-transform 
                                    after:duration-300 after:content-[''] 
                                    peer-checked:after:translate-x-5"
                                    ></div>
                                </label>
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            <div className="flex flex-col gap-y-[2px]">
                <div className="flex flex-col p-[8px] w-full">
                    <span className="text-foreground/80">Admins can: </span>
                    {[...adminPermissions].map(([key, value]) => (
                        <div
                            key={key}
                            className="flex px-[8px] items-center justify-between"
                        >
                            <div className="flex items-center w-fit gap-x-[8px] h-[32px]">
                                <span>{key}</span>
                            </div>

                            <div className="flex">
                                <label className="relative inline-flex cursor-pointer items-center">
                                    <input
                                        aria-label={`Admin permission: ${key}`}
                                        type="checkbox"
                                        className="peer sr-only"
                                        checked={value}
                                        onChange={(e) => {
                                            toggleAdminPermissions(
                                                key,
                                                e.target.checked,
                                            );
                                        }}
                                    />

                                    <div
                                        className="relative h-6 w-11 rounded-full bg-foreground 
                                    transition-colors duration-300 peer-checked:bg-accent 
                                    after:absolute after:left-0.5 after:top-0.5 after:h-5 after:w-5
                                    after:rounded-full after:bg-white after:transition-transform 
                                    after:duration-300 after:content-[''] 
                                    peer-checked:after:translate-x-5"
                                    ></div>
                                </label>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default GroupRolesPermissions;
