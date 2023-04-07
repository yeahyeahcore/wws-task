package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yeahyeahcore/wws-task/internal/utils"
)

func TestGenerateSignature(t *testing.T) {
	signature, timestamp := utils.GenerateSignature("secret")

	assert.NotEmpty(t, signature)
	assert.NotEmpty(t, timestamp)
}
