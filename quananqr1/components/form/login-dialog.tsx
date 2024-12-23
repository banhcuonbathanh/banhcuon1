"use client";

import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Info } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle
} from "@/components/ui/dialog";
import { Form, FormField, FormItem, FormMessage } from "@/components/ui/form";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { LoginBodyType, LoginBody } from "@/schemaValidations/auth.schema";
import { useAuthStore } from "@/zusstand/new_auth/new_auth_controller";
import { usePathname, useRouter } from "next/navigation";
import { handleErrorApi, handleLoginRedirect } from "@/lib/utils";
import { logWithLevel } from "@/lib/log";

const LOG_PATH = "quananqr1/components/form/login-dialog.tsx";

const LoginDialog1 = () => {
  // Log component initialization
  logWithLevel(
    { component: "LoginDialog1", event: "initialization" },
    LOG_PATH,
    "debug",
    1
  );

  const {
    login,
    isLoginDialogOpen,
    closeLoginDialog,
    openRegisterDialog,
    openGuestDialog
  } = useAuthStore();
  const pathname = usePathname();
  const router = useRouter();

  useEffect(() => {
    // Log pathname changes
    logWithLevel({ pathname, event: "pathname_change" }, LOG_PATH, "debug", 1);
  }, [pathname]);

  const form = useForm<LoginBodyType>({
    resolver: zodResolver(LoginBody),
    defaultValues: {
      email: "",
      password: ""
    }
  });

  const onSubmit = async (data: LoginBodyType) => {
    // Log form submission attempt
    logWithLevel(
      {
        event: "form_submission",
        email: data.email,
        hasPassword: !!data.password
      },
      LOG_PATH,
      "info",
      2
    );

    try {
      await login(data);

      // Log successful login
      logWithLevel(
        {
          event: "login_success",
          email: data.email
        },
        LOG_PATH,
        "info",
        3
      );

      handleLoginRedirect(pathname, router);

      // Log navigation attempt
      logWithLevel(
        {
          event: "navigation_redirect",
          pathname
        },
        LOG_PATH,
        "debug",
        6
      );
    } catch (error: any) {
      // Log login error
      logWithLevel(
        {
          event: "login_error",
          error: error.message,
          code: error.code
        },
        LOG_PATH,
        "error",
        5
      );

      handleErrorApi({
        error,
        setError: form.setError
      });
    }
  };

  const handleRegisterClick = () => {
    // Log dialog state change
    logWithLevel(
      {
        event: "dialog_state_change",
        action: "switch_to_register"
      },
      LOG_PATH,
      "debug",
      4
    );

    closeLoginDialog();
    openRegisterDialog();
  };

  const handleGuestClick = () => {
    // Log dialog state change
    logWithLevel(
      {
        event: "dialog_state_change",
        action: "switch_to_guest"
      },
      LOG_PATH,
      "debug",
      4
    );

    closeLoginDialog();
    openGuestDialog();
  };

  // Log dialog visibility changes
  useEffect(() => {
    logWithLevel(
      {
        event: "dialog_visibility_change",
        isOpen: isLoginDialogOpen
      },
      LOG_PATH,
      "debug",
      4
    );
  }, [isLoginDialogOpen]);

  return (
    <Dialog
      open={isLoginDialogOpen}
      onOpenChange={(open) => {
        if (!open) {
          // Log dialog close
          logWithLevel(
            {
              event: "dialog_close",
              trigger: "user_action"
            },
            LOG_PATH,
            "debug",
            4
          );
          closeLoginDialog();
        }
      }}
    >
      <DialogContent className="sm:max-w-[425px] bg-white dark:bg-gray-800 shadow-lg">
        <DialogHeader>
          <DialogTitle className="text-2xl">Đăng nhập</DialogTitle>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Nhập email và mật khẩu của bạn để đăng nhập vào hệ thống
          </p>
        </DialogHeader>
        <Form {...form}>
          <form
            className="space-y-4 w-full"
            noValidate
            onSubmit={form.handleSubmit(onSubmit)}
          >
            <div className="grid gap-4">
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <div className="grid gap-2">
                      <Label htmlFor="email">Email</Label>
                      <Input
                        id="email"
                        type="email"
                        placeholder="m@example.com"
                        required
                        className="border-2 border-gray-300 dark:border-gray-600"
                        {...field}
                      />
                      <FormMessage />
                    </div>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <div className="grid gap-2">
                      <Label htmlFor="password">Password</Label>
                      <Input
                        id="password"
                        type="password"
                        placeholder="••••••••"
                        required
                        className="border-2 border-gray-300 dark:border-gray-600"
                        {...field}
                      />
                      <div className="text-sm text-gray-500 dark:text-gray-400 flex items-center gap-1">
                        <Info size={16} />
                        <span>
                          Password should be at least 8 characters long and
                          include a mix of letters, numbers, and symbols.
                        </span>
                      </div>
                      <FormMessage />
                    </div>
                  </FormItem>
                )}
              />
              <Button type="submit" className="w-full">
                Đăng nhập
              </Button>
              <div className="flex flex-col gap-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={handleRegisterClick}
                  className="w-full"
                >
                  Đăng ký tài khoản mới
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  onClick={handleGuestClick}
                  className="w-full"
                >
                  Đăng nhập với tư cách khách
                </Button>
              </div>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};

export default LoginDialog1;
