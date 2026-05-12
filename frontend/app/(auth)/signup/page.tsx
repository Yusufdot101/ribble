"use client";
import SignupForm from "@/components/SignupForm";

const Login = () => {
    return (
        <div className="flex flex-col gap-y-[4px] h-full w-full items-center">
            <p className="text-center w-full max-[619px]:text-[16px]  min-[620px]:text-[24px]"></p>
            <SignupForm />
        </div>
    );
};

export default Login;
