package anh

import "context"

type Collection[T any] struct {
	items         []T
	fetchToken    string
	fetchNextPage func(ctx context.Context, token string) ([]T, string, error)
}

func (c *Collection[T]) HasItems() bool {
	return c.items != nil && len(c.items) > 0
}

func (c *Collection[T]) Items() []T {
	if c.items != nil {
		return c.items
	}
	return []T{}
}

func (c *Collection[T]) NextPage(ctx context.Context) error {
	var err error

	if c.fetchToken == "" {
		c.items = nil
		return nil
	}

	if c.fetchNextPage == nil {
		return nil
	}

	c.items, c.fetchToken, err = c.fetchNextPage(ctx, c.fetchToken)
	if err != nil {
		return err
	}

	return nil
}
