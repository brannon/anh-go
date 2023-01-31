package anh

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PagedCollection_HasItems(t *testing.T) {
	c := PagedCollection[int]{}
	assert.False(t, c.HasItems())

	c = PagedCollection[int]{items: []int{}}
	assert.False(t, c.HasItems())

	c = PagedCollection[int]{items: []int{1}}
	assert.True(t, c.HasItems())
}

func Test_PagedCollection_Items(t *testing.T) {
	c := PagedCollection[int]{}
	assert.Equal(t, []int{}, c.Items())

	c = PagedCollection[int]{items: []int{}}
	assert.Equal(t, []int{}, c.Items())

	c = PagedCollection[int]{items: []int{1}}
	assert.Equal(t, []int{1}, c.Items())
}

func Test_PagedCollection_NextPage(t *testing.T) {
	var err error

	var tokenParam string
	var fetchNextPageCallCount int

	c := PagedCollection[int]{
		fetchToken: "token",
		fetchNextPage: func(ctx context.Context, token string, count int) ([]int, string, error) {
			fetchNextPageCallCount++
			tokenParam = token
			return []int{1, 2, 3}, "", nil
		},
	}
	err = c.NextPage(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, c.Items())

	err = c.NextPage(context.Background())
	assert.NoError(t, err)
	assert.False(t, c.HasItems())

	assert.Equal(t, "token", tokenParam)
	assert.Equal(t, 1, fetchNextPageCallCount)
}
