import { DialogContentModel } from "@/types";
import { create } from "zustand";
import { persist } from "zustand/middleware";

type DialogStore = {
  isOpen: boolean;
  dialogContentModel: DialogContentModel | null;
  openDialog: (dialogContentModel: DialogContentModel) => void;
  closeDialog: () => void;
  reset: (callback: Function) => void;
};

export const useDialogStorePersist = create<DialogStore>()(
  persist<DialogStore>(
    (set) => ({
      isOpen: true,
      dialogContentModel: null,
      openDialog: (dialogContentModel) =>
        set({ isOpen: true, dialogContentModel }),
      closeDialog: () => set({ isOpen: false, dialogContentModel: null }),
      reset: (callback: Function) => {
        set({ isOpen: false, dialogContentModel: null });
        callback();
      }
    }),
    {
      name: "dialog-store",
      skipHydration: true
    }
  )
);
