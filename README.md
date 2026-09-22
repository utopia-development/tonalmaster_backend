# Tonalmaster Backend

Backend base de Tonalmaster siguiendo la arquitectura y planeación documentadas en `docs/`.

## Fase actual

Fases 1–5 cerradas (base, migraciones, dominio calendárico CASO, API v1, tests de contrato).

Auth por sesión ya cableada (ver `docs/auth-sesiones.md`); forma parte del inicio de la Fase 6.

## Desarrollo local

1. Copia `.env.example` a `.env` y ajusta valores si es necesario.
2. Ejecuta `make up`.
3. `make up` levanta PostgreSQL, ejecuta automáticamente las migraciones y después inicia la API. No necesitas `psql` para el flujo normal.

   Para una instalación limpia desde cero puedes usar `make reset-db` y después `make up`.

4. Comprueba:
   - `GET http://localhost:8080/health`
   - `GET http://localhost:8080/ready`
   - `GET http://localhost:8080/api/v1/calendars`

`/health`, `/ready` y los endpoints de `calendars` no dependen del esquema de base de datos y responderán aunque no hayas aplicado migraciones. **Los endpoints de `auth`, `events` e `interpretations` sí las requieren.**

El flujo normal no requiere aplicar migraciones manualmente. El servicio `migrate` de Compose espera a que PostgreSQL esté saludable y ejecuta las migraciones versionadas antes de arrancar la API.

Para desarrollo avanzado, las migraciones también pueden ejecutarse manualmente:

   ```bash
   psql -h localhost -U ule_user -d ule_tonalmaster -f migrations/000001_init_schema.up.sql
   psql -h localhost -U ule_user -d ule_tonalmaster -f migrations/000002_sessions.up.sql
   ```

   Este es el mismo procedimiento que ejecuta CI (`.github/workflows/backend.yml`). Aún no hay ejecución automática de migraciones al arrancar el servicio (ver `docs/planeacion.md`, Fase 6).

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

Los endpoints de escritura de eventos e interpretaciones requieren autenticación. Las sesiones y sus límites están documentados en `docs/auth-sesiones.md`.

- `GET /api/v1/calendars`
- `GET /api/v1/calendars/{id}`
- `GET /api/v1/calendars/convert?date=YYYY-MM-DD&system=tonalpohualli_caso`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/me`
- `POST /api/v1/auth/logout`
- `POST /api/v1/events` (auth)
- `GET /api/v1/events` (auth)
- `DELETE /api/v1/events/{id}` (auth)
- `POST /api/v1/interpretations` (auth)
- `GET /api/v1/interpretations?system=...&date=YYYY-MM-DD`

Los errores usan el formato JSON `{"error":{"code":"...","message":"..."}}`.

Detalle de sesiones: `docs/auth-sesiones.md`.


## Contenido público Ule

La API expone las rutas públicas `/api/v1/articles`, `/api/v1/bibliography`, `/api/v1/catalogs` y `/api/v1/ads` según `docs/contrato_datos.md`. Estas rutas no requieren sesión. La migración `000003_ule_content` se ejecuta automáticamente mediante el servicio `migrate` de Docker Compose.
