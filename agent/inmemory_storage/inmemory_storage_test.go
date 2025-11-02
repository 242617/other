package inmemory_storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/242617/other/agent/inmemory_storage"
)

func TestInmemoryStorage(t *testing.T) {
	s := inmemory_storage.New(1, 1)
	require.NoError(t, s.Append("system"))
	require.NoError(t, s.Append("first message"))
	require.NoError(t, s.Append("second message"))

	list, err := s.List()
	require.NoError(t, err, "list")
	assert.Len(t, list, 2, "unexpected list length")
	assert.Equal(t, "system", list[0], "unexpected first element")
	assert.Equal(t, "second message", list[1], "unexpected second element")
}
