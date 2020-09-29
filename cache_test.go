package cache

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func parse(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}

func TestGetEmpty(t *testing.T) {
	storage := NewStorage()
	content := storage.Get("MY_KEY")
	require.Nil(t, content)
}

func TestGetEmptyCast(t *testing.T) {
	storage := NewStorage()
	content, ok := storage.Get("MY_KEY").(*testing.T)
	require.False(t, ok)
	require.Nil(t, content)
}

func TestKeys(t *testing.T) {
	storage := Storage{
		items: map[string]Item{
			"aaabbbccc":  {},
			"aaa0123bbb": {},
			"0123aaabbb": {},
		},
		mtx: &sync.RWMutex{},
	}
	require.Equal(t, 0, len(storage.Keys("^[0-9]+$")))
	require.Equal(t, 1, len(storage.Keys("^[a-zA-Z]+$")))
	require.Equal(t, 1, len(storage.Keys("^[a-zA-Z]+[0-9]+[a-zA-Z]+$")))
	require.Equal(t, 1, len(storage.Keys("^[0-9]+[a-zA-Z]+$")))
}

func TestGetValue(t *testing.T) {
	storage := NewStorage()
	storage.Set("MY_KEY", []byte("123456"), parse("5s"))
	content, ok := storage.Get("MY_KEY").([]byte)
	require.True(t, ok)
	require.Equal(t, content, []byte("123456"))
}

func TestGetExpiredValue(t *testing.T) {
	storage := NewStorage()
	storage.Set("MY_KEY", []byte("123456"), parse("1s"))
	time.Sleep(parse("1s200ms"))
	content := storage.Get("MY_KEY")
	require.Nil(t, content)
}

func TestStorage_Flush(t *testing.T) {
	storage := NewStorage()
	storage.Set("MY_KEY", []byte("123456"), parse("1s"))
	require.Equal(t, 1, len(storage.items))
	storage.Flush("")
	require.Equal(t, 0, len(storage.items))
}

func TestStorage_Rem(t *testing.T) {
	storage := NewStorage()
	storage.Set("MY_KEY", "123456", parse("1s"))
	require.Equal(t, 1, len(storage.items))
	v, ok := storage.Rem("MY_KEY").(string)
	require.True(t, ok)
	require.Equal(t, "123456", v)
	v, ok = storage.Rem("NOTHING").(string)
	require.False(t, ok)
	require.Equal(t, string(""), v)
}
