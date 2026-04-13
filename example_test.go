package orchan_test

import (
	"fmt"
	"time"

	"github.com/danilovaalina/orchan"
)

// ExampleOr показывает, как объединить несколько каналов.
// Функция закроется, как только сработает самый быстрый канал (1 секунда).
func ExampleOr() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-orchan.Or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	// Округляем время, чтобы тест всегда проходил успешно
	duration := time.Since(start).Round(time.Second)
	fmt.Printf("done after %v", duration)

	// Output:
	// done after 1s
}
