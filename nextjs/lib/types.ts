import { links } from "../lib/header/data";

export type SectionName = (typeof links)[number]["name"];



export type SideNavItem = {
    title: string;
    path: string;
    icon?: JSX.Element;
    submenu?: boolean;
    subMenuItems?: SideNavItem[];
  };
  
  export type SideNavItemGroup = {
    title: string;
    menuList: SideNavItem[]
  }