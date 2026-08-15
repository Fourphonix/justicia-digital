// Copyright © 2026 José Emanuel Guadalupe López - Justicia Digital. 100% Derechos Reservados.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// GenerarFirmaArchivo lee un archivo al 100% y extrae su huella digital inalterable
func GenerarFirmaArchivo(ruta string) (string, error) {
	archivo, err := os.Open(ruta)
	if err != nil {
		return "", err
	}
	defer archivo.Close()

	motor := sha256.New()
	if _, err := io.Copy(motor, archivo); err != nil {
		return "", err
	}

	// Extraer la firma SHA-256 en formato hexadecimal
	return hex.EncodeToString(motor.Sum(nil)), nil
}

func main() {
	directorioRaiz := "." // Escanea la bóveda actual al 100%
	archivoSalida := "registro_forense_maestro.txt"

	// Crear o sobreescribir el archivo de evidencia
	registro, err := os.Create(archivoSalida)
	if err != nil {
		fmt.Printf("ERROR TÁCTICO: No se pudo crear el registro forense: %v\n", err)
		return
	}
	defer registro.Close()

	// Capturar el tiempo universal al 100%
	tiempoExacto := time.Now().UTC().Format(time.RFC3339)
	encabezado := fmt.Sprintf("========================================================\n" +
		"🛡️  BÓVEDA 'JUSTICIA DIGITAL' - INVENTARIO FORENSE AL 100%%\n" +
		"AUTORIDAD HUMANA: Jose-Emanuel-GL-001\n" +
		"FECHA DEL SELLADO (UTC): %s\n" +
		"ESTADO: 100%% LEGAL. 100%% DEBIDO. 100%% PROPIEDAD INTELECTUAL.\n" +
		"========================================================\n\n", tiempoExacto)

	// Imprimir en consola y guardar en el archivo
	registro.WriteString(encabezado)
	fmt.Println(encabezado)

	// Escanear el 100% de los archivos y subcarpetas
	err = filepath.Walk(directorioRaiz, func(ruta string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Ignorar carpetas y el propio archivo de registro para evitar bucles
		if !info.IsDir() && info.Name() != archivoSalida {
			firma, errFirma := GenerarFirmaArchivo(ruta)
			if errFirma == nil {
				lineaLegal := fmt.Sprintf("[SELLO SHA-256] %s | ARCHIVO: %s\n", firma, ruta)
				
				// Inyectar la evidencia en el documento
				registro.WriteString(lineaLegal)
				// Mostrar la evidencia en la terminal en tiempo real
				fmt.Print(lineaLegal)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("\n⚠️ ADVERTENCIA: Brecha en el escaneo: %v\n", err)
	} else {
		pie := "\n========================================================\n" +
			"✅ ESCANEO COMPLETADO AL 100%. EVIDENCIA GUARDADA FÍSICAMENTE.\n" +
			"========================================================\n"
		registro.WriteString(pie)
		fmt.Println(pie)
	}
}
