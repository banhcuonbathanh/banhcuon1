import Image from "next/image";
import Link from "next/link";
import DashBoardMenu from "./dashboard_component/dashboard_menu";
import DashBoardNavbar from "./dashboard_component/dashboard_navbar";

export default function DashboardLayout({
  children
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="h-screen flex">
      {/* LEFT */}
      <div className="w-[14%] md:w-[8%] lg:w-[16%] xl:w-[14%] p-4">
        <Link
          href="/"
          className="flex items-center justify-center lg:justify-start gap-2"
        >
          <Image src="/logo.png" alt="logo" width={32} height={32} />
          <span className="hidden lg:block font-bold">SchooLama</span>
        </Link>
        <DashBoardMenu />
      </div>
      {/* RIGHT */}
      <div className="w-[86%] md:w-[92%] lg:w-[84%] xl:w-[86%] bg-[#F7F8FA] overflow-scroll flex flex-col">
        <DashBoardNavbar />
        {children}
      </div>
    </div>
  );
}

{
  /* <div className="mx-auto max-w-7xl"> */
}
