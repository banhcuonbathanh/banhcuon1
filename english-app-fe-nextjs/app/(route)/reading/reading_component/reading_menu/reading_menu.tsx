"use client";

import React, { useState } from "react";
import Link from "next/link";
import {
  Accordion,
  AccordionItem,
  AccordionTrigger,
  AccordionContent
} from "@/components/ui/accordion";
import Reading_Lesson from "./reading_lesson";

const reading_skills = [
  {
    title: "Ki nang co ban",
    items: [
      {
        icon: "/home.png",
        reading_skill: "xac dinh noi dung chinh",
        href: "",
        visible: ["admin", "teacher", "student", "parent"],
        lessons: [
          "day 1 Ki nang co ban",
          "day 2 Ki nang co ban",
          "day 3 Ki nang co ban",
          "day 4 Ki nang co ban",
          "day 4 Ki nang co ban"
        ]
      },
      {
        icon: "/teacher.png",
        reading_skill: "tim thong tin chi tiet",
        href: "/list/teachers",
        visible: ["admin", "teacher"],
        lessons: [
          "day 1 tim thong tin chi tiet",
          "day 2 tim thong tin chi tiet",
          "day 3 tim thong tin chi tiet",
          "day 4 tim thong tin chi tiet",
          "day 4 tim thong tin chi tiet"
        ]
      },
      {
        icon: "/student.png",
        reading_skill: "sap xep xu ly thong tin",
        href: "/list/students",
        visible: ["admin", "teacher"],
        lessons: [
          "day 1 sap xep xu ly thong tin",
          "day 2 sap xep xu ly thong tin",
          "day 3",
          "day 4",
          "day 4"
        ]
      },
      {
        icon: "/calendar.png",
        reading_skill: "Events",
        href: "/list/events",
        visible: ["admin", "teacher", "student", "parent"],
        lessons: ["day 1", "day 2", "day 3", "day 4", "day 4"]
      },
      {
        icon: "/announcement.png",
        reading_skill: "Announcements",
        href: "/tien do",
        visible: ["admin", "teacher", "student", "parent"],
        lessons: ["day 1", "day 2", "day 3", "day 4", "day 4"]
      }
    ]
  },
  {
    title: "OTHER",
    items: [
      {
        icon: "/profile.png",
        reading_skill: "Profile",
        href: "/profile",
        visible: ["admin", "teacher", "student", "parent"],
        lessons: ["day 1 Profile", "day 2 Profile", "day 3", "day 4", "day 4"]
      },
      {
        icon: "/setting.png",
        reading_skill: "Settings",
        href: "/settings",
        visible: ["admin", "teacher", "student", "parent"],
        lessons: ["day 1", "day 2", "day 3", "day 4", "day 4"]
      },
      {
        icon: "/logout.png",
        reading_skill: "Logout",
        href: "/logout",
        visible: ["admin", "teacher", "student", "parent"],
        lessons: ["day 1 settings", "day 2 settings", "day 3", "day 4", "day 4"]
      }
    ]
  }
];

const DashBoardMenu = () => {
  const [expandedSkill, setExpandedSkill] = useState<string | null>(null);

  return (
    <div className="mt-4 text-sm">
      <Accordion type="multiple">
        {reading_skills.map((section, sectionIndex) => (
          <AccordionItem value={`item-${sectionIndex}`} key={section.title}>
            <AccordionTrigger className="text-gray-500 transition-colors hover:text-primary no-underline">
              {section.title}
            </AccordionTrigger>
            {section.items.map((item) => (
              <AccordionContent key={item.reading_skill}>
                <div
                  className="cursor-pointer"
                  onClick={() =>
                    setExpandedSkill(
                      expandedSkill === item.reading_skill
                        ? null
                        : item.reading_skill
                    )
                  }
                >
                  <Link
                    href={item.href}
                    className="flex items-center justify-center lg:justify-start gap-4 text-gray-500 py-2 md:px-2 rounded-md hover:bg-lamaSkyLight"
                    onClick={(e) => e.preventDefault()} // Prevent navigation
                  >
                    {/* <Image src={item.icon} alt="" width={20} height={20} /> */}
                    <span className="hidden lg:block">
                      {item.reading_skill}
                    </span>
                  </Link>
                </div>
                {expandedSkill === item.reading_skill && (
                  <Reading_Lesson lessons={item.lessons} />
                )}
              </AccordionContent>
            ))}
          </AccordionItem>
        ))}
      </Accordion>
    </div>
  );
};

export default DashBoardMenu;
