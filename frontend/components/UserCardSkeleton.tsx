"use client";
import { useState } from "react";

interface Props {
    index: number;
}

const nameSizes = [30, 40, 50, 60, 75];
const emailSizes = [10, 20, 30, 50, 60, 75];

const UserCardSkeleton = ({ index }: Props) => {
    const [nameSize] = useState(
        () => nameSizes[Math.floor(Math.random() * nameSizes.length)],
    );
    const [emailSize] = useState(
        () => emailSizes[Math.floor(Math.random() * emailSizes.length)],
    );
    return (
        <div
            className={`${index === 0 ? "" : "border-t-[1px]"} border-foreground p-[4px] pointer-events-none duration-300 flex flex-col gap-y-[4px] h-[64px] duration-300`}
        >
            <div
                style={{ width: `${nameSize}%` }}
                className="h-full bg-foreground/75"
            ></div>
            <div
                style={{ width: `${emailSize}%` }}
                className="h-full bg-foreground/75"
            ></div>
        </div>
    );
};

export default UserCardSkeleton;
