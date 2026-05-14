import { render } from "@testing-library/react";
import { it, beforeEach, expect, vi } from "vitest";
import { useAuthStore } from "@/store/useAuthStore";
import Login from "./page";
import type { ImageProps } from "next/image";

beforeEach(() => {
    useAuthStore.setState({
        accessToken: null,
        isLoggedIn: false,
        userID: null,
    });

    mockPush.mockClear();
});

const mockPush = vi.fn();
vi.mock("next/navigation", () => ({
    useRouter: () => ({
        push: mockPush,
    }),
}));

vi.mock("next/image", () => ({
    default: (props: ImageProps) => <img alt="text image" {...props} />,
}));

it("redirects to home if already logged in", () => {
    useAuthStore.setState({ isLoggedIn: true });
    render(<Login />);
    expect(mockPush).toHaveBeenCalledWith("/");
});

it("does not redirect if not logged in", () => {
    useAuthStore.setState({ isLoggedIn: false });
    render(<Login />);
    expect(mockPush).not.toHaveBeenCalled();
});
