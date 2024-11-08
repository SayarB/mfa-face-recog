import Image from "next/image";

export default function Home() {
  return (
    <main className="">
      <div className="flex flex-col items-center justify-center h-screen">
        <button className="mt-2 py-2 w-[100px] bg-blue-500 text-white rounded-md">
          <a href="/auth/login">Login</a>
        </button>
        <button className="mt-2 py-2 w-[100px] bg-blue-500 text-white rounded-md">
          <a href="/auth/register">Register</a>
        </button>
      </div>
    </main>
  );
}
