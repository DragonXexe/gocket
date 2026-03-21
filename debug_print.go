package gocket

import "net/http"

// debugPrintRequest prints something similair to:
// GET /path/to?query1=var
func debugPrintRequest(req *http.Request) {
	LogInfo("%s %s?%s", req.Method, req.URL.Path, req.URL.RawQuery)
}

func debugPrintMatched(response Response) {
	LogInfo("  >> handler responded: %d", response.StatusCode())
}

func debugPrintMiddlewareBlockRes(reason Response) {
	LogWarning("  >> middlware blocked: %d", reason.StatusCode())
}

func debugPrintNotFound() {
	LogWarning("  >> no matches: returning default 404")
}
