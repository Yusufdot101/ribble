import { socketManager } from "@/utils/socket-manager";
import { create } from "zustand";

type AuthStore = {
    accessToken: string | null;
    isLoggedIn: boolean;
    userId: number | null;
    setAccessToken: (token: string) => void;
    setIsLoggedIn: (isLoggedIn: boolean) => void;
    setUserId: (userId: number) => void;
    clearAccessToken: () => void;
};

export const useAuthStore = create<AuthStore>((set) => ({
    accessToken: null,
    isLoggedIn: false,
    userId: null,
    setAccessToken: (token: string) => set({ accessToken: token }),
    setIsLoggedIn: (isLoggedIn: boolean) => set({ isLoggedIn: isLoggedIn }),
    setUserId: (userId: number) => set({ userId: userId }),
    clearAccessToken: () => {
        try {
            socketManager.close();
        } catch (err) {
            console.error("Failed to close socket:", err);
        } finally {
            set({
                accessToken: null,
                isLoggedIn: false,
                userId: null,
            });
        }
    },
}));
