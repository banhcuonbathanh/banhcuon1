package service

import (
	"context"
	proto "english-ai-full/ecomm-grpc/proto/reading"
	repository "english-ai-full/ecomm-grpc/repository/reading_repository"
	"fmt"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"english-ai-full/ecomm-api/types"
)
type ReadingServericeStruct struct {
    readingRepo *repository.ReadingRepository
    proto.UnimplementedEcommReadingServer
}

func NewReadingServer(readingRepo *repository.ReadingRepository) *ReadingServericeStruct {
    return &ReadingServericeStruct{
        readingRepo: readingRepo,
    }
}


func (rs *ReadingServericeStruct) CreateReading(ctx context.Context, req *proto.ReadingReq) (*proto.ReadingRes, error) {
    log.Println("Creating Reading:",
        "ID:", req.Id,
        "Reading Test:", req.ReadingTest,
    )

    newReading := &types.ReadingReqModel{
        ReadingReqTestType: convertProtoReadingReqTestTypeToModel(req),
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
    }

    createdReading, err := rs.readingRepo.CreateReading(ctx, newReading)
    if err != nil {
        log.Println("Error creating reading:", err)
        return nil, err
    }

    log.Println("Reading created successfully. ID:", createdReading.ID)
    return convertModelReadingToProtoReadingRes(createdReading), nil
}


// func (rs *ReadingServericeStruct) SaveReading(ctx context.Context, res *proto.ReadingRes) (*proto.ReadingRes, error) {
//     reading := convertProtoReadingResToModel(res)
    
  

//     savedReading, err := rs.readingRepo.SaveReading(ctx, reading)
//     if err != nil {
//         return nil, fmt.Errorf("error saving reading: %w", err)
//     }
    
//     return convertModelToProtoReadingRes(savedReading), nil
// }
func (rs *ReadingServericeStruct) SaveReading(ctx context.Context, res *proto.ReadingRes) (*proto.ReadingRes, error) {
    reading := convertProtoReadingResToModel(res)

    savedReading, err := rs.readingRepo.SaveReading(ctx, reading)
    if err != nil {
        return nil, fmt.Errorf("error saving reading: %w", err)
    }
    
    return convertModelToProtoReadingRes(savedReading), nil
}

func (rs *ReadingServericeStruct) UpdateReading(ctx context.Context, req *proto.ReadingReq) (*proto.ReadingRes, error) {
    // Convert the incoming proto request to a ReadingResModel
    updatedReading := &types.ReadingResModel{
        ID: req.Id, // Assuming req.Id exists in proto.ReadingReq
        ReadingResType: convertProtoReadingReqToModel(req).ReadingReqTestType,
        UpdatedAt: time.Now(),
    }
    
    // Call the repository method with the correct arguments
    result, err := rs.readingRepo.UpdateReading(ctx, updatedReading)
    if err != nil {
        return nil, err
    }
    
    // Convert the result back to proto response
    return convertModelReadingToProtoReadingRes(result), nil
}

func (rs *ReadingServericeStruct) DeleteReading(ctx context.Context, res *proto.ReadingRes) (*emptypb.Empty, error) {
    reading := convertProtoReadingResToModel(res)
    err := rs.readingRepo.DeleteReading(ctx, reading)
    return &emptypb.Empty{}, err
}

func (rs *ReadingServericeStruct) FindAllReading(ctx context.Context, _ *emptypb.Empty) (*proto.ReadingResList, error) {
    readingsList, err := rs.readingRepo.FindAllReading(ctx)
    if err != nil {
        return nil, err
    }
    
    // Convert *types.ReadingResList to []types.ReadingResModel
    readings := readingsList.Readings
    
    return convertModelReadingsToProtoReadingResList(readings), nil
}

func (rs *ReadingServericeStruct) FindByID(ctx context.Context, req *proto.ReadingRes) (*proto.ReadingRes, error) {
    reading, err := rs.readingRepo.FindByID(ctx,req.Id)
    if err != nil {
        return nil, err
    }
    return convertModelReadingToProtoReadingRes(reading), nil
}

func (rs *ReadingServericeStruct) FindReadingByPage(ctx context.Context, req *proto.PageRequest) (*proto.ReadingResList, error) {
    result, err := rs.readingRepo.FindReadingByPage(ctx, &types.PageRequestModel{
        PageNumber: int(req.PageNumber),
        PageSize:   int(req.PageSize),
    })
    if err != nil {
        return nil, err
    }
    return &proto.ReadingResList{
        Readings:   convertModelReadingsResListToProtoReadings(result.Readings),
        TotalCount: int32(result.TotalCount),
    }, nil
}

