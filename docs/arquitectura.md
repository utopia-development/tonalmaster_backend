ARQUITECTURA.MD — Blueprint del Backend Unificado (Ule & Tonalmaster)
Principio Rector: El contenido y los servicios cambian; la interfaz no debe enterarse de dónde provienen. Una sola API ligera en Go con PostgreSQL para alimentar el ecosistema educativo y la red social calendárica.

1. Principios SOLID Aplicados al Backend en Go
Para garantizar un código limpio, mantenible y escalable, la arquitectura del backend aplica estrictamente los principios SOLID:

S (Single Responsibility Principle - Principio de Responsabilidad Única): Cada paquete tiene un propósito exclusivo. Los controladores (handlers) solo manejan solicitudes HTTP y validaciones de entrada; los repositorios (repository) solo interactúan con la base de datos; los servicios de conversión (calendars) contienen la lógica matemática pura de los cómputos mesoamericanos.

O (Open/Closed Principle - Principio de Abierto/Cerrado): El sistema está abierto a la extensión (por ejemplo, añadir nuevos sistemas calendáricos como una nueva variante del Tonalpohualli o nuevos tipos de contenido en Ule) sin necesidad de modificar el código existente de los controladores principales.

L (Liskov Substitution Principle - Principio de Sustitución de Liskov): Los accesos a datos se definen mediante interfaces pequeñas y específicas por dominio (`CalendarRepository`, `UserRepository`, etc.). Cualquier implementación (PostgreSQL, mock tests, etc.) puede sustituir a otra sin romper la capa de negocio.

I (Interface Segregation Principle - Principio de Segregación de Interfaces): Se evitan interfaces monolíticas. Se prefieren contratos pequeños y específicos por dominio (ArticleService, CalendarService, AuthService).

D (Dependency Inversion Principle - Principio de Inversión de Dependencias): Los controladores no instancian directamente las conexiones a la base de datos ni los repositorios; estos se inyectan a través de sus constructores utilizando interfaces.

2. Experiencia "Like Vikunja": Despliegue con un Solo Comando
El proyecto se despliega de forma autónoma mediante Docker Compose. No requiere configuraciones complejas en el sistema operativo anfitrión.

El archivo `docker-compose.yml` en la raíz del repositorio define los servicios `db` (PostgreSQL) y `api`. Las variables de entorno reales que consume (`POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_PORT`, `APP_ENV`, `APP_HOST`, `APP_PORT`, `DATABASE_URL`, `CORS_ALLOWED_ORIGINS`) están documentadas en `.env.example`; consulta ese archivo y `docker-compose.yml` como fuente de verdad en lugar de nombres de variables antiguos (`DB_USER`, `DB_HOST`, `PORT`, etc.) que pudieran aparecer en versiones previas de este documento.

Para ponerlo en marcha:
1. Copiar el archivo de entorno `.env.example` a `.env`.
2. Ejecutar en la terminal:

```bash
docker compose up --build -d
```

3. Aplicar las migraciones (aún no se ejecutan automáticamente al arrancar, ver `docs/planeacion.md` § Fase 6):

```bash
psql -h localhost -U ule_user -d ule_tonalmaster -f migrations/000001_init_schema.up.sql
psql -h localhost -U ule_user -d ule_tonalmaster -f migrations/000002_sessions.up.sql
```

La base de datos y la API estarán listas y comunicadas de manera interna y segura; los endpoints de auth/events/interpretations solo responderán correctamente después del paso 3.

3. Esquema de Base de Datos (PostgreSQL)

**Alcance del MVP:** el backend inicial es exclusivamente Tonalmaster. Ule (articles, bibliography, catalogs y ads) queda fuera del MVP y se incorporará en una fase posterior, sin obligar a implementar sus tablas o endpoints durante las Fases 1–5. La arquitectura conserva esos dominios como evolución prevista del mismo backend.
Estructura relacional normalizada que unifica los dominios de educación, catálogos y sistemas calendáricos sociales.

