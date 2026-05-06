import { api, BASE_CHAT_SERVICE_API_URL } from "./api";
import { UserType } from "./users";

export interface ChatType {
    id: number;
    name: string;
    isGroup: boolean;
}

export const getChatByUserIDs = async (
    userIDs?: number[],
    rolePermissions?: Map<string, string[]>,
    userRoles?: Map<number, string>,
    chatName?: string,
    isGroup?: boolean,
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
                name: chatName,
                isGroup: isGroup,
            }),
        });

        if (!res) return;
        const body = await res.json();
        if (body.error) {
            console.error(body.error);
            return;
        }
        const chat = body.data;
        return chat;
    } catch (error) {
        console.error(error);
    }
};

export const getChatByID = async (
    chatID: number,
): Promise<ChatType | undefined> => {
    try {
        const res = await api(`${BASE_CHAT_SERVICE_API_URL}/chats/${chatID}`);
        if (!res) return;
        const body = await res.json();
        if (body.error) {
            console.error(body.error);
            return;
        }
        const chat = body.data;
        return chat;
    } catch (error) {
        console.error(error);
    }
};

export const getChatUsers = async (
    chatID: number,
): Promise<UserType[] | undefined> => {
    try {
        const res = await api(
            `${BASE_CHAT_SERVICE_API_URL}/chats/${chatID}/users`,
        );
        if (!res) return;
        const body = await res.json();
        if (body.error) {
            console.error(body.error);
            return;
        }
        const data = body.data;
        const users = data.users;
        return users;
    } catch (error) {
        console.error(error);
    }
};

export type ConversationDataType = {
    chats: ChatType[];
    groups: ChatType[];
    contacts: UserType[];
};

export const getConversations = async (
    query: string,
): Promise<ConversationDataType | undefined> => {
    try {
        const res = await api(
            `${BASE_CHAT_SERVICE_API_URL}/conversations?q=${encodeURIComponent(query)}`,
        );
        if (!res) return;
        const body = await res.json();
        if (body.error) {
            console.error(body.error);
            return;
        }
        return body.data;
    } catch (error) {
        console.error(error);
    }
};
