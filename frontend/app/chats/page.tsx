import Contacts from "@/components/Contacts";

const Chats = () => {
    return (
        <div className="h-full">
            <div className="min-[899px]:hidden h-full flex">
                {/* only show on mobile */}
                <Contacts />
            </div>
            <div className="max-[900px]:hidden flex justify-center h-full">
                {/* only show on desktop */}
                Welcome to Ripple
            </div>
        </div>
    );
};

export default Chats;
