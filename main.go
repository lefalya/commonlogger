package commonlogger

import (
	"errors"
	"log/slog"
	"regexp"
	"strings"

	"math/rand"
)

type CommonError struct {
	Id          string
	Context     string
	Err         error
	ErrResponse error
}

func generateRandomString(length int) string {

	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = characters[rand.Intn(len(characters))]
	}

	return string(result)
}

func LogInfo(logger *slog.Logger, context string, args ...string) {

	if len(args)%2 != 0 {

		LogInfo(logger, "logInfo.InvalidArgs", "args", "null")
		return
	}

	var keyValues []interface{}
	for i := 0; i < len(args); i += 2 {
		keyValues = append(keyValues, args[i], args[i+1])
	}

	slog.SetDefault(logger)

	argsGroup := slog.Group(
		"args",
		keyValues...,
	)

	slog.Info(
		context,
		argsGroup,
	)
}

func LogError(
	logger *slog.Logger,
	err error,
	errDetail string,
	context string,
	args ...string) *CommonError {

	if err == nil || logger == nil {

		// no point to continue the procedures, hence direct exit
		return nil
	}

	var returnedError CommonError
	var argsGroup slog.Attr

	slog.SetDefault(logger)
	identifier := generateRandomString(10)

	// process args
	if len(args)%2 == 0 {
		// only process args with correct modulo
		var keyValues []interface{}
		for i := 0; i < len(args); i += 2 {
			keyValues = append(keyValues, args[i], args[i+1])
		}

		argsGroup = slog.Group(
			"args",
			keyValues...,
		)
	} else {

		// will notify user the absence of args
		LogInfo(logger, "logError.InvalidArgs", "args", "null")
	}

	// retrieve logsource
	logsource := ""
	re := regexp.MustCompile(`\((.*?)\)`)
	match := re.FindStringSubmatch(err.Error())

	if len(match) > 1 {
		logsource = match[1]
	}

	// compose error group
	errorMessageSplited := strings.Split(err.Error(), ";")

	var errorGroup slog.Attr
	if len(errorMessageSplited) > 1 {
		// Error with code, used by modules.
		// contains indentifier for faster log tracing.

		errorGroup = slog.Group(
			"error",
			"logsource", logsource,
			"code", strings.Split(err.Error(), ";")[0],
			"message", strings.Split(err.Error(), ";")[1],
			"detail", errDetail,
			"context", context,
			"identifier", identifier,
		)

		returnedError = CommonError{
			Id:          identifier,
			Context:     context,
			Err:         err,
			ErrResponse: errors.New(err.Error() + ";" + identifier),
		}
	} else {
		// For error without code, used by dependencies.

		errorGroup = slog.Group(
			"error",
			"logsource", logsource,
			"message", err.Error(),
			"detail", errDetail,
			"context", context,
		)

		returnedError = CommonError{
			Id:      identifier,
			Context: context,
			Err:     err,
		}
	}

	slog.Error(
		err.Error(),
		argsGroup,
		errorGroup,
	)

	return &returnedError
}
