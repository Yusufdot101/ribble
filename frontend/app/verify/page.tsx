"use client";
import { verifyAccount } from "@/utils/users";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect } from "react";

const page = () => {
    const searchParams = useSearchParams();
    const token = searchParams.get("token");
    const identity = searchParams.get("identity");
    if (!token || !identity) return;
    const router = useRouter();
    useEffect(() => {
        verifyAccount(token, identity);
        router.replace("/");
    }, []);
    return <div className="flex justify-center">Activating Account....</div>;
};

export default page;
