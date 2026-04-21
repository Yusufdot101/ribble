import { api, BASE_CHAT_SERVICE_API_URL } from "./api";

const baseURL = `${BASE_CHAT_SERVICE_API_URL}/messages`;

export interface MessageType {
    ID: number;
    ChatID: number;
    SenderID: number;
    Content: string;
    CreatedAt: string;
    UpdatedAt: string;
}

export const getChatMessages = async (
    chatID: number,
): Promise<MessageType[]> => {
    try {
        const res = await api(baseURL, {
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

export const deleteMessage = async (messageID: number) => {
    try {
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

export const editMessage = async (messageID: number, newContent: string) => {
    try {
        const res = await api(`${baseURL}/${messageID}`, {
            method: "PATCH",
            body: JSON.stringify({ newContent }),
            headers: {
                "Content-Type": "application/json",
            },
        });
        if (!res) {
            alert("an error occured editing message");
            return;
        }
        const data = await res.json();
        if (data.error) {
            alert("an error occured editing message: " + data.error);
        }
    } catch (error) {
        console.error(error);
    }
};
