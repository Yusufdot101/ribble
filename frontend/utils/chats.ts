import { api, BASE_CHAT_SERVICE_API_URL } from "./api";

interface ChatType {
    ID: number;
}

export const getChatByUserIDs = async (
    userIDs: number[],
): Promise<ChatType | undefined> => {
    try {
        const res = await api(`${BASE_CHAT_SERVICE_API_URL}/chats/find`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ userIDs }),
        });
        if (!res) return;
        const data = await res.json();
        console.log(data);
        const chat = data.chat;
        return chat;
    } catch (error) {
        console.error(error);
    }
};
