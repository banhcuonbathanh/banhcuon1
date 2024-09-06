import React from "react";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger
} from "@/components/ui/accordion";
const menuItems = [
    {
      title: "Ki nang co ban",
      items: [
        {
          icon: "/home.png",
          label: "xac dinh noi dung chinh",
          href: "/",
          visible: ["admin", "teacher", "student", "parent"]
        },
        {
          icon: "/teacher.png",
          label: "tim thong tin chi tiet",
          href: "/list/teachers",
          visible: ["admin", "teacher"]
        },
        {
          icon: "/student.png",
          label: "sap xep xu ly thong tin",
          href: "/list/students",
          visible: ["admin", "teacher"]
        },
  
        {
          icon: "/calendar.png",
          label: "Events",
          href: "/list/events",
          visible: ["admin", "teacher", "student", "parent"]
        },
  
        {
          icon: "/announcement.png",
          label: "Announcements",
          href: "/tien do",
          visible: ["admin", "teacher", "student", "parent"]
        }
      ]
    },
    {
      title: "OTHER",
      items: [
        {
          icon: "/profile.png",
          label: "Profile",
          href: "/profile",
          visible: ["admin", "teacher", "student", "parent"]
        },
        {
          icon: "/setting.png",
          label: "Settings",
          href: "/settings",
          visible: ["admin", "teacher", "student", "parent"]
        },
        {
          icon: "/logout.png",
          label: "Logout",
          href: "/logout",
          visible: ["admin", "teacher", "student", "parent"]
        }
      ]
    }
  ];
const Reading_Tittle = () => {
  return (
    <Accordion type="single" collapsible>
      {/* item 1 */}
      <AccordionItem value="item-1">
        <AccordionTrigger className="text-gray-500 transition-colors hover:text-primary no-underline">
          Kĩ năng cơ bản
        </AccordionTrigger>
        <AccordionContent className="px-4 border-l border-gray-500 text-gray-500 hover:text-primary">
          xac dinh noi dung chinh
        </AccordionContent>
        <AccordionContent className="px-4 border-l border-gray-500 text-gray-500 hover:text-primary">
          tim thong tin chi tiet
        </AccordionContent>
      </AccordionItem>

      {/* item 2 */}

      <AccordionItem value="item-2">
        <AccordionTrigger className="text-gray-500 transition-colors hover:text-primary no-underline">
          Kĩ năng cơ bản
        </AccordionTrigger>
        <AccordionContent className="px-4 border-l border-gray-500 text-gray-500 hover:text-primary">
          xac dinh noi dung chinh
        </AccordionContent>
        <AccordionContent className="px-4 border-l border-gray-500 text-gray-500 hover:text-primary">
          tim thong tin chi tiet
        </AccordionContent>
      </AccordionItem>
    </Accordion>
  );
};

export default Reading_Tittle;
