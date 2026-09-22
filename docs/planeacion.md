# Planeación — Tonalmaster Backend

## Objetivo

Construir el backend inicial de Tonalmaster con Go + PostgreSQL, siguiendo `arquitectura.md`, principios SOLID y una arquitectura mínima que permita crecer sin complejidad innecesaria.

Regla: **primero hacer funcionar el camino completo más pequeño; después extenderlo.**

---

## Fase 1 — Base ejecutable

**Objetivo:** tener el proyecto arrancando localmente.

- Crear `go.mod`.
- Crear `cmd/server/main.go`.
- Crear configuración desde variables de entorno.
- Crear conexión a PostgreSQL.
- Crear `Dockerfile`.
- Crear `docker-compose.yml`.
- Crear `.env.example`.
- Crear `Makefile`.
- Implementar `GET /health`.
- Implementar `GET /ready`.
- Agregar logs estructurados.
- Configurar CORS explícitamente.

**Resultado:** API + PostgreSQL levantan con Docker y los endpoints de salud responden.

---

## Fase 2 — Modelo mínimo de Tonalmaster

**Objetivo:** persistir los datos esenciales.

Crear migraciones para:

- `users`
- `calendar_systems`
- `interpretations`
- `events`

Agregar:

- claves primarias;
- relaciones;
- índices necesarios;
- timestamps;
- restricciones básicas;
- seed mínimo para el sistema calendárico inicial.

**Resultado:** base de datos reproducible desde cero.

---

## Fase 3 — Dominio calendárico

**Objetivo:** separar las reglas calendáricas de HTTP y PostgreSQL.

- Crear una interfaz pequeña para sistemas calendáricos.
- Implementar el primer sistema: `tonalpohualli_caso`.
- Mantener los cálculos como funciones puras.
- Cubrir los cálculos principales con tests.
- No colocar SQL ni lógica HTTP dentro del dominio.

**Resultado:** conversión de fechas verificable y reutilizable.

---

## Fase 4 — API v1

**Objetivo:** exponer el dominio mediante una API estable.

Implementar:

- `GET /api/v1/calendars`
- `GET /api/v1/calendars/{id}`
- `GET /api/v1/calendars/convert?date=YYYY-MM-DD&system=...`

Después, cuando el núcleo esté estable:

- endpoints de `events`;
- endpoints de `interpretations`.

Aplicar:

- DTOs;
- validación;
- errores JSON consistentes;
- fechas ISO 8601;
- IDs públicos estables;
- paginación donde corresponda;
- CORS;
- request ID.

**Resultado:** contrato HTTP estable sin exponer el esquema interno de PostgreSQL.

---

## Fase 5 — Tests y contrato

**Objetivo:** evitar regresiones.

- Tests unitarios del dominio calendárico.
- Tests de repositorios.
- Tests HTTP de endpoints principales.
- Validar errores y entradas inválidas.
- Verificar que el proyecto arranca desde un entorno limpio.
- Documentar ejemplos mínimos de API.

**Resultado:** una base comprobable antes de agregar funcionalidades.

---

**Estado:** Fases 1–5 cerradas en código (dominio + API calendarios + tests + CI de migraciones). Auth por sesión iniciada (migración `000002`, handlers `/api/v1/auth/*`); documentada en `docs/auth-sesiones.md`. Siguiente foco Fase 6: migraciones al arranque, endpoints autenticados de events/interpretations, endurecer tests de auth HTTP.

## Alcance transversal

El MVP inicial es exclusivamente Tonalmaster. Ule (artículos, bibliografía, catálogos y anuncios) queda fuera de las Fases 1–5 y se planifica como evolución posterior del mismo backend.

La autenticación se implementará después de las Fases 1–5 y antes de habilitar escritura autenticada sobre `events` e `interpretations`.

## Fase 6 — Integración y evolución

**Objetivo:** preparar Tonalmaster para crecer sin rehacer el núcleo.

- Integrar progresivamente el frontend.
- Mantener contratos de API estables.
- Agregar autenticación cuando sea necesaria.
- Incorporar publicaciones/interpretaciones y eventos según prioridad real.
- Revisar rendimiento y consultas N+1.
- Agregar caché HTTP donde aporte valor.
- Mantener migraciones versionadas.

