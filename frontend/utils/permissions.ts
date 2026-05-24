import { api, BASE_CHAT_SERVICE_API_URL } from "./api";

export type PermissionType = {
    id: number;
    name: string;
};

type GetUserPermissionsResponse = {
    permissions: PermissionType[];
};
export const getUserPermissions = async (
    chatId: number,
): Promise<GetUserPermissionsResponse> => {
    try {
        const res = await api(
            `${BASE_CHAT_SERVICE_API_URL}/chats/${chatId}/permissions`,
        );
        if (!res) return { permissions: [] };
        const body = await res.json();
        const data = body.data;
        if (body.error || !data) {
            console.error(body.error);
            return { permissions: [] };
        }
        return data;
    } catch (error) {
        console.error(error);
        return { permissions: [] };
    }
};

type GetRolesPermissionsResponse = {
    rolePermissions: Record<string, PermissionType[]>;
};
export const getRolesPermissions = async (
    chatId: number,
    roles: string[],
): Promise<GetRolesPermissionsResponse> => {
    try {
        const res = await api(
            `${BASE_CHAT_SERVICE_API_URL}/chats/${chatId}/roles/permissions`,
            {
                method: "POST",
                body: JSON.stringify({
                    roles,
                }),
                headers: {
                    "Content-Type": "application/json",
                },
            },
        );
        if (!res) {
            return { rolePermissions: {} };
        }
        const body = await res.json();
        const data = body.data;
        if (body.error || !data) {
            console.error(body.error);
            return { rolePermissions: {} };
        }
        return data;
    } catch (error) {
        console.error(error);
        return { rolePermissions: {} };
    }
};

export const updateRolePermissions = async (
    chatId: number,
    role: string,
    permission: string,
    action: "grant" | "revoke",
): Promise<boolean> => {
    try {
        const res = await api(
            `${BASE_CHAT_SERVICE_API_URL}/chats/${chatId}/roles/${role}/permissions`,
            {
                method: "PATCH",
                body: JSON.stringify({
                    permission,
                    action,
                }),
                headers: {
                    "Content-Type": "application/json",
                },
            },
        );
        if (!res) return false;
        const body = await res.json();
        if (body.message) {
            alert(body.message);
        }
        if (body.error) {
            console.error(body.error);
            return false;
        }
        return true;
    } catch (error) {
        console.error(error);
        return false;
    }
};
