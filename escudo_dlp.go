// Copyright © 2026 José Emanuel Guadalupe López - Justicia Digital. 100% Derechos Reservados.
package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
)

// =====================================================================
// FASE 1: MATRIZ DE PATRONES CLASIFICADOS (100% LEGAL)
// =====================================================================
// El Escudo DLP buscará estas huellas digitales en el tráfico saliente.
var patronesBloqueados = []*regexp.Regexp{
	regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`), // Formato de Número de Seguro Social (SSN)
	regexp.MustCompile(`(?i)privilegio_cliente_abogado`), // Palabra clave altamente confidencial
}

// dlpResponseWriter es el interceptor forense que captura la respuesta antes de que salga a la red
type dlpResponseWriter struct {
	http.ResponseWriter
	bodyCatcher *bytes.Buffer
	statusCode  int
}

func (w *dlpResponseWriter) Write(b []byte) (int, error) {
	return w.bodyCatcher.Write(b) // Almacenar en la bóveda temporal en lugar de enviarlo
}

func (w *dlpResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

// =====================================================================
// FASE 2: MOTOR DEL ESCUDO (MIDDLEWARE DLP)
// =====================================================================
func EscudoDLP(siguiente http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturador := &dlpResponseWriter{
			ResponseWriter: w,
			bodyCatcher:    &bytes.Buffer{},
			statusCode:     http.StatusOK,
		}

		// Ejecutar la ruta solicitada, pero el resultado queda atrapado en el 'capturador'
		siguiente.ServeHTTP(capturador, r)

		textoSaliente := capturador.bodyCatcher.String()

		// Auditar el 100% del texto buscando datos clasificados
		for _, patron := range patronesBloqueados {
			if patron.MatchString(textoSaliente) {
				// 🚨 INFRACCIÓN DETECTADA: Abortar transmisión e imprimir log forense
				log.Printf("🚨 [ALERTA DLP FORENSE] Fuga de datos bloqueada al 100%%. IP: %s | Ruta: %s\n", r.RemoteAddr, r.URL.Path)
				http.Error(w, "ACCESO DENEGADO: Protocolo DLP activado. Datos confidenciales detectados en la transmisión. Operación abortada al 100%.", http.StatusForbidden)
				return
			}
		}

		// ✅ TRÁFICO LIMPIO: Si es 100% seguro, liberar la información
		if capturador.statusCode != http.StatusOK {
			w.WriteHeader(capturador.statusCode)
		}
		w.Write(capturador.bodyCatcher.Bytes())
	})
}

// =====================================================================
// FASE 3: RUTAS DE SIMULACIÓN TÁCTICA
// =====================================================================
func rutaExpedienteClasificado(w http.ResponseWriter, r *http.Request) {
	// Esta ruta intenta enviar información protegida. El DLP la destruirá al instante.
	fmt.Fprint(w, "Expediente 001: El cliente posee el SSN 999-00-1234 bajo la regla de privilegio_cliente_abogado.")
}

func rutaPublica(w http.ResponseWriter, r *http.Request) {
	// Esta ruta no contiene datos clasificados. El DLP la dejará pasar al 100%.
	fmt.Fprint(w, "Información pública de Justicia Digital. 100% autorizada.")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/clasificado", rutaExpedienteClasificado)
	mux.HandleFunc("/api/publico", rutaPublica)

	// Envolver el 100% del enrutador con el Escudo DLP de grado militar
	servidorProtegido := EscudoDLP(mux)

	fmt.Println("=====================================================================")
	fmt.Println("🛡️  MOTOR DLP (ESCUDO DE PÉRDIDA DE DATOS) ACTIVO AL 100%")
	fmt.Printf("Hora de sellado (UTC): %s\n", time.Now().UTC().Format(time.RFC3339))
	fmt.Println("=====================================================================")
	fmt.Println("Vigilando el 100% del tráfico perimetral. Puerto 8080 en Termux.")

	errorServidor := http.ListenAndServe(":8081", servidorProtegido)
	if errorServidor != nil {
		log.Fatalf("ERROR TÁCTICO: %v", errorServidor)
	}
}
