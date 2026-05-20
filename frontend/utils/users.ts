import { useAuthStore } from "@/store/useAuthStore";
import {
    api,
    BASE_CHAT_SERVICE_API_URL,
    BASE_USER_SERVICE_API_URL,
} from "./api";
import { decodeJWT } from "./userIdFromJWT";

export type UserType = {
    id: number;
    sub: string;
    provider: string;
    name: string;
    email: string;
    createdAt: string;
    role: string;
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
        if (!res || !res.ok) {
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
    chatId: number,
    email: string,
): Promise<getAddableChatUsersResponse> => {
    try {
        const res = await api(
            `${BASE_CHAT_SERVICE_API_URL}/chats/${chatId}/addable-users?q=${encodeURIComponent(email)}`,
        );
        if (!res || !res.ok) {
            return { users: [] };
        }
        const body = await res.json();
        const data = body.data;
        return data ?? { users: [] };
    } catch (error) {
        console.error(error);
        return { users: [] };
    }
};

export const verifyAccount = async (
    token: string,
    identityID: string,
): Promise<boolean> => {
    try {
        const url = `${BASE_USER_SERVICE_API_URL}/auth/verify?token=${encodeURIComponent(token)}&identity=${identityID}`;
        const res = await api(url);
        if (!res || !res.ok) {
            alert("an error occurred, please try again");
            return false;
        }
        const body = await res.json();
        if (body.message) {
            alert(body.message);
        }

        const accessToken = body.data;
        if (!accessToken) {
            return false;
        }

        const { payload } = decodeJWT(accessToken);
        if (!payload || !payload.sub) {
            console.error("invalid JWT payload");
            useAuthStore.getState().clearAccessToken();
            return false;
        }

        // Check token expiration
        if (payload.exp && payload.exp * 1000 < Date.now()) {
            console.error("JWT token has expired");
            useAuthStore.getState().clearAccessToken();
            return false;
        }

        const userId = payload.sub;
        if (userId === "") {
            console.error("invalid user ID in JWT");
            useAuthStore.getState().clearAccessToken();
            return false;
        }

        const authStore = useAuthStore.getState();
        authStore.setUserId(userId);
        authStore.setAccessToken(accessToken);
        authStore.setIsLoggedIn(true);
        return true;
    } catch (error) {
        console.error(error);
        return false;
    }
};

export const changeMemberRole = async (
    newRole: string,
    userId: number,
    chatId: number,
): Promise<boolean> => {
    try {
        const res = await api(
            `${BASE_CHAT_SERVICE_API_URL}/chats/${chatId}/users/${userId}`,
            {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ newRole }),
            },
        );
        if (!res || !res.ok) {
            return false;
        }
        const body = await res.json();
        if (body.error) {
            alert(body.error);
            return false;
        }
        if (body.message) {
            alert(body.message);
        }
        return true;
    } catch (error) {
        console.error(error);
        return false;
    }
};
