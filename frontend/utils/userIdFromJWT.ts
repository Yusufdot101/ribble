function base64UrlDecode(str: string) {
    try {
        const padded = str.padEnd(
            str.length + ((4 - (str.length % 4)) % 4),
            "=",
        );
        const base64 = padded.replace(/-/g, "+").replace(/_/g, "/");
        return atob(base64);
    } catch (error) {
        console.error("failed to decode base64: ", error);
        return "";
    }
}

export const decodeJWT = (token: string) => {
    if (!token)
        return {
            header: {},
            payload: {},
            signature: {},
        };

    const [h, p, s] = token.split(".");
    if (!h || !p || !s) {
        console.error("invalid JWT structure");
        return {
            header: {},
            payload: {},
            signature: {},
        };
    }

    try {
        const header = JSON.parse(base64UrlDecode(h));
        const payload = JSON.parse(base64UrlDecode(p));
        return { header, payload, signature: s };
    } catch (error) {
        console.error("failed to parse JWT:", error);
        return {
            header: {},
            payload: {},
            signature: "",
        };
    }
};
