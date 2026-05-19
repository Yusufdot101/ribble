import { api, BASE_CHAT_SERVICE_API_URL } from "./api";

export interface MessageType {
    clientId: string;
    id: number;
    chatId: number;
    senderId: number;
    content: string;
    createdAt: string;
    updatedAt: string;
    deletedAt: string | null;
    deleted: boolean;
    status: "pending" | "delivered" | "failed";
    messageType?: string;
}

type GetChatMessagesResponse = {
    messages: MessageType[];
};

export const getChatMessages = async (
    chatId: number,
): Promise<GetChatMessagesResponse> => {
    try {
        const baseURL = `${BASE_CHAT_SERVICE_API_URL}/chats/${chatId}/messages`;
        const res = await api(baseURL);
        if (!res || !res.ok) return { messages: [] };
        const body = await res.json();
        if (body.error) {
            console.error(body.error);
            return { messages: [] };
        }
        const data = body.data;
        if (!data || !data.messages) {
            console.error("Invalid response structure");
            return { messages: [] };
        }
        return data;
    } catch (error) {
        console.error(error);
        return { messages: [] };
    }
};

export const syncChatMessages = async (
    chatId: number,
    lastMessageID: number,
): Promise<GetChatMessagesResponse> => {
    try {
        const baseURL = `${BASE_CHAT_SERVICE_API_URL}/chats/${chatId}/messages/sync?lastMessageId=${lastMessageID}`;
        const res = await api(baseURL);
        if (!res || !res.ok) return { messages: [] };
        const body = await res.json();
        const data = body.data;
        if (body.error) {
            console.error(body.error);
            return { messages: [] };
        }
        if (!data || !data.messages) {
            console.error("Invalid response structure");
            return { messages: [] };
        }
        return data;
    } catch (error) {
        console.error(error);
        return { messages: [] };
    }
};

export const deleteMessage = async (chatId: number, messageID: number) => {
    try {
        const baseURL = `${BASE_CHAT_SERVICE_API_URL}/chats/${chatId}/messages`;
        const res = await api(`${baseURL}/${messageID}`, {
            method: "DELETE",
        });
        if (!res) {
            alert("an error occured deleting message");
            return;
        }
        const body = await res.json();
        if (body.error) {
            alert("an error occured deleting message");
        }
    } catch (error) {
        console.error(error);
    }
};

export const editMessage = async (
    chatId: number,
    messageID: number,
    newContent: string,
): Promise<boolean | undefined> => {
    try {
        const baseURL = `${BASE_CHAT_SERVICE_API_URL}/chats/${chatId}/messages`;
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
        const body = await res.json();
        if (body.error) {
            alert("an error occured editing message: " + body.error);
            return false;
        }

        return true;
    } catch (error) {
        console.error(error);
    }
};
