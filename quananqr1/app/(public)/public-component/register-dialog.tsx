"use client";

import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { Info } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogTrigger,
  DialogTitle
} from "@/components/ui/dialog";
import { Form, FormField, FormItem, FormMessage } from "@/components/ui/form";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardHeader,
  CardDescription,
  CardContent
} from "@/components/ui/card";
import { handleErrorApi } from "@/lib/utils";
import { useAuthStore } from "@/zusstand/auth/controller/auth-controller";
import {
  RegisterBodyType,
  RegisterBody
} from "@/zusstand/auth/domain/auth.schema";
import axios from "axios";

const RegisterDialog = () => {
  const handleAddUser = async () => {
    console.log(
      "quananqr1/app/(public)/public-component/register-dialog.tsx hander use"
    );
    const userData = {
      name: "Alice ",
      email: "alice.johnson@example.com11111",
      password: "password1231234",
      is_admin: false,
      phone: 1234567890,
      image: "alice.jpg",
      address: "123 Main St, Anytown, USA",
      created_at: "2024-08-19T16:17:16+07:00",
      updated_at: "2024-08-19T16:17:16+07:00"
    };

    try {
      const response = await axios.post(
        "http://localhost:8888/users",
        userData
      );
      console.log("User added successfully:", response.data);
    } catch (error) {
      console.error("Error adding user:", error);
    }
  };

  const [serverStatus, setServerStatus] = useState("");

  // const checkServerConnection = async () => {
  //   console.log("checkServerConnection");
  //   try {
  //     const response = await axios.get("http://localhost:8888/test");
  //     console.log("checkServerConnectio n  done", response);
  //     if (response.status === 200) {
  //       setServerStatus("Connected to server successfully");
  //     } else {
  //       setServerStatus("Failed to connect to server");
  //     }
  //   } catch (error) {
  //     setServerStatus("Error connecting to server");
  //     console.error("Server connection error:", error);
  //   }
  // };
  // console.log("checkServerConnectio n  done");
  // useEffect(() => {
  //   checkServerConnection();
  // }, []);

  const { register, openLoginDialog } = useAuthStore();
  const form = useForm<RegisterBodyType>({
    resolver: zodResolver(RegisterBody),
    defaultValues: {
      name: "",
      email: "",
      password: "",
      is_admin: false,
      phone: 1234,
      image: "",
      address: "",
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }
  });
  const [open, setOpen] = useState(false);
  //   {
  //     "name": "Alice Johnson12341234",
  //     "email": "alice.johnson@example.com12341234",
  //     "password": "password1231234",
  //     "is_admin": false,
  //     "phone": 1234567890,
  //     "image": "alice.jpg",
  //     "address": "123 Main St, Anytown, USA",
  //     "created_at": "2024-08-19T16:17:16+07:00",
  //     "updated_at": "2024-08-19T16:17:16+07:00"

  // }
  const onSubmit = async () => {
    try {
      console.log(
        "onSubmit register form quananqr1/app/(public)/public-component/register-dialog.tsx"
      );

      await register({
        name: "Alice Jo1234f",
        email: "alice.johnson@example.vvvvvvv",
        password: "password123@%$@1234",
        is_admin: false,
        phone: 1234567890,
        image: "alice.jpg",
        address: "123 Main St, Anytown, USA",
        created_at: "2024-08-19T16:17:16+07:00",
        updated_at: "2024-08-19T16:17:16+07:00"
      });
      // setOpen(false);
      // openLoginDialog();
    } catch (error: any) {
      console.log("Error during registration: ", error);
      handleErrorApi({
        error,
        setError: form.setError
      });
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>Đăng ký asdfasdf</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px] bg-white dark:bg-gray-800 shadow-lg">
        <DialogTitle className="text-2xl font-semibold">
          Đăng ký sdFASDFASD
        </DialogTitle>
        <Card className="border-0 shadow-none">
          <CardHeader>
            <CardDescription>
              Điền thông tin của bạn để tạo tài khoản mới
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Form {...form}>
              <form
                className="space-y-4 w-full"
                noValidate
                onSubmit={form.handleSubmit(onSubmit, (err) => {
                  console.log(
                    "Registration err onSubmit: quananqr1/app/(public)/public-component/register-dialog.tsx",
                    err
                  );
                })}
              >
                <div className="grid gap-4">
                  <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                      <FormItem>
                        <div className="grid gap-2">
                          <Label htmlFor="name">Họ và tên</Label>
                          <Input
                            id="name"
                            type="text"
                            placeholder="Alice Johnson"
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
                    name="email"
                    render={({ field }) => (
                      <FormItem>
                        <div className="grid gap-2">
                          <Label htmlFor="email">Email</Label>
                          <Input
                            id="email"
                            type="email"
                            placeholder="alice.johnson@example.com"
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
                          <Label htmlFor="password">Mật khẩu</Label>
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
                              Mật khẩu phải có ít nhất 8 ký tự và bao gồm chữ
                              cái, số và ký tự đặc biệt.
                            </span>
                          </div>
                          <FormMessage />
                        </div>
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={form.control}
                    name="phone"
                    render={({ field }) => (
                      <FormItem>
                        <div className="grid gap-2">
                          <Label htmlFor="phone">Số điện thoại</Label>
                          <Input
                            id="phone"
                            type="tel"
                            placeholder="1234567890"
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
                    name="address"
                    render={({ field }) => (
                      <FormItem>
                        <div className="grid gap-2">
                          <Label htmlFor="address">Địa chỉ</Label>
                          <Input
                            id="address"
                            type="text"
                            placeholder="123 Main St, Anytown, USA"
                            required
                            className="border-2 border-gray-300 dark:border-gray-600"
                            {...field}
                          />
                          <FormMessage />
                        </div>
                      </FormItem>
                    )}
                  />
                  <Button type="submit" className="w-full">
                    Đăng ký
                  </Button>
                  <Button onClick={handleAddUser}>Add User</Button>
                </div>
              </form>
            </Form>
          </CardContent>
        </Card>
      </DialogContent>
    </Dialog>
  );
};

export default RegisterDialog;
