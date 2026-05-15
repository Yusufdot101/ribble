"use client";
import Sidebar from "@/components/Sidebar";

export default function ChatsLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <div className="flex flex-col gap-y-[8px] flex-1 min-h-0">
            <div className="gap-x-[8px] flex-1 flex min-h-0">
                {/* Sidebar */}
                <div className="w-[40%] max-[900px]:hidden h-full flex">
                    <Sidebar />
                </div>
                {/* Main content */}
                <div className="flex-1 min-h-0 flex flex-col">{children}</div>
            </div>
        </div>
    );
}
