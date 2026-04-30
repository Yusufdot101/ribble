import { api, BASE_CHAT_SERVICE_API_URL } from "./api";

export const addUsersToGroup = async (chatID: number, userIDs: number[]) => {
    const res = await api(
        `${BASE_CHAT_SERVICE_API_URL}/chats/${chatID}/addToGroup`,
        {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ userIDs }),
        },
    );

    if (!res) throw new Error("No response from chat service");
    if (!res.ok) {
        const errBody = await res.text();
        throw new Error(
            errBody || `Failed to add users to group (${res.status})`,
        );
    }
};
