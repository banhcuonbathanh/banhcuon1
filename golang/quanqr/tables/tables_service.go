package tables_test

import (
	"context"
	"log"

	"english-ai-full/quanqr/proto_qr/table"

	"google.golang.org/protobuf/types/known/emptypb"
)

type TableServiceStruct struct {
	tableRepo *TableRepository
	table.UnimplementedTableServiceServer
}

func NewTableService(tableRepo *TableRepository) *TableServiceStruct {
	return &TableServiceStruct{
		tableRepo: tableRepo,
	}
}

func (ts *TableServiceStruct) GetTableList(ctx context.Context, _ *emptypb.Empty) (*table.TableListResponse, error) {
	log.Println("Fetching all tables")

	tables, err := ts.tableRepo.GetTableList(ctx)
	if err != nil {
		log.Println("Error fetching tables:", err)
		return nil, err
	}

	return &table.TableListResponse{
		Data:    tables,
		Message: "Tables fetched successfully",
	}, nil
}

func (ts *TableServiceStruct) GetTableDetail(ctx context.Context, req *table.TableNumberRequest) (*table.TableResponse, error) {
	log.Println("Fetching table detail for number:", req.Number)

	tableDetail, err := ts.tableRepo.GetTableDetail(ctx, req.Number)
	if err != nil {
		log.Println("Error fetching table detail:", err)
		return nil, err
	}

	return &table.TableResponse{
		Data:    tableDetail,
		Message: "Table detail fetched successfully",
	}, nil
}

func (ts *TableServiceStruct) CreateTable(ctx context.Context, req *table.CreateTableRequest) (*table.TableResponse, error) {
	log.Print("golang/quanqr/tables/tables_service.go 22 ")

	createdTable, err := ts.tableRepo.CreateTable(ctx, req)
	if err != nil {
		log.Print("golang/quanqr/tables/tables_service.go 22 err", err)
		log.Println("Error creating table:", err)
		return nil, err
	}

	return &table.TableResponse{
		Data:    createdTable,
		Message: "Table created successfully",
	}, nil
}

func (ts *TableServiceStruct) UpdateTable(ctx context.Context, req *table.UpdateTableRequest) (*table.TableResponse, error) {
	log.Println("Updating table:", req.Number)

	updatedTable, err := ts.tableRepo.UpdateTable(ctx, req)
	if err != nil {
		log.Println("Error updating table:", err)
		return nil, err
	}

	return &table.TableResponse{
		Data:    updatedTable,
		Message: "Table updated successfully",
	}, nil
}

func (ts *TableServiceStruct) DeleteTable(ctx context.Context, req *table.TableNumberRequest) (*table.TableResponse, error) {
	log.Println("Deleting table:", req.Number)

	deletedTable, err := ts.tableRepo.DeleteTable(ctx, req.Number)
	if err != nil {
		log.Println("Error deleting table:", err)
		return nil, err
	}

	return &table.TableResponse{
		Data:    deletedTable,
		Message: "Table deleted successfully",
	}, nil
}