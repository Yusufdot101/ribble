"use client";
import Contacts from "@/components/Contacts";
import CreateGroup from "@/components/CreateGroup";
import CreateGroupButton from "@/components/CreateGroupButton";
import { getChatByUserIDs } from "@/utils/chats";
import { UserType } from "@/utils/users";
import { useRouter } from "next/navigation";
import { useState } from "react";

const Chats = () => {
    const router = useRouter();
    const handleClick = async (user: UserType) => {
        const chat = await getChatByUserIDs([user.id]);
        if (!chat) return;
        router.push(`/chats/${chat.ID}`);
    };

    const [isCreatingGroup, setIsCreatingGroup] = useState(false);
    return (
        <div className="flex flex-col gap-y-[8px] flex-1 min-h-0">
            <div className="min-[899px]:hidden flex-1 overflow-x-hidden flex flex-col relative">
                <div
                    className={`${isCreatingGroup ? "-translate-x-full" : "translate-x-0"} h-full transition-transform duration-300 ease-in-out flex flex-col gap-y-[8px]`}
                >
                    {/* only show on mobile */}
                    <CreateGroupButton
                        handleClick={() => setIsCreatingGroup(true)}
                    />
                    <Contacts
                        selectedUsers={[]}
                        handleUserClick={handleClick}
                    />
                </div>

                <CreateGroup
                    handleClose={() => setIsCreatingGroup(false)}
                    createGroupOpen={isCreatingGroup}
                />
            </div>

            <div className="max-[900px]:hidden flex justify-center h-full">
                {/* only show on desktop */}
                Welcome to Ripple
            </div>
        </div>
    );
};

export default Chats;
