package rhttp

import "net/http"

// Determines whether a request should be retried
type RetryPolicy func(resp *http.Response, err error) (bool, error)

// DefaultRetry simply checks for:
//  1. err != nil
//  2. resp.StatusCode >= 500
func DefaultRetry(resp *http.Response, err error) (bool, error) {
	if err != nil {
		return true, err
	}
	return resp.StatusCode >= 500, nil
}
