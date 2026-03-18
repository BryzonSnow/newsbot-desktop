package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"newsbot-desktop/internal/db"
	"newsbot-desktop/internal/models"
	"newsbot-desktop/internal/news"
	"newsbot-desktop/internal/notifier"
)

type App struct {
	ctx      context.Context
	database *sql.DB
	reiniciarReloj chan bool 
}

func NewApp() *App {
	return &App{
		reiniciarReloj: make(chan bool),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	database, err := db.InitDB("newsbot_local.db")
	if err != nil {
		fmt.Println("Error iniciando DB:", err)
		return
	}
	a.database = database
	go a.iniciarDaemon()
}

func (a *App) ObtenerConfiguracion() map[string]interface{} {
	phone, newsKey, waKey, global, local, intervalo, err := db.ObtenerConfig(a.database)
	if err != nil {
		return map[string]interface{}{"error": "No hay configuración guardada"}
	}
	
	return map[string]interface{}{
		"phone":   phone,
		"newsKey": newsKey,
		"waKey":   waKey,
		"global":  global,
		"local":   local,
		"intervalo": intervalo,
	}
}

func (a *App) GuardarConfiguracion(phone, newsKey, waKey, global, local string, intervalo int) string {
	err := db.GuardarConfig(a.database, phone, newsKey, waKey, global, local, intervalo)
	if err != nil {
		return "Error al guardar: " + err.Error()
	}
	
	select {
	case a.reiniciarReloj <- true:
	default:
	}

	return "✅ Configuración guardada. Motor actualizado."
}

func (a *App) iniciarDaemon() {
	_, _, _, _, _, intervalo, _ := db.ObtenerConfig(a.database)
	if intervalo <= 0 { intervalo = 30 }

	ticker := time.NewTicker(time.Duration(intervalo) * time.Minute)
	defer ticker.Stop()

	a.ejecutarCiclo() 

	for {
		select {
		case <-ticker.C:
			a.ejecutarCiclo()
			
		case <-a.reiniciarReloj:
			_, _, _, _, _, nuevoIntervalo, _ := db.ObtenerConfig(a.database)
			if nuevoIntervalo <= 0 { nuevoIntervalo = 30 }
			
			fmt.Printf("🔄 Reiniciando reloj a %d minutos...\n", nuevoIntervalo)
			ticker.Reset(time.Duration(nuevoIntervalo) * time.Minute)
			a.ejecutarCiclo()

		case <-a.ctx.Done():
			return
		}
	}
}

func (a *App) ejecutarCiclo() {
	phone, newsKey, waKey, global, local, _, err := db.ObtenerConfig(a.database)
	if err != nil || phone == "" || newsKey == "" || waKey == "" {
		fmt.Println("Daemon pausado: Faltan credenciales del usuario.")
		return
	}

	fmt.Println("Ejecutando ciclo de noticias...")
	newsClient := news.NewClient(newsKey)
	whatsappClient := notifier.NewCallMeBotClient(waKey)

	globalTopics := strings.Split(global, ",")
	for _, topic := range globalTopics {
		topic = strings.TrimSpace(topic)
		if topic == "" { continue }
		if articulos, err := newsClient.Fetch(topic); err == nil {
			a.procesarYEnviar(whatsappClient, phone, topic, articulos, 2)
		}
	}

	localTopics := strings.Split(local, ",")
	for _, topic := range localTopics {
		topic = strings.TrimSpace(topic)
		if topic == "" { continue }
		if articulosRaw, err := newsClient.Fetch(topic); err == nil {
			palabraClave := strings.ReplaceAll(topic, "\" AND Chile", "")
			palabraClave = strings.ReplaceAll(palabraClave, "\"", "")
			articulosFiltrados := filtrarNoticiasLocales(articulosRaw, palabraClave)
			a.procesarYEnviar(whatsappClient, phone, topic, articulosFiltrados, 2)
		}
	}
}

func filtrarNoticiasLocales(articulos []models.Article, localidad string) []models.Article {
	var filtrados []models.Article
	locLower := strings.ToLower(localidad)
	for _, articulo := range articulos {
		if strings.Contains(strings.ToLower(articulo.Title), locLower) {
			filtrados = append(filtrados, articulo)
		}
	}
	return filtrados
}

func (a *App) procesarYEnviar(client *notifier.CallMeBotClient, phone, topic string, articulos []models.Article, limite int) {
	enviadosEnEsteCiclo := 0
	
	var mensajeGrupal strings.Builder
	mensajeGrupal.WriteString(fmt.Sprintf("📰 *Resumen: %s*\n\n", topic))

	for _, articulo := range articulos {
		if enviadosEnEsteCiclo >= limite { break }
		if db.ArticleExists(a.database, articulo.URL) { continue }

		// Agregamos la noticia al bloque de texto
		noticia := fmt.Sprintf("🔹 *%s*\n🏢 %s\n🔗 %s\n\n", articulo.Title, articulo.Source, articulo.URL)
		mensajeGrupal.WriteString(noticia)
		
		db.MarkAsSent(a.database, articulo.URL)
		enviadosEnEsteCiclo++
	}

	if enviadosEnEsteCiclo > 0 {
		textoFinal := strings.TrimSpace(mensajeGrupal.String())
		
		if err := client.SendMessage(phone, textoFinal); err != nil {
			log.Printf("❌ Error enviando WhatsApp al %s: %v\n", phone, err)
			return
		}

		fmt.Printf("WhatsApp enviado agrupado (Tópico: %s, Noticias: %d)\n", topic, enviadosEnEsteCiclo)
		
		time.Sleep(15 * time.Second) 
	} else {
		fmt.Printf("📭 Nada nuevo para: %s\n", topic)
	}
}