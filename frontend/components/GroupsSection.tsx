import { ChatType } from "@/utils/chats";
import Chats from "./Chats";

type Props = {
    selectedChats: number[];
    handleChatClick: (chatId: number) => void;
    chats: ChatType[];
    isLoading: boolean;
};

const GroupsSection = ({
    selectedChats,
    handleChatClick,
    chats,
    isLoading,
}: Props) => {
    return (
        <div className="flex flex-col">
            <span className="text-foreground/70">groups</span>
            <Chats
                handleChatClick={handleChatClick}
                selectedChats={selectedChats}
                chats={chats}
                isLoading={isLoading}
            />
        </div>
    );
};

export default GroupsSection;
