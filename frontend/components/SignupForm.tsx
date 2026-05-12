"use client";
import Link from "next/link";
import { useEffect, useState } from "react";
import google from "@/assets/google.svg";
import { signup } from "@/utils/signup";
import ShowHide from "./ShowHide";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/useAuthStore";
import { BASE_USER_SERVICE_API_URL } from "@/utils/api";
import Icon from "./Icon";

const googleInfo = {
    src: google,
    href: `${BASE_USER_SERVICE_API_URL}/auth/google`,
    alt: "continue with Google",
};

const SignupForm = () => {
    const router = useRouter();
    const isLoggedIn = useAuthStore((state) => state.isLoggedIn);
    useEffect(() => {
        if (isLoggedIn) {
            router.push("/");
        }
    }, [isLoggedIn, router]);

    const [name, setName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");

    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);

    const [error, setError] = useState("");

    const handleSubmit = async () => {
        if (password !== confirmPassword) {
            return;
        }
        if (
            !name ||
            !email ||
            !password ||
            !confirmPassword ||
            !(password.length >= 8 && password.length <= 72) ||
            confirmPassword !== password
        ) {
            return;
        }
        const success = await signup(
            (error: string) => setError(error),
            name,
            email,
            password,
        );

        if (!success) return;
    };

    return (
        <div className="flex flex-col gap-y-[8px] justify-center border-1 border-muted-foreground p-[20px] rounded-[8px] w-full max-w-[800px] min-[901]:text-[20px] ">
            <form
                onSubmit={(e) => {
                    e.preventDefault();
                    handleSubmit();
                }}
                className="flex flex-col gap-y-[8px] w-full max-w-[800px] min-[901]:text-[20px] "
            >
                <div className="flex flex-col w-full">
                    <label htmlFor="name">Name</label>
                    <input
                        id="name"
                        name="name"
                        type="text"
                        minLength={2}
                        value={name}
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                            setName(e.target.value)
                        }
                        required
                        min={2}
                        className="border-muted-foreground border-1 rounded-[4px] p-[8px] outline-none"
                        placeholder="Choose a name"
                    />
                </div>

                <div className="flex flex-col w-full">
                    <label htmlFor="email">Email</label>
                    <input
                        id="email"
                        name="email"
                        type="email"
                        value={email}
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                            setEmail(e.target.value)
                        }
                        required
                        className="border-muted-foreground border-1 rounded-[4px] p-[8px] outline-none"
                        placeholder="your.email@example.com"
                    />
                </div>

                <div className="flex flex-col w-full">
                    <label htmlFor="password">Password</label>
                    <div className="relative">
                        <input
                            id="password"
                            name="password"
                            type={showPassword ? "text" : "password"}
                            value={password}
                            onChange={(
                                e: React.ChangeEvent<HTMLInputElement>,
                            ) => setPassword(e.target.value)}
                            required
                            minLength={8}
                            maxLength={72}
                            className="border-muted-foreground border-1 rounded-[4px] p-[8px] outline-none w-full pr-[40px]"
                            placeholder="Create a strong password"
                        />
                        <ShowHide
                            show={showPassword}
                            handleClick={() => setShowPassword((prev) => !prev)}
                        />
                    </div>
                </div>

                <div className="flex flex-col w-full">
                    <label htmlFor="confirmPassword">Confirm Password</label>
                    <div className="relative">
                        <input
                            id="confirmPassword"
                            name="confirmPassword"
                            type={showConfirmPassword ? "text" : "password"}
                            value={confirmPassword}
                            onChange={(
                                e: React.ChangeEvent<HTMLInputElement>,
                            ) => setConfirmPassword(e.target.value)}
                            required
                            minLength={8}
                            maxLength={72}
                            pattern={password}
                            title="passwords must match"
                            className="border-muted-foreground border-1 rounded-[4px] p-[8px] outline-none w-full pr-[40px]"
                            placeholder="Re-enter your password"
                        />
                        <ShowHide
                            show={showConfirmPassword}
                            handleClick={() =>
                                setShowConfirmPassword((prev) => !prev)
                            }
                        />
                    </div>
                </div>

                {error && (
                    <div className="bg-red-500 text-white py-[4px] text-center rounded-[4px]">
                        {error}
                    </div>
                )}

                <button className="bg-foreground text-background rounded-[4px] py-[4px] cursor-pointer hover:bg-muted hover:text-foreground hover:bg-accent duration-300 border-1 border-muted-foreground">
                    Create account
                </button>
            </form>

            <section className="flex items-center w-full">
                <hr className="w-full text-muted-foreground" />
                OR
                <hr className="w-full text-muted-foreground" />
            </section>

            <div className="flex flex-col gap-[24px]">
                <div
                    className="flex flex-wrap h-fit flex items-center justify-center border-gray-500 border rounded-[4px] hover:cursor-pointer hover:bg-white/10 active:bg-black duration-300"
                    onClick={() => {
                        router.push(googleInfo.href);
                    }}
                >
                    <span>Continue With</span>
                    <Icon
                        src={googleInfo.src}
                        href={googleInfo.href}
                        alt={googleInfo.alt}
                        width="50px"
                    />
                </div>
            </div>

            <section className="flex gap-x-[4px] w-full justify-center items-center">
                <p className="text-muted-foreground">
                    Already have an account?{" "}
                </p>
                <Link href={"/login"} className="font-bold">
                    Log In
                </Link>
            </section>
        </div>
    );
};

export default SignupForm;
