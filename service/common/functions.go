package common

import (
	"github.com/alioth-center/infrastructure/network/http"
	"strings"
)

// GetApiKey extracts the API key from the Authorization header of the request.
// It removes the "Bearer " prefix from the Authorization header value, allowing
// for the raw API key to be used for further processing or validation.
//
// Parameters:
//   - params: http.Params containing the request headers.
//
// Returns:
//   - A string representing the API key extracted from the Authorization header.
//     If the Authorization header does not contain the "Bearer " prefix, the original
//     header value is returned.
func GetApiKey(params http.Params) string {
	return strings.TrimPrefix(params.GetString("Authorization"), "Bearer ")
}
