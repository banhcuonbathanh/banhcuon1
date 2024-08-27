package types

import "time"

// ReadingRequest represents the structure for requesting a reading test
type ReadingReqModel struct {
    ReadingReqTestType ReadingTestModel `json:"reading_req_type"`
    	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ReadingResponse represents the structure for responding with a reading test
type ReadingResModel struct {
    ID       int64      `json:"id"`
    ReadingResType ReadingTestModel `json:"reading_res_type"`
    CreatedAt time.Time `json:"created_at"`
UpdatedAt time.Time `json:"updated_at"`
}

// ReadingTest structure (as provided)
type ReadingTestModel struct {
    TestNumber int       `json:"testNumber"`
    Sections   []SectionModel `json:"sections"`
}
type ReadingResList struct {
    Readings   []*ReadingResModel `json:"readings"`
    TotalCount int         `json:"total_count"`
}
// Section structure (as provided)
type SectionModel struct {
    SectionNumber int       `json:"sectionNumber"`
    TimeAllowed   int       `json:"timeAllowed"`
    Passages      []PassageModel `json:"passages"`
}

// Passage structure (as provided)
type PassageModel struct {
    PassageNumber int                `json:"passageNumber"`
    Title         string             `json:"title"`
    Content       []ParagraphContentModel `json:"content"`
    Questions     []QuestionModel        `json:"questions"`
}

// ParagraphContent structure (as provided)
type ParagraphContentModel struct {
    ParagraphSummary string `json:"paragraphSummary"`
    KeyWords         string `json:"keyWords"`
    KeySentence      string `json:"keySentence"`
}

// Question structure (as provided)
type QuestionModel struct {
    QuestionNumber int          `json:"questionNumber"`
    Type           QuestionTypeModel `json:"type"`
    Content        string       `json:"content"`
    Options        []string     `json:"options,omitempty"`
    CorrectAnswer  interface{}  `json:"correctAnswer,omitempty"`
}

// QuestionType and constants (as provided)
type QuestionTypeModel string

const (
    MultipleChoice    QuestionTypeModel = "MultipleChoice"
    TrueFalseNotGiven QuestionTypeModel = "TrueFalseNotGiven"
    Matching          QuestionTypeModel = "Matching"
    ShortAnswer       QuestionTypeModel = "ShortAnswer"
)

type PageRequestModel struct {
    PageNumber int          `json:"page_number"`
    PageSize int          `json:"page_size"`

}