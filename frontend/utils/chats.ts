import { api, BASE_CHAT_SERVICE_API_URL } from "./api";

interface ChatType {
    ID: number;
}

export const getChatByUserIDs = async (
    userIDs: number[],
): Promise<ChatType | undefined> => {
    try {
        const rolePermissions = new Map<string, string[]>();
        rolePermissions.set("admin", ["sendMessages"]);
        rolePermissions.set("member", ["sendMessages"]);

        const userRoles = new Map<number, string>();
        for (const user of userIDs) {
            userRoles.set(user, "member");
        }

        const res = await api(`${BASE_CHAT_SERVICE_API_URL}/chats`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                rolePermissions: Object.fromEntries(rolePermissions),
                userRoles: Object.fromEntries(userRoles),
            }),
        });
        if (!res) return;
        const data = await res.json();
        const chat = data.chat;
        return chat;
    } catch (error) {
        console.error(error);
    }
};
