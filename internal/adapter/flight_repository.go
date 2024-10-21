package adapter

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/jiazhen-lin/flight-booking/internal/domain"
)

type flightPostgresRepository struct {
	db *gorm.DB
}

var _ domain.FlightRepository = (*flightPostgresRepository)(nil)

func NewFlightPostgresRepository(db *gorm.DB) *flightPostgresRepository {
	return &flightPostgresRepository{db: db}
}

type FlightRow struct {
	ID                 uuid.UUID           `gorm:"column:id;type:uuid;default:uuid_generate_v4()"`
	CreatedAt          time.Time           `gorm:"column:created_at"`
	UpdatedAt          time.Time           `gorm:"column:updated_at"`
	Number             string              `gorm:"column:number"`
	DepartureAirportID domain.AirportID    `gorm:"column:departure_airport_id"`
	ArrivalAirportID   domain.AirportID    `gorm:"column:arrival_airport_id"`
	DepartureTime      time.Time           `gorm:"column:departure_time"`
	DurationSeconds    int64               `gorm:"column:duration_seconds"`
	Status             domain.FlightStatus `gorm:"column:status"`
	TotalSeats         int                 `gorm:"column:total_seats"`
	OverbookedSeats    int                 `gorm:"column:overbooked_seats"`
	AvailableSeats     int                 `gorm:"column:available_seats"`
	Price              decimal.Decimal     `gorm:"column:price"`
}

func (r FlightRow) TableName() string {
	return "flights"
}

func flightModelToDomain(rows []FlightRow) []domain.Flight {
	flights := make([]domain.Flight, len(rows))
	for i, row := range rows {
		flights[i] = domain.Flight{
			ID:                 row.ID.String(),
			Number:             row.Number,
			DepartureAirportID: row.DepartureAirportID,
			ArrivalAirportID:   row.ArrivalAirportID,
			DepartureTime:      row.DepartureTime,
			DurationSeconds:    row.DurationSeconds,
			Status:             row.Status,
			TotalSeats:         row.TotalSeats,
			OverbookedSeats:    row.OverbookedSeats,
			AvailableSeats:     row.AvailableSeats,
			Price:              row.Price,
		}
	}
	return flights
}

func (r *flightPostgresRepository) List(ctx context.Context, filter domain.ListFilter) ([]domain.Flight, error) {
	var rows []FlightRow
	if err := r.db.WithContext(ctx).
		Where("status = ?", domain.FlightStatusEnabled).
		Where("departure_airport_id = ?", filter.DepartureAirportID).
		Where("arrival_airport_id = ?", filter.ArrivalAirportID).
		Where("departure_time BETWEEN ? AND ?", filter.DepartureTimeFrom, filter.DepartureTimeTo).
		Order("departure_time ASC, id ASC").
		Limit(filter.Limit).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	return flightModelToDomain(rows), nil
}

type BookRow struct {
	ID        uuid.UUID            `gorm:"column:id;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time            `gorm:"column:created_at"`
	UpdatedAt time.Time            `gorm:"column:updated_at"`
	FlightID  uuid.UUID            `gorm:"column:flight_id"`
	Seats     int                  `gorm:"column:seats"`
	Status    domain.BookingStatus `gorm:"column:status"`
}

func (r BookRow) TableName() string {
	return "bookings"
}

func (r *flightPostgresRepository) Book(ctx context.Context, flightID string, seats int) error {
	fid, err := uuid.Parse(flightID)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// check flight id exists
		var row FlightRow
		if err := tx.Where("id = ?", flightID).
			First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrFlightNotFound
			}
			return err
		}
		if row.AvailableSeats < seats {
			return domain.ErrUnavailableFlightSeats
		}

		// update flight available seats
		result := tx.Model(&FlightRow{}).
			Where("id = ?", flightID).
			Where("available_seats >= ?", seats).
			Update("available_seats", gorm.Expr("available_seats - ?", seats))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return domain.ErrUnavailableFlightSeats
		}

		// create booking record
		if err := tx.Create(&BookRow{
			FlightID: fid,
			Seats:    seats,
			Status:   domain.BookingStatusEnabled,
		}).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
