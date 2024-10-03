import { NextResponse } from "next/server";
import {
  CreateDishBody,
  DishListRes,
  DishRes
} from "@/zusstand/dished/domain/dish.schema"; // Adjust the import based on your project structure
import envConfig from "@/config";

export async function GET() {
  console.log("Fetching dishes...");

  try {
    const response = await fetch(
      `${envConfig.NEXT_PUBLIC_API_ENDPOINT}${envConfig.NEXT_PUBLIC_Add_Dished}`,
      {
        method: "GET",
        cache: "no-store"
      }
    );

    // Check for response.ok before parsing JSON
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const dishes = await response.json(); // Parse response once
    console.log("Fetched dishes:", dishes);

    return NextResponse.json({
      data: dishes.data,
      message: "Dishes retrieved successfully"
    });
  } catch (error) {
    console.error("Error fetching dishes:", error);
    return NextResponse.json(
      { message: "Failed to fetch dishes asdfasdf" },
      { status: 500 }
    );
  }
}

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const validatedBody = CreateDishBody.parse(body);

    const response = await fetch(
      `${envConfig.NEXT_PUBLIC_API_ENDPOINT}${envConfig.NEXT_PUBLIC_Add_Dished}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(validatedBody)
      }
    );

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const newDish = await response.json();
    return NextResponse.json(
      DishRes.parse({ data: newDish, message: "Dish created successfully" }),
      { status: 201 }
    );
  } catch (error) {
    console.error("Error creating dish:", error);
    if (error instanceof Error) {
      return NextResponse.json(
        { message: `Failed to create dish: ${error.message}` },
        { status: 400 }
      );
    }
    return NextResponse.json(
      { message: "Failed to create dish" },
      { status: 400 }
    );
  }
}
