package cli

import (
	"errors"

	"github.com/tamnd/v2ex-cli/v2ex"
)

func isNotFound(err error) bool {
	return errors.Is(err, v2ex.ErrNotFound)
}
