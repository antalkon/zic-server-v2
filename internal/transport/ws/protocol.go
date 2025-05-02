package ws

// Структура стандартного WS-сообщения
type Message struct {
	Type    string      `json:"type"`    // например: "ping", "command", "response"
	Tunnel  string      `json:"tunnel"`  // UUID туннеля
	Payload interface{} `json:"payload"` // конкретные данные
}

// Пример полезной нагрузки команды
type CommandPayload struct {
	Command string `json:"command"`
	Params  string `json:"params"` // или можно использовать map[string]interface{}
}
