# Cambios en `feat/bootstrap-backend` — corrección de tests y alineación

Fecha: 2026-09-21

## Errores de compilación / tests corregidos

### 1. Tipos en `TonalpohualliCASO.Convert`

`mod` recibe y opera sobre `int64` y devuelve `int`. Tras la primera llamada:

```go
offset := mod(jdn-s.BaseJDN, 260) // tipo int
```

las llamadas posteriores `mod(offset+4, 20)` y `mod(offset, 13)` no compilaban (`int` no es `int64`).

**Fix:** promover el offset a `int64` y castear solo donde se necesita aritmética de `int` (trecena):

```go
offset := int64(mod(jdn-s.BaseJDN, 260))
Trecena: int(offset)/13 + 1
```

Los tests de dominio (`TestTonalpohualliCASO20260518`, ciclo 260, `TestMod`, JDN) pasan.

### 2. `writeJSON` duplicado en el paquete `handlers`

Había dos definiciones de `writeJSON` (`health.go` y `calendars.go`), lo que impediría compilar el paquete una vez resuelto el dominio.

**Fix:** se eliminó la copia en `health.go`; la implementación canónica queda en `calendars.go` junto con `writeError`.

Los tests HTTP (`List`, `Convert`, fecha inválida, sistema desconocido) pasan.

## Ajustes estructurales (no solo tests)

### Seed de `calendar_systems`

El seed de la migración `000001` tenía `jdn_base = NULL` y un texto que decía que la correlación se definiría después, pero el dominio ya fija el ancla CASO (`CASOAnchorJDN = 2276828`) y la documentación en `arquitectura.md` ya registra el fixture verificado.

**Cambio:**

| Campo | Antes | Después |
|-------|--------|---------|
| `autor_correlacion` | `CASO` | `Alfonso Caso` |
| `jdn_base` | `NULL` | `2276828` |
| `descripcion` | texto provisional | ancla + fixture 2026-05-18 |

Esto alinea migración, dominio y `docs/arquitectura.md` (nota de verificación del caso).

### `go.mod` / `go.sum`

- Se generó `go.sum` (faltaba en el árbol clonado).
- Directiva `go` fijada en **1.23** como mínimo (compatible con `pgx/v5`).
- CI y Dockerfile siguen apuntando a **1.25** (toolchain de build); no es un bloqueo: Go permite toolchain más nueva que el `go` directive.

## Verificación local

```bash
go test ./...
go build ./...
```

Resultado esperado: todos los paquetes de dominio y handlers en verde; `cmd/server` compila.

## Sin cambios de contrato HTTP

Los endpoints y el shape JSON documentados en README / `arquitectura.md` no cambian:

- `GET /api/v1/calendars`
- `GET /api/v1/calendars/{id}`
- `GET /api/v1/calendars/convert?date=YYYY-MM-DD&system=tonalpohualli_caso`
- Errores: `{"error":{"code":"...","message":"..."}}`

## Notas para fases siguientes

- `Registry` aún no expone `List()`; el handler de listado hardcodea el DTO del sistema CASO. Cuando haya más sistemas, conviene iterar el registry.
- `NightLord` / señor de la noche sigue vacío hasta disponer de tabla verificable (ya anotado en arquitectura).
- Ejecución de migraciones en arranque aún no está cableada (planeación Fase 2 / 6).
