"use client";

import React, { useEffect } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import useOrderStore from "@/zusstand/order/order_zustand";
import OrderDetails from "../total-dishes-detail";
import OrderCreationComponent from "./add_order_button";
import { decodeTableToken } from "@/lib/utils";
import { logWithLevel } from "@/lib/logger/log";

interface OrderProps {
  number: string;
  token: string;
}

const LOG_PATH =
  "quananqr1/app/(client)/table/[number]/component/order/order.tsx";

export default function OrderSummary({ number, token }: OrderProps) {
  const decoded = decodeTableToken(token);

  // Log token decode result
  logWithLevel({ decoded, token }, LOG_PATH, "debug", 1);

  const {
    addTableNumber,
    addTableToken,
    getOrderSummary,
    canhKhongRau,
    canhCoRau,
    smallBowl,
    wantChili,
    selectedFilling,
    updateCanhKhongRau,
    updateCanhCoRau,
    updateSmallBowl,
    updateWantChili,
    updateSelectedFilling
  } = useOrderStore();

  useEffect(() => {
    try {
      if (token) {
        addTableToken(token);
      }
      if (number) {
        const tableNumber = addTableNumberconvert(number);
        addTableNumber(tableNumber);
      }

      // Log table number and token updates
      logWithLevel(
        {
          token,
          number,
          tableNumber: number ? addTableNumberconvert(number) : null
        },
        LOG_PATH,
        "info",
        2
      );
    } catch (error) {
      logWithLevel({ error, token, number }, LOG_PATH, "error", 2);
    }
  }, [token, addTableToken, number, addTableNumber]);

  const handleBowlChange = (
    type: "chili" | "noChili" | "small",
    change: number
  ) => {
    const prevValues = {
      canhKhongRau,
      canhCoRau,
      smallBowl
    };

    switch (type) {
      case "chili":
        const newToppingValue = canhKhongRau + change;
        if (newToppingValue >= 0) {
          updateCanhKhongRau(newToppingValue);
          logWithLevel(
            {
              type: "canhKhongRau",
              previousValue: canhKhongRau,
              newValue: newToppingValue,
              change
            },
            LOG_PATH,
            "debug",
            3
          );
        }
        break;
      case "noChili":
        const newNoChiliValue = canhCoRau + change;
        if (newNoChiliValue >= 0) {
          updateCanhCoRau(newNoChiliValue);
          logWithLevel(
            {
              type: "canhCoRau",
              previousValue: canhCoRau,
              newValue: newNoChiliValue,
              change
            },
            LOG_PATH,
            "debug",
            3
          );
        }
        break;
      case "small":
        const newSmallBowlValue = smallBowl + change;
        if (newSmallBowlValue >= 0) {
          updateSmallBowl(newSmallBowlValue);
          logWithLevel(
            {
              type: "smallBowl",
              previousValue: smallBowl,
              newValue: newSmallBowlValue,
              change
            },
            LOG_PATH,
            "debug",
            3
          );
        }
        break;
    }
  };

  const handleChiliUpdate = (newValue: boolean) => {
    updateWantChili(newValue);
    logWithLevel(
      {
        previousValue: wantChili,
        newValue
      },
      LOG_PATH,
      "debug",
      4
    );
  };

  const handleFillingUpdate = (
    fillingType: "mocNhi" | "thit" | "thitMocNhi"
  ) => {
    const prevFilling = { ...selectedFilling };
    updateSelectedFilling(fillingType);
    logWithLevel(
      {
        fillingType,
        previousFilling: prevFilling,
        newFilling: {
          ...selectedFilling,
          [fillingType]: !selectedFilling[fillingType]
        }
      },
      LOG_PATH,
      "debug",
      5
    );
  };

  const orderSummary = getOrderSummary();
  logWithLevel({ orderSummary }, LOG_PATH, "info", 6);

  return (
    <div className="container mx-auto px-4 py-5 space-y-5">
      <Card>
        <CardHeader>
          <CardTitle className="flex justify-between items-center">
            Canh Banh Cuon
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="space-y-4">
            <div className="space-y-3">
              {/* Bowl without vegetables */}
              <div className="flex items-center justify-between">
                <span>Canh không rau</span>
                <div className="flex items-center gap-2">
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("chili", -1)}
                    disabled={canhKhongRau === 0}
                  >
                    -
                  </Button>
                  <span className="mx-2 min-w-[2rem] text-center">
                    {canhKhongRau}
                  </span>
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("chili", 1)}
                  >
                    +
                  </Button>
                </div>
              </div>

              {/* Bowl with vegetables */}
              <div className="flex items-center justify-between">
                <span>Canh rau</span>
                <div className="flex items-center gap-2">
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("noChili", -1)}
                    disabled={canhCoRau === 0}
                  >
                    -
                  </Button>
                  <span className="mx-2 min-w-[2rem] text-center">
                    {canhCoRau}
                  </span>
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("noChili", 1)}
                  >
                    +
                  </Button>
                </div>
              </div>

              {/* Small bowl */}
              <div className="flex items-center justify-between">
                <span>Bát bé</span>
                <div className="flex items-center gap-2">
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("small", -1)}
                    disabled={smallBowl === 0}
                  >
                    -
                  </Button>
                  <span className="mx-2 min-w-[2rem] text-center">
                    {smallBowl}
                  </span>
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("small", 1)}
                  >
                    +
                  </Button>
                </div>
              </div>

              {/* Chili option */}
              <div className="flex items-center justify-between">
                <span>Có ớt</span>
                <div className="flex items-center gap-2">
                  <Button
                    size="sm"
                    variant={wantChili ? "default" : "outline"}
                    onClick={() => handleChiliUpdate(!wantChili)}
                  >
                    {wantChili ? "Selected" : "Select"}
                  </Button>
                </div>
              </div>

              {/* Nhân mọc nhĩ */}
              <div className="flex items-center justify-between">
                <span>Nhân mọc nhĩ</span>
                <div className="flex items-center gap-2">
                  <Button
                    size="sm"
                    variant={selectedFilling.mocNhi ? "default" : "outline"}
                    onClick={() => handleFillingUpdate("mocNhi")}
                  >
                    {selectedFilling.mocNhi ? "Selected" : "Select"}
                  </Button>
                </div>
              </div>

              {/* Nhân thịt */}
              <div className="flex items-center justify-between">
                <span>Nhân thịt</span>
                <div className="flex items-center gap-2">
                  <Button
                    size="sm"
                    variant={selectedFilling.thit ? "default" : "outline"}
                    onClick={() => handleFillingUpdate("thit")}
                  >
                    {selectedFilling.thit ? "Selected" : "Select"}
                  </Button>
                </div>
              </div>

              {/* Nhân thịt và mọc nhĩ */}
              <div className="flex items-center justify-between">
                <span>Nhân thịt và mọc nhĩ</span>
                <div className="flex items-center gap-2">
                  <Button
                    size="sm"
                    variant={selectedFilling.thitMocNhi ? "default" : "outline"}
                    onClick={() => handleFillingUpdate("thitMocNhi")}
                  >
                    {selectedFilling.thitMocNhi ? "Selected" : "Select"}
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <OrderDetails
        dishes={orderSummary.dishes}
        sets={orderSummary.sets}
        totalPrice={orderSummary.totalPrice}
        totalItems={orderSummary.totalItems}
      />

      <OrderCreationComponent table_token={token} table_number={number} />
    </div>
  );
}

function addTableNumberconvert(value: string): number {
  let tableNumber: number;

  try {
    if (typeof value === "string") {
      if (/^\d+$/.test(value)) {
        tableNumber = parseInt(value, 10);
      } else {
        throw new Error("Invalid input: expected a string of digits.");
      }
    } else if (typeof value === "number") {
      tableNumber = value;
    } else {
      throw new Error("Invalid input: expected a string or number.");
    }

    logWithLevel(
      {
        inputValue: value,
        convertedNumber: tableNumber
      },
      LOG_PATH,
      "debug",
      7
    );

    return tableNumber;
  } catch (error) {
    logWithLevel(
      {
        error,
        inputValue: value
      },
      LOG_PATH,
      "error",
      7
    );
    throw error;
  }
}
