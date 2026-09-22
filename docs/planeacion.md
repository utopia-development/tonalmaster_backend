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

**Estado: CERRADA.**

Objetivo cumplido: dejar el núcleo Tonalmaster operativo y preparado para evolucionar.

Implementado y documentado:
- sesiones opacas en PostgreSQL;
- autenticación HTTP;
- middleware reutilizable;
- events e interpretations autenticados;
- migraciones versionadas;
- migraciones automáticas mediante Docker Compose;
- CORS explícito;
- estructura preparada para clientes desacoplados.

La validación final sobre Windows/WSL/red local queda como prueba de aceptación de integración en Fase 8, no como requisito bloqueante de esta fase.

## Fase 7 — Integración de contenido público con ule_educativo

**Estado: IMPLEMENTACIÓN COMPLETADA.**

Objetivo: exponer contenido público de Ule mediante el contrato vinculante docs/contrato_datos.md.

Implementado:
- migración 000003_ule_content;
- articles, bibliography, article_bibliography, catalogs, catalog_items, ads;
- repositorio PostgreSQL;
- servicio de transformación;
- handlers HTTP;
- rutas públicas GET /api/v1/articles, GET /api/v1/articles/{id}, GET /api/v1/bibliography, GET /api/v1/bibliography/{id}, GET /api/v1/catalogs, GET /api/v1/catalogs/{id}, GET /api/v1/ads;
- relaciones artículo↔bibliografía en ambos sentidos;
- filtrado de visibilidad en SQL;
- contrato sincronizado con ule_educativo;
- migración 000003 incorporada al pipeline CI;
- arranque Docker aplica automáticamente las migraciones.

La Fase 7 no incluye migración del contenido editorial real ni cambios en el frontend. Esos pasos pertenecen a la integración del producto.

## Fase 8 — Verificación, compatibilidad y aceptación de integración

**Objetivo:** comprobar que el backend cumple el contrato de Ule y que no rompe Tonalmaster.

### Tests automatizados
- go test ./...;
- go build ./...;
- tests HTTP de las rutas Ule;
- 200 con datos y listas vacías;
- 404 en recursos inexistentes;
- exclusión de visible=false;
- exclusión de anuncios inactivos/vencidos;
- validación de DTOs contra docs/contrato_datos.md;
- relaciones artículo↔bibliografía;
- catalogs/catalog_items y reconstrucción de categorías;
- autenticación y endpoints Tonalmaster existentes;
- migraciones 000001–000003 up/down en CI.

### Compatibilidad frontend
- comparar respuestas reales con ule_educativo/docs/contrato_datos.md;
- ejecutar validadores del frontend con dataSource=api;
- probar 404, 5xx y caída de red;
- comprobar que no se modifican componentes/páginas innecesariamente;
- verificar CORS con el origen real del frontend;
- verificar que las URLs de imágenes sean resolubles desde el navegador.

### Aceptación manual
1. levantar backend en WSL mediante Docker Compose;
2. servirlo hacia la red local;
3. levantar ule_educativo en Windows;
4. apuntar ULE.config.apiBaseUrl al backend;
5. navegar artículos, bibliografía, catálogos y anuncios;
6. comprobar auth/events/interpretations;
7. registrar resultados e incidencias.

**Criterio de cierre:** CI verde + contrato compatible + frontend real funcionando contra el backend desplegado en la red local.

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
