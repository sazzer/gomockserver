package gomockserver

import (
	"encoding/json"
	"net/http"
)

// Response is a representation of the response to send to the client.
type Response struct {
	Status  int
	Headers http.Header
	Body    []byte
}

// Write will write the response details to the provided response writer.
func (r Response) Write(w http.ResponseWriter) {
	for name, values := range r.Headers {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(r.Status)
	_, _ = w.Write(r.Body)
}

// ResponseBuilder is a means to contribute to the response to send to the client.
type ResponseBuilder interface {
	// PopulateResponse is called to allow the builder to update any parts of the response struct.
	PopulateResponse(r *Response, req *http.Request)
}

// ResponseBuilderFunc is function type that implements the `ResponseBuilder` interface.
type ResponseBuilderFunc func(*Response, *http.Request)

func (f ResponseBuilderFunc) PopulateResponse(r *Response, req *http.Request) {
	f(r, req)
}

// ResponseBuilders represents a set of builders that are all called in sequence to populate a response.
type ResponseBuilders []ResponseBuilder

func (b ResponseBuilders) PopulateResponse(r *Response, req *http.Request) {
	for _, builder := range b {
		builder.PopulateResponse(r, req)
	}
}

// ResponseStatus returns a `ResponseBuilder` to specify the status code of the response.
func ResponseStatus(status int) ResponseBuilder {
	return ResponseBuilderFunc(func(r *Response, req *http.Request) {
		r.Status = status
	})
}

// ResponseSetHeader will set a header value to exactly the value provided.
func ResponseSetHeader(name, value string) ResponseBuilder {
	return ResponseBuilderFunc(func(r *Response, req *http.Request) {
		r.Headers.Set(name, value)
	})
}

// ResponseAppendHeader will append a new value for the header name provided.
func ResponseAppendHeader(name, value string) ResponseBuilder {
	return ResponseBuilderFunc(func(r *Response, req *http.Request) {
		r.Headers.Add(name, value)
	})
}

// ResponseBody will indicate the data to return as the response body.
func ResponseBody(data []byte) ResponseBuilder {
	return ResponseBuilderFunc(func(r *Response, req *http.Request) {
		r.Body = data
	})
}

// ResponseJSON will encode the provided value as JSON and use it as the response, also setting the content-type header.
func ResponseJSON(data interface{}) ResponseBuilder {
	bytes, _ := json.Marshal(data)

	return ResponseBuilders{
		ResponseSetHeader("content-type", "application/json"),
		ResponseBody(bytes),
	}
}
