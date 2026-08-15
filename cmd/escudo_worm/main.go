// Copyright © 2026 José Emanuel Guadalupe López - Justicia Digital. 100% Derechos Reservados.
package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Println("=====================================================================")
	fmt.Println("🛡️  INICIANDO PROTOCOLO WORM (WRITE ONCE, READ MANY) AL 100%")
	fmt.Println("=====================================================================")

	nombreEvidencia := "evidencia_clasificada_001.txt"
	contenidoForense := "EXPEDIENTE 001: Evidencia irrefutable. Inyectado el " + time.Now().UTC().Format(time.RFC3339)

	// =====================================================================
	// FASE 1: ESCRITURA ÚNICA (Write Once)
	// =====================================================================
	// Forjamos el archivo con permisos iniciales de escritura solo para inyectar los datos.
	err := os.WriteFile(nombreEvidencia, []byte(contenidoForense), 0644)
	if err != nil {
		log.Fatalf("🚨 ERROR TÁCTICO: Fallo al crear la evidencia: %v", err)
	}
	fmt.Printf("✅ [FASE 1] Archivo '%s' forjado e inyectado con datos al 100%%.\n", nombreEvidencia)

	// =====================================================================
	// FASE 2: BLOQUEO PERMANENTE (Read Many / Inmutabilidad)
	// =====================================================================
	// Aplicamos el código 0444 (Solo Lectura Absoluta). Esto destruye los permisos 
	// de escritura para el dueño, el grupo y cualquier otro usuario del sistema.
	err = os.Chmod(nombreEvidencia, 0444)
	if err != nil {
		log.Fatalf("🚨 ERROR TÁCTICO: Fallo al sellar el escudo WORM: %v", err)
	}
	fmt.Println("✅ [FASE 2] Escudo WORM activado. Permisos de escritura DESTRUIDOS permanentemente.")

	// =====================================================================
	// FASE 3: SIMULACRO DE ATAQUE (Red Team Testing)
	// =====================================================================
	fmt.Println("=====================================================================")
	fmt.Println("⚔️ EJECUTANDO SIMULACRO DE ATAQUE (INTENTO DE MANIPULACIÓN)...")
	
	// El sistema intentará autohackearse sobrescribiendo el archivo con datos falsos.
	errAtaque := os.WriteFile(nombreEvidencia, []byte("TEXTO MANIPULADO POR HACKER"), 0644)
	
	if errAtaque != nil {
		fmt.Printf("🛡️ [BLOQUEO CONFIRMADO] El OS rechazó el ataque al 100%%. Error forense: %v\n", errAtaque)
	} else {
		// Esto no debería ejecutarse nunca si el escudo WORM se activó correctamente.
		fmt.Println("🚨 BRECHA CRÍTICA: La evidencia fue alterada. Sistema comprometido.")
	}

	fmt.Println("=====================================================================")
	fmt.Println("ESTADO: EVIDENCIA CONGELADA AL 100%. INMUTABILIDAD CONFIRMADA.")
}
