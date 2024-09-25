import accountApiRequest from "@/apiRequests/account";
import { get_Account } from "@/zusstand/auth/server/server-auth-controler";
import { cookies } from "next/headers";

export default async function ManageHomePage() {
  // const cookieStore = cookies();
  // const accessToken = cookieStore.get("accessToken")?.value!;
  // let name = "";
  console.log("ManageHomePage quananqr1/app/manage/page.tsx");
  // const resul = await get_Account("alice.johnson@example.com1234");

  // console.log("ManageHomePage quananqr1/app/manage/page.tsx email", resul);
  // try {
  //   const result = await accountApiRequest.sMe(accessToken)
  //   name = result.payload.data.name
  // } catch (error: any) {
  //   if (error.digest?.includes('NEXT_REDIRECT')) {
  //     throw error
  //   }
  // }
  return <div>ManageHomePage </div>;
}
