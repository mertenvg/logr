package logr

import (
	"fmt"
	"os"
)

// listen concurrently works through the buffered messages channel
func listen(ms <-chan *Message) {
	defer close(listenerDone)
	for {
		select {

		case wm := <-addWriter:
			writers[wm.c] = wm.w

		case wm := <-removeWriter:
			delete(writers, wm.c)

		case m, ok := <-ms:
			if !ok {
				return
			}
			if m.Type != None {
				for c, w := range writers {
					if m.Type&c.filter != m.Type {
						continue
					}
					_, err := w.Write(c.format(m))
					if err != nil {
						fmt.Fprintf(os.Stderr, "logr: failed to write message to Writer: %v\n", err)
					}
				}
			}

			close(m.done)

			m.Reset()
			pool.Put(m)

		}
	}
}
