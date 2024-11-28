import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { decodeToken } from "@/lib/utils";
import { Role } from "@/constants/type";

export default async function ManageHomePage() {
  console.log("quananqr1/app/manage/admin/page.tsx ManageHomePage");
  const cookieStore = cookies();
  const accessToken = cookieStore.get("accessToken")?.value;

  console.log("quananqr1/app/manage/admin/page.tsx ManageHomePage");
  // Double-check authorization on server side
  if (!accessToken) {
    redirect("/login");
  }
  const decoded = decodeToken(accessToken);
  console.log(
    "quananqr1/app/manage/admin/page.tsx ManageHomePage 222 decoded",
    decoded
  );
  try {
    const decoded = decodeToken(accessToken);
    if (!(decoded.role === Role.Admin || decoded.role === Role.Employee)) {
      redirect("/manage/employee");
    }
  } catch (error) {
    redirect("/auth");
  }
  console.log("quananqr1/app/manage/admin/page.tsx ManageHomePage 333");
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold">Manage Dashboard</h1>
      {/* Add your management dashboard content here */}
    </div>
  );
}
