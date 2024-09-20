import { SideNavItemGroup } from "@/lib/types";
import { ChevronRight } from "lucide-react";

export const SIDENAV_ITEMS: SideNavItemGroup[] = [
  {
    title: "Dashboards",
    menuList: [
      {
        title: "Dashboard",
        path: "/",
        icon: <ChevronRight size={20} />
      }
    ]
  },
  {
    title: "Manage",
    menuList: [
      {
        title: "Products",
        path: "/products",
        icon: <ChevronRight size={20} />,
        submenu: true,
        subMenuItems: [
          { title: "All", path: "/products" },
          { title: "New", path: "/products/new" }
        ]
      },
      {
        title: "Orders",
        path: "/orders",
        icon: <ChevronRight size={20} />
      },
      {
        title: "Feedbacks",
        path: "/feedbacks",
        icon: <ChevronRight size={20} />
      }
    ]
  },
  {
    title: "Others",
    menuList: [
      {
        title: "Account",
        path: "/account",
        icon: <ChevronRight size={20} />
      },
      {
        title: "Help",
        path: "/help",
        icon: <ChevronRight size={20} />
      }
    ]
  }
];
