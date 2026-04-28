import { api, BASE_CHAT_SERVICE_API_URL } from "./api";

export interface MessageType {
    ID: number;
    ChatID: number;
    SenderID: number;
    Content: string;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt: string | null;
    Deleted: boolean;
}

export const getChatMessages = async (
    chatID: number,
): Promise<MessageType[]> => {
    try {
        const baseURL = `${BASE_CHAT_SERVICE_API_URL}/chats/${chatID}/messages`;
        const res = await api(baseURL);
        if (!res || !res.ok) return [];
        const data = await res.json();
        if (data.error) {
            console.error(data.error);
            return [];
        }
        const messages = data.messages;
        return Array.isArray(messages) ? messages : [];
    } catch (error) {
        console.error(error);
        return [];
    }
};

export const syncChatMessages = async (
    chatID: number,
    lastMessageID: number,
): Promise<MessageType[]> => {
    try {
        const baseURL = `${BASE_CHAT_SERVICE_API_URL}/chats/${chatID}/messages/sync?lastMessageId=${lastMessageID}`;
        const res = await api(baseURL);
        if (!res || !res.ok) return [];
        const data = await res.json();
        if (data.error) {
            console.error(data.error);
            return [];
        }
        const messages = data.messages;
        return Array.isArray(messages) ? messages : [];
    } catch (error) {
        console.error(error);
        return [];
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
