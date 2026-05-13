"use client";
import Link from "next/link";
import { useEffect, useState } from "react";
import google from "@/assets/google.svg";
import ShowHide from "./ShowHide";
import { login } from "@/utils/login";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/useAuthStore";
import { BASE_USER_SERVICE_API_URL } from "@/utils/api";
import Icon from "./Icon";

const googleInfo = {
    src: google,
    href: `${BASE_USER_SERVICE_API_URL}/auth/google`,
    alt: "continue with Google",
};

const LoginForm = () => {
    const router = useRouter();
    const isLoggedIn = useAuthStore((state) => state.isLoggedIn);
    useEffect(() => {
        if (isLoggedIn) {
            router.push("/");
        }
    }, [isLoggedIn, router]);

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    const [showPassword, setShowPassword] = useState(false);

    const [error, setError] = useState("");

    const handleSubmit = async () => {
        if (
            !email ||
            !password ||
            !(password.length >= 8 && password.length <= 72)
        ) {
            return;
        }
        const success = await login(
            (error: string) => setError(error),
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
                    <label htmlFor="email">Email</label>
                    <input
                        id="email"
                        name="email"
                        type="email"
                        required
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        className="border-muted-foreground border-1 rounded-[4px] p-[8px] outline-none"
                        placeholder="your.email@example.com"
                    />
                </div>

                <div className="flex flex-col w-full">
                    <div className="flex justify-between">
                        <label htmlFor="password">Password</label>
                        <Link
                            href={"/forgotpassword"}
                            className="font-light text-muted-foreground cursor-pointer hover:text-foreground duration-300"
                        >
                            Forgot Password?
                        </Link>
                    </div>
                    <div className="relative">
                        <input
                            id="password"
                            name="password"
                            type={showPassword ? "text" : "password"}
                            required
                            minLength={8}
                            maxLength={72}
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            className="border-muted-foreground border-1 rounded-[4px] p-[8px] outline-none w-full pr-[40px]"
                            placeholder="Enter your password"
                        />
                        <ShowHide
                            show={showPassword}
                            handleClick={() => {
                                setShowPassword((prev) => !prev);
                            }}
                        />
                    </div>
                </div>

                {error && (
                    <div className="bg-red-500 text-white py-[4px] text-center rounded-[4px]">
                        {error}
                    </div>
                )}

                <button className="bg-foreground text-background rounded-[4px] py-[4px] cursor-pointer hover:bg-muted hover:text-foreground hover:bg-accent duration-300 border-1 border-muted-foreground">
                    Log In
                </button>
            </form>

            <section className="flex items-center w-full">
                <hr className="w-full text-muted-foreground" />
                OR
                <hr className="w-full text-muted-foreground" />
            </section>

            <div className="flex flex-col gap-[24px]">
                <div
                    className="flex flex-wrap h-fit items-center justify-center border-gray-500 border rounded-[4px] hover:cursor-pointer hover:bg-white/10 active:bg-black duration-300"
                    onClick={() => {
                        window.location.href = googleInfo.href;
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
                    Don&apos;t have an account?{" "}
                </p>
                <Link href={"/signup"} className="font-bold">
                    Sign up
                </Link>
            </section>
        </div>
    );
};

export default LoginForm;
