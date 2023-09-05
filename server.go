package flies

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
)

const (
	defaultLogFormat          = "pretty"
	defaultPort               = "8080"
	defaultResponseStatus     = "421 Misdirected Request"
	defaultResponseStatusCode = 421
)

// Server is the TCP Server used by [flies].
type Server struct {
	Port                string
	ResponseStatus      string
	ResponseStatusCode  int
	ResponseBodyContent string
	ErrWriter           io.Writer
	RawWriter           io.Writer
	ReqWriter           RequestWriter
}

// Listen writes any data received on the TCP connection to RawWriter, and
// writes any errors (followed by a newline) to ErrWriter. The request is
// parsed, the [http.Request.Body] is read and transformed into a [BodyReader],
// and the request is passed to ReqWriter. A response is provided to the client
// and the connection is closed.
func (s *Server) Listen() error {

	l, err := net.Listen("tcp", s.getAddr())
	if err != nil {
		return err
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			s.writeError(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go s.handle(conn)
	}
	return err
}

func (s *Server) handle(c net.Conn) {
	raw := s.getRawWriter()
	tee := io.TeeReader(c, raw)
	req, err := http.ReadRequest(bufio.NewReader(tee))
	if err != nil {
		s.writeError(err)
	}
	// Allow reading the body multiple times
	req.Body = NewBodyReader(req.Body)
	// Read the body once to buffer it, then respond and close the
	// connection. Otherwise we won't know when to close the connection
	s.respond(c)
	raw.Write([]byte("\r\n"))
	c.Close()
	s.getReqWriter().WriteRequest(req)
}

func (s *Server) getRawWriter() io.Writer {
	if s.RawWriter == nil {
		return io.Discard
	}
	return s.RawWriter
}

func (s *Server) getReqWriter() RequestWriter {
	if s.ReqWriter == nil {
		return &RequestWriterDiscard{}
	}
	return s.ReqWriter
}

func (s *Server) respond(c net.Conn) {
	responseStatus := s.ResponseStatus
	if responseStatus == "" {
		responseStatus = defaultResponseStatus
	}
	responseStatusCode := s.ResponseStatusCode
	if responseStatusCode == 0 {
		responseStatusCode = defaultResponseStatusCode
	}
	r := &http.Response{
		Status:        responseStatus,
		StatusCode:    responseStatusCode,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          io.NopCloser(bytes.NewReader([]byte(s.ResponseBodyContent))),
		ContentLength: int64(len(s.ResponseBodyContent)),
	}
	r.Write(c)
}

func (s *Server) writeError(err error) error {
	s.ErrWriter.Write([]byte(err.Error()))
	s.ErrWriter.Write([]byte("\n"))
	return nil
}

func (s *Server) getAddr() string {
	port := s.Port
	if port == "" {
		port = defaultPort
	}
	return net.JoinHostPort("", port)
}
