package response

import (
	"errors"
	"fmt"
	"io"

	"github.com/pavychsasha/httpfromtcp/internal/headers"
)

type WriterState int

const (
	WriterInitialized WriterState = iota
	WriterStatusLineInitialized
	WriterHeadersInitialized
	WriterBodyInitialized
)

type Writer struct {
	WriterState WriterState
	Writer      io.Writer
}

func (w *Writer) Write(p []byte) (int, error) {
	n, err := w.Writer.Write(p)
	if err != nil {
		return n, err
	}
	if n != len(p) {
		return n, io.ErrShortWrite
	}
	return len(p), nil
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.WriterState > WriterStatusLineInitialized {
		return errors.New("Writer's state should be initialized first.")
	}

	writerStatusLine := GetWriterStatusLine(statusCode)
	_, err := fmt.Fprint(w, writerStatusLine)
	if err != nil {
		return err
	}

	w.WriterState = WriterStatusLineInitialized
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	if w.WriterState > WriterHeadersInitialized {
		return errors.New("Writer's needs write headers afer status line initialization.")
	}

	for key, value := range headers {
		_, err := fmt.Fprintf(w, "%s:%s%s", key, value, CRLF)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "%s", CRLF)
	if err != nil {
		return err
	}

	w.WriterState = WriterHeadersInitialized
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {

	n, err := fmt.Fprintf(w, "%s%s", p, CRLF)
	if err != nil {
		return n, err
	}

	w.WriterState = WriterBodyInitialized
	return n, nil
}
