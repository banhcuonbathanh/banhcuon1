"use client";

import { ParagraphContent, Passage } from "@/types";
import React from "react";
import { Editor } from "reactjs-editor";

interface PassageComponentProps {
  passage: Passage;
}

const PassageContend: React.FC<PassageComponentProps> = ({ passage }) => {
  const contentHtml = passage.content
    .map((paragraph: ParagraphContent, index: number) => {
      const [key, content] = Object.entries(paragraph)[0];
      return `
        <div class="mb-4 p-4 border border-gray-200 rounded-lg">
          <h3 class="text-lg font-bold mb-2">Paragraph ${key}</h3>
          <${key}>${content}</${key}>

        </div>
      `;
    })
    .join("");

  return (
    <div>
      <div className="space-y-6 mb-8">
        <Editor
          htmlContent={`
            <main class="bookContainer">
              <aside>
                <h1 class="text-3xl font-bold text-center mb-6">${passage.title}</h1>
                ${contentHtml}
              </aside>
            </main>
          `}
        />
      </div>
    </div>
  );
};

export default PassageContend;
