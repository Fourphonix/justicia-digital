# ==============================================================================
# ETAPA 1: Constructor (Builder) - 100% Aislado
# ==============================================================================
# Reemplazar la versión antigua por la versión correspondiente
FROM golang:1.24-alpine AS builder

# Establecer el directorio de trabajo seguro
WORKDIR /app

# Copiar única y exclusivamente los archivos de dependencias primero (Caché seguro)
COPY go.mod go.sum ./

# Descargar las dependencias y verificar su integridad criptográfica al 100%
RUN go mod download && go mod verify

# Copiar el resto del código fuente
COPY . .

# Compilación 100% Segura y Estática:
# CGO_ENABLED=0 : Desactiva dependencias de C, creando un binario 100% estático.
# -trimpath     : Elimina las rutas absolutas de tu servidor de los rastros de error (Privacidad).
# -ldflags="-w -s": Elimina la tabla de símbolos y la información de depuración (Anti Ingeniería Inversa).
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-w -s" -o /app/secure-backend .

# ==============================================================================
# ETAPA 2: Producción (Distroless) - 100% Reducido y Seguro
# ==============================================================================
# Usamos una imagen 'static-debian' de Distroless que no tiene consola (sh/bash)
# La versión 'nonroot' nos asegura ejecutar el proceso con el menor privilegio posible.
FROM gcr.io/distroless/static-debian12:nonroot

# Etiquetado legal para auditorías
LABEL maintainer="tu-empresa-segura"
LABEL compliance="hipaa,pci-dss,gdpr"

# Establecer el directorio de trabajo
WORKDIR /

# Copiar el binario compilado desde la ETAPA 1. No se copia el código fuente.
COPY --from=builder --chown=nonroot:nonroot /app/secure-backend /secure-backend

# Forzar el uso del usuario 'nonroot' (ID 65532)
# NUNCA ejecutar contenedores como root.
USER nonroot:nonroot

# Exponer el puerto de tu API (ej. 8443 para HTTPS)
EXPOSE 8443

# Ejecutar el binario al iniciar el contenedor
ENTRYPOINT ["/secure-backend"]
