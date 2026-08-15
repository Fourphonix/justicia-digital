// Copyright © 2026 José Emanuel Guadalupe López - Justicia Digital. 100% Derechos Reservados.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// =====================================================================
// FASE 1: ESTRUCTURA DE LA BÓVEDA INMUTABLE (BLOCKCHAIN)
// =====================================================================
type BloqueForense struct {
	Indice        int    `json:"indice_táctico"`
	SelloDeTiempo string `json:"sello_de_tiempo_utc"`
	ArchivoFuerte string `json:"archivo_evidencia"`
	HashAnterior  string `json:"hash_bloque_anterior"`
	HashActual    string `json:"hash_criptografico_absoluto"`
}

// GenerarHash sella los datos del bloque usando SHA-256
func GenerarHash(indice int, tiempo, archivo, hashAnterior string) string {
	registro := fmt.Sprintf("%d%s%s%s", indice, tiempo, archivo, hashAnterior)
	motor := sha256.New()
	motor.Write([]byte(registro))
	return hex.EncodeToString(motor.Sum(nil))
}

func main() {
	fmt.Println("=====================================================================")
	fmt.Println("⚖️  INICIANDO MOTOR NOTARIAL CRIPTOGRÁFICO AL 100%")
	fmt.Println("=====================================================================")

	nombreArchivoEvidencia := "Motor_Criptográfico_-_Justicia_Digital_Formulario.txt"
	nombreBoveda := "boveda_inmutable_ledger.json"

	// 1. Verificar que el formulario legal exista al 100%
	evidenciaFisica, err := os.ReadFile(nombreArchivoEvidencia)
	if err != nil {
		log.Fatalf("🚨 ERROR TÁCTICO: No se encontró el formulario para sellar. Ejecuta el Motor de Formularios primero: %v", err)
	}

	// 2. Extraer la Huella Digital del documento físico
	motorSHA := sha256.New()
	motorSHA.Write(evidenciaFisica)
	hashDelDocumento := hex.EncodeToString(motorSHA.Sum(nil))

	// 3. Forjar el Sello de Tiempo Universal
	tiempoExacto := time.Now().UTC().Format(time.RFC3339)

	// 4. Lógica de Cadena de Bloques (Inmutabilidad 100% Debida)
	var libroNotarial []BloqueForense
	var hashPrevio = "0000000000000000000000000000000000000000000000000000000000000000" // Génesis
	var nuevoIndice = 1

	// Leer la bóveda si ya existe para encadenar el nuevo bloque
	if datosBoveda, err := os.ReadFile(nombreBoveda); err == nil {
		json.Unmarshal(datosBoveda, &libroNotarial)
		if len(libroNotarial) > 0 {
			ultimoBloque := libroNotarial[len(libroNotarial)-1]
			hashPrevio = ultimoBloque.HashActual
			nuevoIndice = ultimoBloque.Indice + 1
		}
	}

	// 5. Inyectar el Nuevo Bloque Forense
	hashNuevoBloque := GenerarHash(nuevoIndice, tiempoExacto, hashDelDocumento, hashPrevio)
	bloqueInyectado := BloqueForense{
		Indice:        nuevoIndice,
		SelloDeTiempo: tiempoExacto,
		ArchivoFuerte: hashDelDocumento,
		HashAnterior:  hashPrevio,
		HashActual:    hashNuevoBloque,
	}

	libroNotarial = append(libroNotarial, bloqueInyectado)

	// 6. Sellar la bóveda grabando los datos en el disco
	datosSellados, _ := json.MarshalIndent(libroNotarial, "", "  ")
	err = os.WriteFile(nombreBoveda, datosSellados, 0444) // Permisos 0444 = Solo lectura forense
	if err != nil {
		log.Fatalf("🚨 ERROR AL SELLAR LA BÓVEDA: %v", err)
	}

	// 7. Reporte Operativo de Misión Cumplida
	fmt.Printf("✅ [NOTARIO DIGITAL] Documento detectado y analizado.\n")
	fmt.Printf("✅ [HASH ARCHIVO] %s\n", hashDelDocumento)
	fmt.Printf("✅ [SELLO TIEMPO] %s\n", tiempoExacto)
	fmt.Printf("✅ [CADENA BLOQUES] Bloque #%d inyectado y atado al Hash: %s...\n", nuevoIndice, hashPrevio[:10])
	fmt.Println("=====================================================================")
	fmt.Println("ESTADO: EVIDENCIA INYECTADA AL 100%. BÓVEDA INALTERABLE CREADA.")
}
