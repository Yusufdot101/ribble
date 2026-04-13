"use client";
import { useParams } from "next/navigation";

const ChatPage = () => {
    // we have;
    //  the chat id

    // we need to;
    //  fetch messages in that chat
    //

    const params = useParams();
    const id = params.id;
    return <div>chat: {id}</div>;
};

export default ChatPage;
