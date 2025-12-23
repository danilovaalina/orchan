package orchan

import "sync"

// Or combines one or more done channels into a single done channel.
// The returned channel is closed as soon as any of the input channels is closed
// or receives a value.
func Or(channels ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	if len(channels) == 0 {
		close(out)
		return out
	}

	var once sync.Once
	closeOut := func() {
		once.Do(func() {
			close(out)
		})
	}

	for _, c := range channels {
		go func() {
			select {
			case <-c:
				closeOut()
			}
		}()
	}

	return out
}
