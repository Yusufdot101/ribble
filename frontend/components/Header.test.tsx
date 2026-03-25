import { render, screen } from "@testing-library/react";
import { expect, test, vi } from "vitest";
import Header from "./Header";
import { useAuthStore } from "@/store/useAuthStore";
import { beforeEach } from "node:test";

test("renders site name", async () => {
    render(<Header />);
    const spanElement = screen.getByText(/ripple/i);
    expect(spanElement).toBeInTheDocument();
});

beforeEach(() => {
    useAuthStore.setState({
        accessToken: null,
        isLoggedIn: false,
        userID: null,
    });
});

test("renders login/logout button", async () => {
    render(<Header />);
    const linkElement = screen.getByText(/^(login|logout)$/i);
    expect(linkElement).toBeInTheDocument();
});

test("shows logout when user is logged in", async () => {
    useAuthStore.setState({ isLoggedIn: true });
    render(<Header />);
    const linkElement = screen.getByText(/^logout$/i);
    expect(linkElement).toBeInTheDocument();
});

test("shows login when user is logged out", async () => {
    useAuthStore.setState({ isLoggedIn: false });
    render(<Header />);
    const linkElement = screen.getByText(/^login$/i);
    expect(linkElement).toBeInTheDocument();
});
