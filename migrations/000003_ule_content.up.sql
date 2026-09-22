CREATE TABLE articles (id VARCHAR(100) PRIMARY KEY, titulo VARCHAR(255) NOT NULL, autor VARCHAR(150), fecha DATE NOT NULL, resumen TEXT NOT NULL, contenido_html TEXT NOT NULL, imagen_destacada TEXT, imagen_alt TEXT, categoria VARCHAR(100), etiquetas TEXT[], visible BOOLEAN NOT NULL DEFAULT TRUE, created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP);

CREATE INDEX articles_visible_fecha_idx ON articles (visible, fecha DESC);

CREATE TABLE bibliography (id VARCHAR(100) PRIMARY KEY, titulo VARCHAR(255) NOT NULL, autores TEXT[], anio INTEGER, tipo VARCHAR(50) NOT NULL CHECK (tipo IN ('libro','capitulo_libro','articulo','articulo_web','web')), editorial TEXT, resumen TEXT, url TEXT, visible BOOLEAN NOT NULL DEFAULT TRUE, created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP);

CREATE INDEX bibliography_visible_anio_idx ON bibliography (visible, anio DESC);

CREATE TABLE article_bibliography (article_id VARCHAR(100) NOT NULL REFERENCES articles(id) ON DELETE CASCADE, bibliography_id VARCHAR(100) NOT NULL REFERENCES bibliography(id) ON DELETE CASCADE, PRIMARY KEY (article_id, bibliography_id));

CREATE INDEX article_bibliography_bibliography_idx ON article_bibliography (bibliography_id);

CREATE TABLE catalogs (id VARCHAR(100) PRIMARY KEY, titulo VARCHAR(255) NOT NULL, descripcion TEXT, imagen_portada TEXT, detalles JSONB NOT NULL DEFAULT '{}'::jsonb, visible BOOLEAN NOT NULL DEFAULT TRUE, created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP);

CREATE INDEX catalogs_visible_idx ON catalogs (visible);

CREATE TABLE catalog_items (id VARCHAR(100) NOT NULL, catalog_id VARCHAR(100) NOT NULL REFERENCES catalogs(id) ON DELETE CASCADE, titulo VARCHAR(255) NOT NULL, imagen TEXT NOT NULL, detalles JSONB NOT NULL DEFAULT '{}'::jsonb, PRIMARY KEY (catalog_id,id));

CREATE INDEX catalog_items_catalog_idx ON catalog_items (catalog_id);

CREATE TABLE ads (id VARCHAR(100) PRIMARY KEY, imagen TEXT, imagen_alt TEXT, contacto TEXT, slogan TEXT, descripcion TEXT, vigencia_inicio DATE, vigencia_fin DATE, activo BOOLEAN NOT NULL DEFAULT TRUE, peso INTEGER NOT NULL DEFAULT 1 CHECK (peso >= 1), enlace TEXT, tipo VARCHAR(50) NOT NULL CHECK (tipo IN ('evento','taller','exhibicion','sponsor','comunidad','otro')), paginas TEXT[], prioridad_slot VARCHAR(50), created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP);

CREATE INDEX ads_active_window_idx ON ads (activo, vigencia_inicio, vigencia_fin);
