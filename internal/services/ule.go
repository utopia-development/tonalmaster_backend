package services
import("context";"encoding/json";"net/url";"time";"github.com/utopia-development/tonalmaster_backend/internal/repository")
type ULEService struct{repo repository.ULERepository}
func NewULEService(r repository.ULERepository)*ULEService{return &ULEService{repo:r}}
func(s *ULEService)Articles(c context.Context)([]repository.Article,error){return s.repo.ListArticles(c)}
func(s *ULEService)Article(c context.Context,id string)(repository.Article,error){return s.repo.GetArticle(c,id)}
func(s *ULEService)Bibliography(c context.Context)([]repository.Bibliography,error){return s.repo.ListBibliography(c)}
func(s *ULEService)Biblio(c context.Context,id string)(repository.Bibliography,error){return s.repo.GetBibliography(c,id)}
func(s *ULEService)Catalogs(c context.Context)([]repository.Catalog,error){return s.repo.ListCatalogs(c)}
func(s *ULEService)Catalog(c context.Context,id string)(repository.Catalog,[]repository.CatalogItem,error){return s.repo.GetCatalog(c,id)}
func(s *ULEService)Ads(c context.Context)([]repository.Ad,error){return s.repo.ListAds(c,time.Now().UTC())}
func text(v *string)any{if v==nil{return nil};return *v}
func articleDTO(x repository.Article)map[string]any{return map[string]any{"id":x.ID,"titulo":x.Titulo,"autor":text(x.Autor),"fecha":x.Fecha.Format("2006-01-02"),"resumen":x.Resumen,"contenido_html":x.ContenidoHTML,"imagen_destacada":imageURL(x.ImagenDestacada),"imagen_alt":text(x.ImagenAlt),"categoria":text(x.Categoria),"etiquetas":x.Etiquetas,"bibliografía_relacionada":[]string{},"visible":x.Visible}}
func biblioDTO(x repository.Bibliography)map[string]any{return map[string]any{"id":x.ID,"titulo":x.Titulo,"autores":x.Autores,"año":x.Anio,"tipo":x.Tipo,"editorial":text(x.Editorial),"resumen":text(x.Resumen),"url":x.URL,"articulos_relacionados":[]string{},"visible":x.Visible}}
func catalogDTO(x repository.Catalog,items []repository.CatalogItem)map[string]any{var meta map[string]any;_ = json.Unmarshal(x.Detalles,&meta);elements:=make([]map[string]any,0,len(items));for _,i:=range items{var d map[string]any;_ = json.Unmarshal(i.Detalles,&d);e:=map[string]any{"id":i.ID,"titulo":i.Titulo,"imagen":imageURL(&i.Imagen)};for k,v:=range d{e[k]=v};elements=append(elements,e)};return map[string]any{"id":x.ID,"titulo":x.Titulo,"descripcion":text(x.Descripcion),"imagen_portada":imageURL(x.ImagenPortada),"categorias_disponibles":valueOrEmpty(meta,"categorias_disponibles"),"elementos":elements,"visible":x.Visible}}
func adDTO(x repository.Ad)map[string]any{return map[string]any{"id":x.ID,"imagen":imageURL(x.Imagen),"imagen_alt":text(x.ImagenAlt),"contacto":text(x.Contacto),"slogan":text(x.Slogan),"descripcion":text(x.Descripcion),"vigencia_inicio":datePtr(x.Inicio),"vigencia_fin":datePtr(x.Fin),"activo":x.Activo,"peso":x.Peso,"enlace":x.Enlace,"tipo":x.Tipo,"paginas":x.Paginas,"prioridad_slot":text(x.PrioridadSlot)}}
func valueOrEmpty(m map[string]any,k string)any{if v,ok:=m[k];ok{return v};return map[string]any{}}
func datePtr(t *time.Time)any{if t==nil{return nil};return t.Format("2006-01-02")}
func imageURL(v *string)any{if v==nil||*v==""{return nil};if u,e:=url.Parse(*v);e==nil&&u.IsAbs(){return *v};return *v}
