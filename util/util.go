package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/rizvn/panics"
)

func GetMethodName(f any) string {
	val := reflect.ValueOf(f)
	pc := val.Pointer()
	funcObj := runtime.FuncForPC(pc)
	return funcObj.Name()
}

func FromJson[T any](r io.Reader) *T {
	var v T
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&v)
	panics.OnError(err, "error decoding JSON")
	return &v
}

func ToJson(w io.Writer, v any) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(v)
	panics.OnError(err, "error encoding JSON")
}

func ToJsonString(v any) string {
	buf := &bytes.Buffer{}
	ToJson(buf, v)
	return buf.String()
}

func FromJsonString[T any](s string) *T {
	var v T
	decoder := json.NewDecoder(bytes.NewBufferString(s))
	err := decoder.Decode(&v)
	panics.OnError(err, "error decoding JSON from string")
	return &v
}

func ExtractJson(input string) string {
	// Extract the first '{' and last '}' inclusively
	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")
	var toolJson string
	if start != -1 && end != -1 && end > start {
		toolJson = input[start : end+1]
	} else {
		toolJson = input
	}
	return toolJson
}

func GetRequiredEnvVar(name string) string {
	value := os.Getenv(name)
	panics.OnBlank(value, "Required env var "+name+" is blank")
	return value
}

func GetEnvVarOrDefault(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetEnvVarOrDefaultBoolean(name string, defaultValue bool) bool {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	parsedValue, err := strconv.ParseBool(value)
	panics.OnError(err, "Error parsing env var "+name+" as boolean")
	return parsedValue
}

// PtrToString safely dereferences a string pointer, returning an empty string if the pointer is nil.
func PtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func StrToPtr(s string) *string {
	return &s
}

// MapGetByType retrieves a dependency from the map by its type.
// It uses the qualified type name as the key.
func MapGetByType[T any](deps *map[string]any) T {
	// CreateAsUser a zero value of type T to get its qualified name
	var i T

	// BuildIndex the qualified type name
	name := GetQualifiedTypeName(i)

	dp := *deps

	if obj, ok := dp[name]; !ok {
		panics.OnFalse(ok, "Dependency with key "+name+" not found in map")
	} else {
		return obj.(T)
	}

	return (*deps)[name].(T)
}

// MapSetByType stores a dependency in the map by its type.
// It uses the qualified type name as the key.
func MapSetByType(deps *map[string]any, i interface{}) {
	name := GetQualifiedTypeName(i)
	(*deps)[name] = i
}

func GetQualifiedTypeName(i interface{}) string {
	t := reflect.TypeOf(i)
	name := t.String()

	return name
}

func RemoveTag(tag, s string) string {
	start := strings.Index(s, "<"+tag+">")
	end := strings.Index(s, "</"+tag+">")
	if start == -1 || end == -1 || end < start {
		return s
	}
	result := s[:start] + s[end+len("</"+tag+">"):]
	return strings.TrimSpace(result)
}

func ExtractTag(tag, s string) string {
	start := strings.Index(s, "<"+tag+">")
	end := strings.Index(s, "</"+tag+">")
	if start == -1 || end == -1 || end < start {
		return ""
	}
	return strings.TrimSpace(s[start+len("<"+tag+">") : end])
}

func RemoveNewLines(s string) string {
	return strings.ReplaceAll(s, "\n", " ")
}

func RemoveStackTrace(s string) string {
	end := strings.Index(s, "Stacktrace: goroutine")
	return s[:end]
}

func GetUserName(rq *http.Request) string {
	if os.Getenv("OIDC_ENABLED") == "true" {
		return rq.Context().Value("username").(string)
	} else {
		return "default_user@example.com"
	}
}

func GetUserRoles(rq *http.Request) string {
	return rq.Context().Value("user_roles").(string)
}
