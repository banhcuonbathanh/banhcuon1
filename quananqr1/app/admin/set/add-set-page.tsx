"use client";

import React, { useState, useRef, useMemo, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormField,
  FormItem,
  FormControl,
  FormMessage
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "@/components/ui/select";
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
import { Upload, PlusCircle } from "lucide-react";
import envConfig from "@/config";
import { DishStatus, DishStatusValues } from "@/constants/type";
import { handleErrorApi, getVietnameseDishStatus } from "@/lib/utils";
import {
  CreateDishBodyType,
  CreateDishBody
} from "@/schemaValidations/dish.schema";
import {
  useDishListQuery,
  useDishStore
} from "@/zusstand/dished/dished-controller";
import { useMediaStore } from "@/zusstand/media/usemediastore";

export default function AddSetPage() {
  const [file, setFile] = useState<File | null>(null);
  const { uploadMedia, isUploading: isUploadingMedia } = useMediaStore();
  const { addDish, isLoading: isAddingDish } = useDishStore();

  const imageInputRef = useRef<HTMLInputElement | null>(null);
  const form = useForm<CreateDishBodyType>({
    resolver: zodResolver(CreateDishBody),
    defaultValues: {
      name: "",
      description: "",
      price: 0,
      image: undefined,
      status: DishStatus.Unavailable
    }
  });
  const image = form.watch("image");
  const name = form.watch("name");
  const previewAvatarFromFile = useMemo(() => {
    if (file) {
      return URL.createObjectURL(file);
    }
    return image;
  }, [file, image]);

  const reset = () => {
    form.reset();
    setFile(null);
  };

  const onSubmit = async (values: CreateDishBodyType) => {
    if (isAddingDish || isUploadingMedia) return;
    try {
      let body = values;
      if (file) {
        const imageUrl = await uploadMedia(
          file,
          envConfig.NEXT_PUBLIC_Folder1_BE + values.name
        );

        console.log(
          "quananqr1/app/admin/test/add-dish.tsx onSubmit imageUrl",
          imageUrl
        );
        body = {
          ...values,
          image:
            envConfig.NEXT_PUBLIC_API_ENDPOINT +
            envConfig.NEXT_PUBLIC_Upload +
            imageUrl.path
        };
      }

      console.log(
        "quananqr1/app/admin/test/add-dish.tsx onSubmit body with link image",
        body
      );
      const result = await addDish(body);
      reset();
      // You might want to add some success feedback here
    } catch (error) {
      handleErrorApi({
        error,
        setError: form.setError
      });
    }
  };

  return (
    <div className="container mx-auto">
      <h1 className="text-2xl font-bold mb-4">Thêm món ăn</h1>
      <Form {...form}>
        <form
          noValidate
          className="grid auto-rows-max items-start gap-4 md:gap-8"
          onSubmit={form.handleSubmit(onSubmit, (e) => {
            console.log(e);
          })}
          onReset={reset}
        >
          <div className="grid gap-4 py-4">
            <FormField
              control={form.control}
              name="image"
              render={({ field }) => (
                <FormItem>
                  <div className="flex gap-2 items-start justify-start">
                    <Avatar className="aspect-square w-[100px] h-[100px] rounded-md object-cover">
                      <AvatarImage src={previewAvatarFromFile} />
                      <AvatarFallback className="rounded-none">
                        {name || "Ảnh món ăn"}
                      </AvatarFallback>
                    </Avatar>
                    <input
                      type="file"
                      accept="image/*"
                      ref={imageInputRef}
                      onChange={(e) => {
                        const file = e.target.files?.[0];
                        if (file) {
                          setFile(file);
                          field.onChange("http://localhost:3000/" + file.name);
                        }
                      }}
                      className="hidden"
                    />
                    <button
                      className="flex aspect-square w-[100px] items-center justify-center rounded-md border border-dashed"
                      type="button"
                      onClick={() => imageInputRef.current?.click()}
                    >
                      <Upload className="h-4 w-4 text-muted-foreground" />
                      <span className="sr-only">Upload</span>
                    </button>
                  </div>
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <div className="grid grid-cols-4 items-center justify-items-start gap-4">
                    <Label htmlFor="name">Tên món ăn</Label>
                    <div className="col-span-3 w-full space-y-2">
                      <Input id="name" className="w-full" {...field} />
                      <FormMessage />
                    </div>
                  </div>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="price"
              render={({ field }) => (
                <FormItem>
                  <div className="grid grid-cols-4 items-center justify-items-start gap-4">
                    <Label htmlFor="price">Giá</Label>
                    <div className="col-span-3 w-full space-y-2">
                      <Input
                        id="price"
                        className="w-full"
                        {...field}
                        type="number"
                      />
                      <FormMessage />
                    </div>
                  </div>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <div className="grid grid-cols-4 items-center justify-items-start gap-4">
                    <Label htmlFor="description">Mô tả sản phẩm</Label>
                    <div className="col-span-3 w-full space-y-2">
                      <Textarea
                        id="description"
                        className="w-full"
                        {...field}
                      />
                      <FormMessage />
                    </div>
                  </div>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="status"
              render={({ field }) => (
                <FormItem>
                  <div className="grid grid-cols-4 items-center justify-items-start gap-4">
                    <Label htmlFor="description">Trạng thái</Label>
                    <div className="col-span-3 w-full space-y-2">
                      <Select
                        onValueChange={field.onChange}
                        defaultValue={field.value}
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Chọn trạng thái" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {DishStatusValues.map((status) => (
                            <SelectItem key={status} value={status}>
                              {getVietnameseDishStatus(status)}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <FormMessage />
                  </div>
                </FormItem>
              )}
            />
          </div>
          <div className="flex justify-end space-x-4">
            <Button type="reset" variant="outline">
              Hủy
            </Button>
            <Button type="submit">Thêm</Button>
          </div>
        </form>
      </Form>

      {/* <DishClient data={dishes} /> */}
    </div>
  );
}
