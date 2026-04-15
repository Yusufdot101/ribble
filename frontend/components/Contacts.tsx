"use client";
import UserCard from "@/components/UserCard";
import UserCardSkeleton from "@/components/UserCardSkeleton";
import { useAuthStore } from "@/store/useAuthStore";
import { getChatByUserIDs } from "@/utils/chats";
import { getUsersByEmail, UserType } from "@/utils/users";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

const Contacts = () => {
    const [users, setUsers] = useState<UserType[]>([]);
    const [activeUser, setActiveUser] = useState<number>();
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        (async () => {
            setIsLoading(true);
            const users = await getUsersByEmail("");
            setUsers(users);
            setIsLoading(false);
        })();
    }, []);

    const router = useRouter();
    const handleClick = async (userID: number) => {
        setActiveUser(userID);
        const chat = await getChatByUserIDs([userID]);
        if (!chat) return;
        router.push(`/chats/${chat.ID}`);
    };

    const loggedInUserID = useAuthStore((state) => state.userID);
    return (
        <div className="flex-1 flex flex-col gap-y-[8px] rounded-[4px] ">
            <div
                className={`${!isLoading ? "hidden" : ""} flex flex-col transition-all duration-300`}
            >
                {[...Array(4).keys()].map((el) => (
                    <UserCardSkeleton key={el} index={el} />
                ))}
            </div>

            {users.length === 0 && !isLoading && (
                <p className="w-full text-center">No users</p>
            )}

            <div
                className={`${users.length === 0 ? "opacity-0 blur-sm" : ""} flex flex-col transition-all duration-300`}
            >
                {users
                    .filter((elem) => elem.id != loggedInUserID)
                    .map((user) => (
                        <UserCard
                            activeUserID={activeUser || -100}
                            key={user.id}
                            user={user}
                            handleClick={handleClick}
                        />
                    ))}
            </div>
        </div>
    );
};

export default Contacts;
