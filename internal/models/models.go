package models

// User representa a un suscriptor del bot
type User struct {
	ID           int
	Phone        string   // Número asociado a CallMeBot
	GlobalTopics []string // Ej: ["Inteligencia Artificial", "Robótica"]
	LocalTopics  []string // Ej: ["Concepción", "Biobío"]
}

// Article representa la estructura estandarizada de una noticia,
// independientemente de la API externa que usemos.
type Article struct {
	Title       string
	URL         string
	Source      string
	PublishedAt string
}