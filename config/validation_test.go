package config_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"

	"phabricator.wikimedia.org/source/blubber/config"
)

func TestIsValidationError(t *testing.T) {
	assert.False(t, config.IsValidationError(nil))
	assert.False(t, config.IsValidationError(errors.New("foo")))
	assert.True(t, config.IsValidationError(validator.ValidationErrors{}))
}
