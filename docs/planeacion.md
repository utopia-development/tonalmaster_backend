# Planeación — Tonalmaster Backend

## Objetivo

Construir una API unificada en Go + PostgreSQL para Tonalmaster y Ule, con arquitectura simple, desacoplada y reproducible mediante Docker.

Regla rectora: **primero hacer funcionar el camino completo más pequeño; después extenderlo.**

## Estado actual

- **Fases 1–6:** CERRADAS.
- **Fase 7:** CERRADA; contenido público Ule implementado.
- **Fase 8:** CERRADA en integración local: backend, Docker, CORS y frontend Ule servido por XAMPP/Windows ya fueron probados con éxito.
- **Siguiente paso:** administración de base de datos y, después, API editorial autenticada.

## Fase 1 — Base ejecutable — CERRADA

- Go, configuración por entorno, PostgreSQL, Docker Compose, Makefile.
- `GET /health` y `GET /ready`.
- CORS explícito.
- Logs estructurados.

## Fase 2 — Modelo mínimo — CERRADA

Migraciones base para usuarios, sistemas calendáricos, interpretaciones y eventos, con relaciones, índices y restricciones.

## Fase 3 — Dominio calendárico — CERRADA

Sistema `tonalpohualli_caso`, lógica separada de HTTP/SQL y tests.

## Fase 4 — API v1 — CERRADA

Calendarios, autenticación, eventos e interpretaciones con DTOs, validación y errores JSON consistentes.

## Fase 5 — Tests y contrato — CERRADA

Tests de dominio/API, CI y verificación de arranque desde entorno limpio.

## Fase 6 — Integración y evolución — CERRADA

- Sesiones opacas en PostgreSQL.
- Middleware de autenticación.
- Events e interpretations autenticados.
- Migraciones automáticas mediante Compose.
- CORS explícito.
- Clientes desacoplados.

## Fase 7 — Contenido público Ule — CERRADA

Implementado:

- migración `000003_ule_content`;
- articles, bibliography, article_bibliography, catalogs, catalog_items y ads;
- repositorios, servicios y handlers separados;
- rutas públicas GET:
  - `/api/v1/articles`
  - `/api/v1/articles/{id}`
  - `/api/v1/bibliography`
  - `/api/v1/bibliography/{id}`
  - `/api/v1/catalogs`
  - `/api/v1/catalogs/{id}`
  - `/api/v1/ads`;
- relaciones artículo↔bibliografía;
- filtrado de visibilidad/vigencia en servidor;
- DTO compatible con `docs/contrato_datos_ule.md`.

La escritura editorial queda deliberadamente fuera de esta fase.

## Fase 8 — Verificación e integración — CERRADA

Se verificó:

- compilación y ejecución Docker;
- migraciones `000001–000003`;
- API Ule pública;
- CORS con `http://localhost`;
- frontend `ule_educativo` servido desde XAMPP/Windows;
- `ULE.config.dataSource = 'api'`;
- `ULE.config.apiBaseUrl = 'http://localhost:8080/api/v1'`;
- consumo real de artículos y anuncios;
- compatibilidad del loader con las rutas inglesas del backend.

Las pruebas de integración demostraron que el backend y frontend pueden operar desacoplados sin modificar los componentes de Ule.

---

# Futuro inmediato

## Fase 9 — Administración de PostgreSQL para desarrollo

**Objetivo:** disponer de una interfaz web local para inspeccionar y administrar PostgreSQL sin convertirla en parte de la API.

### Decisión

Usar **pgAdmin 4** desktop

garantizar la conectabilidad ya sea en entorno local o en producción de una maquina cualquiera con pgadmin a la base de datos usando contraseña

# Fase 10 — API editorial autenticada

**Objetivo:** permitir que usuarios autorizados administren el contenido de Ule mediante la API, manteniendo separado el sitio público de lectura y la escritura editorial.

### 10.1 Registro de usuarios

Agregar un flujo de registro controlado:

```text
POST /api/v1/auth/register
```

El registro solicitará:

- username;
- email;
- password;
- código de verificación.

El código de verificación será **estático y temporalmente almacenado en configuración/código del backend**. Se usará **únicamente durante `register`** para impedir la creación indiscriminada de cuentas. **No se usará en `login`**: el login será normal desde el inicio mediante email/username + contraseña y la sesión existente.

Más adelante el código podrá:

- migrarse a una variable de entorno/DB;
- sustituirse por invitaciones;
- eliminarse cuando el registro público sea permitido.

La contraseña debe almacenarse únicamente como hash seguro; nunca en texto plano.

### 10.2 Login y autorización

Reutilizar la sesión existente de Tonalmaster.

Definir permisos para distinguir como mínimo:

- usuario autenticado;
- editor/contributor;
- administrador.

La API pública Ule seguirá siendo de solo lectura y sin autenticación.

La API editorial requerirá autenticación y autorización.

### 10.3 Artículos

Implementar:

```text
POST   /api/v1/articles
GET    /api/v1/articles/{id}
PUT    /api/v1/articles/{id}
DELETE /api/v1/articles/{id}
```

El payload debe validarse contra el contrato Ule.

### 10.4 Bibliografía

```text
POST   /api/v1/bibliography
PUT    /api/v1/bibliography/{id}
DELETE /api/v1/bibliography/{id}
```

Además, administrar la relación artículo↔bibliografía.

### 10.5 Catálogos

```text
POST   /api/v1/catalogs
PUT    /api/v1/catalogs/{id}
DELETE /api/v1/catalogs/{id}
POST   /api/v1/catalogs/{id}/items
PUT    /api/v1/catalogs/{id}/items/{item_id}
DELETE /api/v1/catalogs/{id}/items/{item_id}
```

Mantener la reconstrucción del contrato Ule desde el modelo interno.

### 10.6 Anuncios

```text
POST   /api/v1/ads
PUT    /api/v1/ads/{id}
DELETE /api/v1/ads/{id}
```

Validar:

- vigencia;
- activo;
- peso;
- tipo;
- páginas;
- prioridad de slot.

### 10.7 Seguridad editorial

No exponer escritura mediante rutas públicas sin middleware.

Todas las operaciones de escritura deben:

- autenticar sesión;
- comprobar rol/permisos;
- validar payload;
- utilizar consultas parametrizadas;
- devolver errores JSON consistentes;
- tener tests HTTP de autorización;
- registrar quién creó/modificó contenido cuando el modelo lo permita.

### 10.8 Criterio de terminado

- registro controlado funcionando;
- login/logout reutilizando sesiones existentes;
- autorización por rol;
- CRUD de artículos;
- CRUD de bibliografía y relaciones;
- CRUD de catálogos y elementos;
- CRUD de anuncios;
- tests de autenticación/autorización;
- tests de contrato;
- regresión completa de Tonalmaster y Ule público;
- documentación de los endpoints editoriales.

---

## Fuera del alcance inmediato

No introducir todavía:

- microservicios;
- Redis;
- Kubernetes;
- CMS genérico;
- marketplace;
- red social completa;
- sistema complejo de permisos;
- almacenamiento avanzado de imágenes.

La administración de PostgreSQL mediante pgAdmin es una **herramienta de desarrollo/operación**, no un panel editorial para usuarios finales.

## Definición de terminado

Una fase se considera terminada cuando:

- compila;
- puede ejecutarse localmente;
- tiene tests correspondientes;
- no rompe fases anteriores;
- mantiene separadas las responsabilidades;
- evita abstracciones innecesarias;
- queda documentada la decisión importante.
