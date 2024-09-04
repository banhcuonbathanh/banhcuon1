"use client";

import { ParagraphContent, Passage } from "@/types";
import React, { useState } from "react";

interface PassageComponentProps {
  passage: Passage;
}

const PassageSummary: React.FC<PassageComponentProps> = ({ passage }) => {
  const [showContent, setShowContent] = useState(false);

  const toggleContent = () => {
    setShowContent(!showContent);
  };

  return (
    <div>
      <button
        className="text-2xl font-bold text-center mb-4 mt-4 bg-blue-500 p-2 rounded-md"
        onClick={toggleContent}
      >
        Summary
      </button>
      {showContent && (
        <div className="space-y-4">
          {passage.content.map((paragraph, index) => (
            <div key={index} className="border p-4 rounded-lg">
              <p className="font-semibold mb-2">
                Paragraph Summary {String.fromCharCode(65 + index)}
              </p>
              <p>{paragraph.paragraphSummary}</p>
              <p className="font-semibold mb-2">Key Words</p>
              <p>{paragraph.keyWords}</p>
              <p className="font-semibold mb-2">Key Sentence</p>
              <p>{paragraph.keySentence}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default PassageSummary;
