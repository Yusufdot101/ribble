"use client";
import { useEffect, useState } from "react";

export const useOnlineStatus = () => {
    const [online, setOnline] = useState(() =>
        typeof navigator !== "undefined" ? navigator.onLine : true,
    );

    const on = () => setOnline(true);
    const off = () => setOnline(false);
    useEffect(() => {
        (() => setOnline(navigator.onLine))();

        window.addEventListener("online", on);
        window.addEventListener("offline", off);

        return () => {
            window.removeEventListener("online", on);
            window.removeEventListener("offline", off);
        };
    }, []);

    return online;
};
