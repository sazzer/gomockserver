package gomockserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nsf/jsondiff"
)

// ResponseJSON will encode the provided value as JSON and use it as the response, also setting the content-type header.
func ResponseJSON(data interface{}) ResponseBuilder {
	bytes, _ := json.Marshal(data)

	return ResponseBuilders{
		ResponseSetHeader("content-type", "application/json"),
		ResponseBody(bytes),
	}
}

func matchJSON(r *http.Request, expected string) jsondiff.Difference {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return jsondiff.NoMatch
	}

	options := jsondiff.DefaultJSONOptions()
	diff, _ := jsondiff.Compare(body, []byte(expected), &options)

	return diff
}

// MatchJSONFull will compare the request body to the provided JSON string and ensure that the two are semantically
// identical.
// The order of keys in JSON objects is not important, but every value must be present.
func MatchJSONFull(expected string) MatchRule {
	return MatchRuleFunc(func(r *http.Request) bool {
		return matchJSON(r, expected) == jsondiff.FullMatch
	})
}

// MatchJSONCompatible will compare the request body to the provided JSON string and ensure that the two are compatible.
// Every value in the request body must be present in the expected body, but the request body may contain keys that are
// not present in the expected body as well.
//
// For example, the following is a valid match:
// * Request Body = {"a": 1, "b": {"c": 2, "d": 3}, "e": 4}
// * Expected = {"a": 1, "b": {"c": 2}}
//
// As with MatchJSONFull, the order of keys is not important.
func MatchJSONCompatible(expected string) MatchRule {
	return MatchRuleFunc(func(r *http.Request) bool {
		match := matchJSON(r, expected)

		return match == jsondiff.FullMatch || match == jsondiff.SupersetMatch
	})
}
