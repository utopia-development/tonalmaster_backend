package repository
import("context";"time";"github.com/jackc/pgx/v5/pgtype")
type Article struct{ID string;Titulo string;Autor *string;Fecha time.Time;Resumen string;ContenidoHTML string;ImagenDestacada *string;ImagenAlt *string;Categoria *string;Etiquetas []string;Visible bool}
type Bibliography struct{ID string;Titulo string;Autores []string;Anio *int;Tipo string;Editorial *string;Resumen *string;URL *string;Visible bool}
type Catalog struct{ID string;Titulo string;Descripcion *string;ImagenPortada *string;Detalles []byte;Visible bool}
type CatalogItem struct{ID string;CatalogID string;Titulo string;Imagen string;Detalles []byte}
type Ad struct{ID string;Imagen *string;ImagenAlt *string;Contacto *string;Slogan *string;Descripcion *string;Inicio *time.Time;Fin *time.Time;Activo bool;Peso int;Enlace *string;Tipo string;Paginas []string;PrioridadSlot *string}
type ULERepository interface{ListArticles(context.Context)([]Article,error);GetArticle(context.Context,string)(Article,error);ListBibliography(context.Context)([]Bibliography,error);GetBibliography(context.Context,string)(Bibliography,error);ListCatalogs(context.Context)([]Catalog,error);GetCatalog(context.Context,string)(Catalog,[]CatalogItem,error);ListAds(context.Context,time.Time)([]Ad,error)}
var _=pgtype.UUID{}
