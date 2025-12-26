package orchan

// Or combines one or more done channels into a single done channel.
// The returned channel is closed as soon as any of the input channels is closed.
func Or(channels ...<-chan interface{}) <-chan interface{} {
	orDone := make(chan interface{})

	if len(channels) == 0 {
		close(orDone)
		return orDone
	}

	if len(channels) == 1 {
		return channels[0]
	}

	go func() {
		defer close(orDone)

		if len(channels) == 2 {
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
			return
		}

		select {
		case <-channels[0]:
		case <-channels[1]:
		case <-Or(append(channels[2:], orDone)...):
		}
	}()
	return orDone
}
