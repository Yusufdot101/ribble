"use client";
import SearchBar from "@/components/SearchBar";
import UserCard from "@/components/UserCard";
import UserCardSkeleton from "@/components/UserCardSkeleton";
import { getUsersByEmail, UserType } from "@/utils/users";
import { useEffect, useState } from "react";

const Home = () => {
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
    return (
        <div className="flex flex-col gap-y-[8px]">
            <SearchBar />
            <div
                className={`${!isLoading ? "hidden" : ""} flex flex-col border-1 border-foreground rounded-[4px] transition-all duration-300`}
            >
                {[...Array(4).keys()].map((el) => (
                    <UserCardSkeleton key={el} index={el} />
                ))}
            </div>

            {users.length === 0 && !isLoading && (
                <p className="w-full text-center">No users</p>
            )}

            <div
                className={`${users.length === 0 ? "opacity-0 blur-sm" : ""} flex flex-col border-1 border-foreground rounded-[4px] transition-all duration-300`}
            >
                {users.map((user, index) => (
                    <UserCard
                        activeUserID={activeUser || -100}
                        index={index}
                        key={user.id}
                        user={user}
                        handleClick={(userID: number) => setActiveUser(userID)}
                    />
                ))}
            </div>
        </div>
    );
};

export default Home;
