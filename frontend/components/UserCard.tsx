import { UserType } from "@/utils/users";

interface Props {
    activeUserID: number;
    user: UserType;
    handleClick: (userID: number) => void;
}

const UserCard = ({ activeUserID, user, handleClick }: Props) => {
    return (
        <div
            onClick={() => handleClick(user.id)}
            className={`${activeUserID === user.id ? "bg-foreground/20" : ""} border-foreground p-[4px] cursor-pointer duration-300 h-[64px]`}
        >
            <p className="min-[620px]:text-[20px]">{user.name}</p>
            <p className="min-[620px]:text-[16px]">{user.email}</p>
        </div>
    );
};

export default UserCard;
