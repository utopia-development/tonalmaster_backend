CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    username VARCHAR(100) NOT NULL UNIQUE,
    avatar_url VARCHAR(255),
    role VARCHAR(50) NOT NULL DEFAULT 'reader'
        CHECK (role IN ('reader', 'contributor', 'admin')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX users_email_lower_idx ON users (LOWER(email));

CREATE TABLE calendar_systems (
    id VARCHAR(50) PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    autor_correlacion VARCHAR(150),
    jdn_base INTEGER,
    descripcion TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE interpretations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    calendar_system_id VARCHAR(50) NOT NULL REFERENCES calendar_systems(id) ON DELETE CASCADE,
    target_date DATE NOT NULL,
    contenido TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX interpretations_calendar_date_idx
    ON interpretations (calendar_system_id, target_date);
CREATE INDEX interpretations_user_date_idx
    ON interpretations (user_id, target_date);

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

CREATE INDEX events_user_date_idx
    ON events (user_id, target_date);
CREATE INDEX events_calendar_date_idx
    ON events (calendar_system_id, target_date);

INSERT INTO calendar_systems (
    id, nombre, autor_correlacion, jdn_base, descripcion
) VALUES (
    'tonalpohualli_caso',
    'Tonalpohualli',
    'CASO',
    NULL,
    'Sistema calendárico inicial de Tonalmaster; la correlación y reglas de cómputo se definirán en la fase de dominio.'
);
