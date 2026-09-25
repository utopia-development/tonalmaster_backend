# Arquitectura del Backend — ULE & Tonalmaster

## 1. Propósito

Este repositorio contiene una API modular en Go que integra los dominios de **ULE** y **Tonalmaster**, con PostgreSQL como persistencia.

La arquitectura busca mantener desacoplados:

- la interfaz HTTP y el contrato externo;
- la lógica de negocio;
- la persistencia;
- los cálculos calendáricos.

El backend expone una sola API y no depende de que el frontend conozca la estructura interna de la aplicación.

---

## 2. Principios arquitectónicos

No se adopta una implementación dogmática de SOLID. Se aplican, de forma práctica, estos principios:

### Responsabilidad única

Cada capa tiene una responsabilidad clara:

- **handlers**: HTTP, decodificación de entrada, validación propia del transporte y adaptación entre DTOs externos y modelos internos.
- **services**: reglas de negocio, autorización y coordinación de operaciones.
- **repository**: acceso a PostgreSQL.
- **config**: configuración de ejecución.
- **calendars**: lógica de cómputo calendárico y funciones relacionadas con el dominio.

### Dependencias hacia abstracciones

Cuando una dependencia necesita sustituirse o aislarse para pruebas, se utilizan interfaces pequeñas en lugar de acoplar la lógica de negocio a una implementación concreta.

### Separación de contrato externo e implementación interna

El contrato HTTP de ULE utiliza las claves JSON definidas por el frontend, aunque los modelos y nombres internos de Go siguen las convenciones idiomáticas del lenguaje.

La conversión se realiza explícitamente en el borde HTTP:

```text
Frontend ULE
    |
    | HTTP
    v
Handler
    |
    | DTO del contrato ULE
    | (JSON en español)
    v
Adaptador
    |
    | modelo interno Go
    v
Service
    |
    v
Repository
    |
    v
PostgreSQL
```

Esto permite evolucionar la implementación interna sin romper el contrato público.

---

## 3. Estructura actual del proyecto

La organización relevante del backend es:

```text
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── repository/
│   ├── services/
│   └── ...
├── migrations/
├── docs/
│   ├── arquitectura.md
│   └── contrato_datos_ule.md
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

El detalle exacto de archivos puede crecer con el proyecto. La separación anterior representa responsabilidades, no una obligación de mantener un número fijo de paquetes.

---

## 4. Flujo de una petición

### Lectura pública

```text
Frontend ULE
    -> HTTP API
    -> Handler
    -> Service
    -> Repository
    -> PostgreSQL
```

Los handlers no contienen consultas SQL ni reglas de negocio complejas.

### Escritura editorial

Las operaciones editoriales pasan por autenticación y autorización antes de llegar a la lógica de negocio:

```text
Frontend / cliente editorial
    -> autenticación
    -> autorización por rol
    -> Handler
    -> DTO -> modelo interno
    -> Service
    -> Repository
    -> PostgreSQL
```

Los roles editoriales actualmente utilizados son `contributor` y `admin`.

---

## 5. Contrato HTTP de ULE

El contrato externo y vinculante vive en:

```text
docs/contrato_datos_ule.md
```

Ese documento es la fuente de verdad para la integración con el frontend ULE.

Reglas importantes:

- las **rutas HTTP son en inglés**;
- las **claves JSON son en español**, exactamente como aparecen en el contrato;
- los nombres con acentos forman parte del contrato cuando así están definidos;
- el frontend no debe depender de nombres internos de Go ni de nombres de columnas PostgreSQL;
- los DTOs HTTP deben adaptarse explícitamente a los modelos internos;
- los cambios de contrato deben actualizarse de forma coordinada con el frontend y sus pruebas.

La API pública de ULE es de lectura y expone únicamente contenido que corresponde a registros visibles.

### Endpoints públicos actuales

```text
GET /api/v1/articles
GET /api/v1/articles/{id}

GET /api/v1/bibliography

GET /api/v1/catalogs
GET /api/v1/catalogs/{id}

GET /api/v1/ads
```

No se debe introducir paginación, filtros u otras extensiones en el contrato únicamente por conveniencia interna: deben responder a una necesidad real y quedar documentadas en el contrato.

---

## 6. API editorial actual

La escritura de contenido ULE ya forma parte de la implementación actual.

Endpoints protegidos:

```text
POST   /api/v1/articles
PUT    /api/v1/articles/{id}
DELETE /api/v1/articles/{id}

POST   /api/v1/bibliography
PUT    /api/v1/bibliography/{id}
DELETE /api/v1/bibliography/{id}

POST   /api/v1/catalogs
PUT    /api/v1/catalogs/{id}
DELETE /api/v1/catalogs/{id}

POST   /api/v1/catalogs/{id}/items
PUT    /api/v1/catalogs/{id}/items/{item_id}
DELETE /api/v1/catalogs/{id}/items/{item_id}

