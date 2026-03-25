import { render, screen } from "@testing-library/react";
import { it, beforeEach, expect, vi } from "vitest";
import { useAuthStore } from "@/store/useAuthStore";
import userEvent from "@testing-library/user-event";
import Login from "./page";
import { BASE_API_URL } from "@/utils/api";

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

vi.mock("@/components/Icon", () => ({
    default: (props: any) => <img {...props} />,
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

it("navigates to google auth on click", async () => {
    render(<Login />);
    await userEvent.click(screen.getByText(/continue with/i));
    expect(mockPush).toHaveBeenCalledWith(`${BASE_API_URL}/auth/google`);
});
