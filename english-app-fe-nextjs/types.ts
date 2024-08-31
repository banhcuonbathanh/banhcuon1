export interface ReadingTest {
    testNumber: number;
    sections: Section[];
  }
  
  export interface Section {
    sectionNumber: number;
    timeAllowed: number;
    passages: Passage[];
  }
  
  export interface Passage {
    passageNumber: number;
    title: string;
    content: ParagraphContent[];
    questions: Question[];
  }
  export interface ParagraphContent {
    [key: string]: string; 
    paragraphSummary: string;
    keyWords: string;
    keySentence: string;

  }
  
  export interface Question {
    questionNumber: number;
    type: QuestionType;
    content: string;
    options?: string[];
    correctAnswer?: string | string[];
  }
  
  export enum QuestionType {
    MultipleChoice = 'MultipleChoice',
    TrueFalseNotGiven = 'TrueFalseNotGiven',
    Matching = 'Matching',
    ShortAnswer = 'ShortAnswer'
  }
  export interface DialogContentModel {
    title: string;
    description: string;
    body: React.ReactNode;
  }
  