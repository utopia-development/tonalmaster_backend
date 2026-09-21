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
