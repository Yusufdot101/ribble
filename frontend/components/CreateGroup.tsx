"use client";
import Contacts from "@/components/Contacts";
import { getUsersByEmail, UserType } from "@/utils/users";
import { useCallback, useEffect, useState } from "react";
import BackArrowButton from "./BackArrowButton";
import XButton from "./XButton";
import { ChatType, getChatByUserIds } from "@/utils/chats";
import { useRouter } from "next/navigation";
import SearchBar from "./SearchBar";
import { useAuthStore } from "@/store/useAuthStore";

interface Props {
    handleClose: () => void;
    handleUpdateGroups: (chat: ChatType) => void;
    createGroupOpen: boolean;
}

const CreateGroup = ({
    handleClose,
    handleUpdateGroups,
    createGroupOpen,
}: Props) => {
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

    const loggedInUserId = useAuthStore((state) => state.userId);
    const [users, setUsers] = useState<UserType[]>([]);
    const searchUsers = useCallback(
        async (email: string = "") => {
            setIsLoading(true);
            try {
                const { users } = await getUsersByEmail(email);
                setUsers(
                    users?.filter((user) => user.id !== loggedInUserId) ?? [],
                );
            } finally {
                setIsLoading(false);
            }
        },
        [loggedInUserId],
    );

    useEffect(() => {
        (() => searchUsers())();
    }, [loggedInUserId, searchUsers]);

    const [showConfigScreen, setShowConfigScreen] = useState(false);
    const [groupName, setGroupName] = useState("");

    const [memberPermissions, setMemberPermissions] = useState({
        "send message": true,
        "add users to group": true,
    });

    const [adminPermissions, setAdminPermissions] = useState({
        "send message": true,
        "add users to group": true,
        "remove users from group": true,
        "delete messages": true,
        "ban users": true,
        "promote members": false,
        "demote admins": false,
    });

    const router = useRouter();
    const handleCreate = async () => {
        const userRoles = new Map<number, string>();
        for (const user of selectedUsers) {
            userRoles.set(user.id, "member");
        }

        const enabledMemberPermissions = Object.entries(memberPermissions)
            .filter(([, value]) => value)
            .map(([key]) => key);

        const enabledAdminPermissions = Object.entries(adminPermissions)
            .filter(([, value]) => value)
            .map(([key]) => key);

        const rolePermissions = new Map<string, string[]>();
        rolePermissions.set("admin", enabledAdminPermissions);
        rolePermissions.set("member", enabledMemberPermissions);

        const chat = await getChatByUserIds(
            undefined,
            rolePermissions,
            userRoles,
            groupName,
            true,
        );
        if (!chat) return;
        handleUpdateGroups(chat);
        handleClose();
        router.push(`/chats/${chat.id}`);
    };

    return (
        <div
            className={`${createGroupOpen ? "translate-x-0" : "translate-x-full invisible"} w-full right-0 absolute h-full transition-transform duration-300 ease-in-out flex-1 flex overflow-x-hidden`}
        >
            <div
                className={`${showConfigScreen ? "-translate-x-full" : "translate-x-0"} h-full transition-transform duration-300 ease-in-out flex flex-1 flex-col gap-y-[8px]`}
            >
                <div className="flex w-full h-[32px] gap-x-[8px] items-center">
                    <BackArrowButton
                        handleClick={handleClose}
                        text="Add group members"
                    />
                </div>

                <div className="flex flex-col space-y-[8px] relative flex-1">
                    {selectedUsers?.length !== 0 && (
                        <div className="flex flex-col space-y-[8px] max-h-[280px] overflow-y-scroll border-b-1 border-foreground transition-all duration-300 ease-in-out">
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
                    <Contacts
                        isLoading={isLoading}
                        users={users}
                        handleUserClick={handleClick}
                        selectedUsers={selectedUsers.map((user) => user.id)}
                        excludeUsers={selectedUsers.map((user) => user.id)}
                    />

                    <div className="flex gap-x-[4px] w-fit absolute bottom-0 right-1/2 translate-x-1/2 w-full">
                        <button
                            aria-label="create group"
                            className=" bg-background border-1 border-foreground text-white p-[4px] rounded-[4px] hover:bg-foreground/10 active:bg-background duration-300 cursor-pointer w-full"
                            onClick={() => setShowConfigScreen(true)}
                        >
                            Group info
                        </button>

                        <button
                            aria-label="create group"
                            className="bg-accent text-white rounded-[4px] hover:bg-accent/80 active:bg-accent duration-300 cursor-pointer w-full"
                            onClick={handleCreate}
                        >
                            Create Group
                        </button>
                    </div>
                </div>
            </div>

            {/*Group configuration screen*/}
            <div
                className={`${showConfigScreen ? "translate-x-0" : "translate-x-full invisible"} w-full right-0 absolute h-full transition-transform duration-300 ease-in-out flex flex-col gap-y-[8px]`}
            >
                <div className="flex w-full h-[32px] gap-x-[8px] items-center">
                    <BackArrowButton
                        handleClick={() => setShowConfigScreen(false)}
                        text="Group info"
                    />
                </div>

                <div className="flex flex-col gap-y-[2px]">
                    <label htmlFor="groupName" className="text-foreground/70">
                        Group name:
                    </label>
                    <input
                        type="text"
                        className="w-full bg-foreground text-background px-[4px] py-[2px] border-none outline-none rounded-[4px]"
                        placeholder="group name"
                        value={groupName}
                        onChange={(e) => {
                            setGroupName(e.target.value);
                        }}
                        onKeyDown={(e) => {
                            if (e.key !== "Enter") return;
                            setShowConfigScreen(false);
                        }}
                        id="groupName"
                    />
                </div>

                <div className="flex flex-col p-[8px] w-full">
                    <span className="text-foreground/80">Members can: </span>
                    {Object.entries(memberPermissions).map(([key, value]) => (
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
                                        type="checkbox"
                                        className="peer sr-only"
                                        checked={value}
                                        onChange={(e) => {
                                            setMemberPermissions((prev) => ({
                                                ...prev,
                                                [key]: e.target.checked,
                                            }));
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

                <div className="flex flex-col gap-y-[2px]">
                    <div className="flex flex-col p-[8px] w-full">
                        <span className="text-foreground/80">Admins can: </span>
                        {Object.entries(adminPermissions).map(
                            ([key, value]) => (
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
                                                    setAdminPermissions(
                                                        (prev) => ({
                                                            ...prev,
                                                            [key]: e.target
                                                                .checked,
                                                        }),
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
                            ),
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default CreateGroup;
