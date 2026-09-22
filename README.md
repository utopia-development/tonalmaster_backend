# Tonalmaster Backend

Backend base de Tonalmaster siguiendo la arquitectura y planeación documentadas en `docs/`.

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


## Fase actual

**Fase 8 — verificación, compatibilidad y aceptación de integración.**

Las Fases 1–6 están cerradas y la implementación de la Fase 7 está completada. La Fase 8 verifica el contrato con `ule_educativo`, el flujo Docker/WSL/XAMPP y las regresiones del backend.

## Integración Windows + WSL 2 + XAMPP

La API escucha en `0.0.0.0:8080` dentro del contenedor. Prueba primero desde Windows:

```text
http://localhost:8080/health
```

En `ule_educativo`, usa `dataSource = 'api'` y:

```text
ULE.config.apiBaseUrl = 'http://localhost:8080/api/v1'
```

Si Windows no alcanza el puerto reenviado por WSL2, obtén la IP con `hostname -I` dentro de Ubuntu y usa esa dirección. En ambos casos, el origen del frontend debe estar permitido en `CORS_ALLOWED_ORIGINS`. Para XAMPP en `http://localhost`:

```text
CORS_ALLOWED_ORIGINS=http://localhost,http://127.0.0.1
```

CORS valida el **origen del navegador**, no la URL de la API. Si XAMPP usa otro puerto, incluye el origen exacto, por ejemplo `http://localhost:8081`.

### Cómo impacta CI en la integración con el frontend

CI no ejecuta las páginas de Ule, pero garantiza la condición previa para que puedan consumir el backend: compila Go, ejecuta los tests y levanta una PostgreSQL limpia para aplicar `000001`, `000002` y `000003`, y después revierte las migraciones en orden inverso.

El flujo correcto es:

```text
cambio backend
   -> go test / go build
   -> CI + migraciones reproducibles
   -> Docker Compose (db -> migrate -> api)
   -> navegador Windows/XAMPP
   -> ULE.loader -> GET /api/v1/... -> PostgreSQL
```

Así evitamos que el frontend tenga que compensar errores de esquema, rutas o DTO. La prueba con XAMPP cubre lo que CI no puede cubrir: CORS real, navegador, 404/5xx, caída de red y resolución de imágenes.

## Datos que se cargan mediante migraciones

`000003_ule_content.up.sql` crea el esquema de contenido; no conviene convertir cada publicación editorial en una migración.

Para **datos fijos de desarrollo/demo**, crea una migración posterior, por ejemplo `000004_ule_seed_dev.up.sql`, con `INSERT ... ON CONFLICT ...` y un `000004_ule_seed_dev.down.sql` que elimine esos datos. Compose la ejecutará automáticamente después de `000003`.

Para **contenido editorial real**, la dirección prevista es usar posteriormente endpoints de escritura autenticados. Así una migración define estructura reproducible y el contenido editorial puede cambiar sin generar una migración por cada artículo.

Ejemplo conceptual:

```sql
INSERT INTO articles (id, titulo, fecha, resumen, contenido_html, visible)
VALUES ('articulo-001', 'Título de prueba', '2026-09-22', 'Resumen', '<p>Contenido</p>', TRUE)
ON CONFLICT (id) DO UPDATE SET
  titulo = EXCLUDED.titulo,
  resumen = EXCLUDED.resumen,
  contenido_html = EXCLUDED.contenido_html,
  visible = EXCLUDED.visible;
```

```bash
make reset-db
make up
```

Con eso el dato queda cargado de forma reproducible desde una base vacía.


## Contenido público Ule

La API expone las rutas públicas `/api/v1/articles`, `/api/v1/bibliography`, `/api/v1/catalogs` y `/api/v1/ads` según `docs/contrato_datos.md`. Estas rutas no requieren sesión. La migración `000003_ule_content` se ejecuta automáticamente mediante el servicio `migrate` de Docker Compose.
