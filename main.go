package commonlogger

import (
	"errors"
	"log/slog"
	"regexp"
	"strings"

	"github.com/lefalya/commonlogger/schema"

	"math/rand"
)

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
	args ...string) *schema.CommonError {

	slog.SetDefault(logger)
	identifier := generateRandomString(10)

	if len(args)%2 != 0 {

		LogInfo(logger, "logError.InvalidArgs", "args", "null")
	} else {

		logsource := ""
		re := regexp.MustCompile(`\((.*?)\)`)
		match := re.FindStringSubmatch(err.Error())

		if len(match) > 1 {
			logsource = match[1]

		}

		var keyValues []interface{}
		for i := 0; i < len(args); i += 2 {
			keyValues = append(keyValues, args[i], args[i+1])
		}

		argsGroup := slog.Group(
			"args",
			keyValues...,
		)

		errorGroup := slog.Group(
			"error",
			"logsource", logsource,
			"code", strings.Split(err.Error(), ";")[0],
			"message", strings.Split(err.Error(), ";")[1],
			"detail", errDetail,
			"context", context,
			"identifier", identifier,
		)

		slog.Error(
			err.Error(),
			argsGroup,
			errorGroup,
		)

	}

	errorResponse := errors.New(err.Error() + " - ID: " + identifier)

	return &schema.CommonError{
		Id:          identifier,
		Context:     context,
		Err:         err,
		ErrResponse: errorResponse,
	}
}
