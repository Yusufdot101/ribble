import { api, BASE_USER_SERVICE_API_URL } from "./api";

export type UserType = {
    id: number;
    sub: string;
    provider: string;
    name: string;
    email: string;
    createdAt: string;
};

export const getUsersByEmail = async (email: string): Promise<UserType[]> => {
    try {
        const res = await api(
            `${BASE_USER_SERVICE_API_URL}/users?email=${encodeURIComponent(email)}`,
        );
        if (!res) {
            return [];
        }
        const data = await res.json();
        return data.users;
    } catch (error) {
        console.error(error);
        return [];
    }
};
