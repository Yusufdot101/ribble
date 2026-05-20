import { UserType } from "@/utils/users";

interface Props {
    activeUsers: number[];
    user: UserType;
    handleClick: (user: UserType) => void;
    handleRightClick?: (user: UserType) => void;
}

const UserCard = ({
    activeUsers,
    user,
    handleClick,
    handleRightClick,
}: Props) => {
    return (
        <div
            tabIndex={0}
            onClick={() => handleClick(user)}
            onContextMenu={(e) => {
                if (!handleRightClick) return;
                e.preventDefault();
                handleRightClick(user);
            }}
            className={`${activeUsers?.includes(user.id) ? "bg-foreground/20" : "hover:bg-foreground/10 focus:bg-foreground/10"} border-foreground p-[4px] cursor-pointer duration-300 h-[64px] flex justify-between items-center`}
            onKeyDown={(e) => {
                if (e.key === "Enter") {
                    handleClick(user);
                }
            }}
        >
            <div>
                <p className="min-[620px]:text-[20px]">{user.name}</p>
                <p className="min-[620px]:text-[16px]">{user.email}</p>
            </div>
            <p className="min-[620px]:text-[16px]">{user.role}</p>
        </div>
    );
};

export default UserCard;
