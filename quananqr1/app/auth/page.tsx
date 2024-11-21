"use client";

import { useSearchParams } from 'next/navigation';
import React from "react";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/zusstand/new_auth/new_auth_controller";
import LoginDialog from "@/components/form/login-dialog";
import GuestLoginDialog from "@/components/form/guest-dialog";
import RegisterDialog from "@/components/form/register-dialog";

const AuthPage = () => {
  const searchParams = useSearchParams();
  const fromPath = searchParams.get('from');

  const {
    openLoginDialog,
    openRegisterDialog,
    openGuestDialog
  } = useAuthStore();

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-100 dark:bg-gray-900 px-4 py-8">
      <div className="w-full max-w-md bg-white dark:bg-gray-800 rounded-xl shadow-lg p-8 space-y-6">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">
            Authentication
          </h1>
          <p className="text-gray-600 dark:text-gray-300 mb-6">
            {fromPath 
              ? `Please log in to access ${fromPath}` 
              : "Choose your preferred login method"}
          </p>
        </div>

        <div className="space-y-4">
          <Button 
            onClick={openLoginDialog} 
            className="w-full"
          >
            Login with Email
          </Button>
          
          <Button 
            onClick={openRegisterDialog} 
            variant="outline" 
            className="w-full"
          >
            Register New Account
          </Button>
          
          <Button 
            onClick={openGuestDialog} 
            variant="secondary" 
            className="w-full"
          >
            Continue as Guest
          </Button>
        </div>

        {/* ... rest of the existing code ... */}
      </div>

      <AuthDialogs fromPath={fromPath} />
    </div>
  );
};

const AuthDialogs = ({ fromPath }: { fromPath: string | null }) => {
  return (
    <>
      <LoginDialog fromPath={fromPath} />
      <GuestLoginDialog fromPath={fromPath} />
      <RegisterDialog />
    </>
  );
};

export default AuthPage;