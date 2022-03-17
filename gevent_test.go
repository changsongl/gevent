package gevent

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

type test struct {
	data int
}

func TestGEventTriggerSyncEvent(t *testing.T) {
	t.Log("Start TestGEventTriggerSyncEvent")
	e := NewGEvent()

	event := "test-event"
	obs1, obs2 := "test-observer1", "test-observer2"

	count := 0

	err := e.AddObserver(event, obs1, func(a test, b *test) {
		count += a.data + b.data

	}, false)

	require.NoErrorf(t, err, "add observer error: %s", obs1)

	err = e.AddObserver(event, "test-observer2", func(a test, b *test) {
		count += 2 * (a.data + b.data)
	}, false)

	require.NoErrorf(t, err, "add observer error: %s", obs2)

	a := test{1}
	b := &test{2}
	e.TriggerEvent(event, a, b)
	require.Equal(t, count, 3*(a.data+b.data), "count not equal")
}

func TestGEventTriggerAsyncEvent(t *testing.T) {
	t.Log("Start TestGEventTriggerAsyncEvent")
	e := NewGEvent()

	event := "test-event"
	obs1, obs2 := "test-observer1", "test-observer2"

	count := 0
	wg := sync.WaitGroup{}
	wg.Add(2)

	err := e.AddObserver(event, obs1, func(a test, b *test) {
		count += a.data + b.data
		wg.Done()
	}, true)

	require.NoErrorf(t, err, "add observer error: %s", obs1)

	err = e.AddObserver(event, "test-observer2", func(a test, b *test) {
		count += 2 * (a.data + b.data)
		wg.Done()
	}, true)

	require.NoErrorf(t, err, "add observer error: %s", obs2)

	a := test{1}
	b := &test{2}
	e.TriggerEvent(event, a, b)

	wg.Wait()
	require.Equal(t, count, 3*(a.data+b.data), "count not equal")
}

func TestGEventRemoveObserver(t *testing.T) {
	t.Log("Start TestGEventRemoveObserver")
	e := NewGEvent()

	event := "test-event"
	obs1, obs2 := "test-observer1", "test-observer2"

	count := 0

	err := e.AddObserver(event, obs1, func(a test, b *test) {
		count += a.data + b.data

	}, false)

	require.NoErrorf(t, err, "add observer error: %s", obs1)

	err = e.AddObserver(event, "test-observer2", func(a test, b *test) {
		count += 2 * (a.data + b.data)
	}, false)

	require.NoErrorf(t, err, "add observer error: %s", obs2)

	result := e.RemoveObserver(event, obs1)
	require.Equal(t, true, result, "remove observer first")

	result = e.RemoveObserver(event, obs1)
	require.Equal(t, false, result, "remove observer second")

	a := test{1}
	b := &test{2}
	e.TriggerEvent(event, a, b)
	require.Equal(t, count, 2*(a.data+b.data), "count not equal")
}

func TestGEventNotFunction(t *testing.T) {
	t.Log("Start TestGEventRemoveObserver")
	e := NewGEvent()
	err := e.AddObserver("test", "test", 8, false)

	require.NotNil(t, err, "non-function error is nil")
}

var testLogCount = 0

type testLog struct {
}

func (t testLog) Error(msg string) {
	testLogCount++
}

func TestGEventLogWhenPanic(t *testing.T) {
	t.Log("Start TestGEventLogWhenPanic")
	e := NewGEvent(NewLogOption(&testLog{}))

	event := "test-event"
	obs1 := "test-observer1"

	err := e.AddObserver(event, obs1, func(a test, b *test) {}, false)
	require.Nil(t, err, "add observer error is not nil")

	require.Equal(t, testLogCount, 0)
	e.TriggerEvent(event, obs1)
	require.Equal(t, testLogCount, 1)
}
