package repository

// import (
// 	"context"
// 	"english-ai-full/ecomm-api/types"
// 	"testing"
// 	"time"

// 	"github.com/jackc/pgx/v4/pgxpool"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func setupTestDB(t *testing.T) *pgxpool.Pool {
// 	// Set up a connection to your test database
// 	// You might want to use an environment variable for the connection string
// 	connString := "postgresql://myuser:mypassword@localhost:5432/testdb"
	
// 	config, err := pgxpool.ParseConfig(connString)
// 	require.NoError(t, err)

// 	pool, err := pgxpool.ConnectConfig(context.Background(), config)
// 	require.NoError(t, err)

// 	return pool
// }

// func TestCreateReading(t *testing.T) {
//     db := setupTestDB(t)
//     defer db.Close()

//     repo := NewReadingRepository(db)

//     ctx := context.Background()
//     req := &types.ReadingReq{
//         ReadingReqType: types.ReadingTest{
//             TestNumber: 1,
//             Sections: []types.Section{
//                 {
//                     SectionNumber: 1,
//                     TimeAllowed:   60,
//                     Passages: []types.Passage{
//                         {
//                             PassageNumber: 1,
//                             Title:         "Passage 1",
//                             Content: []types.ParagraphContent{
//                                 {ParagraphSummary: "This is the first paragraph."},
//                                 {ParagraphSummary: "This is the second paragraph."},
//                             },
//                             Questions: []types.Question{
//                                 {
//                                     QuestionNumber: 1,
//                                     Type:           types.MultipleChoice,
//                                     Content:        "What is the main idea of the first paragraph?",
//                                 },
//                             },
//                         },
//                     },
//                 },
//             },
//         },
//         CreatedAt: time.Now(),
//         UpdatedAt: time.Now(),
//     }

//     res, err := repo.CreateReading(ctx, req)

//     assert.NoError(t, err)
//     assert.NotNil(t, res)
//     assert.NotZero(t, res.ID)
//     assert.Equal(t, req.ReadingReqType.TestNumber, res.ReadingResType.TestNumber)
//     assert.Equal(t, req.ReadingReqType.Sections, res.ReadingResType.Sections)
//     assert.WithinDuration(t, time.Now(), res.CreatedAt, time.Second)
//     assert.WithinDuration(t, time.Now(), res.UpdatedAt, time.Second)
// }


// func TestSaveReading(t *testing.T) {
//     db := setupTestDB(t)
//     defer db.Close()

//     repo := NewReadingRepository(db)

//     ctx := context.Background()
//     req := &types.ReadingRes{
//         ID: 1,
//         ReadingResType: types.ReadingTest{
//             TestNumber: 1,
//             Sections: []types.Section{
//                 {
//                     SectionNumber: 1,
//                     TimeAllowed:   60,
//                     Passages: []types.Passage{
//                         {
//                             PassageNumber: 1,
//                             Title:         "Passage 1",
//                             Content: []types.ParagraphContent{
//                                 {ParagraphSummary: "This is the first paragraph."},
//                                 {ParagraphSummary: "This is the second paragraph."},
//                             },
//                             Questions: []types.Question{
//                                 {
//                                     QuestionNumber: 1,
//                                     Type:           types.MultipleChoice,
//                                     Content:        "What is the main idea of the first paragraph?",
//                                 },
//                             },
//                         },
//                     },
//                 },
//             },
//         },
//         CreatedAt: time.Now(),
//         UpdatedAt: time.Now(),
//     }

//     res, err := repo.SaveReading(ctx, req)

//     assert.NoError(t, err)
//     assert.NotNil(t, res)
//     assert.Equal(t, req.ID, res.ID)
//     assert.Equal(t, req.ReadingResType.TestNumber, res.ReadingResType.TestNumber)
//     assert.Equal(t, req.ReadingResType.Sections, res.ReadingResType.Sections)
//     assert.WithinDuration(t, time.Now(), res.UpdatedAt, time.Second)
// }
// func TestUpdateReading(t *testing.T) {
//     db := setupTestDB(t)
//     defer db.Close()

//     repo := NewReadingRepository(db)

//     ctx := context.Background()
//     req := &types.ReadingRes{
//         ID: 1,
//         ReadingResType: types.ReadingTest{
//             TestNumber: 1,
//             Sections: []types.Section{
//                 {
//                     SectionNumber: 1,
//                     TimeAllowed:   70, // Changed from 60 to 70
//                     Passages: []types.Passage{
//                         {
//                             PassageNumber: 1,
//                             Title:         "Updated Passage 1",
//                             Content: []types.ParagraphContent{
//                                 {ParagraphSummary: "This is the updated first paragraph."},
//                             },
//                             Questions: []types.Question{
//                                 {
//                                     QuestionNumber: 1,
//                                     Type:           types.MultipleChoice,
//                                     Content:        "Updated question?",
//                                 },
//                             },
//                         },
//                     },
//                 },
//             },
//         },
//         UpdatedAt: time.Now(),
//     }

