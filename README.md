# Tonalmaster Backend

Backend base de Tonalmaster siguiendo la arquitectura y planeación documentadas en `docs/`.

## Fase actual

Fase 1: base ejecutable. Incluye configuración por entorno, PostgreSQL, Docker Compose, logs estructurados y endpoints de salud.

## Desarrollo local

1. Copia `.env.example` a `.env` y ajusta valores si es necesario.
2. Ejecuta `make up`.
3. Comprueba:
   - `GET http://localhost:8080/health`
   - `GET http://localhost:8080/ready`

## Base de datos

La Fase 2 incorpora la migración inicial en `migrations/000001_init_schema.up.sql` y su rollback correspondiente. El `docker-compose.yml` prepara PostgreSQL; la ejecución de migraciones se integrará antes de habilitar los repositorios de dominio.

## Siguientes fases

- Fase 2: migraciones y modelo mínimo.
- Fase 3: dominio calendárico.
- Fase 4: API v1.
- Fase 5: pruebas y contrato.
- Fase 6: integración y evolución.

Consulta `docs/planeacion.md` para el alcance.
