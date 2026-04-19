import { useAuthStore } from "@/store/useAuthStore";
import { refreshAccessToken } from "./refreshAccessToken";

export const BASE_USER_SERVICE_API_URL =
    process.env.NEXT_PUBLIC_BASE_USER_SERVICE_API_URL;
export const BASE_CHAT_SERVICE_API_URL =
    process.env.NEXT_PUBLIC_BASE_CHAT_SERVICE_API_URL;

export const api = async (path: string, options: RequestInit = {}) => {
    const { accessToken } = useAuthStore.getState();
    const url = path.startsWith("http") ? path : `${path}`;

    try {
        let res = await fetch(url, {
            ...options,
            credentials: "include",
            headers: {
                ...(options.headers || {}),
                Authorization: accessToken ? `Bearer ${accessToken}` : "",
            },
        });

        if (res.status === 401) {
            await refreshAccessToken();
            const newToken = useAuthStore.getState().accessToken;

            res = await fetch(url, {
                ...options,
                credentials: "include",
                headers: {
                    ...(options.headers || {}),
                    Authorization: newToken ? `Bearer ${newToken}` : "",
                },
            });

            if (!res.ok) {
                useAuthStore.getState().clearAccessToken(); // because the refresh token didn't refresh access token successfully
                // alert("please login to use this feature.");
                return undefined;
            }
        }
        return res;
    } catch (error) {
        console.error(error);
        return undefined;
    }
};