### Fase 6 — Estado de implementación

Implementado en esta fase:
- sesiones opacas en PostgreSQL y autenticación HTTP;
- middleware reutilizable de autenticación;
- repositorio PostgreSQL para `events` e `interpretations`;
- creación/listado/eliminación autenticada de eventos;
- creación autenticada y consulta pública de interpretaciones;
- migraciones versionadas `000001` y `000002` verificadas por CI.

Pendiente antes de cerrar Fase 6: pruebas de integración HTTP + PostgreSQL para auth/content, ejecución automática de migraciones al arranque o mediante un comando operativo reproducible, rate limiting y revisión final de CORS/CSRF para despliegues cross-site.

# Plan de desarrollo — Fase 7: endpoints de contenido para `ule_educativo`

> Objetivo único: exponer `articles`, `bibliography`, `catalogs` y `ads` como rutas **públicas**
> de `tonalmaster_backend`, con el DTO exacto que define `docs/contrato_datos.md` (contrato
> vinculante con `ule_educativo`). Fases 1–5 (Tonalmaster) y Fase 6 (auth/events/interpretations)
> no se tocan. Alcance deliberadamente acotado: solo lectura pública, sin admin/CMS, sin
> paginación, sin comentarios.

## Estado de partida

`docs/arquitectura.md` ya boceta el esquema SQL y un ejemplo de DTO para `articles`; sirve como
punto de partida pero **no está implementado** (no hay migración, ni `internal/*`, ni rutas). El
boceto tiene brechas frente al contrato del frontend — ver las notas "Nota para el backend" en
`docs/contrato_datos.md` §1–§4. Este plan las resuelve.

## Pasos

1. **Migración `000003_content.up/down.sql`** con las 5 tablas (`articles`, `bibliography`,
   `article_bibliography`, `catalogs`, `catalog_items`, `ads`), incorporando ya los ajustes
   detectados en el contrato:
   - `articles.autor` → nullable (no `NOT NULL`).
   - `bibliography`: agregar `visible BOOLEAN DEFAULT TRUE`; usar `autores TEXT[]` en vez de
     `autor` singular; separar `editorial TEXT` y `resumen TEXT` en vez de `referencia` genérico;
     `enlace` → `url`.
   - `catalogs`: agregar `visible BOOLEAN DEFAULT TRUE`.
   - `catalog_items`: mantener `detalles JSONB` (ahí viven `categorias`, `imagen_alt`,
     `descripcion`, `año_descubrimiento`, `ubicacion`); quitar la columna `categoria` singular,
     sobra frente al JSONB.
   - `ads`: ampliar con `imagen_alt`, `contacto`, `slogan`, `descripcion`, `peso INTEGER`,
     `tipo VARCHAR(50)`, `prioridad_slot VARCHAR(50)`; renombrar `fecha_inicio`/`fecha_fin` a
     `vigencia_inicio`/`vigencia_fin` (o mantener el nombre de columna y mapear solo en el DTO,
     lo que sea más barato dado el estado actual — no hay datos en producción que migrar).

2. **Capas Go**, replicando el patrón ya usado por `calendars`/`content` (auth):
   - `internal/repository`: `ArticleRepository`, `BibliographyRepository`, `CatalogRepository`,
     `AdRepository` (interfaz + implementación Postgres), siguiendo el estilo de
     `internal/repository/postgres_auth.go`.
   - `internal/services`: una capa fina que arma los DTOs desde las filas — aquí vive la
     reconstrucción de `categorias_disponibles`/`categorias` desde `detalles JSONB` y el mapeo de
     columnas de `ads` y `bibliography` descrito arriba.
   - `internal/handlers/content_ule.go` (o dividir en `articles.go`, `bibliography.go`,
     `catalogs.go`, `ads.go` como ya sugiere la estructura de `docs/arquitectura.md`): handlers
     HTTP que devuelven `[]DTO` o `DTO` + 404, mismo estilo que `internal/handlers/content.go`
     existente para `interpretations`.

