import { api, BASE_CHAT_SERVICE_API_URL } from "./api";
import { UserType } from "./users";

export const addUsersToGroup = async (chatId: number, userIds: number[]) => {
    const res = await api(
        `${BASE_CHAT_SERVICE_API_URL}/chats/${chatId}/addToGroup`,
        {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ userIds }),
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

export const removeUserFromGroup = async (chatId: number, userId: number) => {
    const res = await api(
        `${BASE_CHAT_SERVICE_API_URL}/chats/${chatId}/users/${userId}`,
        { method: "DELETE" },
    );

    if (!res) throw new Error("No response from chat service");
    if (!res.ok) {
        const errBody = await res.text();
        throw new Error(
            errBody || `Failed to remove user from group (${res.status})`,
        );
    }
};

function addTime(value: number, unit: string): Date {
    const d = new Date();

    switch (unit) {
        case "hours":
            d.setHours(d.getHours() + value);
            break;
        case "days":
            d.setDate(d.getDate() + value);
            break;
        case "weeks":
            d.setDate(d.getDate() + value * 7);
            break;
        case "months":
            d.setMonth(d.getMonth() + value);
            break;
        case "years":
            d.setFullYear(d.getFullYear() + value);
            break;
    }

    return d;
}

export const banUser = async (
    chatId: number,
    userId: number,
    reason: string,
    timeFrame: string,
    timeValue: number,
): Promise<boolean | undefined> => {
    const expiry =
        timeValue !== -1 ? addTime(timeValue, timeFrame).toISOString() : "";
    const res = await api(`${BASE_CHAT_SERVICE_API_URL}/chats/${chatId}/bans`, {
        method: "POST",
        body: JSON.stringify({
            userId: userId,
            reason: reason,
            expiresAt: expiry !== "" ? expiry : undefined,
        }),
        headers: {
            "Content-Type": "application/json",
        },
    });

    if (!res) {
        throw new Error("No response from chat service");
    }
    if (!res.ok) {
        const errBody = await res.text();
        throw new Error(
            errBody || `Failed to ban user from group (${res.status})`,
        );
    }
    return true;
};

type ChatUsers = {
    users: UserType[];
};

export const getBannedUsers = async (
    chatId: number,
    query: string,
): Promise<ChatUsers> => {
    try {
        const res = await api(
            `${BASE_CHAT_SERVICE_API_URL}/chats/${chatId}/bans?q=${encodeURIComponent(query)}`,
        );
        if (!res) {
            throw new Error("No response from chat service");
        }
        if (!res.ok) {
            const errBody = await res.text();
            throw new Error(
                errBody || `Failed to get banned users (${res.status})`,
            );
        }
        const body = await res.json();
        return body.data ?? { users: [] };
    } catch (error) {
        console.error(error);
        return { users: [] };
    }
};

export const unbanUser = async (
    chatId: number,
    userId: number,
): Promise<boolean | undefined> => {
    const res = await api(
        `${BASE_CHAT_SERVICE_API_URL}/chats/${chatId}/bans/${userId}`,
        {
            method: "DELETE",
        },
    );

    if (!res) {
        throw new Error("No response from chat service");
    }
    if (!res.ok) {
        const errBody = await res.text();
        throw new Error(
            errBody || `Failed to unban user from group (${res.status})`,
        );
    }
    return true;
};
