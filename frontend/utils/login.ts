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
        if (body.error) {
            handleError(body.error);
            return false;
        }

        if (body.message) {
            alert(body.message);
            return false;
        }

        const accessToken = body.data;
        if (!accessToken) {
            alert("An error occurred. Please try again later");
            return false;
        }

        const { payload } = decodeJWT(accessToken);
        if (!payload || !payload.sub) {
            console.error("invalid JWT payload");
            useAuthStore.getState().clearAccessToken();
            return false;
        }

        const userID = payload.sub;
        if (userID === "") {
            console.error("invalid user ID in JWT");
            useAuthStore.getState().clearAccessToken();
            return false;
        }

        useAuthStore.getState().setUserID(userID);
        useAuthStore.getState().setAccessToken(accessToken);
        useAuthStore.getState().setIsLoggedIn(true);
        return true;
    } catch (error) {
        alert("An error occurred. Please try again later");
        console.log(error);
        return false;
    }
};
