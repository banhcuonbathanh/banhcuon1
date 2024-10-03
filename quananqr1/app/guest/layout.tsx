import DarkModeToggle from "@/components/dark-mode-toggle";
import NavLinks from "../admin/admin_component/nav-links";
import DropdownAvatar from "../admin/admin_component/dropdown-avatar";
import LoginDialog from "../(public)/public-component/login-dialog";
import RegisterDialog from "../(public)/public-component/register-dialog";
// import NavLinks from './admin_component/nav-links'
// import DropdownAvatar from './admin_component/dropdown-avatar'

export default function Layout({
  children
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="flex min-h-screen w-full flex-col bg-muted/40">
      <div className="flex flex-col sm:gap-4 sm:py-4 sm:pl-14">
        <header className="sticky top-0 z-30 flex h-14 items-center gap-4 border-b bg-background px-4 sm:static sm:h-auto sm:border-0 sm:bg-transparent sm:px-6">
          {/* <MobileNavLinks /> */}
          <div className="relative ml-auto flex-1 md:grow-0">
            <div className="flex justify-end">
              <DarkModeToggle />
            </div>
          </div>
          <DropdownAvatar />
          <LoginDialog />
          <RegisterDialog />
        </header>
        {children}
      </div>
    </div>
  );
}
