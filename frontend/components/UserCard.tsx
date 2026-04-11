import { UserType } from "@/utils/users";

interface Props {
    index: number;
    activeUserID: number;
    user: UserType;
    handleClick: (userID: number) => void;
}

const UserCard = ({ index, activeUserID, user, handleClick }: Props) => {
    console.log(index, user);
    return (
        <div
            onClick={() => handleClick(user.id)}
            className={`${index === 0 ? "" : "border-t-[1px]"} ${activeUserID === user.id ? "bg-foreground/20" : ""} border-foreground p-[4px] cursor-pointer duration-300 h-[64px]`}
        >
            <p className="min-[620px]:text-[20px]">{user.name}</p>
            <p className="min-[620px]:text-[16px]">{user.email}</p>
        </div>
    );
};

export default UserCard;
