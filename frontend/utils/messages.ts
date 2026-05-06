import { api, BASE_CHAT_SERVICE_API_URL } from "./api";

export interface MessageType {
    ClientID?: string;
    ID: number;
    ChatID: number;
    SenderID: number;
    Content: string;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt: string | null;
    Deleted: boolean;
    Status: "pending" | "delivered" | "failed";
    MessageType?: string;
}

type GetChatMessagesResponse = {
    messages: MessageType[];
};

export const getChatMessages = async (
    chatID: number,
): Promise<GetChatMessagesResponse> => {
    try {
        const baseURL = `${BASE_CHAT_SERVICE_API_URL}/chats/${chatID}/messages`;
        const res = await api(baseURL);
        if (!res || !res.ok) return { messages: [] };
        const body = await res.json();
        if (body.error) {
            console.error(body.error);
            return { messages: [] };
        }
        const data = body.data;
        return data;
    } catch (error) {
        console.error(error);
        return { messages: [] };
    }
};

export const syncChatMessages = async (
    chatID: number,
    lastMessageID: number,
): Promise<GetChatMessagesResponse> => {
    try {
        const baseURL = `${BASE_CHAT_SERVICE_API_URL}/chats/${chatID}/messages/sync?lastMessageId=${lastMessageID}`;
        const res = await api(baseURL);
        if (!res || !res.ok) return { messages: [] };
        const body = await res.json();
        const data = body.data;
        if (body.error) {
            console.error(body.error);
            return { messages: [] };
        }
        return data;
    } catch (error) {
        console.error(error);
        return { messages: [] };
    }
};

export const deleteMessage = async (chatID: number, messageID: number) => {
    try {
        const baseURL = `${BASE_CHAT_SERVICE_API_URL}/chats/${chatID}/messages`;
        const res = await api(`${baseURL}/${messageID}`, {
            method: "DELETE",
        });
        if (!res) {
            alert("an error occured deleting message");
            return;
        }
        const data = await res.json();
        if (data.error) {
            alert("an error occured deleting message");
        }
    } catch (error) {
        console.error(error);
    }
};

export const editMessage = async (
    chatID: number,
    messageID: number,
    newContent: string,
): Promise<boolean | undefined> => {
    try {
        const baseURL = `${BASE_CHAT_SERVICE_API_URL}/chats/${chatID}/messages`;
        const res = await api(`${baseURL}/${messageID}`, {
            method: "PATCH",
            body: JSON.stringify({ newContent }),
            headers: {
                "Content-Type": "application/json",
            },
        });
        if (!res) {
            alert("an error occured editing message");
            return false;
        }
        const data = await res.json();
        if (data.error) {
            alert("an error occured editing message: " + data.error);
            return false;
        }

        return true;
    } catch (error) {
        console.error(error);
    }
};
