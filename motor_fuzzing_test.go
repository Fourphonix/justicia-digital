// Copyright © 2026 José Emanuel Guadalupe López - Justicia Digital. 100% Derechos Reservados.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"unicode/utf8"
)

// =====================================================================
// EL NÚCLEO A EVALUAR (La "Mejora" que estás a punto de desplegar)
// =====================================================================
// Esta función simula el procesamiento de evidencia. Debe ser 100% invulnerable.
func ProcesarEvidenciaMejorada(cargaMaliciosa string) string {
	// Si la carga no es texto UTF-8 válido, la rechazamos instantáneamente
	if !utf8.ValidString(cargaMaliciosa) {
		return "ERROR_FORENSE: CODIFICACIÓN_INVÁLIDA"
	}
	
	motor := sha256.New()
	motor.Write([]byte(cargaMaliciosa))
	return hex.EncodeToString(motor.Sum(nil))
}

// =====================================================================
// MOTOR DE FUZZING (Auto-Ataque de Mutación Aleatoria al 100%)
// =====================================================================
// FuzzMotorCriptografico bombardeará la función con millones de combinaciones.
func FuzzMotorCriptografico(f *testing.F) {
	// 1. Semillas Tácticas (Vectores de ataque base)
	f.Add("Expediente_Normal_001") // Caso estándar
	f.Add("'; DROP TABLE boveda; --") // Inyección SQL
	f.Add("<script>alert('hack')</script>") // Cross-Site Scripting
	f.Add("\x00\x00\x00\x00") // Desbordamiento de memoria (Null bytes)
	f.Add("                                    ") // Espacios masivos

	f.Log("⚙️ [INICIANDO FUZZING] Generando mutaciones de ataque al 100%...")

	// 2. El Bucle de Asedio
	f.Fuzz(func(t *testing.T, intentoAtaque string) {
		resultado := ProcesarEvidenciaMejorada(intentoAtaque)

		// Verificación Táctica 1: Prevención de Colapsos
		if resultado == "" {
			t.Errorf("🚨 BRECHA CRÍTICA: Un ataque mutado colapsó el motor y devolvió un vacío. Ataque: %q", intentoAtaque)
		}

		// Verificación Táctica 2: Integridad del Hash
		// Si el resultado no es un error forense, DEBE ser un Hash SHA-256 válido (64 caracteres)
		if resultado != "ERROR_FORENSE: CODIFICACIÓN_INVÁLIDA" && len(resultado) != 64 {
			t.Errorf("🚨 VULNERABILIDAD MATEMÁTICA: La mejora rompió la criptografía. Longitud: %d. Ataque: %q", len(resultado), intentoAtaque)
		}
	})
}
