import { useAuthStore } from "@/store/useAuthStore";
import { BASE_USER_SERVICE_API_URL } from "./api";
import { decodeJWT } from "./userIdFromJWT";

export async function refreshAccessToken() {
    const res = await fetch(`${BASE_USER_SERVICE_API_URL}/auth/refreshtoken`, {
        method: "PUT",
        credentials: "include", // important! sends cookie
    });

    if (!res.ok) {
        console.error(
            `Failed to refresh access token: ${res.status} ${res.statusText}`,
        );
        return;
    }

    const data = await res.json();
    const token = data.access_token;
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
