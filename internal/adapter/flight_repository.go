package adapter

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
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
	ID                 int64               `gorm:"column:id;primary_key;autoIncrement"`
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
			ID:                 row.ID,
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
	query := r.db.WithContext(ctx).Where("status = ?", domain.FlightStatusEnabled)

	if filter.AvailableSeats != nil {
		if *filter.AvailableSeats {
			query = query.Where("available_seats > 0")
		} else {
			query = query.Where("available_seats = 0")
		}
	}
	if filter.CursorID != nil {
		query = query.Where("(departure_time = ? AND id > ?) OR departure_time > ?",
			filter.DepartureTimeFrom, *filter.CursorID, filter.DepartureTimeFrom,
		)
		query = query.Where("departure_time < ?", filter.DepartureTimeTo)
	} else {
		query = query.Where("departure_time BETWEEN ? AND ?", filter.DepartureTimeFrom, filter.DepartureTimeTo)
	}

	query = query.Order("departure_time ASC, id ASC").Limit(filter.Limit)

	var rows []FlightRow
	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}

	return flightModelToDomain(rows), nil
}

type BookRow struct {
	ID        int64                `gorm:"column:id;primary_key;autoIncrement"`
	CreatedAt time.Time            `gorm:"column:created_at"`
	UpdatedAt time.Time            `gorm:"column:updated_at"`
	Status    domain.BookingStatus `gorm:"column:status"`
}

func (r BookRow) TableName() string {
	return "bookings"
}

type BookingSeatsRow struct {
	ID        int64     `gorm:"column:id;primary_key;autoIncrement"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	BookingID int64     `gorm:"column:booking_id"`
	FlightID  int64     `gorm:"column:flight_id"`
	Seats     int       `gorm:"column:seats"`
}

func (r BookingSeatsRow) TableName() string {
	return "booking_seats"
}

func (r *flightPostgresRepository) Book(ctx context.Context, seats []domain.FlightSeats) error {
	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// insert a booking record
		booking := BookRow{
			Status: domain.BookingStatusEnabled,
		}
		if err := tx.Create(&booking).Error; err != nil {
			return err
		}
		bookingID := booking.ID
		logrus.Infof("bookingID: %d", bookingID)

		// check flight id exists
		for _, s := range seats {
			var row FlightRow
			if err := tx.Where("id = ?", s.FlightID).
				First(&row).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return domain.ErrFlightNotFound
				}
				return err
			}
			if row.AvailableSeats < s.Seats {
				return domain.ErrUnavailableFlightSeats
			}

			// update flight available seats
			result := tx.Model(&FlightRow{}).
				Where("id = ?", s.FlightID).
				Where("available_seats >= ?", s.Seats).
				Update("available_seats", gorm.Expr("available_seats - ?", s.Seats))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return domain.ErrUnavailableFlightSeats
			}

			// create booking seats record
			if err := tx.Create(&BookingSeatsRow{
				BookingID: bookingID,
				FlightID:  s.FlightID,
				Seats:     s.Seats,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
