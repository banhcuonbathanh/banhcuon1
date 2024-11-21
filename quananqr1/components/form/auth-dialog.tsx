"use client";
import React from "react";
import LoginDialog from "./login-dialog";
import GuestLoginDialog from "./guest-dialog";
import RegisterDialog from "./register-dialog";


const AuthDialogs = () => {
  return (
    <>
      <LoginDialog fromPath={null} />
      <GuestLoginDialog fromPath={null} />
      <RegisterDialog /> {/* If you have this component */}
    </>
  );
};

export default AuthDialogs;