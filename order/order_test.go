package order

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsServiceOrderStatusValid(t *testing.T) {
	c := require.New(t)

	c.True(isServiceOrderStatusValid(ServiceOrderStatus("pending")))
	c.False(isServiceOrderStatusValid(ServiceOrderStatus("anotherstatus")))
}
