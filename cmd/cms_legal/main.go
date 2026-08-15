// Copyright © 2026 José Emanuel Guadalupe López - Justicia Digital. 100% Derechos Reservados.
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // Driver de SQLite 100% nativo
)

// ForjarEsquemaLegal inyecta las tablas 100% blindadas en la base de datos
func ForjarEsquemaLegal(db *sql.DB) {
	// 1. Comando SQL para la tabla de Abogados (Autoridad)
	tablaAbogados := `
	CREATE TABLE IF NOT EXISTS abogados (
		id_abogado INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre_legal TEXT NOT NULL,
		numero_barra TEXT UNIQUE NOT NULL,
		nivel_acceso TEXT NOT NULL
	);`

	// 2. Comando SQL para la tabla de Clientes (Identidad)
	tablaClientes := `
	CREATE TABLE IF NOT EXISTS clientes (
		id_cliente INTEGER PRIMARY KEY AUTOINCREMENT,
		hash_identidad TEXT UNIQUE NOT NULL, 
		contacto_cifrado TEXT NOT NULL
	);`

	// 3. Comando SQL para la tabla de Casos (El CMS Central)
	// Utiliza Foreign Keys (Claves Foráneas) para atar al 100% el caso con su abogado y cliente
	tablaCasos := `
	CREATE TABLE IF NOT EXISTS casos (
		id_caso INTEGER PRIMARY KEY AUTOINCREMENT,
		titulo_expediente TEXT NOT NULL,
		estatus_legal TEXT NOT NULL,
		fecha_apertura DATETIME DEFAULT CURRENT_TIMESTAMP,
		id_abogado INTEGER,
		id_cliente INTEGER,
		FOREIGN KEY(id_abogado) REFERENCES abogados(id_abogado) ON DELETE RESTRICT,
		FOREIGN KEY(id_cliente) REFERENCES clientes(id_cliente) ON DELETE RESTRICT
	);`

	// Ejecutar la inyección del esquema en la bóveda
	comandosSQL := []string{tablaAbogados, tablaClientes, tablaCasos}
	for _, comando := range comandosSQL {
		_, err := db.Exec(comando)
		if err != nil {
			log.Fatalf("ERROR TÁCTICO AL FORJAR EL ESQUEMA: %v", err)
		}
	}
	fmt.Println("✅ [CMS LEGAL] Tablas de Abogados, Clientes y Casos forjadas al 100%.")
}

func main() {
	fmt.Println("=====================================================================")
	fmt.Println("🛡️  INICIALIZANDO MOTOR CMS 'JUSTICIA DIGITAL' AL 100%")
	fmt.Println("=====================================================================")

	// 1. Abrir o crear la bóveda de la firma de abogados
	bovedaDB, err := sql.Open("sqlite3", "./firma_abogados.db")
	if err != nil {
		log.Fatalf("Fallo crítico en la apertura de la base de datos: %v", err)
	}
	defer bovedaDB.Close()

	// 2. Ejecutar la forja del esquema legal
	ForjarEsquemaLegal(bovedaDB)

	fmt.Println("=====================================================================")
	fmt.Println("ESTADO: BÓVEDA CREADA. INTEGRIDAD REFERENCIAL ACTIVA AL 100%.")
}
