# Go Mock Server

[![CI](https://github.com/sazzer/gomockserver/actions/workflows/ci.yml/badge.svg)](https://github.com/sazzer/gomockserver/actions/workflows/ci.yml)
[![GoDoc](https://godoc.org/github.com/sazzer/gomockserver?status.svg)](https://godoc.org/github.com/sazzer/gomockserver)
[![Coverage Status](https://coveralls.io/repos/github/sazzer/gomockserver/badge.svg?branch=main)](https://coveralls.io/github/sazzer/gomockserver?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/sazzer/gomockserver)](https://goreportcard.com/report/github.com/sazzer/gomockserver)
[![license](https://img.shields.io/badge/license-MIT-blue.svg)]()

Provide a way to mock out an HTTP Server for testing purposes, to ensure that calls are made correctly and return appropriate responses.

## Usage

Firstly a new mock server instance is needed:

```go
server := gomockserver.New(t)
defer server.Close()
```

Closing the server is important, and should be done once it's no longer needed. The `defer` ensures this happens at the end of the current test.

Once created, the server will handle all incoming requests to it. The URL can be determined by using `server.URL()`, which will return a string like `http://127.0.0.1:54681`. This is the base URL to the server, under which all requests can be handled.

Once set up, the server can be configured to match incoming requests and return responses to them. For example, the following will handle a call to `GET /testing/abc` and return an `application/json` response with a body of `"Hello"`:

```go
server.Matches(gomockserver.MatchRequest("GET", "/testing/abc")).
	RespondsWith(gomockserver.ResponseJSON("Hello"))
```

### Matches

The call to `server.Matches()` will take a number of `MatchRule` instances. A `Match` is considered successful if every single `MatchRule` passes for the incoming request. This allows requests to match against as many or as few details as possible.

Standard rules that can be configured are:

- `MatchMethod` - Matches the HTTP Method
- `MatchURLPath` - Matches the full incoming URL
- `MatchURLQuery` - Matches a query parameter with a specific value
- `MatchRequest` - Matches both the HTTP Method and the URL
- `MatchHeader` - Matches a header name with a specific value
- `MatchJSONFull` - Matches the request body in full against a JSON document
- `MatchJSONCompatible` - Ensures the request body is a superset of a given JSON document - i.e. additional fields in the request do not stop this from matching.

Additionally, you can write any custom match rule that you want as long as it fulfils the `MatchRule` interface. There is also a `MatchRuleFunc` function type that already implements the interface, so rules can be written as anonymous functions if desired.

### Responses

The result of calling `server.Matches()` is a `*Match`. This can then be augmented to detail how the response should look, by adding `ResponseBuilder` instances via the `RespondWith` method. As with `Matches()`, this method can take as many `ResponseBuilder` instances as needed, each of which will configure the response in some manner.

Standard response builders that can be configured are:

- `ResponseStatus` - Sets the status code
- `ResponseSetHeader` - Overwrites a response header
- `ResponseAppendHeader` - Append a new value to a response header
- `ResponseBody` - Set the body of the response
- `ResponseJSON` - Set the body of the response to the JSON encoding of the provided object, and set the `Content-Type` header to `application/json`.

Additionally, you can write any custom builder that you want as long as it fulfils the `ResponseBuilder` interface. There is also a `ResponseBuilderFunc` function type that already implements the interface, so rules can be written as anonymous functions if desired.

## Matching Requests

Every request that is received by the mock server is compared to every `Match` that is configured, in the order they were configured, until the first one is a match. At this point,the response from this `Match` is built and sent back to the client.

You can configure as many different `Match`es on the server as you want, but every request will only ever match at most one.

Any incoming requests that do not match a configured `Match` will return an `HTTP 404 Not Found`.

## Counting Requests

Go Mock Server will keep track of the number of times every `Match` has been used to respond to a request. This can be used in tests to assert that a given request was made the correct number of times:

```go
match := server.Matches(gomockserver.MatchRequest("GET", "/testing/abc")).
	RespondsWith(gomockserver.ResponseJSON("Hello"))

// Run tests

is.Equal(match.Count(), 1)
```

We also keep track of the number of times we handled unmatched requests, in case that's interesting. Most often that will be used to assert that this was zero - i.e. that all the requests that we handled were matched:

```go
// Run tsts

is.Equal(server.UnmatchedCount(), 0)
```

## Examples

Examples of how to use this can be found in [server_test.go](https://github.com/sazzer/gomockserver/blob/main/server_test.go).
