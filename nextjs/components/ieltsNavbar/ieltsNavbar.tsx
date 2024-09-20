import React from "react";
import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger
} from "@/components/ui/navigation-menu";
import { ModeToggle } from "../navbar/themes_toogles";

const IeltsNavBar = () => {
  return (
    <nav className="p-4 gap-8">
      <div className="container mx-auto flex justify-evenly items-center gap-5">
        <a href="/" className="text-2xl font-bold">
          {/* Replace with your logo or site name */}
          <span className="sr-only">Home</span>
          🏠
        </a>

        <NavigationMenu className="ml-56">
          <NavigationMenuList className="flex gap-3">
            {"dsfgsdf "}
            {/* Added flex and gap-3 here */}
            <NavigationMenuItem className="ml-56">
              <NavigationMenuTrigger className="hover:bg-[#34495e]">
                IELTS Exam Library
              </NavigationMenuTrigger>
              <NavigationMenuContent className="ml-56 bg-black">
                <ul className="grid gap-3 p-6 md:w-[400px] lg:w-[500px] ml-56">
                  asdf asdf asdf
                </ul>
                <ul className="grid gap-3 p-6 md:w-[400px] lg:w-[500px]">
                  asdf asdf asdf
                </ul>
              </NavigationMenuContent>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <NavigationMenuTrigger className="hover:bg-[#34495e]">
                IELTS Tips
              </NavigationMenuTrigger>
              <NavigationMenuContent>
                <ul className="grid gap-3 p-6 md:w-[400px] lg:w-[500px]">
                  1241234213423
                </ul>
                <ul className="grid gap-3 p-6 md:w-[400px] lg:w-[500px]">
                  123412341234
                </ul>
              </NavigationMenuContent>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <NavigationMenuTrigger className="hover:bg-[#34495e]">
                IELTS Prep
              </NavigationMenuTrigger>
              <NavigationMenuContent>
                <ul className="grid gap-3 p-6 md:w-[400px] lg:w-[500px]">
                  qwerqwer
                </ul>
                <ul className="grid gap-3 p-6 md:w-[400px] lg:w-[500px]">
                  qwerqwer
                </ul>
              </NavigationMenuContent>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <NavigationMenuLink className="hover:bg-[#34495e]">
                Live Lessons
              </NavigationMenuLink>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <NavigationMenuLink className="bg-[#e74c3c] hover:bg-[#c0392b] rounded">
                IELTS Courses
              </NavigationMenuLink>
            </NavigationMenuItem>
            <NavigationMenuItem className=" bg-black">
              <NavigationMenuTrigger className="hover:bg-[#34495e]">
                InterGreat Study Abroad
              </NavigationMenuTrigger>
              <NavigationMenuContent className="bg-black ml-48">
                <ul>zxcvzxcv</ul>
                <ul>zxcvzxcv</ul>
              </NavigationMenuContent>
            </NavigationMenuItem>
          </NavigationMenuList>

          <NavigationMenuList className="flex gap-3">
            {"dsfgsdf "}
            {/* Added flex and gap-3 here */}
            <NavigationMenuItem className="ml-56">
              <NavigationMenuTrigger className="hover:bg-[#34495e]">
                IELTS Exam Library
              </NavigationMenuTrigger>
              <NavigationMenuContent className="ml-56 bg-black">
                <ul className="grid gap-3 p-6 md:w-[400px] lg:w-[500px] ml-56">
                  asdf asdf asdf
                </ul>
                <ul className="grid gap-3 p-6 md:w-[400px] lg:w-[500px]">
                  asdf asdf asdf
                </ul>
              </NavigationMenuContent>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <NavigationMenuTrigger className="hover:bg-[#34495e]">
                IELTS Tips
              </NavigationMenuTrigger>
              <NavigationMenuContent>
                <ul className="grid gap-3 p-6 md:w-[400px] lg:w-[500px]">
                  1241234213423
                </ul>
                <ul className="grid gap-3 p-6 md:w-[400px] lg:w-[500px]">
                  123412341234
                </ul>
              </NavigationMenuContent>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <NavigationMenuTrigger className="hover:bg-[#34495e]">
                IELTS Prep
              </NavigationMenuTrigger>
              <NavigationMenuContent>
                <ul className="grid gap-3 p-6 md:w-[400px] lg:w-[500px]">
                  qwerqwer
                </ul>
                <ul className="grid gap-3 p-6 md:w-[400px] lg:w-[500px]">
                  qwerqwer
                </ul>
              </NavigationMenuContent>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <NavigationMenuLink className="hover:bg-[#34495e]">
                Live Lessons
              </NavigationMenuLink>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <NavigationMenuLink className="bg-[#e74c3c] hover:bg-[#c0392b] rounded">
                IELTS Courses
              </NavigationMenuLink>
            </NavigationMenuItem>
            <NavigationMenuItem className=" bg-black">
              <NavigationMenuTrigger className="hover:bg-[#34495e]">
                InterGreat Study Abroad
              </NavigationMenuTrigger>
              <NavigationMenuContent className="bg-black ml-48">
                <ul>zxcvzxcv</ul>
                <ul>zxcvzxcv</ul>
              </NavigationMenuContent>
            </NavigationMenuItem>
          </NavigationMenuList>
        </NavigationMenu>

        <div className="flex items-center space-x-4">
          <button className=" hover:underline">Sign Up</button>
          <button className=" hover:underline">Log In</button>
        </div>
        <p>test</p>
        <ModeToggle />
      </div>
    </nav>
  );
};

export default IeltsNavBar;
