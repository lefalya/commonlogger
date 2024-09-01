package commonlogger

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/lefalya/commonlogger/interfaces"
	"github.com/lefalya/commonlogger/schema"
	"github.com/stretchr/testify/assert"
)

func TestInfoLog(t *testing.T) {

	t.Run("successfull print loginfo", func(t *testing.T) {

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		context := "submission-created"
		dummyCreatorUUID := uuid.New().String()
		dummyCampaignUUID := uuid.New().String()

		LogInfo(
			logger,
			context,
			"campaignCreatorUUID",
			dummyCreatorUUID,
			"campaignUUID",
			dummyCampaignUUID)
	})

	t.Run("invalid args modulo", func(t *testing.T) {

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		context := "submission-created"
		dummyCreatorUUID := uuid.New().String()

		LogInfo(
			logger,
			context,
			"campaignCreatorUUID",
			dummyCreatorUUID,
			"campaignUUID")
	})
}

func TestErrorLog(t *testing.T) {

	t.Run("successfull print logerror", func(t *testing.T) {

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		DUMMY_ERROR := errors.New("10000;(collection) MySQL fatal error")
		details := "errcon HY2000 mysql host not found!"
		context := "AddCampaign.MYSQL_FATAL_ERROR"

		resultErr := LogError(
			logger,
			DUMMY_ERROR,
			details,
			context,
			"collectionUUID", uuid.New().String(),
			"campaignUUID", uuid.New().String(),
		)

		assert.NotNil(t, resultErr)
		assert.NotNil(t, resultErr.Id)
		assert.NotNil(t, resultErr.Context)
		assert.NotNil(t, resultErr.Err)
		assert.NotNil(t, resultErr.ErrResponse)

		assert.Equal(t, DUMMY_ERROR, resultErr.Err)
		assert.Equal(t, 10, len(resultErr.Id))

		alphanumericRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)

		assert.True(t, alphanumericRegex.MatchString(resultErr.Id))
	})

	t.Run("invalid args modulo", func(t *testing.T) {

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		DUMMY_ERROR := errors.New("10000;(collection) MySQL fatal error")
		details := "errcon HY2000 mysql host not found!"
		context := "AddCampaign.MYSQL_FATAL_ERROR"

		resultErr := LogError(
			logger,
			DUMMY_ERROR,
			details,
			context,
			"collectionUUID", uuid.New().String(),
			"campaignUUID",
		)

		assert.NotNil(t, resultErr)
		assert.NotNil(t, resultErr.Id)
		assert.NotNil(t, resultErr.Context)
		assert.NotNil(t, resultErr.Err)
		assert.NotNil(t, resultErr.ErrResponse)

		fmt.Println(resultErr.ErrResponse)

		assert.Equal(t, DUMMY_ERROR, resultErr.Err)
		assert.Equal(t, 10, len(resultErr.Id))

		alphanumericRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)

		assert.True(t, alphanumericRegex.MatchString(resultErr.Id))
	})

	t.Run("successfull return logerror but invalid logsource", func(t *testing.T) {

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		DUMMY_ERROR := errors.New("10000;(collection MySQL fatal error")
		details := "errcon HY2000 mysql host not found!"
		context := "AddCampaign.MYSQL_FATAL_ERROR"

		resultErr := LogError(
			logger,
			DUMMY_ERROR,
			details,
			context,
			"collectionUUID", uuid.New().String(),
			"campaignUUID", uuid.New().String(),
		)

		assert.NotNil(t, resultErr)
		assert.NotNil(t, resultErr.Id)
		assert.NotNil(t, resultErr.Context)
		assert.NotNil(t, resultErr.Err)
		assert.NotNil(t, resultErr.ErrResponse)

		assert.Equal(t, DUMMY_ERROR, resultErr.Err)
		assert.Equal(t, 10, len(resultErr.Id))

		alphanumericRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)

		assert.True(t, alphanumericRegex.MatchString(resultErr.Id))
	})
}

type Submission struct {
	UUID    string
	Caption string
}

func LogErrorSubmission(
	logger *slog.Logger,
	errorSubject error,
	errorDetail string,
	context string,
	submission any,
) *schema.CommonError {

	submissionAs := submission.(Submission)

	return LogError(
		logger,
		errorSubject,
		errorDetail,
		context,
		"uuid", submissionAs.UUID,
		"caption", submissionAs.Caption,
	)
}

func SetRedis(
	submission interface{},
	loghelper interfaces.LogHelper,
) *schema.CommonError {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return loghelper(
		logger,
		errors.New("10000;(commonlogger) test error"),
		"error detail",
		"TestParticipantSetRedis.Error",
		submission,
	)
}

func TestCallbackLogHelper(t *testing.T) {

	dummySubmission := Submission{
		UUID:    uuid.New().String(),
		Caption: "dummy caption",
	}

	errorResult := SetRedis(
		dummySubmission,
		LogErrorSubmission,
	)

	assert.NotNil(t, errorResult)
}
