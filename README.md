# NewsBot Desktop: Agregador de Noticias Local

Una aplicación de escritorio ligera y descentralizada, construida en Go y Wails, que actúa como tu propio motor de búsqueda de noticias en segundo plano y te notifica directamente por WhatsApp.

En lugar de depender de servicios SaaS centralizados con suscripciones costosas, NewsBot opera bajo una arquitectura **BYOK (Bring Your Own Key)**. El usuario ejecuta el daemon en su propia máquina, gestionando sus propios límites de API y manteniendo su configuración 100% privada gracias a una base de datos SQLite embebida.

## Características Principales

* **Daemon Concurrente:** Utiliza *Goroutines* de Go para escanear noticias en segundo plano a intervalos configurables, sin bloquear el hilo de la interfaz de usuario.
* **Arquitectura BYOK:** Seguridad delegada. El usuario ingresa sus propias credenciales de NewsAPI y CallMeBot, eliminando costos de infraestructura y previniendo la saturación del servidor.
* **Persistencia Local Automática:** Una base de datos SQLite embebida guarda la configuración del usuario y un registro de las URLs ya enviadas para garantizar que **nunca recibas noticias duplicadas**.
* **Algoritmo de Filtrado Estricto:** Separa la lógica de búsqueda global (temas amplios) de la búsqueda local, aplicando validaciones de strings en memoria para evitar falsos positivos geográficos.
* **Agrupación de Notificaciones (Digest):** Agrupa múltiples noticias de un mismo tópico en un solo mensaje de WhatsApp para optimizar el uso de la API gratuita (evitando el error HTTP 403 / Límite de tasa).
* **Interfaz Nativa:** Frontend limpio y responsivo construido con HTML puro y Tailwind CSS, empaquetado nativamente a través de Wails (consumiendo mucha menos RAM que alternativas como Electron).

## Stack Tecnológico

* **Backend:** Go (Golang)
* **Frontend:** Vanilla JavaScript + HTML + Tailwind CSS
* **Desktop Framework:** Wails v2 (Usa el motor web nativo del SO)
* **Base de Datos:** SQLite (Driver `modernc.org/sqlite` en Go puro)

## Instalación y Uso (Modo Portable)

No necesitas instalar Go ni dependencias complejas.

1.  Ve a la pestaña de [Releases](https://github.com/BryzonSnow/newsbot-desktop/releases/tag/v1.0.0.1) y descarga el archivo `newsbot-desktop.exe`.
2.  Ejecuta la aplicación (Portable, sin instalación).
3.  En la interfaz, ingresa tu número de WhatsApp y tus APIs gratuitas:
    * **NewsAPI:** Obtén tu llave en [newsapi.org](https://newsapi.org/).
    * **CallMeBot:** Autoriza tu número y obtén tu código en [callmebot.com](https://www.callmebot.com/).
4.  Configura tus tópicos de interés (ej. "Machine Learning", "Concepción").
5.  Haz clic en Guardar. El motor en Go comenzará a ejecutarse en segundo plano automáticamente.

## Nota sobre los Límites de API
Esta aplicación está diseñada para ser un motor personal silencioso. Ten en cuenta que la API gratuita de CallMeBot tiene un límite de **16 mensajes cada 4 horas**. Se recomienda configurar el intervalo de escaneo en la aplicación a **60 minutos** para mantenerse holgadamente dentro de la cuota gratuita sin ser bloqueado.
