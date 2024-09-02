"use client"
import { useEffect, useRef } from "react";

const MFARegister = () => {
    const ref = useRef<HTMLVideoElement | null>(null)
    useEffect(() => {
        if (ref.current === null) return;
        navigator.mediaDevices.getUserMedia({ video: true, audio: false }).then((stream) => {
            if (ref.current === null) return;
            ref.current.srcObject = stream;
            ref.current.play();
        })
    })

    const handleSubmit = () => {
        // get video frame at the moment and send it to the server
        if (ref.current === null) return;
        const canvas = document.createElement("canvas");
        canvas.width = ref.current?.videoWidth || 0;
        canvas.height = ref.current?.videoHeight || 0;
        const ctx = canvas.getContext("2d");
        if (ctx === null) return;
        ctx.drawImage(ref.current, 0, 0, canvas.width, canvas.height);
        let imageBlob: Blob | null = null;
        canvas.toBlob((blob) => {
            console.log("blob created");
            if (blob !== null) {
                const formData = new FormData();
                const userId = localStorage.getItem("user_id")
                console.log(userId);
                formData.append("user_id", userId ?? '1');
                formData.append("face_image", blob);
                fetch("http://localhost:8000/api/v1/mfa/face/register", {
                    method: "POST",
                    body: formData,
                })
            }
        }, "image/jpeg");
    }

    return (
        <main className="">
            <div className="absolute top-0 left-0 w-[400px] aspect-video bg-black/50 flex flex-col justify-center items-center">
                <video src="" ref={ref} className="w-full"></video>
                <button onClick={handleSubmit}>Submit</button>
            </div>
        </main >
    );
};

export default MFARegister;