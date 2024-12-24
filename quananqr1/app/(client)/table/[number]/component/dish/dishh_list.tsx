"use client";

import { DishInterface } from "@/schemaValidations/interface/type_dish";

import { DishCard } from "./disih_tem";
import React from "react";
import GridContainer from "@/components/general-container-dish";

interface DishSelectionProps {
  dishes: DishInterface[];
}

export function DishSelection({ dishes }: DishSelectionProps) {
  const [isMounted, setIsMounted] = React.useState(false);

  React.useEffect(() => {
    setIsMounted(true);
  }, []);

  if (!isMounted) {
    return null; // or a loading skeleton
  }

  return (
    <GridContainer>
      {dishes.map((dish: DishInterface) => (
        <DishCard key={dish.id} dish={dish} />
      ))}
    </GridContainer>
  );
}
