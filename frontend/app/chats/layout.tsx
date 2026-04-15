import Contacts from "@/components/Contacts";
import SearchBar from "@/components/SearchBar";

export default function ChatsLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <div className="flex flex-col gap-y-[8px] flex-1 h-full">
            <SearchBar />

            <div className="gap-x-[8px] flex-1 flex">
                {/* Sidebar */}
                {/* Hide in mobile */}
                <div className="w-[40%] max-[900px]:hidden nlex border-r-1 border-red">
                    <Contacts />
                </div>

                {/* Main content */}
                <div className="flex-1">{children}</div>
            </div>
        </div>
    );
}
