package main

import (
	"testing"

	"github.com/DoubleWB/website/signatures"

	"github.com/stretchr/testify/assert"
)

const TESTING_STORAGE = "test"
const NAME1 = "alice"
const NAME2 = "bob"
const BAD_NAME1 = ""

func TestValidSignature(t *testing.T) {
	s := signatures.NewSignature(NAME1)
	currentSignatures = []signatures.Signature{s}
	assert.False(t, validSignature(NAME1))
	assert.False(t, validSignature(BAD_NAME1))
	assert.True(t, validSignature(NAME2))
}
