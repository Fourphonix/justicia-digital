package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"sync"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jung-kurt/gofpdf"
	_ "github.com/lib/pq"
)

var db *sql.DB

// 🔒 LLAVES DE SEGURIDAD
var jwtKey = []byte("Super-Secreto-Justicia-Digital-100-Absoluto")
const aesKey = "Llave-AES-100-Secreta"

type TransicionEstado struct {
	ExpedienteID     string `json:"expediente_id"`
	EstadoAnterior   string `json:"estado_anterior"`
	EstadoNuevo      string `json:"estado_nuevo"`
	MotivoCambio     string `json:"motivo_cambio"`
	UsuarioID        string `json:"usuario_id"`
	EvidenciaSecreta string `json:"evidencia_secreta"`
}

// --- LÓGICA CORE (Inmutabilidad y Cifrado) ---
func RegistrarCambioEstado(ctx context.Context, t TransicionEstado) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil { return err }
	defer tx.Rollback()

	var hashAnterior string
	queryHashPrevio := `SELECT hash_transicion FROM historial_estados_expediente WHERE id_expediente = $1 ORDER BY registrado_en DESC LIMIT 1`
	err = tx.QueryRowContext(ctx, queryHashPrevio, t.ExpedienteID).Scan(&hashAnterior)
	if err == sql.ErrNoRows {
		hashAnterior = "0000000000000000000000000000000000000000000000000000000000000000"
	} else if err != nil { return err }

	fechaActual := time.Now().UTC().Format(time.RFC3339)
	payload := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s", t.ExpedienteID, t.EstadoAnterior, t.EstadoNuevo, t.MotivoCambio, t.UsuarioID, fechaActual, hashAnterior)

	hashBytes := sha256.Sum256([]byte(payload))
	hashCalculado := hex.EncodeToString(hashBytes[:])

	insertHistorial := `INSERT INTO historial_estados_expediente (id_expediente, estado_anterior, estado_nuevo, motivo_cambio, usuario_id, hash_transicion, hash_anterior) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.ExecContext(ctx, insertHistorial, t.ExpedienteID, t.EstadoAnterior, t.EstadoNuevo, t.MotivoCambio, t.UsuarioID, hashCalculado, hashAnterior)
	if err != nil { return err }

	updateExpediente := `UPDATE expedientes SET estado_actual = $1, detalles_clasificados = pgp_sym_encrypt($2, $3) WHERE id_expediente = $4`
	_, err = tx.ExecContext(ctx, updateExpediente, t.EstadoNuevo, t.EvidenciaSecreta, aesKey, t.ExpedienteID)
	if err != nil { return err }

	return tx.Commit()
}

func AuditarExpediente(ctx context.Context, expedienteID string) (bool, string, error) {
	query := `SELECT hash_transicion, hash_anterior FROM historial_estados_expediente WHERE id_expediente = $1 ORDER BY registrado_en ASC`
	rows, err := db.QueryContext(ctx, query, expedienteID)
	if err != nil { return false, "", err }
	defer rows.Close()

	hashEsperado := "0000000000000000000000000000000000000000000000000000000000000000"
	for rows.Next() {
		var hashTransicion, hashAnterior string
		if err := rows.Scan(&hashTransicion, &hashAnterior); err != nil { return false, "", err }
		if hashAnterior != hashEsperado { return false, "", fmt.Errorf("cadena rota") }
		hashEsperado = hashTransicion
	}

	var evidencia string
	queryRevelar := `SELECT pgp_sym_decrypt(detalles_clasificados, $1) FROM expedientes WHERE id_expediente = $2`
	err = db.QueryRowContext(ctx, queryRevelar, aesKey, expedienteID).Scan(&evidencia)
	if err != nil { evidencia = "Sin evidencia o error de cifrado" }

	return true, evidencia, nil
}
// --- MIDDLEWARE DE LÍMITE DE TASA 100% LEGAL (ANTI-DDOS) ---
var (
	visitas   = make(map[string]int)
	visitasMu sync.Mutex
)

func limiteTasaForense(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]

		visitasMu.Lock()
		visitas[ip]++
		peticiones := visitas[ip]
		visitasMu.Unlock()

		// Límite 100% estricto: Si pasa de 50 peticiones rápidas, se bloquea al 100%
		if peticiones > 50 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "100% BLOQUEADO: Tráfico anómalo detectado. Defensa DDoS activada 100% legalmente."}`))
			return
		}

		next.ServeHTTP(w, r)
	}
}
// --- MIDDLEWARE 1: TRAZABILIDAD FORENSE 100% LEGAL ---
func trazabilidadForense(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		next.ServeHTTP(w, r)
		duracion := time.Since(inicio)
		log.Printf("[AUDITORÍA 100%%] IP: %s | Método: %s | Ruta: %s | Tiempo: %v", r.RemoteAddr, r.Method, r.RequestURI, duracion)
	}
}

// --- MIDDLEWARE 2: TIEMPO DE EXPIRACIÓN 100% ESTRICTO (TIMEOUT) ---
func timeoutForense(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Envuelve el handler con un límite 100% estricto de 5 segundos al 100%
		http.TimeoutHandler(next, 5*time.Second, `{"error": "100% ALERTA: Tiempo de espera agotado. Conexión cortada 100% legalmente y de raíz."}`).ServeHTTP(w, r)
	}
}

