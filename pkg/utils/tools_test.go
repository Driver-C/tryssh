package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveDuplicate_NormalCase(t *testing.T) {
	input := []string{"a", "b", "a", "c", "b", "d"}
	result := RemoveDuplicate(input)
	assert.Equal(t, []string{"a", "b", "c", "d"}, result)
}

func TestRemoveDuplicate_EmptySlice(t *testing.T) {
	result := RemoveDuplicate([]string{})
	assert.Equal(t, []string{}, result)
}

func TestRemoveDuplicate_AllDuplicates(t *testing.T) {
	input := []string{"x", "x", "x", "x"}
	result := RemoveDuplicate(input)
	assert.Equal(t, []string{"x"}, result)
}

func TestRemoveDuplicate_NoDuplicates(t *testing.T) {
	input := []string{"a", "b", "c", "d"}
	result := RemoveDuplicate(input)
	assert.Equal(t, []string{"a", "b", "c", "d"}, result)
}

func TestToInterfaceSlice_NilInput(t *testing.T) {
	result := ToInterfaceSlice[int](nil)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
	// Should be a non-nil empty slice
	assert.Equal(t, []interface{}{}, result)
}

func TestToInterfaceSlice_EmptySlice(t *testing.T) {
	result := ToInterfaceSlice[string]([]string{})
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
}

func TestToInterfaceSlice_Strings(t *testing.T) {
	input := []string{"hello", "world"}
	result := ToInterfaceSlice(input)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "hello", result[0])
	assert.Equal(t, "world", result[1])
}

func TestToInterfaceSlice_Ints(t *testing.T) {
	input := []int{1, 2, 3}
	result := ToInterfaceSlice(input)
	assert.Equal(t, 3, len(result))
	assert.Equal(t, 1, result[0])
	assert.Equal(t, 2, result[1])
	assert.Equal(t, 3, result[2])
}

func TestToInterfaceSlice_SingleElement(t *testing.T) {
	input := []float64{3.14}
	result := ToInterfaceSlice(input)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, 3.14, result[0])
}
