// internal/datastruct/error.go

package datastruct

import (
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrorCode — тип для кодов ошибок (строковый, но с константами)
type ErrorCode string

// Константы кодов ошибок — все возможные коды в одном месте
const (
	CodeNotFound         ErrorCode = "not_found"
	CodeValidationError  ErrorCode = "validation_error"
	CodeInternalError    ErrorCode = "internal_error"
	CodeUnauthorized     ErrorCode = "unauthorized"
	CodeConflict         ErrorCode = "conflict"
	CodePermissionDenied ErrorCode = "permission_denied"
	CodeRateLimited      ErrorCode = "rate_limited"
)

// AppError — наша основная структура ошибки
type AppError struct {
	Err  error
	Code ErrorCode
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %v", e.Code, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// ── Фабрики ────────────────────────────────────────────────────────────────

func NewNotFound(msg string) *AppError {
	return &AppError{Err: errors.New(msg), Code: CodeNotFound}
}

func NewValidationError(msg string) *AppError {
	return &AppError{Err: errors.New(msg), Code: CodeValidationError}
}

func NewInternalError(msg string) *AppError {
	return &AppError{Err: errors.New(msg), Code: CodeInternalError}
}

func NewUnauthorized(msg string) *AppError {
	return &AppError{Err: errors.New(msg), Code: CodeUnauthorized}
}

func NewConflict(msg string) *AppError {
	return &AppError{Err: errors.New(msg), Code: CodeConflict}
}

func NewRateLimited(msg string) *AppError {
	return &AppError{Err: errors.New(msg), Code: CodeRateLimited}
}

func NewPermissionDenied(msg string) *AppError {
	return &AppError{Err: errors.New(msg), Code: CodePermissionDenied}
}

// WrapPgError — маппинг ошибок PostgreSQL → наши коды
func WrapPgError(err error) *AppError {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return NewInternalError(err.Error())
	}

	switch pgErr.Code {
	case pgerrcode.UniqueViolation:
		return &AppError{Err: err, Code: CodeConflict}

	case pgerrcode.NotNullViolation,
		pgerrcode.ForeignKeyViolation,
		pgerrcode.CheckViolation,
		pgerrcode.InvalidTextRepresentation:
		return &AppError{Err: err, Code: CodeValidationError}

	case pgerrcode.UndefinedTable,
		pgerrcode.UndefinedColumn:
		return &AppError{Err: err, Code: CodeNotFound}

	default:
		return &AppError{Err: err, Code: CodeInternalError}
	}
}
