import React from "react";
import Image from "next/image";

const ImageContainer = () => {
  return (
    <Image
      src="/next.svg"
      alt="logo"
      width={70}
      height={50}
      style={{ maxWidth: "auto", maxHeight: "auto" }}
    />
  );
};

export default ImageContainer;