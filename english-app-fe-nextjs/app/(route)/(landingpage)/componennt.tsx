"use client";

import { ParagraphContent, Passage } from "@/types";
import React, { useState } from "react";
import { Editor } from "reactjs-editor";

interface PassageComponentProps {
  passage: Passage;
}

const PassageComponent: React.FC<PassageComponentProps> = ({ passage }) => {
  const [showDetails, setShowDetails] = useState(false);

  const toggleDetails = () => {
    setShowDetails(!showDetails);
  };

  return (
    <div>
      <div className="space-y-6 mb-8">
        <button
          onClick={toggleDetails}
          className="mb-4 px-4 py-2 bg-blue-500 text-white rounded"
        >
          {showDetails ? "Hide" : "Show"} All Details
        </button>
        <Editor
          htmlContent={`
            <main class="bookContainer">
              <aside>
                <h1 class="text-3xl font-bold text-center mb-6">${passage.title}</h1>
              </aside>
            </main>
          `}
        />
        {passage.content.map((paragraph: ParagraphContent, index: number) => {
          const [key, content] = Object.entries(paragraph)[0];
          return (
            <div
              key={index}
              className="mb-4 p-4 border border-gray-200 rounded-lg"
            >
              <h3 className="text-lg font-bold mb-2">Paragraph {key}</h3>
              <div
                dangerouslySetInnerHTML={{
                  __html: `<${key}>${content}</${key}>`
                }}
              />
              {showDetails && (
                <div>
                  <p>
                    <strong>Summary:</strong> {paragraph.paragraphSummary}
                  </p>
                  <p>
                    <strong>Key Words:</strong> {paragraph.keyWords}
                  </p>
                  <p>
                    <strong>Key Sentence:</strong> {paragraph.keySentence}
                  </p>
                </div>
              )}
            </div>
          );
        })}
      </div>
      {/* Questions section remains unchanged */}
    </div>
  );
};

export default PassageComponent;
