# Tonalmaster Backend

Backend unificado de Tonalmaster + Ule, construido con Go + PostgreSQL y desplegable con Docker Compose.

## Estado

- Fases 1–6: cerradas.
- Fase 7: contenido público Ule cerrada.
- Fase 8: integración Windows/WSL/XAMPP cerrada.
- Siguiente: **Fase 9 — pgAdmin para administración de PostgreSQL en desarrollo**.
- Después: **Fase 10 — API editorial autenticada**.

La planeación completa está en `docs/planeacion.md`.

## Desarrollo local

Requisitos: Docker y Docker Compose. En Windows, el backend puede ejecutarse dentro de WSL2.

```bash
cp .env.example .env
make up
```

Compose levanta:

```text
db -> migrate -> api
 \-> pgadmin :5050
```

pgAdmin es una herramienta de administración de desarrollo. No forma parte de la API ni del dominio.

La base se crea y las migraciones versionadas se aplican automáticamente.

Comprobaciones:

```text
GET http://localhost:8080/health
GET http://localhost:8080/ready
GET http://localhost:8080/api/v1/calendars
```

Para reconstruir PostgreSQL desde cero:

```bash
make reset-db
make up
```

No es necesario ejecutar `psql` manualmente para el flujo normal.

## Windows + WSL2 + XAMPP

La API se publica en el puerto `8080`.

En `ule_educativo`:

```js
ULE.config = ULE.config || {};
ULE.config.dataSource = 'api';
ULE.config.apiBaseUrl = 'http://localhost:8080/api/v1';
```

Si Windows no alcanza el reenvío de WSL2, puede utilizarse la IP obtenida con:

```bash
hostname -I
```

El origen del frontend debe estar en `CORS_ALLOWED_ORIGINS`. Para XAMPP en `http://localhost`:

```text
CORS_ALLOWED_ORIGINS=http://localhost,http://127.0.0.1
```

CORS valida el origen del navegador, no la URL de la API.

### Flujo de integración

```text
XAMPP / navegador
      |
      | GET /api/v1/...
      v
Docker -> API Go -> PostgreSQL
      |
      +--> JSON del contrato Ule
```

El frontend no conoce PostgreSQL y no necesita modificar sus componentes para cambiar de datos locales a API.

## API pública Ule

Las siguientes rutas son públicas y de solo lectura:

```text
GET /api/v1/articles
GET /api/v1/articles/{id}
GET /api/v1/bibliography
GET /api/v1/bibliography/{id}
GET /api/v1/catalogs
GET /api/v1/catalogs/{id}
GET /api/v1/ads
```

El contrato vinculante está en `docs/contrato_datos_ule.md`.

Las URLs de las rutas son inglesas; las claves JSON conservan los nombres españoles definidos por el contrato, incluidos `bibliografía_relacionada`, `año` y `año_descubrimiento`.

## API Tonalmaster existente

```text
GET  /api/v1/calendars
GET  /api/v1/calendars/{id}
GET  /api/v1/calendars/convert?date=YYYY-MM-DD&system=tonalpohualli_caso

POST /api/v1/auth/register
POST /api/v1/auth/login
GET  /api/v1/auth/me
POST /api/v1/auth/logout

POST /api/v1/events
GET  /api/v1/events
DELETE /api/v1/events/{id}

POST /api/v1/interpretations
GET  /api/v1/interpretations?system=...&date=YYYY-MM-DD
```

Los endpoints que modifican datos requieren autenticación. Las sesiones están documentadas en `docs/auth-sesiones.md`.

## Datos y migraciones

Las migraciones definen estructura reproducible. La migración `000003_ule_content` crea el esquema de contenido Ule.

Para datos fijos de desarrollo/demo, una migración posterior puede contener `INSERT ... ON CONFLICT`. El contenido editorial real no debe convertirse en una migración por publicación: la futura API editorial será la vía de escritura.

## Administración de PostgreSQL

**Fase 9 — pgAdmin 4:** el proyecto incluye pgAdmin como servicio Docker de desarrollo, accesible desde Windows mediante `http://localhost:5050`.

En el primer acceso usa las credenciales definidas en `.env`:

```text
PGADMIN_DEFAULT_EMAIL
PGADMIN_DEFAULT_PASSWORD
```

Para registrar la base en pgAdmin:

```text
Host: db
Port: 5432
Database: ${POSTGRES_DB}
Username: ${POSTGRES_USER}
Password: ${POSTGRES_PASSWORD}
```

El host `db` es correcto dentro de Docker Compose; no uses `localhost` para la conexión de pgAdmin a PostgreSQL.

PostgreSQL actualmente conserva su publicación local para desarrollo. No debe exponerse al LAN como parte del despliegue normal. pgAdmin se comunica con PostgreSQL por la red interna de Compose.

## API editorial — próxima fase

La futura API editorial será autenticada y autorizada. Incluirá:

```text
POST/PUT/DELETE /api/v1/articles
POST/PUT/DELETE /api/v1/bibliography
POST/PUT/DELETE /api/v1/catalogs
POST/PUT/DELETE /api/v1/catalogs/{id}/items
POST/PUT/DELETE /api/v1/ads
```

El registro de usuarios tendrá inicialmente un código de verificación estático controlado por configuración/código del backend. Es una medida temporal para impedir registros anónimos indiscriminados; posteriormente podrá migrarse a DB, invitaciones o retirarse cuando el registro público sea intencional.

No se habilitará escritura anónima para el contenido Ule.

## Arquitectura

La arquitectura separa:

```text
handlers -> services -> repositories -> PostgreSQL
```

El dominio no depende de HTTP ni SQL. Las interfaces pequeñas y la inyección de dependencias se mantienen como principios del proyecto.

Consulta:

- `docs/arquitectura.md`
- `docs/planeacion.md`
- `docs/contrato_datos_ule.md`
- `docs/auth-sesiones.md`

## CI

CI ejecuta:

```text
go test ./...
go build ./...
migrations 000001 -> 000002 -> 000003
rollback  000003 -> 000002 -> 000001
```

El objetivo es garantizar que un checkout limpio pueda reproducir el esquema antes de la integración con el frontend.
