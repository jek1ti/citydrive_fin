package domain

import "errors"

var (
    ErrCarNotFound      = errors.New("car not found")
    ErrInvalidTimeRange = errors.New("invalid time range: from must be less than to")
    ErrInvalidCarID     = errors.New("invalid car id")
)