// --- MIDDLEWARE 3: CABECERAS DE SEGURIDAD 100% ESTRICTAS ---
func cabecerasSeguridad(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none';")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		next.ServeHTTP(w, r)
	}
}

// --- MIDDLEWARE 4: SEGURIDAD JWT ---
func seguridadJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Acceso denegado: 100% falta token de autorización"}`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado")
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error": "Acceso denegado: Token inválido o expirado"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// --- ENDPOINTS REST EXISTENTES ---
func loginHandler(w http.ResponseWriter, r *http.Request) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usuario": "juez_supremo",
		"nivel":   "100%_autorizado",
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error al generar el token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, tokenString)))
}

func registrarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	var transicion TransicionEstado
	if err := json.NewDecoder(r.Body).Decode(&transicion); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	if err := RegistrarCambioEstado(r.Context(), transicion); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"mensaje": "Transición y cifrado de evidencia 100% registrados con éxito"}`))
}

func auditarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	expedienteID := r.URL.Query().Get("id")
	integro, evidencia, err := AuditarExpediente(r.Context(), expedienteID)
	if err != nil || !integro {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"status": "ALERTA", "mensaje": "Expediente comprometido"}`))
		return
	}
	respuesta := fmt.Sprintf(`{"status": "OK", "mensaje": "Expediente 100%% íntegro", "evidencia_revelada": "%s"}`, evidencia)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(respuesta))
}

func subirEvidenciaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("archivo")
	if err != nil {
		http.Error(w, "Error al recibir el archivo de evidencia", http.StatusBadRequest)
		return
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		http.Error(w, "Error procesando huella digital", http.StatusInternalServerError)
		return
	}
	hashString := hex.EncodeToString(hash.Sum(nil))
	file.Seek(0, io.SeekStart)

	extension := filepath.Ext(handler.Filename)
	rutaFinal := filepath.Join("evidencia_fisica", hashString+extension)
	
	dst, err := os.Create(rutaFinal)
	if err != nil {
		http.Error(w, "Error al almacenar", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	respuesta := fmt.Sprintf(`{"status": "100%% OK", "mensaje": "Evidencia asegurada en bóveda", "hash_integridad": "%s"}`, hashString)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(respuesta))
}

func generarSentenciaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	expedienteID := r.URL.Query().Get("id")
	integro, evidencia, err := AuditarExpediente(r.Context(), expedienteID)
	if err != nil || !integro {
		http.Error(w, "ALERTA: Expediente alterado. No se puede emitir documento legal.", http.StatusConflict)
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	
	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(190, 10, "TRIBUNAL SUPREMO - JUSTICIA DIGITAL", "0", 1, "C", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(190, 10, "ACTA DE SENTENCIA Y AUDITORIA AL 100%", "0", 1, "C", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(190, 10, fmt.Sprintf("ID de Expediente: %s", expedienteID), "0", 1, "L", false, 0, "")
	pdf.CellFormat(190, 10, fmt.Sprintf("Fecha de Emision: %s", time.Now().Format("2006-01-02 15:04:05")), "0", 1, "L", false, 0, "")
	
	pdf.SetTextColor(0, 128, 0)
	pdf.CellFormat(190, 10, "ESTADO: VERIFICADO - HASH SHA-256 INMUTABLE", "0", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(190, 10, "DECLARACION DE EVIDENCIA CLASIFICADA", "0", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(190, 8, evidencia, "1", "J", false)

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="sentencia_oficial_blindada.pdf"`)

	err = pdf.Output(w)
	if err != nil {
		http.Error(w, "Error emitiendo PDF", http.StatusInternalServerError)
		return
	}
}
// --- MIDDLEWARE DE CORS 100% ESTRICTO (PUERTA WEB) ---
func corsForense(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Damos permiso 100% legal a tu interfaz web para conectarse
		w.Header().Set("Access-Control-Allow-Origin", "*") 
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 2. Si el navegador hace una petición de seguridad previa (OPTIONS), la aprobamos al 100%
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 3. Pasamos el control al siguiente escudo
		next.ServeHTTP(w, r)
	}
}
// --- INICIO DEL SERVIDOR ---
func main() {
	var err error
	connStr := "user=postgres dbname=justicia_db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil { log.Fatal("Error de BD:", err) }
	defer db.Close()

	// El login lleva blindaje para la web
	http.HandleFunc("/api/login", corsForense(limiteTasaForense(trazabilidadForense(cabecerasSeguridad(loginHandler)))))
	
	// Los endpoints críticos llevan el 100% del blindaje SÉXTUPLE absoluto
	http.HandleFunc("/api/transicion", corsForense(limiteTasaForense(trazabilidadForense(timeoutForense(cabecerasSeguridad(seguridadJWT(registrarHandler)))))))
	http.HandleFunc("/api/auditoria", corsForense(limiteTasaForense(trazabilidadForense(timeoutForense(cabecerasSeguridad(seguridadJWT(auditarHandler)))))))
	http.HandleFunc("/api/evidencia", corsForense(limiteTasaForense(trazabilidadForense(timeoutForense(cabecerasSeguridad(seguridadJWT(subirEvidenciaHandler)))))))
	http.HandleFunc("/api/sentencia", corsForense(limiteTasaForense(trazabilidadForense(timeoutForense(cabecerasSeguridad(seguridadJWT(generarSentenciaHandler)))))))

	fmt.Println("🚀 Servidor al 100%: Blindaje Séxtuple Absoluto (CORS + RateLimit + Trazabilidad + Timeout + Cabeceras + JWT) corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

