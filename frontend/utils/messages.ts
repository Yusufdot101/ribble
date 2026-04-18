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
