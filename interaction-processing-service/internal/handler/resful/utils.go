package resful

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// readIDParamFromPath trích xuất ID từ path URL kiểu /resource/{id}
// Ví dụ: path = "/users/123" => id = "123"
func readIDParamFromPath(r *http.Request) (string, error) {
	// Tách path thành các phần
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// Kiểm tra độ dài tối thiểu để có thể lấy ID
	if len(parts) < 2 {
		return "", errors.New("invalid path format, expected /{resource}/{id}")
	}

	// Trả về phần tử cuối là ID
	id := parts[len(parts)-1]

	if id == "" {
		return "", errors.New("missing 'id' in path")
	}

	return id, nil
}

// readInt trích xuất query param dạng int, fallback nếu lỗi
func readInt(query url.Values, key string, defaultValue int) int {
	if value := query.Get(key); value != "" {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	}
	return defaultValue
}

// readString trích xuất query param dạng string, fallback nếu không có
func readString(query url.Values, key, defaultValue string) string {
	if v := query.Get(key); v != "" {
		return v
	}
	return defaultValue
}
