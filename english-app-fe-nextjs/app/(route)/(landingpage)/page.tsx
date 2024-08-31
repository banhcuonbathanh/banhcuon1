"use client";

import { readingTest1 } from "@/data";
import PassageComponent from "./presentation_component/passage_component";
import PassageSummary from "./presentation_component/passage_summary";
import PassageQuestion from "./presentation_component/passage_question";
import PassageContend from "./presentation_component/passage_component";

import { useState } from "react";
import Timer from "./presentation_component/timer";
import ExampleDialog from "./landingpage_dialog/landdingpage_dialog";
import { useDialogStorePersist } from "./landing_page_zustand/landding_page_zustand";

const ReadingTestPage: React.FC = () => {
  const [timerMinutes, setTimerMinutes] = useState(0);
  const [timerSeconds, setTimerSeconds] = useState(0);

  const handleTimerEnd = () => {
    // Add any logic to handle the timer ending here
    console.log("Timer ended");
  };
  const openDialog = useDialogStorePersist((state) => state.openDialog);

  const handleOpenDialog = () => {
    openDialog({
      title: "Example Dialog",
      description: "This is an example dialog using our custom store",
      body: (
        <p>This is the body of the dialog. You can put any React node here.</p>
      )
    });
  };
  return (
    <div className="container mx-auto py-8">
      <button onClick={handleOpenDialog}>Open Dialog</button>
      {/* <ExampleDialog /> */}
      <h1 className="text-4xl font-bold text-center mb-8">
        Reading Test {readingTest1.testNumber}
      </h1>
      <div className="mb-8">
        <label htmlFor="timer-minutes" className="mr-2">
          Set timer minutes:
        </label>
        <input
          type="number"
          id="timer-minutes"
          min="0"
          className="border border-gray-300 px-2 py-1 rounded-md bg-black"
          onChange={(e) => setTimerMinutes(parseInt(e.target.value))}
        />
        <label htmlFor="timer-seconds" className="ml-4 mr-2">
          Set timer seconds:
        </label>
        <input
          type="number"
          id="timer-seconds"
          min="0"
          max="59"
          className="border border-gray-300 px-2 py-1 rounded-md bg-black"
          onChange={(e) => setTimerSeconds(parseInt(e.target.value))}
        />
      </div>
      <Timer
        initialMinutes={timerMinutes}
        initialSeconds={timerSeconds}
        onTimerEnd={handleTimerEnd}
      />

      {readingTest1.sections.map((section) => (
        <div key={section.sectionNumber} className="mb-12">
          <h2 className="text-2xl font-semibold mb-4">
            Section {section.sectionNumber}
          </h2>

          {section.passages.map((passage) => (
            <PassageContend key={passage.passageNumber} passage={passage} />
          ))}
          {section.passages.map((passage) => (
            <PassageQuestion key={passage.passageNumber} passage={passage} />
          ))}
          {section.passages.map((passage) => (
            <PassageSummary key={passage.passageNumber} passage={passage} />
          ))}
        </div>
      ))}
      <div></div>
    </div>
  );
};

export default ReadingTestPage;
