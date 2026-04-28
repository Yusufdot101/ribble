"use client";
import Image from "next/image";
import searchIcon from "../assets/searchIcon.svg";
import { useState } from "react";

interface Props {
    handleEnter: (value: string) => void;
    placeholder?: string;
}

const SearchBar = ({
    handleEnter,
    placeholder = "Search on start new chat",
}: Props) => {
    const [value, setValue] = useState("");
    return (
        <div
            onKeyDown={(e) => {
                if (e.key === "Enter") {
                    handleEnter(value);
                }
            }}
            className="shrink-0 border-[1px] border-solid border-[#ffffff] rounded-[8px] py-[12px] px-[12px] flex gap-x-[4px] h-[50px] cursor-pointer opacity-80"
        >
            <Image
                role="button"
                tabIndex={0}
                aria-label="search"
                src={searchIcon}
                alt="search icon"
                className="h-full w-[40px] rounded-[8px] "
                onClick={() => {
                    handleEnter(value);
                }}
            />

            <input
                type="text"
                placeholder={placeholder}
                className="border-none outline-none h-full w-full"
                value={value}
                onChange={(e) => {
                    setValue(e.target.value);
                }}
            />
        </div>
    );
};

export default SearchBar;
