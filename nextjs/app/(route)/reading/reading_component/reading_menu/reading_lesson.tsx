// const ReadingLesson = ({ lessons }: LessonProps) => {
//   if (!lessons || lessons.length === 0) {
//     return <p>No lessons available.</p>;
//   }

//   return (
//     <ul>
//       {lessons.map((lesson) => (
//         <li key={lesson}>{lesson}</li>
//       ))}
//     </ul>
//   );
// };

// export default ReadingLesson;
import React from "react";
import { AccordionContent } from "@/components/ui/accordion";
import Link from "next/link";

interface LessonProps {
  lessons: string[];
}

const Reading_Lesson = ({ lessons }: LessonProps) => {
  return (
    <div className="pl-4">
      {lessons.map((lesson, index) => (
        <AccordionContent key={lesson}>
          <Link
            href={`/lesson/${index + 1}`}
            className="flex items-center justify-start gap-2 text-gray-500 py-1 px-2 rounded-md hover:bg-lamaSkyLight"
          >
            {/* <span className="text-xs">{`Lesson ${index + 1}:`}</span> */}
            <span className="text-sm">{lesson}</span>
          </Link>
        </AccordionContent>
      ))}
    </div>
  );
};

export default Reading_Lesson;
