package encoder

import (
	"bytes"
	"compress/gzip"
	"io"
)

var SupportedEncoding = map[string]bool{
	"gzip": true,
}

func EncoderFactory(encoding string) Encoder {
	supported, ok := SupportedEncoding[encoding]
	if ok && supported {
		switch encoding {
		case "gzip":
			return GZipEncoder{name: "gzip"}
		}
	}
	return nil
}

type Encoder interface {
	Encode([]byte) (io.Reader, error)
	Decode([]byte) (io.Reader, error)
}

type GZipEncoder struct {
	name string
}

func (gz GZipEncoder) Encode(p []byte) (io.Reader, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	defer zw.Close()
	_, err := zw.Write(p)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

func (gz GZipEncoder) Decode(p []byte) (io.Reader, error) {
	buf := bytes.NewBuffer(p)
	zr, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	return zr, nil
}
