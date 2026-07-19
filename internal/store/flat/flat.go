package flat

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/vadimfilimonov/house/internal/models"
	"github.com/vadimfilimonov/house/internal/storage/pg"
	"github.com/vadimfilimonov/house/internal/store/pgerr"
)

var (
	ErrFlatAlreadyExists  = errors.New("flat already exists")
	ErrFlatNotFound       = errors.New("flat is not found")
	ErrFlatStatusConflict = errors.New("flat status transition conflict")
	defaultTimeout        = 5 * time.Second
)

const updateStatusQuery = `UPDATE flats SET status = $2 WHERE id = $1 AND status = $3 RETURNING id, number, house_id, price, rooms, status`

var requiredStatusByNextStatus = map[models.Status]models.Status{
	models.OnModerationStatus: models.CreatedStatus,
	models.ApprovedStatus:     models.OnModerationStatus,
	models.DeclinedStatus:     models.OnModerationStatus,
}

type Store struct {
	storage *pg.Storage
}

func New(storage *pg.Storage) *Store {
	return &Store{storage: storage}
}

func (s *Store) Add(ctx context.Context, number, houseID, price, rooms int) (*models.Flat, error) {
	var err error

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.storage.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("transaction rollback failed: %w", rollbackErr)
			}
		}
	}()

	flatsQuery := `INSERT INTO flats (number, house_id, price, rooms, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var flatID int
	if err = tx.QueryRowContext(ctx, flatsQuery, number, houseID, price, rooms, models.CreatedStatus).Scan(&flatID); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgerr.UniqueViolation {
			return nil, ErrFlatAlreadyExists
		}

		return nil, fmt.Errorf("cannot add flat to database: %w", err)
	}

	housesQuery := `UPDATE houses SET update_at = NOW() WHERE id = $1`
	if _, err = tx.ExecContext(ctx, housesQuery, houseID); err != nil {
		return nil, fmt.Errorf("cannot update houses table: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return &models.Flat{
		ID:      flatID,
		Number:  number,
		HouseID: models.HouseID(houseID),
		Price:   price,
		Rooms:   rooms,
		Status:  models.CreatedStatus,
	}, nil
}

func (s *Store) ListByHouseID(ctx context.Context, houseID int, includeAllStatuses bool) ([]models.Flat, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `SELECT id, number, house_id, price, rooms, status FROM flats WHERE house_id = $1 ORDER BY id`
	args := []any{houseID}
	if !includeAllStatuses {
		query = `SELECT id, number, house_id, price, rooms, status FROM flats WHERE house_id = $1 AND status = $2 ORDER BY id`
		args = append(args, models.ApprovedStatus)
	}

	rows, err := s.storage.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot list flats: %w", err)
	}
	defer rows.Close()

	flats := make([]models.Flat, 0)
	for rows.Next() {
		flat, err := scanFlat(rows)
		if err != nil {
			return nil, err
		}

		flats = append(flats, *flat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot read flats rows: %w", err)
	}

	return flats, nil
}

func (s *Store) UpdateStatus(ctx context.Context, flatID int, status models.Status) (*models.Flat, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	requiredStatus, ok := requiredStatusByNextStatus[status]
	if !ok {
		return nil, ErrFlatStatusConflict
	}

	flat, err := scanFlat(s.storage.QueryRowContext(ctx, updateStatusQuery, flatID, status.String(), requiredStatus.String()))
	if err == nil {
		return flat, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("cannot update flat status: %w", err)
	}

	exists, err := s.exists(ctx, flatID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrFlatNotFound
	}

	return nil, ErrFlatStatusConflict
}

func (s *Store) exists(ctx context.Context, flatID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM flats WHERE id = $1)`

	var exists bool
	if err := s.storage.QueryRowContext(ctx, query, flatID).Scan(&exists); err != nil {
		return false, fmt.Errorf("cannot check flat existence: %w", err)
	}

	return exists, nil
}

type flatScanner interface {
	Scan(dest ...any) error
}

func scanFlat(scanner flatScanner) (*models.Flat, error) {
	var flat models.Flat
	var status string
	var houseID int
	if err := scanner.Scan(&flat.ID, &flat.Number, &houseID, &flat.Price, &flat.Rooms, &status); err != nil {
		return nil, err
	}

	flat.HouseID = models.HouseID(houseID)
	flat.Status = models.Status(status)

	return &flat, nil
}
