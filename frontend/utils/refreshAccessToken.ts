import { useAuthStore } from "@/store/useAuthStore";
import { BASE_USER_SERVICE_API_URL } from "./api";
import { decodeJWT } from "./userIdFromJWT";

export async function refreshAccessToken() {
    const res = await fetch(`${BASE_USER_SERVICE_API_URL}/auth/refreshtoken`, {
        credentials: "include",
    });

    if (!res.ok) {
        console.error(
            `Failed to refresh access token: ${res.status} ${res.statusText}`,
        );
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

    const userId = +payload.sub;
    if (isNaN(userId)) {
        console.error("invalid user ID in JWT");
        useAuthStore.getState().clearAccessToken();
        return;
    }

    useAuthStore.getState().setAccessToken(token);
    useAuthStore.getState().setUserID(userId);
    useAuthStore.getState().setIsLoggedIn(true);
}
