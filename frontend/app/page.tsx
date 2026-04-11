"use client";
import SearchBar from "@/components/SearchBar";
import UserCard from "@/components/UserCard";
import { getUsersByEmail, UserType } from "@/utils/users";
import { useEffect, useState } from "react";

const Home = () => {
    const [users, setUsers] = useState<UserType[]>([]);
    const [activeUser, setActiveUser] = useState<number>();
    useEffect(() => {
        (async () => {
            const users = await getUsersByEmail("");
            console.log(users);
            setUsers(users);
        })();
    }, []);
    return (
        <div className="flex flex-col gap-y-[8px]">
            <SearchBar />
            <div className="flex flex-col border-1 border-foreground rounded-[4px]">
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
