package anh

import "context"

type PagedCollection[T any] struct {
	items         []T
	fetchToken    string
	fetchNextPage func(ctx context.Context, token string, count int) ([]T, string, error)
	pageSize      int
}

func (c *PagedCollection[T]) HasItems() bool {
	return c.items != nil && len(c.items) > 0
}

func (c *PagedCollection[T]) Items() []T {
	if c.items != nil {
		return c.items
	}
	return []T{}
}

func (c *PagedCollection[T]) NextPage(ctx context.Context) error {
	var err error

	if c.fetchToken == "" {
		c.items = nil
		return nil
	}

	if c.fetchNextPage == nil {
		return nil
	}

	c.items, c.fetchToken, err = c.fetchNextPage(ctx, c.fetchToken, c.pageSize)
	if err != nil {
		return err
	}

	return nil
}
