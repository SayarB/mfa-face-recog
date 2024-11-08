"use client"
const Register = () => {
    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        console.log(e.target);
        const formData = new FormData(e.target as HTMLFormElement);
        const formDataJson = Object.fromEntries(formData) as { name: string, email: string, password: string, confirmPassword: string };

        if (formDataJson.password !== formDataJson.confirmPassword) {
            alert("Passwords do not match");
            return;
        }
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(formDataJson.email)) {
            alert("Invalid email");
            return;
        }
        const requestBody = {
            name: formDataJson.name,
            email: formDataJson.email,
            password: formDataJson.password,
        };
        const res = await fetch("http://localhost:8000/api/v1/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Accept": "application/json",
            },

            body: JSON.stringify(requestBody),
        })
        console.log(res);
        if (res.status.toString().startsWith("2")) {
            const data = await res.json();
            alert("Registration successful user id = " + data.id);
            localStorage.setItem("user_id", JSON.stringify(data.id));
            window.location.href = "/auth/login";
        } else {
            alert("Registration failed");
        }
    };
    return (
        <main className="">
            <div className="flex flex-col items-center justify-center h-screen">
                <h1>Register</h1>
                <form onSubmit={handleSubmit}>
                    <div className="flex flex-col w-[300px]">
                        <input type="text" name="name" placeholder="Name" className="border-2 my-2" />
                        <input type="text" name="email" placeholder="Email" className="border-2 my-2" />
                        <input type="password" name="password" placeholder="Password" className="border-2 my-2" />
                        <input type="password" name="confirmPassword" placeholder="Confirm Password" className="border-2 my-2" />
                        <button className="mt-2 py-2 w-full bg-blue-500 text-white rounded-md" type="submit">Submit</button>
                    </div>
                </form>
            </div>
        </main >
    );
};

export default Register;