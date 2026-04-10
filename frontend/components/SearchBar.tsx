"use client";
import Image from "next/image";
import searchIcon from "../assets/searchIcon.svg";
import { useState } from "react";
// import { getUsersByEmail, UserType } from "@/utils/users";

const SearchBar = () => {
    const [email, setEmail] = useState("");
    // const [users, setUsers] = useState<UserType[]>([]);
    const handleEnter = async () => {
        // const users = await getUsersByEmail(email);
    };
    return (
        <div
            onKeyDown={(e) => {
                if (e.key === "Enter") {
                    handleEnter();
                }
            }}
            className="border-[1px] border-solid border-[#ffffff] rounded-[8px] py-[12px] px-[12px] flex gap-x-[4px] h-[50px] cursor-pointer opacity-80"
        >
            <Image
                role="button"
                tabIndex={0}
                aria-label="search"
                src={searchIcon}
                alt="search icon"
                className="h-full w-[40px] rounded-[8px]"
            />

            <input
                type="text"
                placeholder="Search or start new chat"
                className="border-none outline-none h-full w-full"
                value={email}
                onChange={(e) => {
                    setEmail(e.target.value);
                }}
            />
        </div>
    );
};

export default SearchBar;