SQL
-- Usuarios unificados (para Tonalmaster y panel administrativo de Ule)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    avatar_url VARCHAR(255),
    role VARCHAR(50) NOT NULL DEFAULT 'reader' CHECK (role IN ('reader', 'contributor', 'admin')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ULE: Artículos educativos
CREATE TABLE articles (
    id VARCHAR(100) PRIMARY KEY,
    titulo VARCHAR(255) NOT NULL,
    autor VARCHAR(150) NOT NULL,
    fecha DATE NOT NULL,
    resumen TEXT NOT NULL,
    contenido_html TEXT NOT NULL,
    imagen_destacada VARCHAR(255),
    categoria VARCHAR(100) NOT NULL,
    etiquetas TEXT[],
    visible BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ULE: Bibliografía
CREATE TABLE bibliography (
    id VARCHAR(100) PRIMARY KEY,
    titulo VARCHAR(255) NOT NULL,
    autor VARCHAR(150) NOT NULL,
    anio INTEGER,
    tipo VARCHAR(50),
    referencia TEXT,
    enlace VARCHAR(255)
);

-- ULE: Relación muchos a muchos (Artículo ↔ Bibliografía)
CREATE TABLE article_bibliography (
    article_id VARCHAR(100) REFERENCES articles(id) ON DELETE CASCADE,
    bibliography_id VARCHAR(100) REFERENCES bibliography(id) ON DELETE CASCADE,
    PRIMARY KEY (article_id, bibliography_id)
);

-- ULE: Catálogos y piezas
CREATE TABLE catalogs (
    id VARCHAR(100) PRIMARY KEY,
    titulo VARCHAR(255) NOT NULL,
    descripcion TEXT
);

CREATE TABLE catalog_items (
    id VARCHAR(100) PRIMARY KEY,
    catalog_id VARCHAR(100) REFERENCES catalogs(id) ON DELETE CASCADE,
    titulo VARCHAR(255) NOT NULL,
    categoria VARCHAR(100),
    imagen VARCHAR(255),
    detalles JSONB
);

-- ULE: Anuncios contextuales
CREATE TABLE ads (
    id VARCHAR(100) PRIMARY KEY,
    titulo VARCHAR(255),
    imagen VARCHAR(255),
    enlace VARCHAR(255),
    paginas TEXT[],
    fecha_inicio DATE,
    fecha_fin DATE,
    activo BOOLEAN DEFAULT TRUE
);

-- TONALMASTER: Sistemas de Cómputo Calendárico
CREATE TABLE calendar_systems (
    id VARCHAR(50) PRIMARY KEY, -- ej: 'tonalpohualli_caso', 'tonalpohualli_mesa', 'tzolkin_gmt'
    nombre VARCHAR(100) NOT NULL,
    autor_correlacion VARCHAR(150),
    jdn_base INTEGER,
    descripcion TEXT
);

-- TONALMASTER: Interpretaciones sociales en fechas específicas
CREATE TABLE interpretations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    calendar_system_id VARCHAR(50) NOT NULL REFERENCES calendar_systems(id) ON DELETE CASCADE,
    target_date DATE NOT NULL,
    contenido TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- TONALMASTER: Eventos de usuarios en los calendarios
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    calendar_system_id VARCHAR(50) NOT NULL REFERENCES calendar_systems(id) ON DELETE CASCADE,
    target_date DATE NOT NULL,
    titulo VARCHAR(255) NOT NULL,
    descripcion TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
4. Contrato de Datos (API Contracts)
Para que el frontend estático de Ule (y la futura interfaz de Tonalmaster) consuman los datos de forma idéntica sin importar si la fuente es un archivo JSON local o esta API en Go, los endpoints devolverán estrictamente los siguientes contratos JSON.

Contrato: Artículo (GET /api/v1/articles/:id)
JSON
{
  "id": "articulo-001",
  "titulo": "Arquitectura de los templos del sol",
  "autor": "Felipe González",
  "fecha": "2026-03-15",
  "resumen": "Análisis geométrico y astronómico...",
  "contenido_html": "<p>Contenido detallado en HTML...</p>",
  "imagen_destacada": "assets/images/articulos/templo.jpg",
  "categoria": "Arquitectura",
  "etiquetas": ["maya", "clásico", "astronomía"],
  "bibliografia_relacionada": ["biblio-001"],
  "visible": true
}
Contrato: Conversión Calendárica (GET /api/v1/calendars/convert?date=2026-05-18&system=tonalpohualli_caso)
JSON
{
  "fecha_gregoriana": "2026-05-18",
  "sistema": "tonalpohualli_caso",
  "jdn": 2461179,
  "resultado": {
    "trecena": 1,
    "signo": "Cozcacuauhtli",
    "numero_dia": 12,
    "senor_de_la_noche": null
  }
}
5. Estructura de Directorios del Repositorio (`tonalmaster_backend`)
Plaintext
tonalmaster_backend/
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
├── cmd/
│   └── server/
│       └── main.go
|── docs/
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── db.go
│   ├── domain/\n│   │   └── calendars/\n│   ├── services/\n│   ├── handlers/
│   │   ├── auth.go
│   │   ├── articles.go
│   │   ├── bibliography.go
│   │   ├── catalogs.go
│   │   ├── ads.go
│   │   ├── calendars.go
│   │   └── interpretations.go
│   ├── models/
│   │   └── models.go
│   └── repository/
│       ├── postgres.go
│       └── repository.go
└── migrations/
    ├── 000001_init_schema.up.sql
    └── 000001_init_schema.down.sql


6. Seguridad y sesiones (posterior al MVP inicial)

La autenticación no bloquea las Fases 1–5 de calendarios. Antes de habilitar escritura sobre `events` e `interpretations`, se implementará un bloque de autenticación con sesiones opacas almacenadas en PostgreSQL: token persistido como hash, expiración y revocación. La sesión se transportará preferentemente mediante cookie `HttpOnly; Secure`; `SameSite=Lax` para frontends bajo el mismo sitio y `SameSite=None` + protección CSRF cuando sean cross-site. CORS usará orígenes explícitos y credenciales, nunca `*`. Las contraseñas se almacenarán con Argon2id o bcrypt. Los roles iniciales serán `reader`, `contributor` y `admin`. Redis y JWT quedan fuera del MVP salvo necesidad futura documentada.


**Nota de verificación del caso:** el cálculo del sistema CASO usa como ancla 13 de agosto de 1521 (calendario juliano), JDN 2276828, identificado como 1-Cóatl y comienzo de la trecena 1 en la metodología consultada. La conversión de 2026-05-18 produce JDN 2461179; a partir de esa ancla, el dominio prueba 12-Cozcacuauhtli y trecena 1. El señor de la noche permanece sin valor hasta incorporar una tabla/fuente específica y verificable.
