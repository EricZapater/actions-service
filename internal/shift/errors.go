package shift

import "errors"

var (
	ErrShiftNotFound = errors.New("shift not found")
	ErrShiftDetailNotFound = errors.New("shift detail not found")
)