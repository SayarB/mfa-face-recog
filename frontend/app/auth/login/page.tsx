"use client"

import { useKeyPair } from "@/app/context/keyContext";
import { useState } from "react";

const LoginPage = () => {
    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        console.log(e.target);
        const formData = new FormData(e.target as HTMLFormElement);
        const formDataJson = Object.fromEntries(formData) as { email: string, password: string, public_key: string | null }
        const res = await fetch("http://localhost:8000/api/v1/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(formDataJson),
        })
        if (res.status.toString().startsWith("2")) {
            const data = await res.json();
            console.log(data);
            if (data.success === "true") {
                localStorage.setItem("mfa_token", data.token)
                if (data.mfa_enabled) {
                    window.location.href = "/mfa/verify";
                } else {
                    window.location.href = "/mfa/register/face";
                }
            }
        }
    }
    return <main className="">
        <div className="flex flex-col items-center justify-center h-screen">
            <h1>Login</h1>
            <form onSubmit={handleSubmit}>
                <div className="flex flex-col w-[300px]">
                    <input type="text" name="email" placeholder="Email" className="border-2 my-2" />
                    <input type="password" name="password" placeholder="Password" className="border-2 my-2" />
                    <button className="mt-2 py-2 w-full bg-blue-500 text-white rounded-md" type="submit">Submit</button>
                </div>
            </form>
        </div>
    </main >
}

export default LoginPage