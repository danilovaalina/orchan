package orchan

import (
	"testing"
	"time"
)

// Or с одним уже закрытым каналом должен немедленно закрыть выходной канал.
func TestOr_Single(t *testing.T) {
	t.Parallel()

	in := make(chan interface{})
	close(in)

	result := Or(in)

	select {
	case <-result:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Or did not close the output channel when the input channel was already closed")
	}
}

// Or с пустым списком аргументов должен возвращать немедленно закрытый канал.
func TestOr_Empty(t *testing.T) {
	t.Parallel()

	resultChan := Or()

	_, ok := <-resultChan
	if ok {
		t.Error("Or with no channels should return a closed channel immediately")
	}
}

// Or с несколькими каналами должен закрыть выходной канал при закрытии любого из входных.
func TestOr_MultipleChannels(t *testing.T) {
	t.Parallel()

	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})

	result := Or(ch1, ch2, ch3)

	// Закрываем только второй канал
	close(ch2)

	select {
	case <-result:
		// Or отреагировал на закрытие ch2
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Or did not close the output channel when one of the inputs was closed")
	}
}

// Or при одновременном закрытии нескольких каналов не должен паниковать.
func TestOr_RaceCondition(t *testing.T) {
	t.Parallel()

	ch1 := make(chan interface{})
	ch2 := make(chan interface{})

	result := Or(ch1, ch2)

	// Закрываем оба канала почти одновременно в разных горутинах
	go close(ch1)
	go close(ch2)

	<-result
}

// Or при передаче одного и того же канала несколько раз не должен паниковать.
func TestOr_DuplicateChannels(t *testing.T) {
	t.Parallel()

	ch := make(chan interface{})

	result := Or(ch, ch)

	close(ch)

	select {
	case <-result:
		// Канал закрылся без паники
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Or did not close the output channel within timeout")
	}
}
