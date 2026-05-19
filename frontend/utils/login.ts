import { useAuthStore } from "../store/useAuthStore";
import { BASE_USER_SERVICE_API_URL } from "./api";
import { decodeJWT } from "./userIdFromJWT";

export const login = async (
    handleError: (error: string) => void,
    email: string,
    password: string,
): Promise<boolean> => {
    try {
        const res = await fetch(BASE_USER_SERVICE_API_URL + "/auth/login", {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ email, password }),
        });
        const body = await res.json();
        if (!res.ok) {
            handleError(
                body.error || `Request failed with status ${res.status}`,
            );
            return false;
        }
        if (body.error) {
            handleError(body.error);
            return false;
        }

        if (body.message) {
            handleError(body.message);
            return false;
        }

        const accessToken = body.data;
        if (!accessToken) {
            handleError("An error occurred. Please try again later");
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
        handleError("An error occurred. Please try again later");
        console.error("Login error:", error);
        return false;
    }
};
