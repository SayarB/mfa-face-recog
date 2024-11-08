"use client"; // Ensure this code only runs on the client

import { createContext, ReactNode, useContext, useState, useEffect, useCallback, ContextType } from "react";

// Define the shape of the key pair state
export interface KeyPair {
    publicKey: string | null;
    privateKey: CryptoKey | null;
}



// Define context type
const KeyContext = createContext<KeyPair & { generateKeyPair: () => Promise<KeyPair>, decryptMessage: (msg: string) => Promise<string>, loadKeyPair: () => Promise<void> } | undefined>(undefined);

// Custom hook to access the KeyPair context
export function useKeyPair(): KeyPair & { generateKeyPair: () => Promise<KeyPair>, decryptMessage: (msg: string) => Promise<string>, loadKeyPair: () => Promise<void> } {
    const context = useContext(KeyContext);
    if (!context) {
        throw new Error("useKeyPair must be used within a KeyProvider");
    }
    return context;
}

// KeyProvider component
interface KeyProviderProps {
    children: ReactNode;
}

export function KeyProvider({ children }: KeyProviderProps) {
    const [keyPair, setKeyPair] = useState<KeyPair>({ publicKey: null, privateKey: null });

    useEffect(() => {
        console.log("private key  = ", keyPair.privateKey)
    }, [keyPair.privateKey])
    const decryptMessage = useCallback(async (cipher: string): Promise<string> => {
        if (!keyPair.privateKey) {
            throw new Error("Private key is not available for decryption.");
        }

        // Decode the Base64 string to ArrayBuffer
        const ciphertext = Uint8Array.from(cipher, c => c.charCodeAt(0)).buffer;

        const decrypted = await window.crypto.subtle.decrypt(
            {
                name: "RSA-OAEP",
            },
            keyPair.privateKey,
            ciphertext
        );

        const decoder = new TextDecoder();
        return decoder.decode(decrypted);
    }, [keyPair.privateKey]);

    const generateKeys = useCallback(async () => {

        const generatedKeyPair = await window.crypto.subtle.generateKey(
            {
                name: "RSA-OAEP",
                modulusLength: 2048,
                publicExponent: new Uint8Array([1, 0, 1]),
                hash: { name: "SHA-256" },
            },
            true, // Keys can be exported
            ["encrypt", "decrypt"]
        );

        const exportedPublicKey = await window.crypto.subtle.exportKey("spki", generatedKeyPair.publicKey);
        const publicKeyPem = convertToPem(exportedPublicKey);

        setKeyPair({ publicKey: publicKeyPem, privateKey: generatedKeyPair.privateKey });
        saveKeyPair({ publicKey: publicKeyPem, privateKey: generatedKeyPair.privateKey });
        console.log({ publicKey: publicKeyPem, privateKey: generatedKeyPair.privateKey })
        return { publicKey: publicKeyPem, privateKey: generatedKeyPair.privateKey }
    }, [setKeyPair])

    const convertToBase64 = (buffer: ArrayBuffer): string => {

        return btoa(String.fromCharCode(...new Uint8Array(buffer).values().toArray()));
    };

    const saveKeyPair = useCallback(async (keys: KeyPair) => {
        if (!keys.privateKey) return;
        if (!keys.publicKey) return;

        // Export private key
        const exportedPrivateKey = await window.crypto.subtle.exportKey("pkcs8", keys.privateKey);
        const base64PrivateKey = convertToBase64(exportedPrivateKey);
        localStorage.setItem("privateKey", base64PrivateKey);

        // Export public key
        const exportedPublicKey = await window.crypto.subtle.exportKey("spki", await window.crypto.subtle.importKey(
            "spki",
            new Uint8Array(Buffer.from(keys.publicKey.split("\n").slice(1, -1).join(""), "base64")),
            {
                name: "RSA-OAEP",
                hash: { name: "SHA-256" },
            },
            true,
            []
        ));
        const base64PublicKey = convertToBase64(exportedPublicKey);
        localStorage.setItem("publicKey", base64PublicKey);
    }, [keyPair]);

    const loadKeyPair = useCallback(async () => {
        if (keyPair.privateKey && keyPair.publicKey) return
        const base64PrivateKey = localStorage.getItem("privateKey");
        const base64PublicKey = localStorage.getItem("publicKey");

        if (base64PrivateKey && base64PublicKey) {
            const privateKeyBuffer = Uint8Array.from(atob(base64PrivateKey), c => c.charCodeAt(0)).buffer;
            const publicKeyBuffer = Uint8Array.from(atob(base64PublicKey), c => c.charCodeAt(0)).buffer;

            const privateKey = await window.crypto.subtle.importKey(
                "pkcs8",
                privateKeyBuffer,
                {
                    name: "RSA-OAEP",
                    hash: { name: "SHA-256" },
                },
                true,
                ["decrypt"]
            );

            const publicKey = await window.crypto.subtle.importKey(
                "spki",
                publicKeyBuffer,
                {
                    name: "RSA-OAEP",
                    hash: { name: "SHA-256" },
                },
                true,
                ["encrypt"]
            );

            const exportedPublicKey = await window.crypto.subtle.exportKey("spki", publicKey);
            const publicKeyPem = convertToPem(exportedPublicKey);

            setKeyPair({ publicKey: publicKeyPem, privateKey });
        }
    }, []);


    const convertToPem = (exportedKey: ArrayBuffer): string => {
        const arr = new Uint8Array(exportedKey).values().toArray()
        const base64String = window.btoa(String.fromCharCode(...arr));
        return `-----BEGIN PUBLIC KEY-----\n${base64String.match(/.{1,64}/g)?.join("\n")}\n-----END PUBLIC KEY-----`;
    };

    return <KeyContext.Provider value={{ ...keyPair, generateKeyPair: generateKeys, decryptMessage, loadKeyPair }}>{children}</KeyContext.Provider>;
}