package http

type path string

type Mux struct {
	handlers map[path]func(ResponseWriter, *Request)
}

func (m *Mux) Handle(rw ResponseWriter, r *Request) {
	handleFunc, ok := m.handlers[path(r.URL.Path)]
	if !ok {
		rw.WriteStatus(StatusNotFound)
		rw.Write([]byte{})
		return
	}
	handleFunc(rw, r)
}

func (m *Mux) HandleFunc(endPath string, handler func(ResponseWriter, *Request)) {
	m.handlers[path(endPath)] = handler
}

func NewMux() *Mux {
	handlers := make(map[path]func(ResponseWriter, *Request))
	return &Mux{handlers: handlers}
}
