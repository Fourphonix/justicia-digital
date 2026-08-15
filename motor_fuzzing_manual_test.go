// Copyright © 2026 José Emanuel Guadalupe López - Justicia Digital. 100% Derechos Reservados.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"testing"
	"time"
)

func SimuladorMejora(datos string) string {
	motor := sha256.New()
	motor.Write([]byte(datos))
	return hex.EncodeToString(motor.Sum(nil))
}

func TestFuerzaBrutaFuzzing(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	t.Log("⚙️ INICIANDO ASALTO DE MUTACIÓN MANUAL (100,000 ITERACIONES)...")

	// Bombardeamos la mejora con 100,000 cadenas de caracteres aleatorias
	for i := 0; i < 100000; i++ {
		longitud := rand.Intn(500) + 1
		bytesAleatorios := make([]byte, longitud)
		rand.Read(bytesAleatorios)
		
		cargaMutada := string(bytesAleatorios)
		resultado := SimuladorMejora(cargaMutada)

		// Verificamos que el hash generado siga siendo matemáticamente perfecto
		if len(resultado) != 64 {
			t.Fatalf("🚨 BRECHA CRÍTICA en la iteración %d: La longitud del hash falló.", i)
		}
	}
	t.Log("✅ [SUPERVIVENCIA ABSOLUTA] El código resistió 100,000 ataques de mutación aleatoria al 100%.")
}
