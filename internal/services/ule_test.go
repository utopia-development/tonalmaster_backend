package services

import (
	"context"
	"testing"
	"time"

	"github.com/utopia-development/tonalmaster_backend/internal/repository"
)

type fakeULERepository struct {
	articles []repository.Article
	bibliography []repository.Bibliography
	catalog repository.Catalog
	items []repository.CatalogItem
	ads []repository.Ad
}

func (f fakeULERepository) ListArticles(context.Context) ([]repository.Article,error) { return f.articles,nil }
func (f fakeULERepository) GetArticle(context.Context,string) (repository.Article,error) { return f.articles[0],nil }
func (f fakeULERepository) ListBibliography(context.Context) ([]repository.Bibliography,error) { return f.bibliography,nil }
func (f fakeULERepository) GetBibliography(context.Context,string) (repository.Bibliography,error) { return f.bibliography[0],nil }
func (f fakeULERepository) ListCatalogs(context.Context) ([]repository.Catalog,error) { return []repository.Catalog{f.catalog},nil }
func (f fakeULERepository) GetCatalog(context.Context,string) (repository.Catalog,[]repository.CatalogItem,error) { return f.catalog,f.items,nil }
func (f fakeULERepository) ListAds(context.Context,time.Time) ([]repository.Ad,error) { return f.ads,nil }

func TestULEServiceContractShape(t *testing.T) {
	article := repository.Article{ID:"articulo-001",Titulo:"Título",Fecha:time.Date(2026,9,22,0,0,0,time.UTC),Resumen:"Resumen",ContenidoHTML:"<p>ok</p>",Etiquetas:[]string{"uno"},Visible:true,BibliographyIDs:[]string{"biblio-001"}}
	biblio := repository.Bibliography{ID:"biblio-001",Titulo:"Referencia",Autores:[]string{"Autor"},Anio:ptrInt(2026),Tipo:"libro",ArticleIDs:[]string{"articulo-001"}}
	catalog := repository.Catalog{ID:"catalogo-001",Titulo:"Catálogo",Detalles:[]byte(`{"categorias_disponibles":{"periodo":["preclasico"]}}`),Visible:true}
	item := repository.CatalogItem{ID:"pieza-001",CatalogID:"catalogo-001",Titulo:"Pieza",Imagen:"https://example.test/pieza.jpg",Detalles:[]byte(`{"descripcion":"Desc","categorias":{"periodo":"preclasico"},"año_descubrimiento":1980}`)}
	ad := repository.Ad{ID:"ad-001",Tipo:"evento",Activo:true,Peso:1,Inicio:ptrTime(time.Date(2026,9,1,0,0,0,time.UTC)),Fin:ptrTime(time.Date(2026,10,1,0,0,0,time.UTC))}

	s := NewULEService(fakeULERepository{articles:[]repository.Article{article},bibliography:[]repository.Bibliography{biblio},catalog:catalog,items:[]repository.CatalogItem{item},ads:[]repository.Ad{ad}})
	arts,_:=s.Articles(context.Background()); if arts[0]["bibliografía_relacionada"]==nil { t.Fatal("missing bibliografía_relacionada") }
	bibs,_:=s.Bibliography(context.Background()); if bibs[0]["articulos_relacionados"]==nil { t.Fatal("missing articulos_relacionados") }
	cats,_:=s.Catalog(context.Background(),"catalogo-001"); if len(cats["elementos"].([]map[string]any))!=1 { t.Fatal("catalog elements not reconstructed") }
	if cats["categorias_disponibles"]==nil { t.Fatal("missing categorias_disponibles") }
	ads,_:=s.Ads(context.Background()); if ads[0]["vigencia_inicio"]!="2026-09-01" || ads[0]["vigencia_fin"]!="2026-10-01" { t.Fatal("ad dates not normalized") }
}

func ptrTime(v time.Time)*time.Time{return &v}\nfunc ptrInt(v int)*int{return &v}
