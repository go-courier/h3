package h3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getRes0Indexes(t *testing.T) {
	count := res0IndexCount()
	indexes := make([]H3Index, count)
	getRes0Indexes(indexes)

	require.True(t, indexes[0] == 0x8001fffffffffff, "correct first basecell")
	require.True(t, indexes[121] == 0x80f3fffffffffff, "correct last basecell")
}
