import { readingTest1 } from "@/data";
import PassageComponent from "./passage_component";
import SkimmingGuidance from "./SkimmingGuidance";

const ReadingTestPage: React.FC = () => {
  return (
    <div className="container mx-auto py-8">
      <h1 className="text-4xl font-bold text-center mb-8">
        Reading Test {readingTest1.testNumber}
      </h1>
      {readingTest1.sections.map((section) => (
        <div key={section.sectionNumber} className="mb-12">
          <h2 className="text-2xl font-semibold mb-4">
            Section {section.sectionNumber}
          </h2>

          {section.passages.map((passage) => (
            <PassageComponent key={passage.passageNumber} passage={passage} />
          ))}
        </div>
      ))}

      <SkimmingGuidance />
    </div>
  );
};

export default ReadingTestPage;
