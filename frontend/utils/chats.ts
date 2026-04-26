import { use } from "react";
import { api, BASE_CHAT_SERVICE_API_URL } from "./api";

interface ChatType {
    ID: number;
}

export const getChatByUserIDs = async (
    userIDs?: number[],
    rolePermissions?: Map<string, string[]>,
    userRoles?: Map<number, string>,
): Promise<ChatType | undefined> => {
    try {
        if (!userIDs && (!rolePermissions || !userRoles)) return;

        if (!rolePermissions) {
            rolePermissions = new Map<string, string[]>();
            rolePermissions.set("admin", ["send message"]);
            rolePermissions.set("member", ["send message"]);
        }

        if (!userRoles) {
            userRoles = new Map<number, string>();
            if (!userIDs) return;
            for (const user of userIDs) {
                userRoles.set(user, "member");
            }
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
        if (data.error) {
            console.error(data.error);
            return;
        }
        const chat = data.chat;
        return chat;
    } catch (error) {
        console.error(error);
    }
};
