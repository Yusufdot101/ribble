import { socketManager } from "@/utils/socket-manager";
import { create } from "zustand";

type AuthStore = {
    accessToken: string | null;
    isLoggedIn: boolean;
    userID: number | null;
    setAccessToken: (token: string) => void;
    setIsLoggedIn: (isLoggedIn: boolean) => void;
    setUserID: (userID: number) => void;
    clearAccessToken: () => void;
};

export const useAuthStore = create<AuthStore>((set) => ({
    accessToken: null,
    isLoggedIn: false,
    userID: null,
    setAccessToken: (token: string) => set({ accessToken: token }),
    setIsLoggedIn: (isLoggedIn: boolean) => set({ isLoggedIn: isLoggedIn }),
    setUserID: (userID: number) => set({ userID: userID }),
    clearAccessToken: () => {
        try {
            socketManager.close();
        } catch (err) {
            console.error("Failed to close socket:", err);
        } finally {
            set({
                accessToken: null,
                isLoggedIn: false,
                userID: null,
            });
        }
    },
}));
