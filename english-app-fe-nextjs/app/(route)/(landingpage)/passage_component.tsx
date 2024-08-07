"use client";

import { ParagraphContent, Passage } from "@/types";
import React from "react";
import { Editor } from "reactjs-editor";

interface PassageComponentProps {
  passage: Passage;
}

const PassageComponent: React.FC<PassageComponentProps> = ({ passage }) => {
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
      <h3 className="text-2xl font-bold text-center mb-4">Questions</h3>
      <div className="space-y-4">
        {passage.questions.map((question) => (
          <div key={question.questionNumber} className="border p-4 rounded-lg">
            <p className="font-semibold mb-2">
              Question {question.questionNumber}
            </p>
            <p>{question.content}</p>
            {question.options && (
              <ul className="list-disc list-inside mt-2">
                {question.options.map((option, index) => (
                  <li key={index}>{option}</li>
                ))}
              </ul>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

export default PassageComponent;
