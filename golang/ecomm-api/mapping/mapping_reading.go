package mapping_user


import (
 	"english-ai-full/ecomm-api/types"
	pb "english-ai-full/ecomm-grpc/proto/reading"

)

func ToPBReadingReq(reading types.ReadingReqModel) *pb.ReadingReq {
    return &pb.ReadingReq{
        ReadingTest: toProtoReadingTest(reading.ReadingReqTestType),
    }
}

func ToReadingResFromPbReadingRes(createdReading *pb.ReadingRes) types.ReadingResModel {
    return types.ReadingResModel{
        ID:             createdReading.Id,
        ReadingResType: fromProtoReadingTest(createdReading.ReadingTest),
        CreatedAt:      createdReading.CreatedAt.AsTime(),
        UpdatedAt:      createdReading.UpdatedAt.AsTime(),
    }
}

func ToReadingResListFromPbReadingResList(readings []*pb.ReadingRes) []types.ReadingResModel {
    result := make([]types.ReadingResModel, len(readings))
    for i, r := range readings {
        result[i] = ToReadingResFromPbReadingRes(r)
    }
    return result
}

func toProtoReadingTest(rt types.ReadingTestModel) *pb.ReadingTestProto {
    return &pb.ReadingTestProto{
        TestNumber: int32(rt.TestNumber),
        Sections:   toProtoSections(rt.Sections),
    }
}

func fromProtoReadingTest(rt *pb.ReadingTestProto) types.ReadingTestModel {
    return types.ReadingTestModel{
        TestNumber: int(rt.TestNumber),
        Sections:   fromProtoSections(rt.Sections),
    }
}

func toProtoSections(sections []types.SectionModel) []*pb.SectionProto {
    result := make([]*pb.SectionProto, len(sections))
    for i, s := range sections {
        result[i] = &pb.SectionProto{
            SectionNumber: int32(s.SectionNumber),
            TimeAllowed:   int32(s.TimeAllowed),
            Passages:      toProtoPassages(s.Passages),
        }
    }
    return result
}

func fromProtoSections(sections []*pb.SectionProto) []types.SectionModel {
    result := make([]types.SectionModel, len(sections))
    for i, s := range sections {
        result[i] = types.SectionModel{
            SectionNumber: int(s.SectionNumber),
            TimeAllowed:   int(s.TimeAllowed),
            Passages:      fromProtoPassages(s.Passages),
        }
    }
    return result
}

func toProtoPassages(passages []types.PassageModel) []*pb.PassageProto {
    result := make([]*pb.PassageProto, len(passages))
    for i, p := range passages {
        result[i] = &pb.PassageProto{
            PassageNumber: int32(p.PassageNumber),
            Title:         p.Title,
            Content:       toProtoParagraphContents(p.Content),
            Questions:     toProtoQuestions(p.Questions),
        }
    }
    return result
}

func fromProtoPassages(passages []*pb.PassageProto) []types.PassageModel {
    result := make([]types.PassageModel, len(passages))
    for i, p := range passages {
        result[i] = types.PassageModel{
            PassageNumber: int(p.PassageNumber),
            Title:         p.Title,
            Content:       fromProtoParagraphContents(p.Content),
            Questions:     fromProtoQuestions(p.Questions),
        }
    }
    return result
}

func toProtoParagraphContents(contents []types.ParagraphContentModel) []*pb.ParagraphContentProto {
    result := make([]*pb.ParagraphContentProto, len(contents))
    for i, c := range contents {
        result[i] = &pb.ParagraphContentProto{
            ParagraphSummary: c.ParagraphSummary,
            KeyWords:         c.KeyWords,
            KeySentence:      c.KeySentence,
        }
    }
    return result
}

func fromProtoParagraphContents(contents []*pb.ParagraphContentProto) []types.ParagraphContentModel {
    result := make([]types.ParagraphContentModel, len(contents))
    for i, c := range contents {
        result[i] = types.ParagraphContentModel{
            ParagraphSummary: c.ParagraphSummary,
            KeyWords:         c.KeyWords,
            KeySentence:      c.KeySentence,
        }
    }
    return result
}

func toProtoQuestions(questions []types.QuestionModel) []*pb.QuestionProto {
    result := make([]*pb.QuestionProto, len(questions))
    for i, q := range questions {
        result[i] = &pb.QuestionProto{
            QuestionNumber: int32(q.QuestionNumber),
            Type:           toProtoQuestionType(q.Type),
            Content:        q.Content,
            Options:        q.Options,
        }
        switch answer := q.CorrectAnswer.(type) {
        case string:
            result[i].CorrectAnswer = &pb.QuestionProto_StringAnswer{StringAnswer: answer}
        case []string:
            result[i].CorrectAnswer = &pb.QuestionProto_StringArrayAnswer{
                StringArrayAnswer: &pb.StringArrayProto{Values: answer},
            }
        }
    }
    return result
}

func fromProtoQuestions(questions []*pb.QuestionProto) []types.QuestionModel {
    result := make([]types.QuestionModel, len(questions))
    for i, q := range questions {
        result[i] = types.QuestionModel{
            QuestionNumber: int(q.QuestionNumber),
            Type:           fromProtoQuestionType(q.Type),
            Content:        q.Content,
            Options:        q.Options,
        }
        switch answer := q.CorrectAnswer.(type) {
        case *pb.QuestionProto_StringAnswer:
            result[i].CorrectAnswer = answer.StringAnswer
        case *pb.QuestionProto_StringArrayAnswer:
            result[i].CorrectAnswer = answer.StringArrayAnswer.Values
        }
    }
    return result
}

func toProtoQuestionType(qt types.QuestionTypeModel) pb.QuestionTypeProto {
    switch qt {
    case types.MultipleChoice:
        return pb.QuestionTypeProto_MULTIPLE_CHOICE
    case types.TrueFalseNotGiven:
        return pb.QuestionTypeProto_TRUE_FALSE_NOT_GIVEN
    case types.Matching:
        return pb.QuestionTypeProto_MATCHING
    case types.ShortAnswer:
        return pb.QuestionTypeProto_SHORT_ANSWER
    default:
        return pb.QuestionTypeProto_MULTIPLE_CHOICE // Default case
    }
}

func fromProtoQuestionType(qt pb.QuestionTypeProto) types.QuestionTypeModel {
    switch qt {
    case pb.QuestionTypeProto_MULTIPLE_CHOICE:
        return types.MultipleChoice
    case pb.QuestionTypeProto_TRUE_FALSE_NOT_GIVEN:
        return types.TrueFalseNotGiven
    case pb.QuestionTypeProto_MATCHING:
        return types.Matching
    case pb.QuestionTypeProto_SHORT_ANSWER:
        return types.ShortAnswer
    default:
        return types.MultipleChoice // Default case
    }
}