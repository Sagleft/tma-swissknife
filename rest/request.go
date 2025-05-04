package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	apiEndpointFormat = "http://%s:%s/api"
	methodPOST        = "POST"
	reqTimeout        = 10 * time.Second
)

type Request struct {
	Method string         `json:"method"`
	Data   map[string]any `json:"data"`
}

// SendRequest выполняет POST-запрос к указанному хосту:порту
func SendRequest(host, port string, req Request) (Message, error) {
	// Сериализация запроса в JSON
	requestBody, err := json.Marshal(req)
	if err != nil {
		return Message{}, fmt.Errorf("marshal error: %w", err)
	}

	// Формирование URL и создание запроса
	url := fmt.Sprintf(apiEndpointFormat, host, port)
	httpReq, err := http.NewRequest(methodPOST, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return Message{}, fmt.Errorf("create request error: %w", err)
	}
	httpReq.Header.Set(HeaderContentType, HeaderValueJSON)

	// Настройка HTTP-клиента с таймаутом
	client := &http.Client{
		Timeout: reqTimeout,
	}

	// Выполнение запроса
	resp, err := client.Do(httpReq)
	if err != nil {
		return Message{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Message{}, fmt.Errorf("read response error: %w", err)
	}

	// Десериализация ответа
	var msg Message
	if err := json.Unmarshal(body, &msg); err != nil {
		return Message{}, fmt.Errorf("parse response error: %w", err)
	}

	// Проверка статуса ответа
	if msg.Status != StatusSuccess {
		return msg, fmt.Errorf("API error: %s", msg.Error)
	}

	return msg, nil
}
