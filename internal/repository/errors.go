package repository

import "errors"

var (
	ErrProductNotFound              = errors.New("product not found")
	ErrInsufficientStock            = errors.New("insufficient stock")
	ErrInvalidQuantity              = errors.New("invalid quantity")
	ErrInvalidLimit                 = errors.New("invalid limit")
	ErrReservationNotFound          = errors.New("reservation not found")
	ErrInvalidReservationTransition = errors.New("invalid reservation transition")
	ErrReservationExpired           = errors.New("reservation expired")
)
