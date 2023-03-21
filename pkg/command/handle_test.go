package command

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrettyError(t *testing.T) {
	err := fmt.Errorf("test")
	assert.ErrorIs(t, err, PrettyError(err))
}
