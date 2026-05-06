"use client";

import { useEffect, useState } from "react";

export const useOnlineStatus = () => {
    const [online, setOnline] = useState(true);

    const update = async () => {
        const nav = navigator.onLine;

        if (!nav) return setOnline(false);

        try {
            await fetch("https://www.google.com/generate_204", {
                cache: "no-store",
                mode: "no-cors",
            });
            setOnline(true);
        } catch {
            setOnline(false);
        }
    };

    useEffect(() => {
        (() => update())();

        window.addEventListener("online", update);
        window.addEventListener("offline", update);

        const interval = setInterval(update, 10000);

        return () => {
            window.removeEventListener("online", update);
            window.removeEventListener("offline", update);
            clearInterval(interval);
        };
    }, []);

    return online;
};
