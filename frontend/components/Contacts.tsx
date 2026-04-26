"use client";
import UserCard from "@/components/UserCard";
import UserCardSkeleton from "@/components/UserCardSkeleton";
import { useAuthStore } from "@/store/useAuthStore";
import { getUsersByEmail, UserType } from "@/utils/users";
import { useEffect, useState } from "react";
import SearchBar from "./SearchBar";

interface Props {
    selectedUsers: number[];
    handleUserClick: (user: UserType) => void;
    excludeUsers?: number[];
}

const Contacts = ({ selectedUsers, handleUserClick, excludeUsers }: Props) => {
    const [users, setUsers] = useState<UserType[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        (async () => {
            setIsLoading(true);
            const users = await getUsersByEmail("");
            setUsers(users);
            setIsLoading(false);
        })();
    }, []);

    const loggedInUserID = useAuthStore((state) => state.userID);
    const visibleUsers = (users ?? []).filter(
        (elem) =>
            elem.id !== loggedInUserID && !excludeUsers?.includes(elem.id),
    );

    return (
        <div className="flex-1 flex flex-col gap-y-[8px] rounded-[4px] gap-y-[8px]">
            <SearchBar
                handleEnter={async (email: string) => {
                    setUsers([]);
                    setIsLoading(true);
                    const users = await getUsersByEmail(email.trim());
                    setUsers(users);
                    setIsLoading(false);
                }}
            />
            <div
                className={`${!isLoading ? "hidden" : ""} flex flex-col transition-all duration-300`}
            >
                {[...Array(4).keys()].map((el) => (
                    <UserCardSkeleton key={el} index={el} />
                ))}
            </div>

            {visibleUsers.length === 0 && !isLoading ? (
                <p className="w-full text-center">No users</p>
            ) : null}

            <div
                className={`${visibleUsers?.length === 0 ? "opacity-0 blur-sm" : ""} flex flex-col transition-all duration-300`}
            >
                {visibleUsers.map((user) => (
                    <UserCard
                        activeUsers={selectedUsers}
                        key={user.id}
                        user={user}
                        handleClick={(userID: number) =>
                            handleUserClick(
                                users.filter((user) => user.id === userID)[0],
                            )
                        }
                    />
                ))}
            </div>
        </div>
    );
};

export default Contacts;
