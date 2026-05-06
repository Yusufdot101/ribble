import { api, BASE_CHAT_SERVICE_API_URL } from "./api";

export type PermissionType = {
    id: number;
    name: string;
};

type GetUserPermissionsResponse = {
    permissions: PermissionType[];
};
export const getUserPermissions = async (
    chatID: number,
): Promise<GetUserPermissionsResponse> => {
    try {
        // /chats/:chatId/permissions
        const res = await api(
            `${BASE_CHAT_SERVICE_API_URL}/chats/${chatID}/permissions`,
        );
        if (!res) return { permissions: [] };
        const body = await res.json();
        const data = body.data;
        return data;
    } catch (error) {
        console.error(error);
        return { permissions: [] };
    }
};
