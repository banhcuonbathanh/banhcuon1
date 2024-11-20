import { Role } from "@/constants/type";
import { decodeToken } from "@/lib/utils";
import { NextRequest, NextResponse } from "next/server";

// Define path configurations
const managePaths = ["/manage"];
const guestPaths = ["/guest"];
const privatePaths = [...managePaths, ...guestPaths];
const unAuthPaths = ["/login"];

// Define allowed roles for different paths - updated to match Role enum casing
const pathRoleConfig = {
  "/manage": [Role.Admin, Role.Employee],
  "/guest": [Role.Guest, Role.Admin, Role.Employee]
};

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const accessToken = request.cookies.get("accessToken")?.value;
  const refreshToken = request.cookies.get("refreshToken")?.value;

  // If trying to access login page while already authenticated
  if (unAuthPaths.includes(pathname) && accessToken) {
    return NextResponse.redirect(new URL("/manage", request.url));
  }

  // If trying to access protected route without authentication
  if (privatePaths.some((path) => pathname.startsWith(path)) && !accessToken) {
    const url = new URL("/login", request.url);
    url.searchParams.set("from", pathname);
    return NextResponse.redirect(url);
  }

  // If authenticated, check role-based access
  if (accessToken) {
    try {
      const decoded = decodeToken(accessToken);
      const userRole = decoded.role;

      // Check if the current path requires specific roles
      for (const [path, allowedRoles] of Object.entries(pathRoleConfig)) {
        if (pathname.startsWith(path)) {
          // if (!allowedRoles.includes(userRole)) {
          //   // Redirect to appropriate page based on role - updated to match Role enum casing
          //   if (userRole === Role.Guest) {
          //     return NextResponse.redirect(new URL("/guest", request.url));
          //   }
          //   return NextResponse.redirect(new URL("/unauthorized", request.url));
          // }
        }
      }
    } catch (error) {
      // Token decode failed - redirect to login
      const url = new URL("/login", request.url);
      url.searchParams.set("from", pathname);
      return NextResponse.redirect(url);
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"]
};
// See "Matching Paths" below to learn more
// export const config = {
//   matcher: ["/manage/:path*", "/guest/:path*", "/login"]
// };

//   const { pathname } = request.nextUrl;
//   // pathname: /manage/dashboard
//   const accessToken = request.cookies.get("accessToken")?.value;
//   const refreshToken = request.cookies.get("refreshToken")?.value;
//   // 1. Chưa đăng nhập thì không cho vào private paths
//   if (privatePaths.some((path) => pathname.startsWith(path)) && !refreshToken) {
//     const url = new URL("/login", request.url);
//     url.searchParams.set("clearTokens", "true");
//     return NextResponse.redirect(url);
//   }

//   // 2. Trường hợp đã đăng nhập
//   if (refreshToken) {
//     // 2.1 Nếu cố tình vào trang login sẽ redirect về trang chủ
//     if (unAuthPaths.some((path) => pathname.startsWith(path))) {
//       return NextResponse.redirect(new URL("/", request.url));
//     }

//     // 2.2 Nhưng access token lại hết hạn
//     if (
//       privatePaths.some((path) => pathname.startsWith(path)) &&
//       !accessToken
//     ) {
//       const url = new URL("/refresh-token", request.url);
//       url.searchParams.set("refreshToken", refreshToken);
//       url.searchParams.set("redirect", pathname);
//       return NextResponse.redirect(url);
//     }

//     // 2.3 Vào không đúng role, redirect về trang chủ
//     const role = decodeToken(refreshToken).role;
//     // Guest nhưng cố vào route owner
//     const isGuestGoToManagePath =
//       role === Role.Guest &&
//       managePaths.some((path) => pathname.startsWith(path));
//     // Không phải Guest nhưng cố vào route guest
//     const isNotGuestGoToGuestPath =
//       role !== Role.Guest &&
//       guestPaths.some((path) => pathname.startsWith(path));
//     if (isGuestGoToManagePath || isNotGuestGoToGuestPath) {
//       return NextResponse.redirect(new URL("/", request.url));
//     }

//     return NextResponse.next();
//   }

//

// import { NextRequest, NextResponse } from 'next/server';
// import { jwtVerify } from 'jose';

// const JWT_SECRET = process.env.JWT_SECRET || 'your-secret-key';

// export async function middleware(request: NextRequest) {
//   const { pathname } = request.nextUrl;
//   const accessToken = request.cookies.get("accessToken")?.value;

//   // Public routes that don't require authentication
//   const publicRoutes = ['/login', '/register', '/forgot-password'];
//   if (publicRoutes.includes(pathname)) {
//     return NextResponse.next();
//   }

//   if (!accessToken) {
//     return NextResponse.redirect(new URL('/login', request.url));
//   }

//   try {
//     const { payload } = await jwtVerify(accessToken, new TextEncoder().encode(JWT_SECRET));
//     const role = payload.role as string;

//     // Define role-based route permissions
//     const routePermissions = {
//       '/admin': ['Admin'],
//       '/manager': ['Manager', 'Admin'],
//       '/employee': ['Employee', 'Manager', 'Admin'],
//     };

//     // Check if the current path requires specific roles
//     for (const [route, allowedRoles] of Object.entries(routePermissions)) {
//       if (pathname.startsWith(route) && !allowedRoles.includes(role)) {
//         return NextResponse.redirect(new URL('/unauthorized', request.url));
//       }
//     }

//     // For all other cases, allow the request to proceed
//     return NextResponse.next();
//   } catch (error) {
//     console.error('Token verification failed:', error);
//     return NextResponse.redirect(new URL('/login', request.url));
//   }
// }

// export const config = {
//   matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"]
// };
