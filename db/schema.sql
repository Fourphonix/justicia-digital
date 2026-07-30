-- 1. Habilitar extensión para UUIDs
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 2. Tabla principal de Expedientes
CREATE TABLE IF NOT EXISTS expedientes (
    id_expediente UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    numero_causa VARCHAR(50) UNIQUE NOT NULL,
    fuero VARCHAR(30) NOT NULL,
    estado_actual VARCHAR(30) NOT NULL DEFAULT 'INICIADO',
    fecha_ultimo_estado TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    es_archivado BOOLEAN DEFAULT FALSE,
    creado_en TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Tabla inmutable de Historial de Estados (Append-Only)
CREATE TABLE IF NOT EXISTS historial_estados_expediente (
    id_historial UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_expediente UUID NOT NULL REFERENCES expedientes(id_expediente),
    estado_anterior VARCHAR(30),
    estado_nuevo VARCHAR(30) NOT NULL,
    motivo_cambio TEXT,
    usuario_id UUID NOT NULL, -- ID del Juez, Secretario o Sistema
    hash_transicion CHAR(64) NOT NULL, -- Hash SHA-256 para integridad
    registrado_en TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 4. Índice de alto rendimiento para consultas históricas
CREATE INDEX IF NOT EXISTS idx_historial_expediente_fecha 
ON historial_estados_expediente (id_expediente, registrado_en DESC);

-- 5. Función de Inmutabilidad para la tabla Historial (Regla WORM)
CREATE OR REPLACE FUNCTION denegar_modificacion_historial()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'OPERACIÓN RECHAZADA: La tabla historial_estados_expediente es inalterable (Modelo WORM). No se permiten UPDATE ni DELETE.';
END;
$$ LANGUAGE plpgsql;

-- 6. Disparador (Trigger) de Inviolabilidad
DROP TRIGGER IF EXISTS trg_proteger_historial ON historial_estados_expediente;
CREATE TRIGGER trg_proteger_historial
BEFORE UPDATE OR DELETE ON historial_estados_expediente
FOR EACH ROW
EXECUTE FUNCTION denegar_modificacion_historial();
