import { api, BASE_CHAT_SERVICE_API_URL } from "./api";
import { UserType } from "./users";

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

export const removeUserFromGroup = async (chatID: number, userID: number) => {
    const res = await api(
        `${BASE_CHAT_SERVICE_API_URL}/chats/${chatID}/users/${userID}`,
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
    chatID: number,
    userID: number,
    reason: string,
    timeFrame: string,
    timeValue: number,
): Promise<boolean | undefined> => {
    const expiry =
        timeValue !== -1 ? addTime(timeValue, timeFrame).toISOString() : "";
    const res = await api(`${BASE_CHAT_SERVICE_API_URL}/chats/${chatID}/bans`, {
        method: "POST",
        body: JSON.stringify({
            userId: userID,
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

export const getBannedUsers = async (
    chatID: number,
    query: string,
): Promise<UserType[]> => {
    try {
        const res = await api(
            `${BASE_CHAT_SERVICE_API_URL}/chats/${chatID}/bans?q=${encodeURIComponent(query)}`,
        );
        if (!res) {
            return [];
        }
        const data = await res.json();
        return data.users;
    } catch (error) {
        console.error(error);
        return [];
    }
};

export const unbanUser = async (
    chatID: number,
    userID: number,
): Promise<boolean | undefined> => {
    const res = await api(
        `${BASE_CHAT_SERVICE_API_URL}/chats/${chatID}/bans/${userID}`,
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