// Helper functions for converting between proto and model types



func convertModelReadingToProtoReadingRes(modelReading *types.ReadingResModel) *proto.ReadingRes {
    return &proto.ReadingRes{
        Id:          modelReading.ID,
        ReadingTest: convertModelReadingTestToProto(modelReading.ReadingResType),
        CreatedAt:   timestamppb.New(modelReading.CreatedAt),
        UpdatedAt:   timestamppb.New(modelReading.UpdatedAt),
    }
}



func convertModelReadingsToProtoReadingResList(readings []*types.ReadingResModel) *proto.ReadingResList {
    protoReadings := make([]*proto.ReadingRes, len(readings))
    for i, reading := range readings {
        protoReadings[i] = convertModelReadingToProtoReadingRes(reading)
    }
    return &proto.ReadingResList{
        Readings:   protoReadings,
        TotalCount: int32(len(readings)),
    }
}

func convertModelReadingsResListToProtoReadings(readings []*types.ReadingResModel) []*proto.ReadingRes {
    protoReadings := make([]*proto.ReadingRes, len(readings))
    for i, reading := range readings {
        protoReadings[i] = convertModelReadingToProtoReadingRes(reading)
    }
    return protoReadings
}

func convertModelReadingTestToProto(modelTest types.ReadingTestModel) *proto.ReadingTestProto {
    protoTest := &proto.ReadingTestProto{
        TestNumber: int32(modelTest.TestNumber),
        Sections:   make([]*proto.SectionProto, len(modelTest.Sections)),
    }

    for i, section := range modelTest.Sections {
        protoSection := &proto.SectionProto{
            SectionNumber: int32(section.SectionNumber),
            TimeAllowed:   int32(section.TimeAllowed),
            Passages:      make([]*proto.PassageProto, len(section.Passages)),
        }

        for j, passage := range section.Passages {
            protoPassage := &proto.PassageProto{
                PassageNumber: int32(passage.PassageNumber),
                Title:         passage.Title,
                Content:       make([]*proto.ParagraphContentProto, len(passage.Content)),
                Questions:     make([]*proto.QuestionProto, len(passage.Questions)),
            }

            for k, content := range passage.Content {
                protoPassage.Content[k] = &proto.ParagraphContentProto{
                    ParagraphSummary: content.ParagraphSummary,
                    KeyWords:         content.KeyWords,
                    KeySentence:      content.KeySentence,
                }
            }

            for k, question := range passage.Questions {
                protoQuestion := &proto.QuestionProto{
                    QuestionNumber: int32(question.QuestionNumber),
                    Type:           convertQuestionTypeToProto(question.Type),
                    Content:        question.Content,
                    Options:        question.Options,
                }

                switch answer := question.CorrectAnswer.(type) {
                case string:
                    protoQuestion.CorrectAnswer = &proto.QuestionProto_StringAnswer{StringAnswer: answer}
                case []string:
                    protoQuestion.CorrectAnswer = &proto.QuestionProto_StringArrayAnswer{
                        StringArrayAnswer: &proto.StringArrayProto{Values: answer},
                    }
                }

                protoPassage.Questions[k] = protoQuestion
            }

            protoSection.Passages[j] = protoPassage
        }

        protoTest.Sections[i] = protoSection
    }

    return protoTest
}

func convertQuestionTypeToProto(questionType types.QuestionTypeModel) proto.QuestionTypeProto {
    switch questionType {
    case types.MultipleChoice:
        return proto.QuestionTypeProto_MULTIPLE_CHOICE
    case types.TrueFalseNotGiven:
        return proto.QuestionTypeProto_TRUE_FALSE_NOT_GIVEN
    case types.Matching:
        return proto.QuestionTypeProto_MATCHING
    case types.ShortAnswer:
        return proto.QuestionTypeProto_SHORT_ANSWER
    default:
        // You might want to handle this case differently, perhaps with an error or a default value
        return proto.QuestionTypeProto_MULTIPLE_CHOICE
    }
}

// func convertProtoReadingReqToModel(protoReading *proto.ReadingReq) *types.ReadingReqModel {
//     if protoReading == nil {
//         return nil
//     }

//     modelReading := &types.ReadingReqModel{
//         ID: protoReading.Id, // Assuming the proto has an Id field
//         ReadingResType: types.ReadingTestModel{
//             TestNumber: int(protoReading.ReadingTest.TestNumber),
//             Sections:   make([]types.SectionModel, len(protoReading.ReadingTest.Sections)),
//         },
//         CreatedAt: time.Now(), // Or use a field from protoReading if available
//         UpdatedAt: time.Now(), // Or use a field from protoReading if available
//     }

