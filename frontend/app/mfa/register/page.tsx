"use client"
import { useKeyPair } from "@/app/context/keyContext";
import { useSearchParams } from "next/navigation";
import { useEffect, useRef, useState } from "react";

const MFARegister = () => {
    const ref = useRef<HTMLVideoElement | null>(null)
    const sessionToken = useSearchParams().get("session_token")
    const [loading, setLoading] = useState(true)
    const { generateKeyPair } = useKeyPair()
    const [userName, setUserName] = useState("")
    let i = 0
    useEffect(() => {
        if (i === 1) return;
        const userid = localStorage.getItem("user_id")
        if (userid === null) return;
        getUser(userid)
        if (ref.current === null) return;
        navigator.mediaDevices.getUserMedia({ video: true, audio: false }).then((stream) => {
            if (ref.current === null) return;
            ref.current.srcObject = stream;
            ref.current.play();
        })
        i++
    }, [])

    const getUser = async (userid: string) => {
        const res = await fetch(`http://localhost:8000/api/v1/user`, {
            method: "GET",
            headers: {
                "Authorization": `BEARER ${localStorage.getItem("mfa_token")}`
            }
        })
        const data = await res.json();
        if (res.status.toString().startsWith("2")) {
            setUserName(data.name);
        }
        setLoading(false)
    }

    const handleSubmit = () => {
        // get video frame at the moment and send it to the server
        if (ref.current === null) return;
        ref.current.pause()
        const canvas = document.createElement("canvas");
        canvas.width = ref.current?.videoWidth || 0;
        canvas.height = ref.current?.videoHeight || 0;
        const ctx = canvas.getContext("2d");
        if (ctx === null) return;
        ctx.drawImage(ref.current, 0, 0, canvas.width, canvas.height);
        let imageBlob: Blob | null = null;
        canvas.toBlob(async (blob) => {
            console.log("blob created");
            if (blob !== null) {
                const formData = new FormData();
                const userId = localStorage.getItem("user_id")
                console.log(userId);
                const keyPair = await generateKeyPair()
                formData.append("user_id", userId ?? '1');
                formData.append("face_image", blob);
                formData.append("public_key", keyPair.publicKey ?? "");
                const res = await fetch("http://localhost:8000/api/v1/mfa/face/register/image", {
                    method: "POST",
                    headers: {
                        "Authorization": `BEARER ${sessionToken}`
                    },
                    body: formData,
                })

                if (res.status.toString().startsWith("2")) {
                    console.log("success")
                }
            }
        }, "image/jpeg");
    }

    return (
        <main className="w-screen h-screen flex justify-center items-center">
            <div className="absolute top-[50%] left-[50%] -translate-x-[50%] -translate-y-[50%] h-[400px] aspect-video flex flex-col justify-center items-center">
                {/* <h1>Account Id: {localStorage.getItem("user_id")}</h1> */}
                <h1>Registering: {loading ? "loading" : userName}</h1>
                <video src="" ref={ref} className="w-full"></video>
                <button className="mt-2 py-2 px-4 bg-blue-500 text-white rounded-md" onClick={handleSubmit}>Submit</button>
            </div>
        </main >
    );
};

export default MFARegister;