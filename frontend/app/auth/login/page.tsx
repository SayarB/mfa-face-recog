"use client"
const LoginPage = () => {
    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        console.log(e.target);
        const formData = new FormData(e.target as HTMLFormElement);
        const formDataJson = Object.fromEntries(formData) as { email: string, password: string };
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
                localStorage.setItem("user_id", data.id);
                window.location.href = "/mfa/verify";
            }
        }
    }
    return <main className="">
        <div className="flex flex-col items-center justify-center">
            <h1>Login</h1>
            <form onSubmit={handleSubmit}>
                <div className="flex flex-col w-[300px]">
                    <input type="text" name="email" placeholder="Email" className="border-2 my-2" />
                    <input type="password" name="password" placeholder="Password" className="border-2 my-2" />
                    <button type="submit">Submit</button>
                </div>
            </form>
        </div>
    </main >
}

export default LoginPage