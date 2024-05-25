package http

type Status int

const (
	StatusOk                  Status = 200
	StatusCreated             Status = 201
	StatusBadRequest          Status = 400
	StatusNotFound            Status = 404
	StatusInternalServerError Status = 500
)

var StatusString map[Status]string = map[Status]string{
	StatusOk:                  "OK",
	StatusNotFound:            "Not Found",
	StatusBadRequest:          "Bad Request",
	StatusInternalServerError: "Internal Server Error",
	StatusCreated:             "Created",
}
