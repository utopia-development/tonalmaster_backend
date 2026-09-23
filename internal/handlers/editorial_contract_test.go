package handlers

import (
    "encoding/json"
    "testing"
)

func TestULEEditorialContractJSONTags(t *testing.T) {
    tests := []struct { name string; value any; required []string; forbidden []string }{
        {"article", articleRequest{ID:"a", Titulo:"T", Autor:strptr("A"), Fecha:"2026-09-23", Resumen:"R", ContenidoHTML:"<p>x</p>", BibliographyIDs:[]string{"b"}}, []string{"id","titulo","autor","fecha","resumen","contenido_html","bibliografía_relacionada"}, []string{"bibliography_ids"}},
        {"bibliography", bibliographyRequest{ID:"b", Titulo:"B", Anio:intptr(2026), Tipo:"libro", ArticleIDs:[]string{"a"}}, []string{"id","titulo","autores","año","tipo","articulos_relacionados"}, []string{"anio","article_ids"}},
        {"catalog", catalogRequest{ID:"c", Titulo:"C", CategoriasDisponibles:map[string][]string{"periodo":[]string{"preclasico"}}}, []string{"id","titulo","categorias_disponibles"}, []string{"detalles"}},
        {"catalog item", itemRequest{ID:"i", Titulo:"I", Imagen:"https://example.org/i.jpg", Categorias:map[string]any{"periodo":"preclasico"}, AnioDescubrimiento:intptr(2026)}, []string{"id","titulo","imagen","categorias","año_descubrimiento"}, []string{"detalles","anio_descubrimiento"}},
        {"ad", adRequest{ID:"ad", Inicio:strptr("2026-10-01"), Fin:strptr("2026-10-31")}, []string{"id","vigencia_inicio","vigencia_fin"}, []string{"inicio","fin"}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            b, err := json.Marshal(tt.value); if err != nil { t.Fatal(err) }
            var got map[string]any; if err := json.Unmarshal(b, &got); err != nil { t.Fatal(err) }
            for _, key := range tt.required { if _, ok := got[key]; !ok { t.Fatalf("missing contract JSON key %q", key) } }
            for _, key := range tt.forbidden { if _, ok := got[key]; ok { t.Fatalf("forbidden/internal JSON key %q leaked", key) } }
        })
    }
}

func strptr(v string) *string { return &v }
func intptr(v int) *int { return &v }