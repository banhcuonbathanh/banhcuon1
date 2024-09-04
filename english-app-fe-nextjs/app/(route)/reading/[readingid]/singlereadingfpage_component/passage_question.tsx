"use client";

import { ParagraphContent, Passage } from "@/types";
import React, { useState } from "react";
import { Editor } from "reactjs-editor";

interface PassageComponentProps {
  passage: Passage;
}

const PassageQuestion: React.FC<PassageComponentProps> = ({ passage }) => {
  const [visibleAnswers, setVisibleAnswers] = useState<{
    [key: number]: boolean;
  }>({});

  const toggleContent = (questionNumber: number) => {
    setVisibleAnswers((prev) => ({
      ...prev,
      [questionNumber]: !prev[questionNumber]
    }));
  };

  return (
    <div>
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

            <button
              className="text-2xl font-bold text-center mb-4 mt-4 bg-blue-500 p-2 rounded-md"
              onClick={() => toggleContent(question.questionNumber)}
            >
              Answer
            </button>

            {visibleAnswers[question.questionNumber] && (
              <div className="space-y-4">
                <p>{question.correctAnswer}</p>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

export default PassageQuestion;
