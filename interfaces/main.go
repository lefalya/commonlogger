package interfaces

import (
	"log/slog"

	"github.com/lefalya/commonlogger/schema"
)

type LogHelper func(
	logger *slog.Logger,
	errorSubject error,
	errorDetail string,
	context string,
	schema any,
) *schema.CommonError
