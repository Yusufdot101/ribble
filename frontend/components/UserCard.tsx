import { UserType } from "@/utils/users";

interface Props {
    activeUsers: number[];
    user: UserType;
    handleClick: (userID: number) => void;
}

const UserCard = ({ activeUsers, user, handleClick }: Props) => {
    return (
        <div
            tabIndex={0}
            onClick={() => handleClick(user.id)}
            className={`${activeUsers?.includes(user.id) ? "bg-foreground/20" : ""} border-foreground p-[4px] cursor-pointer duration-300 h-[64px]`}
            onKeyDown={(e) => {
                if (e.key === "Enter") {
                    handleClick(user.id);
                }
            }}
        >
            <p className="min-[620px]:text-[20px]">{user.name}</p>
            <p className="min-[620px]:text-[16px]">{user.email}</p>
        </div>
    );
};

export default UserCard;
