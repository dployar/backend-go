package handler

import (
	"encoding/json"
	"net/http"
)

// HealthHandler es el manejador HTTP encargado de verificar
// que la aplicación se encuentra activa y responde.
//
// Se utiliza principalmente para:
// - Monitorear la disponibilidad del servidor.
// - Realizar health checks desde Docker o Kubernetes.
// - Permitir verificaciones desde balanceadores de carga.
// - Facilitar el diagnóstico básico de la aplicación.
//
// Se declara como una estructura vacía porque actualmente
// no necesita almacenar dependencias ni estado interno.
type HealtHandler struct{}

// NewHealthHandler crea y devuelve una nueva instancia
// del manejador de salud.
//
// Se utiliza el patrón constructor para centralizar
// la inicialización del handler.
//
// Retornar un puntero permite mantener una interfaz
// consistente si en el futuro se agregan dependencias,
// como un logger, configuración o servicios.
func NewHealthHandler() *HealtHandler {
	return &HealtHandler{}
}

// HealthResponse define la estructura de la respuesta
// que será enviada al cliente en formato JSON.
//
// Los tags `json` determinan los nombres de las
// propiedades al serializar la estructura.
//
// Ejemplo de respuesta:
//
//	{
//	    "status": "ok",
//	    "version": "0.0.1"
//	}
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// ServeHTTP procesa las solicitudes HTTP recibidas
// por el endpoint de salud.
//
// Este método implementa la interfaz http.Handler,
// definida por la biblioteca estándar de Go.
//
// Gracias a esto, HealthHandler puede utilizarse
// directamente con routers como Chi.
//
// Parámetros:
// - w: permite escribir la respuesta HTTP.
// - r: contiene la información de la solicitud.
//
// El receptor (*HealthHandler) asocia este método
// con la estructura HealthHandler.
func (handle *HealtHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Status:  "Ok",
		Version: "0.0.1",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return
	}
}
