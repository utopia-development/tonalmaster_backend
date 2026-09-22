# Fase 7 — Contenido público para ule_educativo

## Objetivo
Exponer contenido editorial de Ule mediante una API pública estable, sin autenticación, usando exactamente `docs/contrato_datos.md`.

## Implementado
- Migración `000003_ule_content`: articles, bibliography, article_bibliography, catalogs, catalog_items y ads.
- Filtrado de visibilidad en SQL.
- Repositorios PostgreSQL y servicio/handlers separados.
- Rutas públicas:
  - GET /api/v1/articles
  - GET /api/v1/articles/{id}
  - GET /api/v1/bibliography
  - GET /api/v1/bibliography/{id}
  - GET /api/v1/catalogs
  - GET /api/v1/catalogs/{id}
  - GET /api/v1/ads
- Relaciones artículo↔bibliografía en ambos sentidos.
- Reconstrucción de `categorias_disponibles`, `elementos[].categorias` y campos del anuncio desde JSONB/columnas internas.
- Contrato sincronizado desde `ule_educativo/docs/contrato_datos.md`.

## Pendiente para cierre
1. Tests unitarios de DTO/mapeos.
2. Tests HTTP de 200/404/lista vacía/visibilidad.
3. Carga de datos editoriales reales.
4. Ajuste del loader de Ule a las rutas inglesas y configuración `dataSource=api`.
5. Prueba integrada desde Windows/WSL sobre la red local.
6. Confirmar que las URLs de imágenes servidas por API son resolubles desde el navegador y cumplen el contrato.

## Criterio de aceptación
El frontend debe poder cambiar de JSON local a API sin modificar sus componentes ni el contrato de datos.
