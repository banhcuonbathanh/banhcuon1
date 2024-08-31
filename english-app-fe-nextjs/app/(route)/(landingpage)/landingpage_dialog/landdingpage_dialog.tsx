"use client";

import React from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle
} from "@/components/ui/dialog";

import { useDialogStorePersist } from "../landing_page_zustand/landding_page_zustand";

// Define the interface for the dialog content

// Define the props interface for the ExampleDialog component

const ExampleDialog: React.FC = () => {
  const { isOpen, dialogContentModel, closeDialog } = useDialogStorePersist();

  if (!dialogContentModel) return null;

  return (
    <Dialog open={isOpen} onOpenChange={closeDialog}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{dialogContentModel.title}</DialogTitle>
          <DialogDescription>
            {dialogContentModel.description}
          </DialogDescription>
        </DialogHeader>
        <div className="py-4">{dialogContentModel.body}</div>
      </DialogContent>
    </Dialog>
  );
};

export default ExampleDialog;