//     for i, protoSection := range protoReading.ReadingTest.Sections {
//         modelSection := types.SectionModel{
//             SectionNumber: int(protoSection.SectionNumber),
//             TimeAllowed:   int(protoSection.TimeAllowed),
//             Passages:      make([]types.PassageModel, len(protoSection.Passages)),
//         }

//         for j, protoPassage := range protoSection.Passages {
//             modelPassage := types.PassageModel{
//                 PassageNumber: int(protoPassage.PassageNumber),
//                 Title:         protoPassage.Title,
//                 Content:       make([]types.ParagraphContentModel, len(protoPassage.Content)),
//                 Questions:     make([]types.QuestionModel, len(protoPassage.Questions)),
//             }

//             for k, protoContent := range protoPassage.Content {
//                 modelPassage.Content[k] = types.ParagraphContentModel{
//                     ParagraphSummary: protoContent.ParagraphSummary,
//                     KeyWords:         protoContent.KeyWords,
//                     KeySentence:      protoContent.KeySentence,
//                 }
//             }

//             for k, protoQuestion := range protoPassage.Questions {
//                 modelQuestion := types.QuestionModel{
//                     QuestionNumber: int(protoQuestion.QuestionNumber),
//                     Type:           types.QuestionTypeModel(protoQuestion.Type.String()),
//                     Content:        protoQuestion.Content,
//                     Options:        protoQuestion.Options,
//                 }

//                 switch answer := protoQuestion.CorrectAnswer.(type) {
//                 case *proto.QuestionProto_StringAnswer:
//                     modelQuestion.CorrectAnswer = answer.StringAnswer
//                 case *proto.QuestionProto_StringArrayAnswer:
//                     modelQuestion.CorrectAnswer = answer.StringArrayAnswer.Values
//                 }

//                 modelPassage.Questions[k] = modelQuestion
//             }

//             modelSection.Passages[j] = modelPassage
//         }

//         modelReading.ReadingResType.Sections[i] = modelSection
//     }

//     return modelReading
// }
func convertProtoReadingReqTestTypeToModel(protoReading *proto.ReadingReq) types.ReadingTestModel {
    if protoReading == nil || protoReading.ReadingTest == nil {
        return types.ReadingTestModel{}
    }

    modelReading := types.ReadingTestModel{
        TestNumber: int(protoReading.ReadingTest.TestNumber),
        Sections:   make([]types.SectionModel, len(protoReading.ReadingTest.Sections)),
    }

    for i, protoSection := range protoReading.ReadingTest.Sections {
        modelSection := types.SectionModel{
            SectionNumber: int(protoSection.SectionNumber),
            TimeAllowed:   int(protoSection.TimeAllowed),
            Passages:      make([]types.PassageModel, len(protoSection.Passages)),
        }

        for j, protoPassage := range protoSection.Passages {
            modelPassage := types.PassageModel{
                PassageNumber: int(protoPassage.PassageNumber),
                Title:         protoPassage.Title,
                Content:       make([]types.ParagraphContentModel, len(protoPassage.Content)),
                Questions:     make([]types.QuestionModel, len(protoPassage.Questions)),
            }

            for k, protoContent := range protoPassage.Content {
                modelPassage.Content[k] = types.ParagraphContentModel{
                    ParagraphSummary: protoContent.ParagraphSummary,
                    KeyWords:         protoContent.KeyWords,
                    KeySentence:      protoContent.KeySentence,
                }
            }

            for k, protoQuestion := range protoPassage.Questions {
                modelQuestion := types.QuestionModel{
                    QuestionNumber: int(protoQuestion.QuestionNumber),
                    Type:         types.QuestionTypeModel(protoQuestion.Type.String()),
                    Content:        protoQuestion.Content,
                    Options:        protoQuestion.Options,
                }

                switch answer := protoQuestion.CorrectAnswer.(type) {
                case *proto.QuestionProto_StringAnswer:
                    modelQuestion.CorrectAnswer = answer.StringAnswer
                case *proto.QuestionProto_StringArrayAnswer:
                    modelQuestion.CorrectAnswer = answer.StringArrayAnswer.Values
                }

                modelPassage.Questions[k] = modelQuestion
            }

            modelSection.Passages[j] = modelPassage
        }

        modelReading.Sections[i] = modelSection
    }

    return modelReading
}


//----------------------------------------


func convertProtoReadingResToModel(res *proto.ReadingRes) *types.ReadingResModel {
    return &types.ReadingResModel{
        ID:             res.Id,
        ReadingResType: convertProtoReadingTestToModel(res.ReadingTest),
        CreatedAt:      res.CreatedAt.AsTime(),
        UpdatedAt:      res.UpdatedAt.AsTime(),
    }
}

