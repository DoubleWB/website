package signatures

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const TESTING_STORAGE = "test"
const NAME1 = "alice"
const NAME2 = "bob"
const NAME3 = "charlie"

func setup() {
	os.Create(TESTING_STORAGE)
}

func cleanup() {
	os.Remove(TESTING_STORAGE)
}

func TestNewSignature(t *testing.T) {
	now := time.Now()
	s := NewSignature(NAME1)
	assert.Equal(t, s.Name, NAME1, "name is not properly instantiated")
	assert.True(t, s.CreatedAt.After(now))
}

func TestReadAndWriteFile(t *testing.T) {
	setup()
	s := NewSignature(NAME1)
	assert.Nil(t, s.WriteToFile(TESTING_STORAGE))
	s1 := NewSignature(NAME2)
	assert.Nil(t, s1.WriteToFile(TESTING_STORAGE))
	s2, err := WriteNewSignature(NAME3, TESTING_STORAGE)
	assert.Nil(t, err)
	sArr, err := ReadFromFile(TESTING_STORAGE)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(sArr))
	assert.Equal(t, s.Name, sArr[0].Name)
	assert.Equal(t, s.CreatedAt.Location(), sArr[1].CreatedAt.Location())
	assert.Equal(t, s1.Name, sArr[1].Name)
	assert.Equal(t, s1.CreatedAt.Location(), sArr[1].CreatedAt.Location())
	assert.Equal(t, s2.Name, sArr[2].Name)
	assert.Equal(t, s2.CreatedAt.Location(), sArr[2].CreatedAt.Location())
	cleanup()
}

func TestSaveOverFile(t *testing.T) {
	setup()
	WriteNewSignature(NAME1, TESTING_STORAGE)
	sArr, err := ReadFromFile(TESTING_STORAGE)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(sArr))
	assert.Equal(t, NAME1, sArr[0].Name)
	newSArr := []Signature{NewSignature(NAME2), NewSignature(NAME3)}
	err = SaveOverFile(newSArr, TESTING_STORAGE)
	sArr, err = ReadFromFile(TESTING_STORAGE)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(sArr))
	assert.Equal(t, NAME2, sArr[0].Name)
	assert.Equal(t, NAME3, sArr[1].Name)
	cleanup()
}
