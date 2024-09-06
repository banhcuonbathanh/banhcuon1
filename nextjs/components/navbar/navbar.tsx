import React from "react";
import ImageContainer from "./image_container";
import NavLinks from "./nav_link";
import { ModeToggle } from "./themes_toogles";

export function MainNav({
  className,
  ...props
}: React.HTMLAttributes<HTMLElement>) {
  return (
    <div className="flex items-center gap-x-6 fixed top-0 left-1/2 transform -translate-x-1/2  max-w-screen-xl bg-background rounded-lg p-4">
      {/* <ImageContainer /> */}
      <NavLinks />
      <ModeToggle />
    </div>
  );
}
