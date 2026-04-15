const ChatPage = () => {
    // we have;
    //  the chat id

    // we need to;
    //  fetch messages in that chat
    //
    return (
        <div className="flex flex-col h-full">
            <input
                type="text"
                placeholder="Type a message here..."
                className="bg-foreground w-full text-background mt-auto outline-none border-none p-[4px] rounded-[4px]"
            />
        </div>
    );
};

export default ChatPage;
