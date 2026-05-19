import { render, screen } from "@testing-library/react";
import { it, beforeEach, expect, vi } from "vitest";
import { useAuthStore } from "@/store/useAuthStore";
import userEvent from "@testing-library/user-event";
import Signup from "./page";
import { BASE_USER_SERVICE_API_URL } from "@/utils/api";
import type { ImageProps } from "next/image";

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

vi.mock("next/image", () => ({
    default: (props: ImageProps) => <img alt="text image" {...props} />,
}));

it("redirects to home if already logged in", () => {
    useAuthStore.setState({ isLoggedIn: true });
    render(<Signup />);
    expect(mockPush).toHaveBeenCalledWith("/");
});

it("does not redirect if not logged in", () => {
    useAuthStore.setState({ isLoggedIn: false });
    render(<Signup />);
    expect(mockPush).not.toHaveBeenCalled();
});

it("navigates to google auth on click", async () => {
    render(<Signup />);
    await userEvent.click(screen.getByText(/continue with/i));
    expect(mockPush).toHaveBeenCalledWith(
        `${BASE_USER_SERVICE_API_URL}/auth/google`,
    );
});
