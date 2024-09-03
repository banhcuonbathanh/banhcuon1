import React from "react";
import ImageContainer from "./image_container";
import NavLinks from "./nav_link";
import { ModeToggle } from "./themes_toogles";

export function MainNav({
  className,
  ...props
}: React.HTMLAttributes<HTMLElement>) {
  return (
    <div className="flex items-center gap-x-6 pt-3">
      <ImageContainer />
      <NavLinks />
      <ModeToggle />
    </div>
  );
}
