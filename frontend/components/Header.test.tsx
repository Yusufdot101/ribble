import { render, screen } from "@testing-library/react";
import { expect, test } from "vitest";
import Header from "./Header";

test("renders site name", async () => {
    render(<Header />);
    const spanElement = screen.getByTitle("site name");
    expect(spanElement).toBeInTheDocument();
});
