package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type DataExportsRepository interface {
	CreateDataExport(ctx context.Context, req domain.CreateDataExportAction) (*domain.DataExport, error)
	GetDataExportByID(ctx context.Context, id, userID int64) (*domain.DataExport, error)
	GetDataExportsByUserID(ctx context.Context, userID int64) ([]domain.DataExport, error)
	UpdateDataExportStatus(ctx context.Context, req domain.UpdateDataExportStatusAction) (*domain.DataExport, error)
	UpdateDataExportFile(ctx context.Context, req domain.UpdateDataExportFileAction) (*domain.DataExport, error)
	DeleteDataExport(ctx context.Context, id, userID int64) error
	DeleteExpiredDataExports(ctx context.Context) error
	GetPendingDataExports(ctx context.Context) ([]domain.DataExport, error)
}

type dataExportsRepository struct {
	store db.Store
}

func NewDataExportsRepository(store db.Store) DataExportsRepository {
	return &dataExportsRepository{
		store: store,
	}
}

func (r *dataExportsRepository) CreateDataExport(ctx context.Context, req domain.CreateDataExportAction) (*domain.DataExport, error) {
	dbExport, err := r.store.CreateDataExport(ctx, db.CreateDataExportParams{
		UserID:     req.UserID,
		ExportType: req.ExportType,
		Status:     domain.DataExportStatusPending,
		ExpiresAt:  req.ExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbExport), nil
}

func (r *dataExportsRepository) GetDataExportByID(ctx context.Context, id, userID int64) (*domain.DataExport, error) {
	dbExport, err := r.store.GetDataExportByID(ctx, db.GetDataExportByIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbExport), nil
}

func (r *dataExportsRepository) GetDataExportsByUserID(ctx context.Context, userID int64) ([]domain.DataExport, error) {
	dbExports, err := r.store.GetDataExportsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	exports := make([]domain.DataExport, len(dbExports))
	for i, export := range dbExports {
		exports[i] = *r.toDomain(export)
	}

	return exports, nil
}

func (r *dataExportsRepository) UpdateDataExportStatus(ctx context.Context, req domain.UpdateDataExportStatusAction) (*domain.DataExport, error) {
	dbExport, err := r.store.UpdateDataExportStatus(ctx, db.UpdateDataExportStatusParams{
		ID:      req.ID,
		Column2: req.Status,
		UserID:  req.UserID,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbExport), nil
}

func (r *dataExportsRepository) UpdateDataExportFile(ctx context.Context, req domain.UpdateDataExportFileAction) (*domain.DataExport, error) {
	var filePath pgtype.Text
	if req.FilePath != "" {
		filePath = pgtype.Text{String: req.FilePath, Valid: true}
	}

	var fileSize pgtype.Int8
	if req.FileSize != 0 {
		fileSize = pgtype.Int8{Int64: req.FileSize, Valid: true}
	}

	dbExport, err := r.store.UpdateDataExportFile(ctx, db.UpdateDataExportFileParams{
		ID:       req.ID,
		UserID:   req.UserID,
		FilePath: filePath,
		FileSize: fileSize,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbExport), nil
}

func (r *dataExportsRepository) DeleteDataExport(ctx context.Context, id, userID int64) error {
	return r.store.DeleteDataExport(ctx, db.DeleteDataExportParams{
		ID:     id,
		UserID: userID,
	})
}

func (r *dataExportsRepository) DeleteExpiredDataExports(ctx context.Context) error {
	return r.store.DeleteExpiredDataExports(ctx)
}

func (r *dataExportsRepository) GetPendingDataExports(ctx context.Context) ([]domain.DataExport, error) {
	dbExports, err := r.store.GetPendingDataExports(ctx)
	if err != nil {
		return nil, err
	}

	exports := make([]domain.DataExport, len(dbExports))
	for i, export := range dbExports {
		exports[i] = *r.toDomain(export)
	}

	return exports, nil
}

func (r *dataExportsRepository) toDomain(dbExport db.DataExport) *domain.DataExport {
	export := &domain.DataExport{
		ID:         dbExport.ID,
		UserID:     dbExport.UserID,
		ExportType: dbExport.ExportType,
		Status:     dbExport.Status,
		ExpiresAt:  dbExport.ExpiresAt,
		CreatedAt:  dbExport.CreatedAt,
	}

	if dbExport.FilePath.Valid {
		export.FilePath = &dbExport.FilePath.String
	}
	if dbExport.FileSize.Valid {
		export.FileSize = &dbExport.FileSize.Int64
	}
	if dbExport.CompletedAt != nil {
		export.CompletedAt = dbExport.CompletedAt
	}

	return export
}
