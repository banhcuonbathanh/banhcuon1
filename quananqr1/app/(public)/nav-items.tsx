"use client";

import Link from "next/link";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger
} from "@/components/ui/alert-dialog";
import { cn, handleErrorApi } from "@/lib/utils";
import { useAuthStore } from "./zustand-public";
import { Role } from "@/constants/type";
import { RoleType } from "@/types/jwt.types";
import { useRouter } from "next/navigation";

const menuItems: {
  title: string;
  href: string;
  role?: RoleType[];
  hideWhenLogin?: boolean;
}[] = [
  {
    title: "Trang chủ",
    href: "/"
  },
  {
    title: "Menu",
    href: "/guest/menu",
    role: ["Guest"]
  },
  {
    title: "Đơn hàng",
    href: "/guest/orders",
    role: ["Guest"]
  },
  {
    title: "Đăng nhập",
    href: "/login",
    hideWhenLogin: true
  },
  {
    title: "Quản lý",
    href: "/manage/dashboard",
    role: ["Owner", "Employee"]
  }
];

// Type guard function
function isValidRole(role: string): role is RoleType {
  return Object.values(Role).includes(role as RoleType);
}

export default function NavItems({ className }: { className?: string }) {
  const { account, logout: logoutAction } = useAuthStore();
  const router = useRouter();

  const logout = async () => {
    try {
      await logoutAction();
      router.push("/");
    } catch (error: any) {
      handleErrorApi({
        error
      });
    }
  };

  return (
    <>
      {menuItems.map((item) => {
        // Check if the user is authenticated and has the required role
        const isAuth =
          item.role &&
          account &&
          isValidRole(account.role) &&
          item.role.includes(account.role);
        // Check if the item should be shown based on login status
        const canShow =
          (item.role === undefined && !item.hideWhenLogin) ||
          (!account && item.hideWhenLogin);

        if (isAuth || canShow) {
          return (
            <Link href={item.href} key={item.href} className={className}>
              {item.title}
            </Link>
          );
        }
        return null;
      })}
      {account && (
        <AlertDialog>
          <AlertDialogTrigger asChild>
            <div className={cn(className, "cursor-pointer")}>Đăng xuất</div>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Bạn có muốn đăng xuất không?</AlertDialogTitle>
              <AlertDialogDescription>
                Việc đăng xuất có thể làm mất đi hóa đơn của bạn
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Thoát</AlertDialogCancel>
              <AlertDialogAction onClick={logout}>OK</AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      )}
    </>
  );
}
