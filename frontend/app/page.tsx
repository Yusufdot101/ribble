"use client";
import { useAuthStore } from "@/store/useAuthStore";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

const Home = () => {
    const router = useRouter();
    const isLoggedIn = useAuthStore((state) => state.isLoggedIn);
    useEffect(() => {
        if (isLoggedIn) {
            router.push("/chats");
        }
    }, [isLoggedIn, router]);
    return <div className="flex justify-center">Home</div>;
};

export default Home;