//     res, err := repo.UpdateReading(ctx, req)

//     assert.NoError(t, err)
//     assert.NotNil(t, res)
//     assert.Equal(t, req.ID, res.ID)
//     assert.Equal(t, req.ReadingResType.TestNumber, res.ReadingResType.TestNumber)
//     assert.Equal(t, req.ReadingResType.Sections, res.ReadingResType.Sections)
//     assert.WithinDuration(t, time.Now(), res.UpdatedAt, time.Second)
// }

// func TestDeleteReading(t *testing.T) {
//     db := setupTestDB(t)
//     defer db.Close()

//     repo := NewReadingRepository(db)

//     ctx := context.Background()
//     req := &types.ReadingRes{
//         ID: 1,
//     }

//     err := repo.DeleteReading(ctx, req)

//     assert.NoError(t, err)

//     // Verify deletion
//     _, err = repo.FindByID(ctx, req.ID)
//     assert.Error(t, err) // Expect an error as the reading should not be found
// }

// func TestFindAllReading(t *testing.T) {
//     db := setupTestDB(t)
//     defer db.Close()

//     repo := NewReadingRepository(db)

//     ctx := context.Background()

//     // Insert some test data
//     for i := 1; i <= 3; i++ {
//         req := &types.ReadingRes{
//             ID: int64(i),
//             ReadingResType: types.ReadingTest{
//                 TestNumber: i,
//                 Sections: []types.Section{
//                     {
//                         SectionNumber: 1,
//                         TimeAllowed:   60,
//                     },
//                 },
//             },
//             CreatedAt: time.Now(),
//             UpdatedAt: time.Now(),
//         }
//         _, err := repo.SaveReading(ctx, req)
//         assert.NoError(t, err)
//     }

//     res, err := repo.FindAllReading(ctx)

//     assert.NoError(t, err)
//     assert.NotNil(t, res)
//     assert.Equal(t, 3, res.TotalCount)
//     assert.Len(t, res.Readings, 3)
// }

// func TestFindByID(t *testing.T) {
//     db := setupTestDB(t)
//     defer db.Close()

//     repo := NewReadingRepository(db)

//     ctx := context.Background()

//     // Insert test data
//     req := &types.ReadingRes{
//         ID: 1,
//         ReadingResType: types.ReadingTest{
//             TestNumber: 1,
//             Sections: []types.Section{
//                 {
//                     SectionNumber: 1,
//                     TimeAllowed:   60,
//                 },
//             },
//         },
//         CreatedAt: time.Now(),
//         UpdatedAt: time.Now(),
//     }
//     _, err := repo.SaveReading(ctx, req)
//     assert.NoError(t, err)

//     res, err := repo.FindByID(ctx, 1)

//     assert.NoError(t, err)
//     assert.NotNil(t, res)
//     assert.Equal(t, req.ID, res.ID)
//     assert.Equal(t, req.ReadingResType.TestNumber, res.ReadingResType.TestNumber)
//     assert.Equal(t, req.ReadingResType.Sections, res.ReadingResType.Sections)
// }

// func TestFindReadingByPage(t *testing.T) {
//     db := setupTestDB(t)
//     defer db.Close()

//     repo := NewReadingRepository(db)

//     ctx := context.Background()

//     // Insert some test data
//     for i := 1; i <= 15; i++ {
//         req := &types.ReadingRes{
//             ID: int64(i),
//             ReadingResType: types.ReadingTest{
//                 TestNumber: i,
//                 Sections: []types.Section{
//                     {
//                         SectionNumber: 1,
//                         TimeAllowed:   60,
//                     },
//                 },
//             },
//             CreatedAt: time.Now(),
//             UpdatedAt: time.Now(),
//         }
//         _, err := repo.SaveReading(ctx, req)
//         assert.NoError(t, err)
//     }

//     pageReq := &types.PageRequest{
//         PageNumber: 2,
//         PageSize:   5,
//     }

//     res, err := repo.FindReadingByPage(ctx, pageReq)

//     assert.NoError(t, err)
//     assert.NotNil(t, res)
//     assert.Equal(t, 15, res.TotalCount)
//     assert.Len(t, res.Readings, 5)
//     assert.Equal(t, int64(11), res.Readings[0].ID) // Assuming descending order by created_at
// }