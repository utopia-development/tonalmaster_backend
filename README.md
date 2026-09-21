# Tonalmaster Backend

Backend base de Tonalmaster siguiendo la arquitectura y planeación documentadas en `docs/`.

## Fase actual

Fases 1–5 cerradas (base, migraciones, dominio calendárico CASO, API v1, tests de contrato).

Auth por sesión ya cableada (ver `docs/auth-sesiones.md`); forma parte del inicio de la Fase 6.

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

### API v1

- `GET /api/v1/calendars`
- `GET /api/v1/calendars/{id}`
- `GET /api/v1/calendars/convert?date=YYYY-MM-DD&system=tonalpohualli_caso`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/me`
- `POST /api/v1/auth/logout`

Los errores usan el formato JSON `{"error":{"code":"...","message":"..."}}`.

Detalle de sesiones: `docs/auth-sesiones.md`.
