"use client";
import { useEffect, useState } from "react";

export const useOnlineStatus = () => {
    const [online, setOnline] = useState(() =>
        typeof navigator !== "undefined" ? navigator.onLine : true,
    );

<<<<<<< HEAD
    useEffect(() => {
        setOnline(navigator.onLine);
        const on = () => setOnline(true);
        const off = () => setOnline(false);
=======
    const on = () => setOnline(true);
    const off = () => setOnline(false);
    useEffect(() => {
        (() => setOnline(navigator.onLine))();
>>>>>>> feature/usernames-to-group-messages

        window.addEventListener("online", on);
        window.addEventListener("offline", off);

        return () => {
            window.removeEventListener("online", on);
            window.removeEventListener("offline", off);
        };
    }, []);

    return online;
};
