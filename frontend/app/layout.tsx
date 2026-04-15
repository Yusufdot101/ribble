import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import Header from "@/components/Header";

const geistSans = Geist({
    variable: "--font-geist-sans",
    subsets: ["latin"],
});

const geistMono = Geist_Mono({
    variable: "--font-geist-mono",
    subsets: ["latin"],
});

export const metadata: Metadata = {
    title: "Ripple",
    description: "A site to chat",
};

export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html
            lang="en"
            className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}
        >
            <body
                suppressHydrationWarning
                className="min-h-screen flex flex-col px-[24px] py-[16px] gap-y-[12px]"
            >
                <Header />
                <div className="flex-1 flex flex-col">{children}</div>
            </body>
        </html>
    );
}
