
import LoginDialog from "./login-dialog";
import GuestLoginDialog from "./guest-dialog";
import RegisterDialog from "./register-dialog";


const AuthDialogs = () => {


  return (
    <>
      <LoginDialog  />
      <GuestLoginDialog />
      <RegisterDialog /> {/* If you have this component */}
    </>
  );
};

export default AuthDialogs;