package kvp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvOS(t *testing.T) {
	in := []string{"foo=bar", "BAZ=BUM"}

	res := envVarsFromOS(in)

	expected := envvars(map[string]string{
		"foo": "bar",
		"BAZ": "BUM",
	})

	assert.Equal(t, res, expected, "These maps should be equal")
}

func TestEnvPairs(t *testing.T) {

}
