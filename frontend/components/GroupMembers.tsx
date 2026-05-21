"use client";
import Contacts from "@/components/Contacts";
import { changeMemberRole, UserType } from "@/utils/users";
import { useCallback, useEffect, useState } from "react";
import BackArrowButton from "./BackArrowButton";
import SearchBar from "./SearchBar";
import { getChatUsers } from "@/utils/chats";
import { useAuthStore } from "@/store/useAuthStore";
import {
    banUser,
    getBannedUsers,
    removeUserFromGroup,
    unbanUser,
} from "@/utils/groups";
import { useRouter } from "next/navigation";
import { useSocket } from "@/providers/socket-provider";

interface Props {
    handleClose: () => void;
    groupMembersIsOpen: boolean;
    chatId: number;
    hasPermission: (permissionName: string) => boolean;
    reloadPermissions: () => void;
}

const GroupMembers = ({
    handleClose,
    groupMembersIsOpen,
    chatId,
    hasPermission,
    reloadPermissions,
}: Props) => {
    const [isLoading, setIsLoading] = useState(false);
    const [users, setUsers] = useState<UserType[]>([]);
    const [bannedUsers, setBannedUsers] = useState<UserType[]>([]);
    const [clickedUser, setClickedUser] = useState<UserType>();
    const [clickedUserIsBanned, setClickedUserIsBanned] = useState(false);
    const [menuIsOpen, setMenuIsOpen] = useState(false);
    const [banMenuIsOpen, setBanMenuIsOpen] = useState(false);

    const [banReason, setBanReason] = useState("");
    const [banDuration, setBanDuration] = useState<number>(-1);
    const [banDurationFrame, setBanDurationFrame] = useState("days");

    const handleBanUser = async () => {
        if (!clickedUser) return;

        if (!confirm(`are you sure you want to ban ${clickedUser?.name}`))
            return;

        const success = await banUser(
            chatId,
            clickedUser.id,
            banReason,
            banDurationFrame,
            banDuration,
        );
        if (!success) return;
        setMenuIsOpen(false);
        setBanMenuIsOpen(false);
        handleClose();
    };

    const handleUnBanUser = async () => {
        if (!clickedUser) return;
        if (!confirm(`are you sure you want to unban ${clickedUser?.name}`))
            return;
        const success = await unbanUser(chatId, clickedUser.id);
        if (!success) return;
        setMenuIsOpen(false);
        setBanMenuIsOpen(false);
        handleClose();
    };

    const searchUsers = useCallback(
        async (value: string = "") => {
            setIsLoading(true);
            try {
                let users = await getChatUsers(chatId);
                const { users: bannedUsers } = await getBannedUsers(
                    chatId,
                    value,
                );

                users = users?.filter(
                    (user) =>
                        user.name.includes(value) || user.email.includes(value),
                );
                setUsers(users ?? []);
                setBannedUsers(bannedUsers);
            } finally {
                setIsLoading(false);
            }
        },
        [chatId],
    );

    const hasAnyPermission = (() => {
        const usedPermissions = [
            "remove users from group",
            "ban users",
            "promote members",
            "demote admins",
        ];
        for (const permission of usedPermissions) {
            if (hasPermission(permission)) return true;
        }
        return false;
    })();
    const handleRightClick = (user: UserType) => {
        if (!hasAnyPermission || user.role === "creator") return;
        setClickedUser(user);
        setMenuIsOpen(true);
    };

    useEffect(() => {
        (() => searchUsers())();
    }, [searchUsers]);

    const handleRemove = async (userId: number) => {
        if (
            userId === loggedInUserId &&
            !confirm(`are you sure you want to exit this group`)
        ) {
            return;
        } else if (
            userId !== loggedInUserId &&
            (!clickedUser ||
                !confirm(
                    `are you sure you want to remove ${clickedUser.name} from this group`,
                ))
        ) {
            return;
        }
        await removeUserFromGroup(chatId, userId);
        setUsers((prev) => prev.filter((user) => user.id !== userId));
        setMenuIsOpen(false);
        handleClose();
    };

    const loggedInUserId = useAuthStore((state) => state.userId);
    const router = useRouter();

    const socket = useSocket();

    const handleMessage = useCallback(
        (event: MessageEvent) => {
            const data = JSON.parse(event.data);

            if (data.type === "error") {
                console.error(data.message);
                return;
            }

            const payload = data.payload;
            if (data.subType === "updatedUserRole") {
                const userId = payload.userId;
                if (loggedInUserId === userId) {
                    reloadPermissions();
                    return;
                }
                setUsers((prev) => {
                    return prev.map((user) =>
                        user.id === userId
                            ? {
                                  ...user,
                                  role: payload.newRole,
                              }
                            : user,
                    );
                });
            }

            if (data.subType === "usersAdded") {
                const users = payload.addedUsers as UserType[];
                setUsers((prev) => [
                    ...prev,
                    ...users
                        .filter((user) => user.id !== loggedInUserId)
                        .map((user) => ({
                            ...user,
                            role: "member",
                        })),
                ]);
            }

            if (["userBanned"].includes(data.subType)) {
                const bannedUser = payload.target as UserType;
                if (bannedUser.id === loggedInUserId) {
                    alert("you have been banned from the group");
                    router.replace("/chats");
                    return;
                }
                setUsers((prev) =>
                    prev.filter((user) => user.id !== bannedUser.id),
                );
                setBannedUsers((prev) => [...(prev ?? []), bannedUser]);
            }

            if (["userUnbanned"].includes(data.subType)) {
                const unbannedUser = payload.target as UserType;
                setBannedUsers((prev) =>
                    prev.filter((user) => user.id !== unbannedUser.id),
                );
            }

            if (["userRemoved"].includes(data.subType)) {
                const removedUser = payload.target as UserType;
                const actor = payload.actor as UserType;
                if (removedUser.id === loggedInUserId) {
                    if (actor.id == loggedInUserId) {
                        alert("you have left the group");
                    } else {
                        alert("you have been removed from the group");
                    }
                    router.replace("/chats");
                    return;
                }
                setUsers((prev) =>
                    prev.filter((user) => user.id !== removedUser.id),
                );
            }
        },
        [loggedInUserId, reloadPermissions, router],
    );

    useEffect(() => {
        socket?.addEventListener("message", handleMessage);
        return () => {
            socket?.removeEventListener("message", handleMessage);
        };
    }, [socket, handleMessage]);

    return (
        <div
            className={`${groupMembersIsOpen ? "translate-x-0" : "translate-x-full"} transition-transform absolute w-full h-screen overflow-y-hidden bg-background duration-300 flex-1 flex overflow-x-hidden z-10 flex-col gap-y-[1px]`}
            onClick={() => setMenuIsOpen(false)}
            onKeyDown={(e) => {
                if (e.key !== "Escape") return;
                setMenuIsOpen(false);
            }}
        >
            <div className="h-full transition-transform duration-300 ease-in-out flex flex-1 flex-col gap-y-[8px]">
                <div className="flex w-full h-[32px] gap-x-[8px] items-center">
                    <BackArrowButton
                        handleClick={handleClose}
                        text="Group members"
                    />
                </div>

                <div className="flex flex-col gap-y-[8px] relative flex-1 h-full">
                    <SearchBar
                        placeholder="Search group members"
                        handleEnter={searchUsers}
                    />

                    <button
                        onClick={() => {
                            if (!loggedInUserId) return;
                            handleRemove(loggedInUserId);
                        }}
                        className="w-full bg-red-500 p-[4px] rounded-[4px] cursor-pointer hover:bg-red-600 active:bg-red-500 duration-300"
                    >
                        Exit group
                    </button>

                    <div className="flex-1 overflow-y-auto h-full">
                        <Contacts
                            isLoading={isLoading}
                            users={users}
                            handleUserClick={() => {}}
                            handleUserRightClick={(user: UserType) => {
                                handleRightClick(user);
                                setClickedUserIsBanned(false);
                            }}
                            selectedUsers={[]}
                            excludeUsers={
                                loggedInUserId ? [loggedInUserId] : []
                            }
                        />
                    </div>

                    <div className="flex-1 flex flex-col justify-between overflow-y-auto">
                        <p className="w-full text-center text-red-500">
                            Banned Users
                        </p>
                        <div className="h-full flex-1 min-h-0 overflow-y-auto">
                            <Contacts
                                isLoading={isLoading}
                                users={bannedUsers}
                                handleUserClick={() => {}}
                                handleUserRightClick={(user: UserType) => {
                                    handleRightClick(user);
                                    setClickedUserIsBanned(true);
                                }}
                                selectedUsers={[]}
                                excludeUsers={[]}
                            />
                        </div>
                    </div>
                </div>
            </div>

            <div
                className={`${menuIsOpen ? "p-[4px] min-h-screen" : "max-h-0 invisible p-0"} z-1 duration-0 flex justify-center items-center absolute overflow-hidden w-full bg-background/80`}
            >
                <div className="bg-background w-80 border-1 border-foreground rounded-[4px] flex flex-col justify-center">
                    {!clickedUserIsBanned &&
                        hasPermission("remove users from group") &&
                        clickedUser?.role !== "creator" && (
                            <button
                                onClick={(e) => {
                                    e.stopPropagation();
                                    if (!clickedUser) return;
                                    handleRemove(clickedUser.id);
                                }}
                                onKeyDown={(e) => {
                                    e.stopPropagation();
                                    if (e.key !== "Enter") return;
                                    if (!clickedUser) return;
                                    handleRemove(clickedUser.id);
                                    router.push("/chats");
                                }}
                                className="cursor-pointer hover:bg-foreground/20 active:bg-background duration-300 p-[4px]"
                            >
                                Remove {clickedUser?.name}
                            </button>
                        )}

                    {!clickedUserIsBanned &&
                        hasPermission("ban users") &&
                        clickedUser?.role !== "creator" && (
                            <button
                                onClick={(e) => {
                                    e.stopPropagation();
                                    setMenuIsOpen(false);
                                    setBanMenuIsOpen(true);
                                }}
                                onKeyDown={(e) => {
                                    if (e.key !== "Enter") return;
                                    e.stopPropagation();
                                    setMenuIsOpen(false);
                                    setBanMenuIsOpen(true);
                                }}
                                className="cursor-pointer hover:bg-foreground/20 active:bg-background duration-300 p-[4px]"
                            >
                                Ban {clickedUser?.name}
                            </button>
                        )}

                    {!clickedUserIsBanned &&
                        hasPermission("promote members") &&
                        clickedUser?.role === "member" && (
                            <button
                                onClick={async (e) => {
                                    e.stopPropagation();
                                    const success = await changeMemberRole(
                                        "admin",
                                        clickedUser.id,
                                        chatId,
                                    );
                                    if (!success) return;
                                    setMenuIsOpen(false);
                                    handleClose();
                                }}
                                onKeyDown={async (e) => {
                                    if (e.key !== "Enter") return;
                                    e.stopPropagation();
                                    const success = await changeMemberRole(
                                        "admin",
                                        clickedUser.id,
                                        chatId,
                                    );
                                    if (!success) return;
                                    setMenuIsOpen(false);
                                    handleClose();
                                }}
                                className="cursor-pointer hover:bg-foreground/20 active:bg-background duration-300 p-[4px]"
                            >
                                Assign {clickedUser?.name} as admin
                            </button>
                        )}

                    {!clickedUserIsBanned &&
                        hasPermission("demote admins") &&
                        clickedUser?.role === "admin" && (
                            <button
                                onClick={async (e) => {
                                    e.stopPropagation();
                                    const success = await changeMemberRole(
                                        "member",
                                        clickedUser.id,
                                        chatId,
                                    );
                                    if (!success) return;
                                    setMenuIsOpen(false);
                                    handleClose();
                                }}
                                onKeyDown={async (e) => {
                                    if (e.key !== "Enter") return;
                                    e.stopPropagation();
                                    const success = await changeMemberRole(
                                        "member",
                                        clickedUser.id,
                                        chatId,
                                    );
                                    if (!success) return;
                                    setMenuIsOpen(false);
                                    handleClose();
                                }}
                                className="cursor-pointer hover:bg-foreground/20 active:bg-background duration-300 p-[4px]"
                            >
                                Dismiss {clickedUser?.name} as admin
                            </button>
                        )}

                    {clickedUserIsBanned && hasPermission("ban users") && (
                        <button
                            onClick={(e) => {
                                e.stopPropagation();
                                handleUnBanUser();
                            }}
                            onKeyDown={(e) => {
                                if (e.key !== "Enter") return;
                                e.stopPropagation();
                                e.preventDefault();
                                handleUnBanUser();
                            }}
                            className="cursor-pointer hover:bg-foreground/20 active:bg-background duration-300 p-[4px]"
                        >
                            Unban {clickedUser?.name}
                        </button>
                    )}
                </div>
            </div>

            <div
                className={`${banMenuIsOpen ? "max-h-96 p-[4px]" : "max-h-0 invisible p-0"} z-1 duration-300 flex justify-center items-center absolute overflow-hidden h-full w-full bg-background/80`}
                onClick={() => {
                    setBanMenuIsOpen(false);
                }}
            >
                <div className="bg-background w-full border-1 border-foreground rounded-[4px] flex flex-col justify-center">
                    {hasPermission("ban users") && (
                        <form
                            className="flex flex-col gap-y-[4px] p-[4px]"
                            onClick={(e) => {
                                e.preventDefault();
                                e.stopPropagation();
                            }}
                            onSubmit={(e) => {
                                e.stopPropagation();
                                e.preventDefault();
                                handleBanUser();
                            }}
                            onKeyDown={(e) => {
                                if (e.key !== "Enter") return;
                                e.stopPropagation();
                                e.preventDefault();
                                handleBanUser();
                            }}
                        >
                            <div className="flex flex-col gap-y-[2px]">
                                <label htmlFor="banReason">Reason</label>
                                <input
                                    id="banReason"
                                    type="text"
                                    required
                                    value={banReason}
                                    onChange={(e) =>
                                        setBanReason(e.target.value)
                                    }
                                    className="bg-foreground text-background border-none outline-none p-[4px]"
                                />
                            </div>
                            <div className="flex flex-col gap-y-[2px]">
                                <label htmlFor="banDuration">
                                    Expiry (-1 for indefinitely)
                                </label>
                                <div className="flex gap-x-[2px]">
                                    <select
                                        name="banDurationFrame"
                                        id="banDurationFrame"
                                        className="bg-foreground text-background border-none outline-none p-[4px]"
                                        onChange={(e) =>
                                            setBanDurationFrame(e.target.value)
                                        }
                                        value={banDurationFrame}
                                    >
                                        <option value="hours">Hours</option>
                                        <option value="days">Days</option>
                                        <option value="weeks">Weeks</option>
                                        <option value="months">Months</option>
                                        <option value="years">Years</option>
                                    </select>
                                    <input
                                        type="number"
                                        className="bg-foreground w-full text-background border-none outline-none p-[4px]"
                                        min={-1}
                                        value={banDuration}
                                        onChange={(e) => {
                                            setBanDuration(+e.target.value);
                                        }}
                                    />
                                </div>
                            </div>
                            <div className="flex gap-x-[4px]">
                                <button
                                    onClick={(e) => {
                                        e.preventDefault();
                                        setMenuIsOpen(false);
                                        setBanMenuIsOpen(false);
                                    }}
                                    onKeyDown={(e) => {
                                        if (e.key !== "Enter") return;
                                        e.preventDefault();
                                        e.stopPropagation();
                                        setMenuIsOpen(false);
                                        setBanMenuIsOpen(false);
                                    }}
                                    className="cursor-pointer w-full bg-green-700 hover:bg-green-600 active:bg-green-700 duration-300 p-[4px]"
                                >
                                    Cancel
                                </button>
                                <button
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        e.preventDefault();
                                        handleBanUser();
                                    }}
                                    onKeyDown={(e) => {
                                        e.stopPropagation();
                                        if (e.key !== "Enter") return;
                                        e.preventDefault();
                                        handleBanUser();
                                    }}
                                    className="cursor-pointer w-full bg-red-700 hover:bg-red-600 active:bg-red-700 duration-300 p-[4px]"
                                >
                                    Ban {clickedUser?.name}
                                </button>
                            </div>
                        </form>
                    )}
                </div>
            </div>
        </div>
    );
};

export default GroupMembers;
