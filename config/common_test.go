package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"phabricator.wikimedia.org/source/blubber/config"
)

func TestCommonConfig(t *testing.T) {
	cfg, err := config.ReadConfig([]byte(`---
    base: fooimage
    sharedvolume: true
    entrypoint: ["/bin/foo"]
    variants:
      build: {}`))

	assert.Nil(t, err)

	assert.Equal(t, "fooimage", cfg.Base)
	assert.Equal(t, true, cfg.SharedVolume.True)
	assert.Equal(t, []string{"/bin/foo"}, cfg.EntryPoint)

	variant, err := config.ExpandVariant(cfg, "build")

	assert.Equal(t, "fooimage", variant.Base)
	assert.Equal(t, true, variant.SharedVolume.True)
	assert.Equal(t, []string{"/bin/foo"}, variant.EntryPoint)
}