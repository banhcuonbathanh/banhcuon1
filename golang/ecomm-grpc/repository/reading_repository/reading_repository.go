package repository

import (
	"context"

	"encoding/json"
	"fmt"
	"log"
	"time"

	"english-ai-full/ecomm-api/types"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ReadingRepository struct {
	db *pgxpool.Pool
}

func NewReadingRepository(db *pgxpool.Pool) *ReadingRepository {
	return &ReadingRepository{
		db: db,
	}
}

func (r *ReadingRepository) CreateReading(ctx context.Context, req *types.ReadingReqModel) (*types.ReadingResModel, error) {
    log.Println("Inserting new reading test into database", req)

    query := `
        INSERT INTO reading_test_models (test_number, sections)
        VALUES ($1, $2)
        RETURNING id, test_number, sections
    `

    sectionsJSON, err := json.Marshal(req.ReadingReqTestType.Sections)
    if err != nil {
        return nil, fmt.Errorf("error marshaling sections: %w", err)
    }

    var returnedID int
    var testNumber int
    var sectionsJSONReturn []byte

    err = r.db.QueryRow(ctx, query, req.ReadingReqTestType.TestNumber, sectionsJSON).Scan(&returnedID, &testNumber, &sectionsJSONReturn)
    if err != nil {
        log.Println("Error inserting reading test:", err)
        return nil, fmt.Errorf("error inserting reading test: %w", err)
    }

    var sections []types.SectionModel
    err = json.Unmarshal(sectionsJSONReturn, &sections)
    if err != nil {
        return nil, fmt.Errorf("error unmarshaling sections: %w", err)
    }

    now := time.Now()

    return &types.ReadingResModel{
        ID: returnedID,
        ReadingResType: types.ReadingTestModel{
            TestNumber: testNumber,
            Sections:   sections,
        },
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}
func (r *ReadingRepository) SaveReading(ctx context.Context, req *types.ReadingResModel) (*types.ReadingResModel, error) {
	log.Println("Saving reading test in database")

	query := `INSERT INTO reading_tests (id, test_number, sections, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5) 
			  ON CONFLICT (id) DO UPDATE 
			  SET test_number = $2, sections = $3, updated_at = $5`

	sectionsJSON, err := json.Marshal(req.ReadingResType.Sections)
	if err != nil {
		return nil, fmt.Errorf("error marshaling sections: %w", err)
	}

	now := time.Now()
	_, err = r.db.Exec(ctx, query, req.ID, req.ReadingResType.TestNumber, sectionsJSON, now, now)
	if err != nil {
		log.Println("Error saving reading test:", err)
		return nil, fmt.Errorf("error saving reading test: %w", err)
	}

	return req, nil
}


func (r *ReadingRepository) UpdateReading(ctx context.Context, req *types.ReadingResModel) (*types.ReadingResModel, error) {
	log.Println("Updating reading test in database")

	query := `UPDATE reading_tests 
			  SET test_number = $2, sections = $3, updated_at = $4
			  WHERE id = $1
			  RETURNING created_at`

	sectionsJSON, err := json.Marshal(req.ReadingResType.Sections)
	if err != nil {
		return nil, fmt.Errorf("error marshaling sections: %w", err)
	}

	now := time.Now()
	var createdAt time.Time
	err = r.db.QueryRow(ctx, query, req.ID, req.ReadingResType.TestNumber, sectionsJSON, now).Scan(&createdAt)
	if err != nil {
		log.Println("Error updating reading test:", err)
		return nil, fmt.Errorf("error updating reading test: %w", err)
	}

	return &types.ReadingResModel{
        ID:             req.ID, // Note: You might need to convert this to int64 if that's what your struct expects
        ReadingResType: req.ReadingResType,
        CreatedAt:      now,
        UpdatedAt:      now,
    }, nil
}

func (r *ReadingRepository) DeleteReading(ctx context.Context, req *types.ReadingResModel) error {
	log.Println("Deleting reading test from database")

	query := `DELETE FROM reading_tests WHERE id = $1`

	_, err := r.db.Exec(ctx, query, req.ID)
	if err != nil {
		log.Println("Error deleting reading test:", err)
		return fmt.Errorf("error deleting reading test: %w", err)
	}

	return nil
}


func (r *ReadingRepository) FindAllReading(ctx context.Context) (*types.ReadingResList, error) {
    log.Println("Fetching all reading tests from database")

    query := `SELECT id, test_number, sections, created_at, updated_at FROM reading_res_models`

    rows, err := r.db.Query(ctx, query)
    if err != nil {
        log.Println("Error fetching reading tests: FindAllReading ReadingRepository", err)
        return nil, fmt.Errorf("error fetching reading tests: %w", err)
    }
    defer rows.Close()

    var readings []*types.ReadingResModel
    for rows.Next() {
        var id int
        var testNumber int
        var sectionsJSON []byte
        var createdAt, updatedAt time.Time

        err := rows.Scan(&id, &testNumber, &sectionsJSON, &createdAt, &updatedAt)
        if err != nil {
            log.Println("Error scanning reading test row: FindAllReading ReadingRepository", err)
            return nil, fmt.Errorf("error scanning reading test row: %w", err)
        }

        var sections []types.SectionModel
        err = json.Unmarshal(sectionsJSON, &sections)
        if err != nil {
            log.Println("Error unmarshaling sections: FindAllReading ReadingRepository", err)
            return nil, fmt.Errorf("error unmarshaling sections: %w", err)
        }

        readings = append(readings, &types.ReadingResModel{
            ID: id,
            ReadingResType: types.ReadingTestModel{
                TestNumber: testNumber,
                Sections:   sections,
            },
            CreatedAt: createdAt,
            UpdatedAt: updatedAt,
        })
    }

    return &types.ReadingResList{
        Readings:   readings,
        TotalCount: int(len(readings)),
    }, nil
}

func (r *ReadingRepository) FindByID(ctx context.Context, ID int) (*types.ReadingResModel, error) {
	log.Println("Fetching reading test by ID from database")

	query := `SELECT test_number, sections, created_at, updated_at FROM reading_tests WHERE id = $1`

	var testNumber int
	var sectionsJSON []byte
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(ctx, query, ID).Scan(&testNumber, &sectionsJSON, &createdAt, &updatedAt)
	if err != nil {
		log.Println("Error fetching reading test:", err)
		return nil, fmt.Errorf("error fetching reading test: %w", err)
	}

	var sections []types.SectionModel
	err = json.Unmarshal(sectionsJSON, &sections)
	if err != nil {
		log.Println("Error unmarshaling sections:", err)
		return nil, fmt.Errorf("error unmarshaling sections: %w", err)
	}

	return &types.ReadingResModel{
		ID:             ID,
		ReadingResType: types.ReadingTestModel{
			TestNumber: testNumber,
			Sections:   sections,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}


func (r *ReadingRepository) FindReadingByPage(ctx context.Context, req *types.PageRequestModel) (*types.ReadingResList, error) {
	log.Println("Fetching reading tests by page from database")

	query := `SELECT id, test_number, sections, created_at, updated_at 
			  FROM reading_tests 
			  ORDER BY created_at DESC 
			  LIMIT $1 OFFSET $2`

	offset := (req.PageNumber - 1) * req.PageSize

	rows, err := r.db.Query(ctx, query, req.PageSize, offset)
	if err != nil {
		log.Println("Error fetching reading tests:", err)
		return nil, fmt.Errorf("error fetching reading tests: %w", err)
	}
	defer rows.Close()

	var readings []*types.ReadingResModel
	for rows.Next() {
        var id int
        var testNumber int
        var sectionsJSON []byte
        var createdAt, updatedAt time.Time

		err := rows.Scan(&id, &testNumber, &sectionsJSON, &createdAt, &updatedAt)
		if err != nil {
			log.Println("Error scanning reading test row:", err)
			return nil, fmt.Errorf("error scanning reading test row: %w", err)
		}

		var sections []types.SectionModel
		err = json.Unmarshal(sectionsJSON, &sections)
		if err != nil {
			log.Println("Error unmarshaling sections:", err)
			return nil, fmt.Errorf("error unmarshaling sections: %w", err)
		}

        readings = append(readings, &types.ReadingResModel{
            ID: id,
            ReadingResType: types.ReadingTestModel{
                TestNumber: testNumber,
                Sections:   sections,
            },
            CreatedAt: createdAt,
            UpdatedAt: updatedAt,
        })
	}

	// Get total count
	var totalCount int
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM reading_tests").Scan(&totalCount)
	if err != nil {
		log.Println("Error getting total count:", err)
		return nil, fmt.Errorf("error getting total count: %w", err)
	}

	return &types.ReadingResList{
		Readings:   readings,
		TotalCount: totalCount,
	}, nil
}