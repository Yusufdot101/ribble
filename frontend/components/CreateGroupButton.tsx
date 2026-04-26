import Link from "next/link";

type Props = {
    handleClick: () => void;
};

const CreateGroupButton = ({ handleClick }: Props) => {
    return (
        <div className="group relative inline-block ml-auto w-fit">
            <svg
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
                fill="currentColor"
                className="text-foreground hover:text-accent active:text-foreground  duration-300 cursor-pointer border-1 border-foreground rounded-[4px]  w-[44px] h-[44px] max-[899px]:w-[28px] max-[899px]:h-[28px]"
                role="button"
                aria-label="create group chat"
                tabIndex={0}
                onClick={handleClick}
                onKeyDown={(e) => {
                    if (e.key !== "Enter") return;
                    handleClick();
                }}
            >
                <g id="SVGRepo_bgCarrier" strokeWidth="0"></g>
                <g
                    id="SVGRepo_tracerCarrier"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                ></g>
                <g id="SVGRepo_iconCarrier">
                    <g id="Complete">
                        <g>
                            {" "}
                            <g>
                                <line
                                    fill="none"
                                    stroke="currentColor"
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth="2"
                                    x1="12"
                                    x2="12"
                                    y1="19"
                                    y2="5"
                                ></line>
                                <line
                                    fill="none"
                                    stroke="currentColor"
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth="2"
                                    x1="5"
                                    x2="19"
                                    y1="12"
                                    y2="12"
                                ></line>
                            </g>
                            t
                        </g>
                    </g>
                </g>
            </svg>

            <div className="absolute opacity-0 invisible group-hover:opacity-100 group-hover:visible transition duration-300 text-nowrap border-1 border-foreground p-[4px] rounded-[4px] min-[899px]:bottom-full max-[899px]:bottom-1/2 min-[899px]:right-1/2 max-[899px]:right-full min-[899px]:translate-x-1/2 max-[899px]:translate-y-1/2 min-[899px]:mb-[4px] max-[899px]:mr-[4px]">
                create group chat
            </div>
        </div>
    );
};

export default CreateGroupButton;
