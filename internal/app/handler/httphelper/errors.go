package httphelper

import (
	"errors"
	"net/http"

	"github.com/pedropiedade7/mundo-invest-test/internal/domain"
)

func StatusFromError(err error) int {
	switch {
	case errors.Is(err, domain.ErrClientAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, domain.ErrClientNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrDuplicateEvent):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
