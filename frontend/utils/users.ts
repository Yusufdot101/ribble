import {
    api,
    BASE_CHAT_SERVICE_API_URL,
    BASE_USER_SERVICE_API_URL,
} from "./api";

export type UserType = {
    id: number;
    sub: string;
    provider: string;
    name: string;
    email: string;
    createdAt: string;
};

type getUsersByEmailResponse = {
    users: UserType[];
};
export const getUsersByEmail = async (
    email: string,
): Promise<getUsersByEmailResponse> => {
    try {
        const res = await api(
            `${BASE_USER_SERVICE_API_URL}/users?email=${encodeURIComponent(email)}`,
        );
        if (!res) {
            return { users: [] };
        }
        const data = await res.json();
        return data;
    } catch (error) {
        console.error(error);
        return { users: [] };
    }
};

type getAddableChatUsersResponse = {
    users: UserType[];
};
export const getAddableChatUsers = async (
    chatID: number,
    email: string,
): Promise<getAddableChatUsersResponse> => {
    try {
        const res = await api(
            `${BASE_CHAT_SERVICE_API_URL}/chats/${chatID}/addable-users?q=${encodeURIComponent(email)}`,
        );
        if (!res) {
            return { users: [] };
        }
        const body = await res.json();
        const data = body.data;
        return data;
    } catch (error) {
        console.error(error);
        return { users: [] };
    }
};
