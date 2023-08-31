package flies

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
)

const (
	defaultLogFormat          = "pretty"
	defaultPort               = "8080"
	defaultResponseStatus     = "421 Misdirected Request"
	defaultResponseStatusCode = 421
)

// TCPServer is the default Server used by flies.
type TCPServer struct {
	Port                string
	ResponseStatus      string
	ResponseStatusCode  int
	ResponseBodyContent string
}

// Init sets up a TCPServer based on environment variables.
func (s *TCPServer) Init() {
	s.Port = os.Getenv("FLIES_PORT")

	s.ResponseStatus = os.Getenv("FLIES_RESPONSE_STATUS")

	if rc, err := strconv.Atoi(os.Getenv("FLIES_RESPONSE_STATUS_CODE")); err == nil {
		s.ResponseStatusCode = rc
	}

	s.ResponseBodyContent = os.Getenv("FLIES_RESPONSE_BODY_CONTENT")
}

// Listen writes any data received on the TCP connection to rawWriter, and
// writes any errors (followed by a newline) to errWriter. The request is
// parsed, the body is read and transformed into a BodyReader, and the request
// is passed to reqWriter. A response is provided to the client and the
// connection is closed.
func (s *TCPServer) Listen(errWriter, rawWriter io.Writer, reqWriter RequestWriter) error {

	l, err := net.Listen("tcp", s.getAddr())
	if err != nil {
		return err
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			s.writeError(errWriter, err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			tee := io.TeeReader(c, rawWriter)
			req, err := http.ReadRequest(bufio.NewReader(tee))
			if err != nil {
				s.writeError(errWriter, err)
			}
			// Allow reading the body multiple times
			req.Body = NewBodyReader(req.Body)
			// Read the body once to buffer it, then respond and close the
			// connection. Otherwise we won't know when to close the connection
			s.respond(c)
			rawWriter.Write([]byte("\r\n"))
			c.Close()
			reqWriter.WriteRequest(req)
		}(conn)
	}
	return err
}

func (s *TCPServer) respond(c net.Conn) {
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

func (s *TCPServer) writeError(out io.Writer, err error) error {
	out.Write([]byte(err.Error()))
	out.Write([]byte("\n"))
	return nil
}

func (s *TCPServer) getAddr() string {
	port := s.Port
	if port == "" {
		port = defaultPort
	}
	return net.JoinHostPort("", port)
}
