"use client"
import { useKeyPair } from '@/app/context/keyContext';
import { useCallback, useEffect, useRef, useState } from 'react';
import QRCode from 'react-qr-code';

const useLongPolling = (url: string, options: any, enabled: boolean) => {
    const [data, setData] = useState<{ isComplete: boolean, isSuccess: boolean, isFailure: boolean } | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<any>(null);
    const timeoutRef = useRef<NodeJS.Timeout | null>(null);
    const fetchData = useCallback(async () => {
        try {
            setIsLoading(true);
            const response = await fetch(url, options);
            const json = await response.json();
            setData(json);
        } catch (error) {
            setError(error);
        } finally {
            setIsLoading(false);
        }
    }, [url, options]);

    let i = 0;
    useEffect(() => {
        if (i === 1) return
        if (!enabled) return
        if (!!timeoutRef.current) return
        timeoutRef.current = setInterval(fetchData, 1000);
        i++
        return () => {
            if (timeoutRef.current) {
                clearInterval(timeoutRef.current);
            }
        };
    }, [enabled]);

    return { data, isLoading, error };
};


const FaceVerify = () => {
    const [sessionToken, setSessionToken] = useState<string>("");
    const [sessionId, setSessionId] = useState<string>("")
    const { data, isLoading, error } = useLongPolling(`http://localhost:8000/api/v1/mfa/session/${sessionId}/status`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        }
    }, sessionId !== "")
    const getSessionToken = async (mfa_token: string) => {
        setSessionToken("")
        const res = await fetch(`http://localhost:8000/api/v1/mfa/sessiontoken`, {
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
                <h1>Face Verify</h1>
                {sessionToken !== "" ? <QRCode value={`https://68e7-182-66-218-123.ngrok-free.app/mfa/verify?session_token=${encodeURIComponent(sessionToken)}`} /> : <p>Loading...</p>}

                <a href={`http://localhost:3000/mfa/verify?session_token=${encodeURIComponent(sessionToken)}`} target="_blank" rel="noopener noreferrer" className='mt-2 py-2 px-4 bg-blue-500 text-white rounded-md'>Visit</a></> :
                <>
                    {data.isSuccess ?
                        <h1>Face Verified</h1>
                        :
                        <h1>Face Not Verified</h1>
                    }
                </>
            }
        </div>
    )
}
export default FaceVerify;