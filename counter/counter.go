package counter

// counter is a thread-safe counter with a maximum value and an increment step.
// It's used to safely increment a counter value in a multi-threaded environment.
type counter struct {
	val         int
	maxVal      int
	step        int
	incremented chan bool
	committed   chan int
	done        chan struct{}
}

// Increment increments counter by its step value. Returns counter value.
func (c *counter) Increment() int {
	c.incremented <- true
	return <-c.committed
}

// Val returns current counter value.
func (c *counter) Val() int {
	c.incremented <- false
	return <-c.committed
}

// Close stops the counter goroutine and closes all channels.
// This should be used when you don't need the counter anymore.
func (c *counter) Close() {
	close(c.done)
	close(c.committed)
	close(c.incremented)
}

func New(maxVal, step int) *counter {
	c := &counter{
		val:         0,
		maxVal:      maxVal,
		step:        step,
		incremented: make(chan bool),
		committed:   make(chan int),
		done:        make(chan struct{}),
	}

	// Start a separate goroutine to handle incrementing the counter value
	// and sending the updated value on the committed channel.
	go func() {
		for ok := range c.incremented {
			select {
			case <-c.done:
				return
			default:
				if ok && c.val+c.step <= c.maxVal {
					c.val += c.step
				}
				c.committed <- c.val
			}
		}
	}()

	return c
}
