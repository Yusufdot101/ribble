"use client";
import Icon from "@/components/Icon";
import google from "@/assets/google.svg";
import { useRouter } from "next/navigation";
import { BASE_USER_SERVICE_API_URL } from "@/utils/api";
import { useAuthStore } from "@/store/useAuthStore";
import { useEffect } from "react";

const googleInfo = {
    src: google,
    href: `${BASE_USER_SERVICE_API_URL}/auth/google`,
    alt: "continue with Google",
};

const Login = () => {
    const router = useRouter();
    const isLoggedIn = useAuthStore((state) => state.isLoggedIn);
    useEffect(() => {
        if (isLoggedIn) {
            router.push("/");
        }
    }, [isLoggedIn, router]);

    return (
        <div className="flex flex-col gap-y-[4px] h-full">
            <p className="text-center w-full max-[619px]:text-[16px]  min-[620px]:text-[24px]"></p>
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
        </div>
    );
};

export default Login;
