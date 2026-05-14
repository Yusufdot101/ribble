"use client";
import { verifyAccount } from "@/utils/users";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { Suspense } from "react";

const VerifyContent = () => {
    const searchParams = useSearchParams();
    const router = useRouter();
    const token = searchParams.get("token");
    const identity = searchParams.get("identity");
    const [status, setStatus] = useState("Activating Account....");

    useEffect(() => {
        if (!token || !identity) {
            (() => setStatus("Invalid activation link"))();
            return;
        }
        let cancelled = false;
        (async () => {
            const ok = await verifyAccount(token, identity);
            if (cancelled) return;
            if (ok) {
                router.replace("/");
            } else {
                (() => setStatus("Activation failed. Please try again."))();
            }
        })();
        return () => {
            cancelled = true;
        };
    }, [token, identity, router]);

    return <div className="flex justify-center">{status}</div>;
};

const Page = () => (
    <Suspense fallback={<div className="flex justify-center">Loading...</div>}>
        <VerifyContent />
    </Suspense>
);

export default Page;
