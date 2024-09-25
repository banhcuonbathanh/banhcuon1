"use client";

import React from "react";
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useForm } from "react-hook-form";
import { Form, FormField, FormItem, FormMessage } from "@/components/ui/form";
import { LoginBody, LoginBodyType } from "@/schemaValidations/auth.schema";
import { zodResolver } from "@hookform/resolvers/zod";
import { handleErrorApi } from "@/lib/utils";
import { useRouter, useSearchParams } from "next/navigation";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from "@/components/ui/card";
import { useAuth } from "@/zusstand/auth/useauth";
import { Info } from "lucide-react";

const LoginDialog = () => {
  const { login, user } = useAuth();
  const searchParams = useSearchParams();
  const clearTokens = searchParams.get("clearTokens");

  const form = useForm<LoginBodyType>({
    resolver: zodResolver(LoginBody),
    defaultValues: {
      email: "",
      password: ""
    }
  });

  const router = useRouter();

  const onSubmit = async (data: LoginBodyType) => {
    try {
      const result = await login(data);
      router.push("/manage/dashboard");
    } catch (error: any) {
      handleErrorApi({
        error,
        setError: form.setError
      });
    }
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button>Log In</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px] bg-white dark:bg-gray-800 shadow-lg">
        <CardHeader>
          <CardTitle className="text-2xl">Đăng nhập</CardTitle>
          <CardDescription>
            Nhập email và mật khẩu của bạn để đăng nhập vào hệ thống
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form
              className="space-y-4 w-full"
              noValidate
              onSubmit={form.handleSubmit(onSubmit, (err) => {
                console.log(err);
              })}
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
                          className="border-2 border-gray-300  dark:border-gray-600 "
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
                          placeholder="password"
                          id="password"
                          type="password"
                          required
                          className="border-2 border-gray-300  dark:border-gray-600"
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
                <Button variant="outline" className="w-full" type="button">
                  Đăng nhập bằng Google
                </Button>
              </div>
            </form>
          </Form>
        </CardContent>
      </DialogContent>
    </Dialog>
  );
};

export default LoginDialog;
