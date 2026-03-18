package notifier

import (
	"fmt"
	"net/http"
	"net/url"
)

// CallMeBotClient maneja el envío de mensajes por WhatsApp
type CallMeBotClient struct {
	APIKey string
}

// NewCallMeBotClient inicializa el notificador
func NewCallMeBotClient(apiKey string) *CallMeBotClient {
	return &CallMeBotClient{
		APIKey: apiKey,
	}
}

// SendMessage formatea y envía el texto al número especificado
func (c *CallMeBotClient) SendMessage(phone string, message string) error {
	// 1. CallMeBot requiere que el texto esté codificado para URL (ej: espacios como %20)
	encodedMessage := url.QueryEscape(message)

	// 2. Armamos la URL exacta según la documentación de CallMeBot
	reqURL := fmt.Sprintf("https://api.callmebot.com/whatsapp.php?phone=%s&text=%s&apikey=%s",
		phone,
		encodedMessage,
		c.APIKey,
	)

	// 3. Hacemos la petición GET
	resp, err := http.Get(reqURL)
	if err != nil {
		return fmt.Errorf("error enviando WhatsApp: %w", err)
	}
	defer resp.Body.Close()

	// 4. Validamos que CallMeBot haya aceptado el mensaje
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("CallMeBot respondió con código de error: %d", resp.StatusCode)
	}

	return nil
}