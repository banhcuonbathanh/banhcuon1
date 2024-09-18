


import DashBoard2Header from "./component/dashboard2_header";
import PageWrapper from "./component/pagewrapper";
import { Dashboard2SideBar } from "./component/sidebar";

export default function UserLayout({ children }: { children: React.ReactNode }) {
    return (
        <>
            <Dashboard2SideBar />
            <div className="flex flex-col h-full w-full">
                <DashBoard2Header />
                <PageWrapper children={children} />
            </div>
        </>
    )
}