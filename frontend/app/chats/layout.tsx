"use client";
import Contacts from "@/components/Contacts";
import CreateGroup from "@/components/CreateGroup";
import CreateGroupButton from "@/components/CreateGroupButton";
import { getChatByUserIDs } from "@/utils/chats";
import { UserType } from "@/utils/users";
import { useRouter } from "next/navigation";
import { useState } from "react";
export default function ChatsLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    const [activeUser, setActiveUser] = useState<number>();
    const router = useRouter();
    const handleClick = async (clickedUser: UserType) => {
        setActiveUser(clickedUser.id);
        const chat = await getChatByUserIDs([clickedUser.id]);
        if (!chat) return;
        router.push(`/chats/${chat.ID}`);
    };

    const [isCreatingGroup, setIsCreatingGroup] = useState(false);
    return (
        <div className="flex flex-col gap-y-[8px] flex-1 min-h-0">
            <div className="gap-x-[8px] flex-1 flex min-h-0 ">
                {/* Sidebar */}
                {/* Hide in mobile */}
                <div className="w-[40%] max-[900px]:hidden flex flex-col gap-y-[8px] relative overflow-x-hidden">
                    <div
                        className={`${isCreatingGroup ? "-translate-x-full" : "translate-x-0"} h-full transition-transform duration-300 ease-in-out flex flex-col gap-y-[8px]`}
                    >
                        <CreateGroupButton
                            handleClick={() => setIsCreatingGroup(true)}
                        />
                        <div className="flex-1 border-r-1 border-foreground">
                            <Contacts
                                selectedUsers={activeUser ? [activeUser] : []}
                                handleUserClick={handleClick}
                            />
                        </div>
                    </div>

                    <CreateGroup
                        handleClose={() => setIsCreatingGroup(false)}
                        createGroupOpen={isCreatingGroup}
                    />
                </div>

                {/* Main content */}
                <div className="flex-1 min-h-0 flex flex-col">{children}</div>
            </div>
        </div>
    );
}
