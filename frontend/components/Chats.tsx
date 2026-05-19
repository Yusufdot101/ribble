import UserCardSkeleton from "@/components/UserCardSkeleton";
import { ChatType } from "@/utils/chats";
import ChatCard from "./ChatCard";

interface Props {
    selectedChats: number[];
    handleChatClick: (chatId: number) => void;
    chats: ChatType[];
    isLoading: boolean;
}

const Chats = ({ selectedChats, handleChatClick, chats, isLoading }: Props) => {
    return (
        <div className="flex-1 flex flex-col gap-y-[8px] rounded-[4px] gap-y-[8px]">
            <div
                className={`${!isLoading ? "hidden" : ""} flex flex-col transition-all duration-300`}
            >
                {[...Array(4).keys()].map((el) => (
                    <UserCardSkeleton key={el} index={el} />
                ))}
            </div>

            {chats.length === 0 && !isLoading ? (
                <p className="w-full text-center">No Data</p>
            ) : null}

            <div
                className={`${chats?.length === 0 ? "opacity-0 blur-sm" : ""} flex flex-col transition-all duration-300`}
            >
                {chats.map((chat) => (
                    <ChatCard
                        handleClick={handleChatClick}
                        activeChats={selectedChats}
                        key={chat.id}
                        chat={chat}
                    />
                ))}
            </div>
        </div>
    );
};

export default Chats;
