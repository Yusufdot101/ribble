import { api } from "./api";

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
        const res = await api(`/users?email=${email}`);
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
