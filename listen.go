package logr

import (
	"fmt"
	"os"
)

// listen concurrently works through the buffered messages channel
func listen(ms <-chan *Message) {
	defer close(listenerDone)
	for m := range ms {
		if m.Type != None {
			mutex.RLock()
			for c, w := range writers {
				if m.Type&c.filter != m.Type {
					continue
				}
				_, err := w.Write(c.format(m))
				if err != nil {
					fmt.Fprintf(os.Stderr, "logr: failed to write message to Writer: %v\n", err)
				}
			}
			mutex.RUnlock()
		}

		close(m.done)

		m.Reset()
		pool.Put(m)
	}
}
