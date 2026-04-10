import { useAuthStore } from "@/store/useAuthStore";
import { BASE_API_URL } from "./api";

export const logout = async () => {
    try {
        const res = await fetch(`${BASE_API_URL}/auth/logout`, {
            method: "POST",
            credentials: "include",
        });
        if (!res.ok) {
            console.error("error loging out", res);
        }
        useAuthStore.getState().clearAccessToken();
    } catch (error) {
        console.error(error);
    }
};
