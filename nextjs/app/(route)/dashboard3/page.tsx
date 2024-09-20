import React from "react";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup
} from "@/components/ui/resizable";
import { Dashboard2SideBar } from "./component/sidebar";

const Dashboard3 = () => {
  return (
    <div className="mt-40">
      <ResizablePanelGroup direction="horizontal">
        <ResizablePanel>
          <p>tsts</p>
        </ResizablePanel>
        <ResizableHandle />
        <ResizablePanel>Two</ResizablePanel>

        <ResizableHandle />

        <ResizablePanel>three</ResizablePanel>
      </ResizablePanelGroup>
    </div>
  );
};

export default Dashboard3;
