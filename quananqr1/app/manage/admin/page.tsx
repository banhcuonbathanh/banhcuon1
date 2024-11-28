import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { decodeToken } from "@/lib/utils";
import { Role } from "@/constants/type";

export default async function ManageHomePage() {
  console.log("quananqr1/app/manage/admin/page.tsx ManageHomePage");
  const cookieStore = cookies();
  const accessToken = cookieStore.get("accessToken")?.value;

  // Double-check authorization on server side
  if (!accessToken) {
    redirect("/login");
  }
  console.log("quananqr1/app/manage/admin/page.tsx ManageHomePage 222");
  try {
    const decoded = decodeToken(accessToken);
    if (!(decoded.role === Role.Admin || decoded.role === Role.Employee)) {
      redirect("/unauthorized");
    }
  } catch (error) {
    redirect("/login");
  }
  console.log("quananqr1/app/manage/admin/page.tsx ManageHomePage 333");
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold">Manage Dashboard</h1>
      {/* Add your management dashboard content here */}
    </div>
  );
}
