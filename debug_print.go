package gocket

import "net/http"

// debugPrintRequest prints something similair to:
// GET /path/to?query1=var
func debugPrintRequest(req *http.Request) {
	LogInfo("%s %s?%s", req.Method, req.URL.Path, req.URL.RawQuery)
}

func debugPrintMatched(response Response) {
	if response != nil {
		LogInfo("  >> handler responded: %d", response.StatusCode())
	} else {
		LogInfo("  >> handler responded: <unkown-res>")
	}
}

func debugPrintMiddlewareBlockRes(reason Response) {
	LogWarning("  >> middlware blocked: %d", reason.StatusCode())
}

func debugPrintNotFound() {
	LogWarning("  >> no matches: returning default 404")
}
