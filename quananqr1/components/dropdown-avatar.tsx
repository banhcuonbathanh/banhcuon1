"use client";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import Link from "next/link";
import { useAuthStore } from "@/zusstand/new_auth/new_auth_controller";

export default function DropdownAvatar() {
  const { user, guest, logout, guestLogout, loading, isGuest } = useAuthStore();

  const handleLogout = async () => {
    try {
      // if (isGuest && guest) {
      //   // Handle guest logout with required data
      //   await guestLogout({
      //     body: {
      //       refresh_token: Cookies.get("refreshToken") || ""
      //     }
      //   });
      // } else {
      //   // Handle regular user logout
      //   await logout();
      // }
    } catch (error) {
      console.error("Logout failed:", error);
    }
  };

  // Get display name based on whether it's a guest or regular user
  const getDisplayName = () => {
    if (isGuest && guest) {
      return `Khách ${guest.name} - Bàn ${guest.table_number}`;
    }
    return user?.name || "User";
  };

  // Get avatar initials
  const getAvatarInitials = () => {
    if (isGuest && guest) {
      return `K${guest.table_number}`;
    }
    return user?.name ? user.name.slice(0, 2).toUpperCase() : "U";
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="outline"
          size="icon"
          className="overflow-hidden rounded-full"
        >
          <Avatar>
            <AvatarImage
              src={isGuest ? undefined : user?.image ?? undefined}
              alt={getDisplayName()}
            />
            <AvatarFallback>{getAvatarInitials()}</AvatarFallback>
          </Avatar>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuLabel>{getDisplayName()}</DropdownMenuLabel>
        <DropdownMenuSeparator />

        {/* Only show settings for regular users */}
        {!isGuest && (
          <DropdownMenuItem asChild>
            <Link href="/manage/setting" className="cursor-pointer">
              Cài đặt
            </Link>
          </DropdownMenuItem>
        )}

        <DropdownMenuItem>Hỗ trợ</DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={handleLogout} disabled={loading}>
          {loading ? "Đang đăng xuất..." : "Đăng xuất"}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
