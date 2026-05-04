"use client";
import Contacts from "@/components/Contacts";
import { UserType } from "@/utils/users";
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

interface Props {
    handleClose: () => void;
    groupMembersIsOpen: boolean;
    chatID: number;
    hasPermission: (permissionName: string) => boolean;
}

const GroupMembers = ({
    handleClose,
    groupMembersIsOpen,
    chatID,
    hasPermission,
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
            chatID,
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
        const success = await unbanUser(chatID, clickedUser.id);
        if (!success) return;
        setMenuIsOpen(false);
        setBanMenuIsOpen(false);
        handleClose();
    };

    const searchUsers = useCallback(
        async (value: string = "") => {
            setIsLoading(true);
            try {
                let users = await getChatUsers(chatID);
                const bannedUsers = await getBannedUsers(chatID, value);
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
        [chatID],
    );

    const handleRightClick = (user: UserType) => {
        setClickedUser(user);
        setMenuIsOpen(true);
    };

    useEffect(() => {
        (() => searchUsers())();
    }, [searchUsers]);

    const handleRemove = async (userID: number) => {
        if (
            (userID === loggedInUserID &&
                !confirm(`are you sure you want to exist this group`)) ||
            !clickedUser ||
            !confirm(
                `are you sure you want to remove ${clickedUser.name} from this group`,
            )
        ) {
            return;
        }
        await removeUserFromGroup(chatID, userID);
        setUsers((prev) => prev.filter((user) => user.id !== userID));
        setMenuIsOpen(false);
        handleClose();
    };

    const loggedInUserID = useAuthStore((state) => state.userID);
    const router = useRouter();

    return (
        <div
            className={`${groupMembersIsOpen ? "translate-x-0" : "translate-x-full"} transition-transform absolute w-full bg-background top-1/2 translate-y-1/2 duration-300 flex-1 flex overflow-x-hidden z-10 flex-col gap-y-[1px]`}
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

                <div className="flex flex-col gap-y-[8px] relative flex-1">
                    <SearchBar
                        placeholder="Search group members"
                        handleEnter={searchUsers}
                    />

                    <div className="h-full max-h-[200px] overflow-y-scroll">
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
                                loggedInUserID ? [loggedInUserID] : []
                            }
                        />
                    </div>
                </div>
            </div>

            <div className="h-full max-h-[200px] overflow-y-scroll">
                <p className="w-full text-center">Banned Users</p>
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

            <div
                className={`${menuIsOpen ? "max-h-96 p-[4px]" : "max-h-0 invisible p-0"} z-1 duration-300 flex justify-center items-center absolute overflow-hidden h-full w-full bg-background/80`}
            >
                <div className="bg-background w-80 border-1 border-foreground rounded-[4px] flex flex-col justify-center">
                    {!clickedUserIsBanned &&
                        hasPermission("remove users from group") && (
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

                    {!clickedUserIsBanned && hasPermission("ban users") && (
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

            <button
                onClick={() => {
                    if (!loggedInUserID) return;
                    handleRemove(loggedInUserID);
                }}
                className="w-full bg-red-500 p-[4px] rounded-[4px] cursor-pointer hover:bg-red-600 active:bg-red-500 duration-300"
            >
                Exit group
            </button>
        </div>
    );
};

export default GroupMembers;
