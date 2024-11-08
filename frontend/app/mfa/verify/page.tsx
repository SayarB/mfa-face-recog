"use client"
import { KeyPair, useKeyPair } from "@/app/context/keyContext";
import { useSearchParams } from "next/navigation";
import { useCallback, useEffect, useRef, useState } from "react";

const decryptAndEncrypt = (message: string, keyPair: KeyPair) => {

}


const MFAVerify = () => {
    const sessionToken = useSearchParams().get("session_token")
    const ref = useRef<HTMLVideoElement | null>(null)
    const timeoutRef = useRef<NodeJS.Timeout | null>(null)
    const [userLoading, setUserLoading] = useState(true)
    const [userName, setUserName] = useState("")
    const [verifiedPreds, setVerifiedPreds] = useState<number>(0)
    const [notVerifiedPreds, setNotVerifiedPreds] = useState<number>(0)
    const [decSessionToken, setDecSessionToken] = useState("")
    const [loading, setLoading] = useState(true)
    const { decryptMessage, loadKeyPair, privateKey } = useKeyPair()
    let i = 0
    useEffect(() => {
        if (i === 1) return;
        // const token = localStorage.getItem("access_token")
        // if (token === null) return;
        // getUser(token)
        loadKeyPair()
        i++
    }, [])
    i = 0
    useEffect(() => {
        if (i === 1) return
        if (decSessionToken === "") return
        if (ref.current === null) return;
        navigator.mediaDevices.getUserMedia({ video: true, audio: false }).then(async (stream) => {
            if (ref.current === null) return;
            ref.current.srcObject = stream;
            if (ref.current.paused) await ref.current.play();
            if (!timeoutRef.current) sendFrameToServer()
        })
        i++
    }, [decSessionToken])
    useEffect(() => {
        console.log("change in session token or private key")
        if (sessionToken && privateKey) decryptSessionToken()
    }, [sessionToken, privateKey])

    const decryptSessionToken = useCallback(async () => {
        if (!sessionToken) return

        if (!privateKey) {
            throw new Error("Private key is not available for decryption.");
        }
        console.log("session", sessionToken)
        // Decode the Base64 string to ArrayBuffer
        const ciphertext = Uint8Array.from(atob(sessionToken), c => c.charCodeAt(0)).buffer;

        const decrypted = await window.crypto.subtle.decrypt(
            {
                name: "RSA-OAEP",
            },
            privateKey,
            ciphertext
        );

        const decoder = new TextDecoder();
        setDecSessionToken(decoder.decode(decrypted));
    }, [sessionToken, privateKey])

    i = 0
    useEffect(() => {
        if (i === 1) return
        console.log(verifiedPreds, notVerifiedPreds)
        if (verifiedPreds >= 3 || notVerifiedPreds >= 5) {
            setLoading(false)
            window.location.href = verifiedPreds >= 3 ? "/verified" : "/not_verified";
        } else if (verifiedPreds > 0 || notVerifiedPreds > 0) {
            sendFrameToServer()
        }
        i++

    }, [verifiedPreds, notVerifiedPreds])

    // const getUser = async (token: string) => {
    //     setUserLoading(true)
    //     const res = await fetch(`http://localhost:8000/api/v1/user`, {
    //         method: "GET",
    //         headers: {
    //             'Authorization': "BEARER " + token
    //         },
    //     })
    //     const data = await res.json();
    //     if (res.status.toString().startsWith("2")) {
    //         setUserName(data.name);
    //     }
    //     setUserLoading(false)
    // }

    const handlePreds = useCallback(async (res: Response) => {
        const data = await res.json();
        if (res.status.toString().startsWith("2")) {
            if (data.verified === "true") {
                setVerifiedPreds(pred => pred + 1)
                setNotVerifiedPreds(0)
            } else {
                setVerifiedPreds(0)
                setNotVerifiedPreds(pred => pred + 1)
            }
            console.log(data)

        } else if (res.status === 400) {
            console.log("expired")
            if (data.success === "false") {
                window.location.href = "/expired";
            }
        }
    }, [verifiedPreds, notVerifiedPreds, decSessionToken])

    const sendFrameToServer = useCallback(() => {
        // get video frame at the moment and send it to the server
        if (decSessionToken === "") {
            console.log("Decrypted Session token is not available.")
            return
        }
        if (ref.current === null) return;
        const canvas = document.createElement("canvas");
        canvas.width = ref.current?.videoWidth || 0;
        canvas.height = ref.current?.videoHeight || 0;
        const ctx = canvas.getContext("2d");
        if (ctx === null) return;
        ctx.drawImage(ref.current, 0, 0, canvas.width, canvas.height);
        canvas.toBlob(async (blob) => {
            console.log("blob created");
            if (blob !== null) {
                const formData = new FormData();
                const userId = localStorage.getItem("user_id")
                console.log(userId);
                formData.append("user_id", userId ?? '1');
                formData.append("face_image", blob);

                const res = await fetch("http://localhost:8000/api/v1/mfa/face/verify", {
                    method: "POST",
                    headers: {
                        'Authorization': "BEARER " + decSessionToken
                    },
                    body: formData,
                })
                handlePreds(res)

            }
        }, "image/jpeg");

    }, [decSessionToken])



    return (
        <main className="">
            <div className="absolute top-[50%] left-[50%] -translate-x-[50%] -translate-y-[50%] w-[800px] aspect-video flex flex-col justify-center items-center">
                {/* <h1>Account Id: {localStorage.getItem("user_id")}</h1> */}
                <h1>Loading: {loading ? "loading" : "not loading"}</h1>
                <video src="" ref={ref} className="w-full"></video>
                <p>Positive Predictions: {verifiedPreds}</p>
                <p>Negative Predictions: {notVerifiedPreds}</p>
            </div>
        </main >
    );
};

export default MFAVerify;