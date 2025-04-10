"use client"
import { useLongPolling } from "@/hooks/useLongPolling";
import { useEffect, useState } from "react";
import QRCode from "react-qr-code";

const FaceRegister = () => {
    const [sessionToken, setSessionToken] = useState<string>("");
    const [sessionId, setSessionId] = useState<string>("")
    const { data, isLoading, error } = useLongPolling(`http://localhost:8000/api/v1/mfa/register/session/${sessionId}/status`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        }
    }, sessionId !== "")
    const getSessionToken = async (mfa_token: string) => {
        setSessionToken("")
        const res = await fetch(`http://localhost:8000/api/v1/mfa/register/sessiontoken`, {
            headers: {
                "Content-Type": "application/json",
                "Authorization": `BEARER ${mfa_token}`
            }
        })

        if (res.status === 401) {
            window.location.href = "/auth/login";
            return
        }

        if (res.status.toString().startsWith("2")) {
            const data = await res.json();
            console.log(data.session_token)
            if (data.success) {
                setSessionToken(data.session_token)
                setSessionId(data.session_id)
            }
        }
    }


    let i = 0
    useEffect(() => {
        if (i === 1) return
        const mfa_token = localStorage.getItem("mfa_token");
        if (!mfa_token) {
            window.location.href = "/auth/login";
            return
        }
        getSessionToken(mfa_token)
        i++
    }, [])
    console.log(encodeURIComponent(sessionToken))
    return (
        <div className='flex flex-col items-center justify-center h-screen'>
            {!data?.isComplete ? <>
                <h1 className="text-3xl font-bold">Scan Qr Code</h1>
                {sessionToken !== "" ? <QRCode value={`https://68e7-182-66-218-123.ngrok-free.app/mfa/verify?session_token=${encodeURIComponent(sessionToken)}`} /> : <p>Loading...</p>}

                <a href={`http://localhost:3000/mfa/register?session_token=${encodeURIComponent(sessionToken)}`} target="_blank" rel="noopener noreferrer" className='mt-2 py-2 px-4 bg-blue-500 text-white rounded-md'>Visit</a></> :
                <>
                    <h1>Registration Complete</h1>
                    <a className="mt-2 py-2 px-4 bg-blue-500 text-white rounded-md" href="/mfa/verify/face">Verification</a>
                </>
            }
        </div>
    )
}
export default FaceRegister;