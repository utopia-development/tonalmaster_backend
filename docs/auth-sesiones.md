# Autenticación por sesión (post Fase 5)

## Endpoints

| Método | Ruta | Auth | Respuesta |
|--------|------|------|-----------|
| `POST` | `/api/v1/auth/register` | no | `201` + usuario + cookie |
| `POST` | `/api/v1/auth/login` | no | `200` + usuario + cookie |
| `GET` | `/api/v1/auth/me` | cookie o Bearer | `200` usuario |
| `POST` | `/api/v1/auth/logout` | cookie o Bearer | `204` |

Cookie: `tonalmaster_session` (HttpOnly, SameSite=Lax, Secure fuera de `development`, Max-Age 30 días).

Alternativa: header `Authorization: Bearer <token>`.

### Register body

```json
{"email":"a@example.com","username":"alice","password":"password12"}
```

Password mínimo 8 caracteres (validado en servicio).

### Login body

```json
{"email":"a@example.com","password":"password12"}
```

### Usuario (DTO)

```json
{"id":"<uuid>","email":"...","username":"...","role":"reader"}
```

## Persistencia

- `users` (migración `000001`)
- `sessions` (migración `000002`): `token_hash` (SHA-256 del token opaco), `expires_at`, `revoked_at`

El token en claro **nunca** se guarda; solo el hash.

## Capas

- `internal/repository` — contrato + Postgres
- `internal/services` — bcrypt + emisión/revocación de sesión
- `internal/handlers` — HTTP + cookie

## CORS

Orígenes permitidos reciben `Access-Control-Allow-Credentials: true` para que el frontend pueda enviar cookies en cross-origin (localhost:3000 / 5173).

## Uso en Fase 6

`events` y la creación de `interpretations` requieren sesión válida mediante cookie `tonalmaster_session` o `Authorization: Bearer <token>`. El middleware `RequireAuth` carga el usuario autenticado en el contexto de la petición.

Las consultas públicas de interpretaciones requieren `system` y `date` y no requieren sesión.

## Deuda conocida

- Sin middleware genérico de autenticación (solo `/me` consume sesión).
- Sin rate-limit en login/register.
- CI aplica migraciones 000001 y 000002; el proceso **aún no** ejecuta migraciones al arrancar (hay que aplicarlas a mano o vía pipeline/ops).
- Tests de auth a nivel servicio con mock; falta integración HTTP+Postgres en CI.
