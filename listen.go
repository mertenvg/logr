package logr

// listen concurrently works through the buffered messages channel
func listen(ms <-chan *Message) {
	for {
		select {

		case wm := <-addWriter:
			writers[wm.c] = wm.w

		case wm := <-removeWriter:
			delete(writers, wm.c)

		case m := <-ms:
			for c, w := range writers {
				if m.Type&c.filter != m.Type {
					continue
				}
				_, err := w.Write(c.format(m))
				if err != nil {
					Errorf("failed to write message to Writer: %v", err)
				}
			}

			close(m.done)

			m.Reset()
			pool.Put(m)

		}
	}
}