func convertModelToProtoReadingRes(model *types.ReadingResModel) *proto.ReadingRes {
    return &proto.ReadingRes{
        Id:          model.ID,
        ReadingTest: convertModelToProtoReadingTest(model.ReadingResType),
        CreatedAt:   timestamppb.New(model.CreatedAt),
        UpdatedAt:   timestamppb.New(model.UpdatedAt),
    }
}

func convertProtoReadingTestToModel(protoTest *proto.ReadingTestProto) types.ReadingTestModel {
    sections := make([]types.SectionModel, len(protoTest.Sections))
    for i, protoSection := range protoTest.Sections {
        sections[i] = convertProtoSectionToModel(protoSection)
    }

    return types.ReadingTestModel{
        TestNumber: int(protoTest.TestNumber),
        Sections:   sections,
    }
}

func convertModelToProtoReadingTest(model types.ReadingTestModel) *proto.ReadingTestProto {
    sections := make([]*proto.SectionProto, len(model.Sections))
    for i, section := range model.Sections {
        sections[i] = convertModelToProtoSection(section)
    }

    return &proto.ReadingTestProto{
        TestNumber: int32(model.TestNumber),
        Sections:   sections,
    }
}

func convertProtoSectionToModel(protoSection *proto.SectionProto) types.SectionModel {
    passages := make([]types.PassageModel, len(protoSection.Passages))
    for i, protoPassage := range protoSection.Passages {
        passages[i] = convertProtoPassageToModel(protoPassage)
    }

    return types.SectionModel{
        SectionNumber: int(protoSection.SectionNumber),
        TimeAllowed:   int(protoSection.TimeAllowed),
        Passages:      passages,
    }
}

func convertModelToProtoSection(model types.SectionModel) *proto.SectionProto {
    passages := make([]*proto.PassageProto, len(model.Passages))
    for i, passage := range model.Passages {
        passages[i] = convertModelToProtoPassage(passage)
    }

    return &proto.SectionProto{
        SectionNumber: int32(model.SectionNumber),
        TimeAllowed:   int32(model.TimeAllowed),
        Passages:      passages,
    }
}

func convertProtoPassageToModel(protoPassage *proto.PassageProto) types.PassageModel {
    content := make([]types.ParagraphContentModel, len(protoPassage.Content))
    for i, protoContent := range protoPassage.Content {
        content[i] = convertProtoParagraphContentToModel(protoContent)
    }

    questions := make([]types.QuestionModel, len(protoPassage.Questions))
    for i, protoQuestion := range protoPassage.Questions {
        questions[i] = convertProtoQuestionToModel(protoQuestion)
    }

    return types.PassageModel{
        PassageNumber: int(protoPassage.PassageNumber),
        Title:         protoPassage.Title,
        Content:       content,
        Questions:     questions,
    }
}

func convertModelToProtoPassage(model types.PassageModel) *proto.PassageProto {
    content := make([]*proto.ParagraphContentProto, len(model.Content))
    for i, paragraphContent := range model.Content {
        content[i] = convertModelToProtoParagraphContent(paragraphContent)
    }

    questions := make([]*proto.QuestionProto, len(model.Questions))
    for i, question := range model.Questions {
        questions[i] = convertModelToProtoQuestion(question)
    }

    return &proto.PassageProto{
        PassageNumber: int32(model.PassageNumber),
        Title:         model.Title,
        Content:       content,
        Questions:     questions,
    }
}

func convertProtoParagraphContentToModel(protoContent *proto.ParagraphContentProto) types.ParagraphContentModel {
    return types.ParagraphContentModel{
        ParagraphSummary: protoContent.ParagraphSummary,
        KeyWords:         protoContent.KeyWords,
        KeySentence:      protoContent.KeySentence,
    }
}

func convertModelToProtoParagraphContent(model types.ParagraphContentModel) *proto.ParagraphContentProto {
    return &proto.ParagraphContentProto{
        ParagraphSummary: model.ParagraphSummary,
        KeyWords:         model.KeyWords,
        KeySentence:      model.KeySentence,
    }
}

func convertProtoQuestionToModel(protoQuestion *proto.QuestionProto) types.QuestionModel {
    questionModel := types.QuestionModel{
        QuestionNumber: int(protoQuestion.QuestionNumber),
        Type:           convertProtoQuestionTypeToModel(protoQuestion.Type),
        Content:        protoQuestion.Content,
        Options:        protoQuestion.Options,
    }

    switch answer := protoQuestion.CorrectAnswer.(type) {
    case *proto.QuestionProto_StringAnswer:
        questionModel.CorrectAnswer = answer.StringAnswer
    case *proto.QuestionProto_StringArrayAnswer:
        questionModel.CorrectAnswer = answer.StringArrayAnswer.Values
    }

    return questionModel
}