3. **Rutas en `cmd/server/main.go`**, montadas **sin** `RequireAuth`:
   ```
   GET /api/v1/articles
   GET /api/v1/articles/{id}
   GET /api/v1/bibliography
   GET /api/v1/bibliography/{id}
   GET /api/v1/catalogs
   GET /api/v1/catalogs/{id}
   GET /api/v1/ads
   ```
   Confirmar explícitamente en el router (y con un test) que ninguna de estas rutas exige cookie
   ni Bearer token.

4. **Filtro de visibilidad en la query**, no solo en el DTO: `WHERE visible = TRUE` para
   articles/bibliography/catalogs; para `ads`, `WHERE activo = TRUE AND (vigencia_fin IS NULL OR
   vigencia_fin >= CURRENT_DATE)`.

5. **CORS**: confirmar que `CORS_ALLOWED_ORIGINS` en el entorno de despliegue incluye el dominio
   real de `ule_educativo` (GitHub Pages u otro). Sin esto, el navegador bloquea las llamadas de
   `loader.js` aunque el backend responda bien.

6. **Tests**: unitarios de servicio (reconstrucción de DTO desde JSONB, mapeo de `ads`) +
   handler tests por recurso (200 con datos, 200 con lista vacía, 404 en detalle, visible:false
   excluido), siguiendo el estilo de `internal/handlers/calendars_test.go`.

7. **Congelar `docs/contrato_datos.md`** (ya actualizado) como archivo idéntico en ambos repos;
   es el criterio de aceptación de esta fase, no el código de `ule_educativo`.

## Fuera de este plan (a cargo del equipo de frontend, después de que el paso 3 esté desplegado)

- Cambiar en `js/loader.js` las 4 llamadas `apiJSON('/articulos'|...)` a
  `apiJSON('/articles'|...)` y fijar `ULE.config.dataSource='api'` +
  `ULE.config.apiBaseUrl`.
- Correr `scripts/validar_datos.py` y `scripts/pruebas_navegador.py` con `dataSource='api'`
  contra el backend ya desplegado.
- Migrar el contenido editorial real (`data/articulos/*.json`, `data/bibliografia/*.json`,
  `data/catalogos/*.json`, `data/anuncios.json`) a filas en Postgres — no es parte de este plan
  de backend; requiere un script de carga puntual (`INSERT`) o un endpoint interno de import, a
  decidir con el equipo editorial cuando este plan esté cerrado.

## No perder de vista

Esta fase agrega superficie pública nueva sin tocar los pendientes de Fase 6 ya documentados en
`docs/planeacion.md` (migraciones automáticas al arranque, rate limiting, revisión CORS/CSRF,
tests de integración HTTP+Postgres). Conviene resolverlos en paralelo o inmediatamente después,
porque el sitio público (`ule_educativo`) va a exponer estas rutas a tráfico real antes que
cualquier cliente de Tonalmaster.

### Fuera del MVP inicial

No implementar todavía:

- panel administrativo;
- red social completa;
- marketplace;
- CMS genérico;
- sistema complejo de permisos;
- microservicios;
- colas;
- Redis;
- Kubernetes;
- almacenamiento avanzado de imágenes.

---

## Estructura objetivo

```
.
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
├── README.md
├── cmd/
│   └── server/
│       └── main.go
|── docs/
├── internal/
│   ├── config/
│   ├── database/
│   ├── domain/
│   │   └── calendars/
│   ├── services/
│   ├── repository/
│   ├── handlers/
│   └── models/
└── migrations/
```

---

## Definición de terminado

Una fase se considera terminada cuando:

- compila;
- puede ejecutarse localmente;
- tiene los tests correspondientes;
- no rompe las fases anteriores;
- mantiene separadas las responsabilidades;
- no agrega abstracciones sin necesidad;
- queda documentada la decisión importante.

---

## Principios de implementación

1. **SOLID**, sin sobreingeniería.
2. **Una responsabilidad por componente.**
3. **Dominio independiente de HTTP y SQL.**
4. **Inyección de dependencias mediante constructores.**
5. **Interfaces pequeñas y sólo cuando aporten desacoplamiento real.**
6. **API estable; la base de datos no se expone al frontend.**
7. **Migraciones reproducibles.**
8. **Código simple antes que arquitectura compleja.**
