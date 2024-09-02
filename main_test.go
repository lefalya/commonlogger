package commonlogger

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func TestInfoLog(t *testing.T) {

	t.Run("successfull print loginfo", func(t *testing.T) {

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

	t.Run("error message with no error code", func(t *testing.T) {

		DUMMY_ERROR := errors.New("(commonpagination) Redis fatal error!")
		details := "redis timeout"
		context := "commonpagination.redis_fatal_error"

		resultErr := LogError(
			logger,
			DUMMY_ERROR,
			details,
			context,
			"collectionUUID", uuid.New().String(),
			"campaignUUID", uuid.New().String(),
		)

		assert.NotNil(t, resultErr)
	})

	t.Run("invalid args modulo", func(t *testing.T) {

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

	t.Run("no error present", func(t *testing.T) {

		details := "errcon HY2000 mysql host not found!"
		context := "AddCampaign.MYSQL_FATAL_ERROR"

		resultErr := LogError(
			logger,
			nil,
			details,
			context,
			"collectionUUID", uuid.New().String(),
			"campaignUUID", uuid.New().String(),
		)

		assert.Nil(t, resultErr)
	})
}
