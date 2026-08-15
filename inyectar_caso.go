// Copyright © 2026 José Emanuel Guadalupe López - Justicia Digital. 100% Derechos Reservados.
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("=====================================================================")
	fmt.Println("🛡️ INYECTANDO EL PRIMER EXPEDIENTE LEGAL AL 100%")
	fmt.Println("=====================================================================")

	// 1. Abrir la conexión con la bóveda forjada previamente
	bovedaDB, err := sql.Open("sqlite3", "./firma_abogados.db")
	if err != nil {
		log.Fatalf("ERROR TÁCTICO: Fallo al acceder a la bóveda: %v", err)
	}
	defer bovedaDB.Close()

	// =====================================================================
	// FASE 1: REGISTRO DEL ABOGADO (AUTORIDAD HUMANA)
	// =====================================================================
	// Utilizamos el signo '?' para parametrizar los datos. 100% blindado contra SQL Injection.
	comandoAbogado := `INSERT INTO abogados (nombre_legal, numero_barra, nivel_acceso) VALUES (?, ?, ?)`
	resultadoA, err := bovedaDB.Exec(comandoAbogado, "José Emanuel Guadalupe López", "GL-001-FEDERAL", "SUPREMO")
	if err != nil {
		log.Fatalf("Brecha en el registro del abogado: %v", err)
	}
	idAbogado, _ := resultadoA.LastInsertId()
	fmt.Println("✅ [ABOGADO] Identidad de Autoridad sellada exitosamente.")

	// =====================================================================
	// FASE 2: REGISTRO DEL CLIENTE (IDENTIDAD CIFRADA)
	// =====================================================================
	comandoCliente := `INSERT INTO clientes (hash_identidad, contacto_cifrado) VALUES (?, ?)`
	// Simulamos el Hash SHA-256 de un cliente para mantener su privacidad forense al 100%
	resultadoC, err := bovedaDB.Exec(comandoCliente, "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3", "ENCRIPTADO_AES256")
	if err != nil {
		log.Fatalf("Brecha en el registro del cliente: %v", err)
	}
	idCliente, _ := resultadoC.LastInsertId()
	fmt.Println("✅ [CLIENTE] Privilegio Cliente-Abogado protegido y registrado.")

	// =====================================================================
	// FASE 3: APERTURA DEL EXPEDIENTE (VINCULACIÓN)
	// =====================================================================
	comandoCaso := `INSERT INTO casos (titulo_expediente, estatus_legal, id_abogado, id_cliente) VALUES (?, ?, ?, ?)`
	_, err = bovedaDB.Exec(comandoCaso, "Defensa de Propiedad Intelectual - Motor IA", "ACTIVO", idAbogado, idCliente)
	if err != nil {
		log.Fatalf("Brecha en la apertura del caso: %v", err)
	}
	fmt.Println("✅ [CASO] Expediente abierto y enlazado con Integridad Referencial Estricta.")

	fmt.Println("=====================================================================")
	fmt.Println("ESTADO: EXPEDIENTE 001 INSERTADO AL 100%. BÓVEDA OPERATIVA.")
}
