package http

type Status int

const (
	StatusOk       Status = 200
	StatusNotFound Status = 404
)

var StatusString map[Status]string = map[Status]string{
	StatusOk:       "OK",
	StatusNotFound: "Not Found",
}
