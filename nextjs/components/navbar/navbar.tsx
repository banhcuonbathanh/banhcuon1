import React from "react";
import ImageContainer from "./image_container";
import NavLinks from "./nav_link";
import { ModeToggle } from "./themes_toogles";

export function MainNav({
  className,
  ...props
}: React.HTMLAttributes<HTMLElement>) {
  return (
    <div className=" flex  fixed  w-full justify-center gap-6 top-0 bg-background items-center py-6">
      {/* <ImageContainer /> */}
      <NavLinks />
      <ModeToggle />
    </div>
  );
}
