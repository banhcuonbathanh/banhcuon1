import { Role, RoleType } from "@/constants/type";
import { decodeToken } from "@/lib/utils";
import { NextRequest, NextResponse } from "next/server";

// Define path configurations
const adminPaths = ["/manage/admin"];
const employeePaths = ["/manage/employee"];

const privatePaths = [...adminPaths, ...employeePaths];
const unAuthPaths = ["/auth"];
const wellComaePage = ["/"];

// Define allowed roles for different paths
const pathRoleConfig: Record<string, RoleType[]> = {
  "/manage/admin": [Role.Admin],
  "/manage/employee": [Role.Employee, Role.Admin]
};

export function middleware(request: NextRequest) {
  // console.log("quananqr1/middleware.ts 1111");
  const { pathname } = request.nextUrl;
  const accessToken = request.cookies.get("accessToken")?.value;
  const refreshToken = request.cookies.get("refreshToken")?.value;

  // If trying to access login pages while already authenticated
  if (unAuthPaths.includes(pathname) && accessToken) {
    // Check if there's a 'from' parameter to redirect after login
    const fromPath = request.nextUrl.searchParams.get("from");

    // If there's a specific path to redirect to, use that
    if (fromPath && fromPath !== "/") {
      return NextResponse.redirect(new URL(fromPath, request.url));
    }
    // console.log("quananqr1/middleware.ts 1111 aaaaaaa");
    // Otherwise, redirect to welcome page
    return NextResponse.redirect(new URL(`${wellComaePage}`, request.url));
  }
  // console.log("quananqr1/middleware.ts 222222");
  // If trying to access protected route without authentication
  if (privatePaths.some((path) => pathname.startsWith(path)) && !accessToken) {
    const url = new URL(`${unAuthPaths}`, request.url);
    url.searchParams.set("from", pathname);

    // console.log("quananqr1/middleware.ts 222222 bbbbbbb");
    return NextResponse.redirect(url);
  }
  // console.log("quananqr1/middleware.ts 33333");
  // If authenticated, check role-based access
  if (accessToken) {
    try {
      const decoded = decodeToken(accessToken);
      const userRole = decoded.role as RoleType;

      // Check if the current path requires specific roles
      for (const [path, allowedRoles] of Object.entries(pathRoleConfig)) {
        if (pathname.startsWith(path) && !allowedRoles.includes(userRole)) {
          // Redirect to unauthorized page or default dashboard

          // console.log("quananqr1/middleware.ts 33333 ccccccc");
          return NextResponse.redirect(new URL("/", request.url));
        }
      }
    } catch (error) {
      // Token decode failed - redirect to login
      const url = new URL("/auth", request.url);
      url.searchParams.set("from", pathname);
      return NextResponse.redirect(url);
    }
  }
  // console.log("quananqr1/middleware.ts 444444");
  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"]
};
