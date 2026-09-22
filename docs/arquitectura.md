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

El archivo `docker-compose.yml` en la raíz define únicamente `db` (PostgreSQL), `migrate` (migraciones versionadas) y `api`. Las variables de entorno que consume (`POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_PORT`, `APP_ENV`, `APP_HOST`, `APP_PORT`, `DATABASE_URL`, `CORS_ALLOWED_ORIGINS`) están documentadas en `.env.example`; ese archivo y `docker-compose.yml` son la fuente de verdad.

Para ponerlo en marcha:
1. Copiar el archivo de entorno `.env.example` a `.env`.
2. Ejecutar en la terminal:

```bash
docker compose up --build -d
```

3. El servicio `migrate` aplica automáticamente todas las migraciones pendientes antes de iniciar la API.

La base de datos y la API estarán listas y comunicadas de manera interna mediante la red de Compose.

### Acceso externo a PostgreSQL

PostgreSQL se publica en el host mediante `POSTGRES_PORT` (por defecto `5432`) con el mapeo:

```text
0.0.0.0:${POSTGRES_PORT:-5432} -> PostgreSQL:5432
```

Esto permite que **pgAdmin 4 Desktop, instalado fuera del proyecto**, se conecte al PostgreSQL del host desde otra máquina usando la IP o nombre de red del servidor, el puerto publicado, la base, el usuario y la contraseña.

El hostname `db` es únicamente un nombre DNS interno de Docker Compose. No debe utilizarse desde una máquina externa.

La publicación del puerto debe complementarse con las reglas de firewall/red del host. La arquitectura del proyecto no incluye un servicio pgAdmin ni una interfaz web de administración de PostgreSQL.

3. Esquema de Base de Datos (PostgreSQL)

**Estado:** Tonalmaster y el contenido público Ule ya están implementados en el backend. La escritura editorial Ule queda para la Fase 10 y requerirá autenticación/autorización.
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
4. Contrato de datos Ule

El contrato externo y vinculante del frontend Ule vive exclusivamente en `docs/contrato_datos_ule.md`. Este repositorio no duplica allí sus DTOs. Si cambia el contrato, se actualiza ese archivo y se verifica el adaptador del servicio Ule.

La API pública Ule usa rutas en inglés y claves JSON en español según ese contrato. El contenido público es de solo lectura y no requiere autenticación.

5. Seguridad y sesiones

La autenticación existente usa sesiones opacas almacenadas en PostgreSQL. Las contraseñas se almacenan como hash seguro y los endpoints protegidos pasan por middleware de autenticación/autorización. CORS usa orígenes explícitos y nunca `*`.

El futuro registro editorial tendrá una particularidad temporal: `POST /api/v1/auth/register` requerirá un código de verificación estático. **Ese código solo controla la creación de cuentas. No participa en el login.** El login desde el inicio será el flujo normal de autenticación existente. El código podrá migrarse posteriormente a configuración/DB o retirarse cuando el registro público sea intencional.

6. Integración frontend

Ule se sirve desacoplado del backend. En desarrollo, XAMPP puede servir el frontend y Docker/WSL la API. `ULE.config.apiBaseUrl` apunta al backend y CORS permite explícitamente el origen del frontend.

7. Evolución editorial

La API pública Ule permanece de lectura. La futura API editorial autenticada añadirá escritura para artículos, bibliografía, catálogos, elementos de catálogo y anuncios. Las operaciones editoriales requieren sesión y autorización por rol; no se convierten en endpoints públicos.

8. Estructura conceptual

```text
handlers -> services -> repositories -> PostgreSQL
                         ^
                    interfaces

frontend Ule -> HTTP API -> services -> repositories

pgAdmin 4 Desktop (externo)
              |
              | TCP POSTGRES_PORT
              v
        PostgreSQL del host
```

Las migraciones son la fuente reproducible del esquema. Los datos editoriales normales entrarán por la API editorial; solo los datos fijos de infraestructura/demo justifican seeds mediante migración.

9. Decisiones fuera del alcance inmediato

No introducir todavía microservicios, Redis, Kubernetes, CMS genérico ni un sistema complejo de permisos. Mantener una API modular y pequeña hasta que una necesidad real justifique otra abstracción.

La administración mediante pgAdmin 4 es externa al proyecto: no es un servicio Docker, no forma parte de la API y no es un panel editorial para usuarios finales.

10. Verificación del dominio CASO

El sistema CASO usa como ancla 13 de agosto de 1521 (calendario juliano), JDN 2276828. La conversión de 2026-05-18 produce JDN 2461179 y el resultado documentado por los tests es 12-Cozcacuauhtli, trecena 1. El señor de la noche permanece sin valor hasta incorporar una tabla/fuente verificable.
