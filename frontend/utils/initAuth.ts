import { useAuthStore } from "@/store/useAuthStore";
import { BASE_USER_SERVICE_API_URL } from "./api";
import { decodeJWT } from "./userIdFromJWT";

export const initAuth = async () => {
    try {
        const res = await fetch(
            `${BASE_USER_SERVICE_API_URL}/auth/refreshtoken`,
            {
                method: "GET",
                credentials: "include",
            },
        );
        if (!res.ok) {
            useAuthStore.getState().setIsLoggedIn(false);
            return;
        }
        const data = await res.json();
        const token = data.accessToken;
        const { payload } = decodeJWT(token ?? "");

        if (!payload || !payload.sub) {
            console.error("invalid JWT payload");
            useAuthStore.getState().clearAccessToken();
            return;
        }

        const expired = new Date(payload.exp * 1000) < new Date();
        if (expired) {
            console.error("expired JWT");
            useAuthStore.getState().clearAccessToken();
            return;
        }

        const userId = +payload.sub;
        if (isNaN(userId)) {
            console.error("invalid user ID in JWT");
            useAuthStore.getState().clearAccessToken();
            return;
        }

        useAuthStore.getState().setAccessToken(token);
        useAuthStore.getState().setUserID(userId);
        useAuthStore.getState().setIsLoggedIn(true);
    } catch (error) {
        console.error(error);
    }
};
