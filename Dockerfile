# ==============================================================================
# ETAPA 1: Constructor (Builder) - 100% Aislado
# ==============================================================================
# Reemplazar la versión antigua por la versión correspondiente
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias y descargarlas
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el código fuente
COPY . .

# Compilar el binario del backend
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-w -s" -o /app/secure-backend ./backend/main.go

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
