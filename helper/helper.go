package helper

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type ResponseMsg struct {
	Meta Meta        `json:"meta"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	// Memulai waktu eksekusi
	start := time.Now()
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}
	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}
	// Menandai fungsi ini akan selesai diukur waktu eksekusinya setelah return
	defer func() {
		// Menghentikan waktu eksekusi
		end := time.Now()

		// Hitung durasi eksekusi
		duration := end.Sub(start)

		// Tampilkan durasi eksekusi
		fmt.Printf("Durasi eksekusi APIResponse(): %s\n", duration)
	}()

	// Simulasi operasi yang memakan waktu
	time.Sleep(2 * time.Second)
	return jsonResponse
}

func APIResponseMessage(message string, code int, status string) ResponseMsg {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}
	jsonResponse := ResponseMsg{
		Meta: meta,
	}
	return jsonResponse
}

func FormatValidationError(err error) []string {
	var errors []string

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}
	return errors
}