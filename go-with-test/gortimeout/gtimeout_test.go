package gortimeout

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTimeout(t *testing.T) {
	err := timeout(doSomething)
	assert.Error(t, err)
	//assert.Equal(t, nil, err)
}
