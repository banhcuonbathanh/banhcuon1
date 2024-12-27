"use client";

import React, { useEffect, useState } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import useCartStore from "@/zusstand/new-order/new-order-zustand";
import { logWithLevel } from "@/lib/logger/log";

const LOG_PATH =
  "quananqr1/app/(client)/table/[number]/component/order/order.tsx";

interface FillingState {
  mocNhi: boolean;
  thit: boolean;
  thitMocNhi: boolean;
}

export default function ChoosingTopping() {
  // Local state management
  const [canhKhongRau, setCanhKhongRau] = useState(0);
  const [canhCoRau, setCanhCoRau] = useState(0);
  const [smallBowl, setSmallBowl] = useState(0);
  const [wantChili, setWantChili] = useState(false);
  const [selectedFilling, setSelectedFilling] = useState<FillingState>({
    mocNhi: false,
    thit: false,
    thitMocNhi: false
  });

  const { current_order, updateTopping } = useCartStore();

  // Function to generate and update topping string
  const updateOrderTopping = () => {
    const toppingParts = [];

    // Add bowl quantities
    if (canhKhongRau > 0) toppingParts.push(`Canh không rau: ${canhKhongRau}`);
    if (canhCoRau > 0) toppingParts.push(`Canh rau: ${canhCoRau}`);
    if (smallBowl > 0) toppingParts.push(`Bát bé: ${smallBowl}`);

    // Add chili preference
    if (wantChili) toppingParts.push("Có ớt");

    // Add filling type
    if (selectedFilling.mocNhi) toppingParts.push("Nhân mọc nhĩ");
    if (selectedFilling.thit) toppingParts.push("Nhân thịt");
    if (selectedFilling.thitMocNhi) toppingParts.push("Nhân thịt và mọc nhĩ");

    const toppingString = toppingParts.join(", ");
    if (current_order) {
      updateTopping(toppingString);
    }

    logWithLevel(
      {
        type: "updateOrderTopping",
        toppingString
      },
      LOG_PATH,
      "debug",
      8
    );
  };

  // Update order topping whenever any topping-related state changes
  useEffect(() => {
    updateOrderTopping();
  }, [canhKhongRau, canhCoRau, smallBowl, wantChili, selectedFilling]);

  const handleBowlChange = (
    type: "chili" | "noChili" | "small",
    change: number
  ) => {
    switch (type) {
      case "chili":
        const newToppingValue = Math.max(0, canhKhongRau + change);
        setCanhKhongRau(newToppingValue);
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
        break;
      case "noChili":
        const newNoChiliValue = Math.max(0, canhCoRau + change);
        setCanhCoRau(newNoChiliValue);
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
        break;
      case "small":
        const newSmallBowlValue = Math.max(0, smallBowl + change);
        setSmallBowl(newSmallBowlValue);
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
        break;
    }
  };

  const handleChiliUpdate = (newValue: boolean) => {
    setWantChili(newValue);
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
    const newFilling = {
      mocNhi: false,
      thit: false,
      thitMocNhi: false,
      [fillingType]: true
    };
    setSelectedFilling(newFilling);
    logWithLevel(
      {
        fillingType,
        previousFilling: prevFilling,
        newFilling
      },
      LOG_PATH,
      "debug",
      5
    );
  };

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

              {/* Filling options */}
              {[
                { key: "mocNhi", label: "Nhân mọc nhĩ" },
                { key: "thit", label: "Nhân thịt" },
                { key: "thitMocNhi", label: "Nhân thịt và mọc nhĩ" }
              ].map(({ key, label }) => (
                <div key={key} className="flex items-center justify-between">
                  <span>{label}</span>
                  <div className="flex items-center gap-2">
                    <Button
                      size="sm"
                      variant={
                        selectedFilling[key as keyof typeof selectedFilling]
                          ? "default"
                          : "outline"
                      }
                      onClick={() =>
                        handleFillingUpdate(
                          key as "mocNhi" | "thit" | "thitMocNhi"
                        )
                      }
                    >
                      {selectedFilling[key as keyof typeof selectedFilling]
                        ? "Selected"
                        : "Select"}
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