POST   /api/v1/ads
PUT    /api/v1/ads/{id}
DELETE /api/v1/ads/{id}
```

Estas operaciones requieren una sesión válida y autorización editorial.

La existencia de estos endpoints no cambia el contrato de lectura pública: el frontend público continúa consumiendo los endpoints GET.

---

## 7. Persistencia

PostgreSQL es la base de datos principal del backend.

Las **migraciones de `migrations/` son la fuente de verdad del esquema**. Este documento no duplica las sentencias SQL de creación de tablas para evitar que la documentación y la implementación diverjan.

Los repositorios son responsables de:

- ejecutar consultas;
- mapear resultados de PostgreSQL a modelos internos;
- manejar errores de persistencia;
- mantener fuera de los handlers los detalles de SQL.

Los servicios no deben depender de nombres de columnas ni construir consultas SQL directamente.

---

## 8. Autenticación, sesiones y autorización

La autenticación utiliza sesiones opacas almacenadas en PostgreSQL.

Las contraseñas se almacenan mediante hash seguro y no se conservan en texto plano.

El flujo general es:

```text
registro/login
    -> autenticación
    -> sesión
    -> cookie/sesión HTTP
    -> middleware
    -> usuario autenticado
    -> autorización por rol
```

El registro utiliza el código configurado mediante `REGISTRATION_CODE`. Este código controla la creación de cuentas y no forma parte del proceso normal de login.

Las operaciones editoriales requieren rol `contributor` o `admin`.

Los secretos, contraseñas y tokens no deben escribirse en logs.

---

## 9. Despliegue actual

El despliegue local utiliza Docker Compose.

Los servicios actuales son:

```text
db       PostgreSQL
migrate  aplicación de migraciones
api      servidor Go
```

El servicio `migrate` aplica las migraciones antes de que la API quede disponible.

El puerto de PostgreSQL se publica en el host mediante `POSTGRES_PORT`. El nombre DNS `db` solamente es válido dentro de la red de Docker Compose.

### Administración de PostgreSQL

**pgAdmin 4 no forma parte del proyecto ni de Docker Compose.**

Puede utilizarse como aplicación externa para conectarse al puerto publicado de PostgreSQL.

La arquitectura no incluye un panel web de administración de base de datos.

---

## 10. Integración con el frontend

ULE se mantiene desacoplado del backend.

El frontend conoce únicamente:

- la URL base de la API;
- el contrato HTTP;
- los datos definidos por dicho contrato;
- el mecanismo de autenticación necesario para las operaciones protegidas.

No debe conocer:

- consultas SQL;
- tablas de PostgreSQL;
- nombres internos de modelos Go;
- detalles de repositorios;
- estructura interna de servicios.

El adaptador HTTP del backend es el punto responsable de traducir entre el contrato externo y la implementación interna.

---

## 11. Guía para continuar el desarrollo

Las decisiones de implementación deben seguir estas reglas:

1. **Primero revisar el contrato** si el cambio afecta datos o endpoints de ULE.
2. **Mantener separados DTOs y modelos internos** cuando sus nombres o formas tengan objetivos diferentes.
3. **Mantener la lógica de negocio en services**, no en handlers.
4. **Mantener SQL y detalles de PostgreSQL en repositories**.
5. **Actualizar las migraciones** cuando cambie persistentemente el esquema.
6. **Agregar pruebas de contrato** cuando una modificación pueda afectar la integración con ULE.
7. **Evitar duplicar información estructural** entre documentación y migraciones.
8. **Preferir cambios pequeños y verificables** sobre refactorizaciones amplias.
9. **No introducir infraestructura adicional** (microservicios, Redis, Kubernetes, CMS u otras capas) sin una necesidad técnica concreta.
10. **Documentar las decisiones permanentes**, no el plan temporal utilizado para llegar a ellas.

Las planeaciones de implementación, tareas temporales y pasos de desarrollo pertenecen a documentos de trabajo separados y no forman parte de este documento arquitectónico.

---

## 12. Dominio Tonalmaster

El dominio calendárico mantiene su lógica separada de HTTP y persistencia.

Los cálculos deben permanecer en componentes de dominio independientes de handlers y consultas SQL. Las correlaciones, constantes y resultados verificables deben estar respaldados por código y pruebas.

El sistema CASO actualmente documenta y prueba su ancla de conversión y sus resultados de referencia. Los datos calendáricos que todavía requieran una fuente verificable no deben representarse como valores definitivos.

---

## 13. Principio de evolución

La arquitectura actual debe entenderse como una API modular, no como una arquitectura distribuida.

Mientras la separación entre HTTP, negocio, persistencia y dominio siga resolviendo las necesidades reales del proyecto, no se debe añadir complejidad estructural sin una razón concreta.

Cuando una nueva necesidad requiera modificar esta arquitectura, el cambio debe evaluarse contra:

- el contrato público existente;
- la separación de responsabilidades;
- la facilidad de prueba;
- la simplicidad operativa;
- la compatibilidad con el frontend;
- la necesidad real que motiva la nueva abstracción.
