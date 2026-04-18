import { api, BASE_CHAT_SERVICE_API_URL } from "./api";

export interface MessageType {
    ID: number;
    ChatID: number;
    SenderID: number;
    Content: string;
    CreatedAt: string;
}

export const getChatMessages = async (
    chatID: number,
): Promise<MessageType[]> => {
    try {
        const res = await api(`${BASE_CHAT_SERVICE_API_URL}/messages`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ chatID }),
        });
        if (!res) return [];
        const data = await res.json();
        const messages = data.messages;
        return messages;
    } catch (error) {
        console.error(error);
        return [];
    }
};
