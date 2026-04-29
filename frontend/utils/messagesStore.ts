import { WebsocketMsg } from "@/components/MessageInput";

const DB_NAME = "chatDB";
const STORE_NAME = "messages";
const DB_VERSION = 1;

let dbPromise: Promise<IDBDatabase> | null = null;

function openDB(): Promise<IDBDatabase> {
    if (dbPromise) return dbPromise;

    dbPromise = new Promise((resolve, reject) => {
        const request = indexedDB.open(DB_NAME, DB_VERSION);

        request.onupgradeneeded = () => {
            const db = request.result;

            if (!db.objectStoreNames.contains(STORE_NAME)) {
                const store = db.createObjectStore(STORE_NAME, {
                    keyPath: "clientID",
                });

                store.createIndex("chatId", "chatID", { unique: false });
            }
        };

        request.onsuccess = () => {
            resolve(request.result);
        };

        request.onerror = () => {
            reject(request.error);
        };
    });

    return dbPromise;
}

async function add(msg: WebsocketMsg): Promise<void> {
    const db = await openDB();

    return new Promise((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, "readwrite");
        const store = tx.objectStore(STORE_NAME);

        const req = store.put(msg);

        req.onsuccess = () => resolve();
        req.onerror = () => reject(req.error);
    });
}

async function getAll(): Promise<WebsocketMsg[]> {
    const db = await openDB();

    return new Promise((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, "readonly");
        const store = tx.objectStore(STORE_NAME);

        const req = store.getAll();

        req.onsuccess = () => resolve(req.result);
        req.onerror = () => reject(req.error);
    });
}

async function getByChat(chatID: number): Promise<WebsocketMsg[]> {
    const db = await openDB();

    return new Promise((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, "readonly");
        const store = tx.objectStore(STORE_NAME);
        const index = store.index("chatId");

        const req = index.getAll(chatID);

        req.onsuccess = () => resolve(req.result);
        req.onerror = () => reject(req.error);
    });
}

async function update(msg: WebsocketMsg): Promise<void> {
    return add(msg);
}

async function remove(clientID: string): Promise<void> {
    const db = await openDB();

    return new Promise((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, "readwrite");
        const store = tx.objectStore(STORE_NAME);

        const req = store.delete(clientID);

        req.onsuccess = () => resolve();
        req.onerror = () => reject(req.error);
    });
}

async function clear(): Promise<void> {
    const db = await openDB();

    return new Promise((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, "readwrite");
        const store = tx.objectStore(STORE_NAME);

        const req = store.clear();

        req.onsuccess = () => resolve();
        req.onerror = () => reject(req.error);
    });
}

export const messageStore = {
    add,
    getAll,
    getByChat,
    update,
    delete: remove,
    clear,
};
