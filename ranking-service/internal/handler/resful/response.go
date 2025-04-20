package resful

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Success gửi phản hồi JSON thành công với data, status code và header tuỳ chỉnh
func Success(w http.ResponseWriter, data interface{}, status int, headers http.Header) error {
	payload := &Response{
		Data: data,
	}

	js, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	for key, values := range headers {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(js)
	return nil
}

// Error gửi phản hồi lỗi dưới dạng JSON
func Error(w http.ResponseWriter, err error, status int) {
	payload := &Response{
		Message: err.Error(),
	}

	js, jsonErr := json.Marshal(payload)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(js)
	return
}
