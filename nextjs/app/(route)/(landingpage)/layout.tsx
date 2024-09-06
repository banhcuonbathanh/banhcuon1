import Image from "next/image";
import Link from "next/link";

import LandingPageMenu from "./landding_page_component/landding_page_menu";
import LandingPageNavbar from "./landding_page_component/landding_page_navbar";
import Header from "./landding_page_component/header";
import ActiveSectionContextProvider from "@/app/context/active-section-context";
import DashboardMail from "../dashboard/dashboard_component/dashboard_mail/dashboard_mail";

export default function DashboardLayout({
  children
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <ActiveSectionContextProvider>
      <div className="h-screen flex mt-20">
        {/* LEFT */}
        <div className="w-[14%] md:w-[8%] lg:w-[16%] xl:w-[14%] p-4 flex flex-col">
          {/* <Link
          href="/"
          className="flex items-center justify-center lg:justify-start gap-2 mb-4"
        >
          <Image src="/logo.png" alt="logo" width={32} height={32} />
          <span className="hidden lg:block font-bold">SchooLama</span>
        </Link> */}
          {/* <Header /> */}

          <LandingPageMenu />
        </div>
        {/* RIGHT */}
        <div className="w-[86%] md:w-[92%] lg:w-[84%] xl:w-[86%] flex flex-col">
          {children}
        </div>

     
      </div>
    </ActiveSectionContextProvider>
  );
}

{
  /* <div className="mx-auto max-w-7xl"> */
}
