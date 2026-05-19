import { render } from "@testing-library/react";
import { it, beforeEach, expect, vi } from "vitest";
import { useAuthStore } from "@/store/useAuthStore";
import Logout from "./page";
import { logout } from "@/utils/logout";
import { ImageProps } from "next/image";

beforeEach(() => {
    useAuthStore.setState({
        accessToken: null,
        isLoggedIn: false,
        userId: null,
    });

    mockPush.mockClear();
});

const mockPush = vi.fn();
vi.mock("next/navigation", () => ({
    useRouter: () => ({
        push: mockPush,
    }),
}));

vi.mock("@/utils/logout", () => ({
    logout: vi.fn(),
}));

vi.mock("next/image", () => ({
    default: (props: ImageProps) => <img alt="text image" {...props} />,
}));

it("redirects to home if not logged in", () => {
    useAuthStore.setState({ isLoggedIn: false });
    render(<Logout />);
    expect(mockPush).toHaveBeenCalledWith("/");
});

it("does not redirect if logged in", () => {
    useAuthStore.setState({ isLoggedIn: true });
    render(<Logout />);
    expect(mockPush).not.toHaveBeenCalled();
});

it("calls logout function on logout", async () => {
    useAuthStore.setState({ isLoggedIn: true });
    render(<Logout />);
    expect(logout).toHaveBeenCalled();
});
