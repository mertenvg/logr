package logr

func M(key string, value any) Meta {
	return Meta{}.With(key, value)
}

type Meta map[string]any

func (m Meta) With(key string, value any) Meta {
	n := m.Copy()
	n[key] = value
	return n
}

func (m Meta) Copy() Meta {
	meta := Meta{}
	for k, v := range m {
		meta[k] = v
	}
	return meta
}