func convertModelToProtoQuestion(model types.QuestionModel) *proto.QuestionProto {
    protoQuestion := &proto.QuestionProto{
        QuestionNumber: int32(model.QuestionNumber),
        Type:           convertModelQuestionTypeToProto(model.Type),
        Content:        model.Content,
        Options:        model.Options,
    }

    switch answer := model.CorrectAnswer.(type) {
    case string:
        protoQuestion.CorrectAnswer = &proto.QuestionProto_StringAnswer{StringAnswer: answer}
    case []string:
        protoQuestion.CorrectAnswer = &proto.QuestionProto_StringArrayAnswer{
            StringArrayAnswer: &proto.StringArrayProto{Values: answer},
        }
    }

    return protoQuestion
}

func convertProtoQuestionTypeToModel(protoType proto.QuestionTypeProto) types.QuestionTypeModel {
    switch protoType {
    case proto.QuestionTypeProto_MULTIPLE_CHOICE:
        return types.MultipleChoice
    case proto.QuestionTypeProto_TRUE_FALSE_NOT_GIVEN:
        return types.TrueFalseNotGiven
    case proto.QuestionTypeProto_MATCHING:
        return types.Matching
    case proto.QuestionTypeProto_SHORT_ANSWER:
        return types.ShortAnswer
    default:
        return ""
    }
}

func convertModelQuestionTypeToProto(modelType types.QuestionTypeModel) proto.QuestionTypeProto {
    switch modelType {
    case types.MultipleChoice:
        return proto.QuestionTypeProto_MULTIPLE_CHOICE
    case types.TrueFalseNotGiven:
        return proto.QuestionTypeProto_TRUE_FALSE_NOT_GIVEN
    case types.Matching:
        return proto.QuestionTypeProto_MATCHING
    case types.ShortAnswer:
        return proto.QuestionTypeProto_SHORT_ANSWER
    default:
        return proto.QuestionTypeProto_MULTIPLE_CHOICE // Default to MULTIPLE_CHOICE
    }
}

func convertProtoReadingReqToModel(protoReading *proto.ReadingReq) types.ReadingReqModel {
    if protoReading == nil || protoReading.ReadingTest == nil {
        return types.ReadingReqModel{}
    }

    modelReading := types.ReadingReqModel{
        ReadingReqTestType: types.ReadingTestModel{
            TestNumber: int(protoReading.ReadingTest.TestNumber),
            Sections:   make([]types.SectionModel, len(protoReading.ReadingTest.Sections)),
        },
    }

    for i, protoSection := range protoReading.ReadingTest.Sections {
        modelSection := types.SectionModel{
            SectionNumber: int(protoSection.SectionNumber),
            TimeAllowed:   int(protoSection.TimeAllowed),
            Passages:      make([]types.PassageModel, len(protoSection.Passages)),
        }

        for j, protoPassage := range protoSection.Passages {
            modelPassage := types.PassageModel{
                PassageNumber: int(protoPassage.PassageNumber),
                Title:         protoPassage.Title,
                Content:       make([]types.ParagraphContentModel, len(protoPassage.Content)),
                Questions:     make([]types.QuestionModel, len(protoPassage.Questions)),
            }

            for k, protoContent := range protoPassage.Content {
                modelPassage.Content[k] = types.ParagraphContentModel{
                    ParagraphSummary: protoContent.ParagraphSummary,
                    KeyWords:         protoContent.KeyWords,
                    KeySentence:      protoContent.KeySentence,
                }
            }

            for k, protoQuestion := range protoPassage.Questions {
                modelQuestion := types.QuestionModel{
                    QuestionNumber: int(protoQuestion.QuestionNumber),
                    Type:           types.QuestionTypeModel(protoQuestion.Type.String()),
                    Content:        protoQuestion.Content,
                    Options:        protoQuestion.Options,
                }

                switch answer := protoQuestion.CorrectAnswer.(type) {
                case *proto.QuestionProto_StringAnswer:
                    modelQuestion.CorrectAnswer = answer.StringAnswer
                case *proto.QuestionProto_StringArrayAnswer:
                    modelQuestion.CorrectAnswer = answer.StringArrayAnswer.Values
                }

                modelPassage.Questions[k] = modelQuestion
            }

            modelSection.Passages[j] = modelPassage
        }

        modelReading.ReadingReqTestType.Sections[i] = modelSection
    }

    return modelReading
}