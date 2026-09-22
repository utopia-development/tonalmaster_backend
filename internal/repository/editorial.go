package repository

import (
 "context"
 "time"
)

type Article struct { ID,Titulo,ContenidoHTML string; Autor *string; Fecha time.Time; Resumen string; ImagenDestacada,ImagenAlt,Categoria *string; Etiquetas []string; Visible bool; BibliographyIDs []string }
type Bibliography struct { ID,Titulo string; Autores []string; Anio *int; Tipo string; Editorial,Resumen,URL *string; Visible bool; ArticleIDs []string }
type Catalog struct { ID,Titulo string; Descripcion,ImagenPortada *string; Detalles []byte; Visible bool }
type CatalogItem struct { ID,CatalogID,Titulo,Imagen string; Detalles []byte }
type Ad struct { ID string; Imagen,ImagenAlt,Contacto,Slogan,Descripcion *string; Inicio,Fin *time.Time; Activo bool; Peso int; Enlace *string; Tipo string; Paginas []string; PrioridadSlot *string }

type EditorialRepository interface {
 CreateArticle(context.Context,Article) (Article,error); UpdateArticle(context.Context,string,Article) (Article,error); DeleteArticle(context.Context,string) error
 CreateBibliography(context.Context,Bibliography) (Bibliography,error); UpdateBibliography(context.Context,string,Bibliography) (Bibliography,error); DeleteBibliography(context.Context,string) error
 SetArticleBibliography(context.Context,string,[]string) error
 CreateCatalog(context.Context,Catalog) (Catalog,error); UpdateCatalog(context.Context,string,Catalog) (Catalog,error); DeleteCatalog(context.Context,string) error
 CreateCatalogItem(context.Context,CatalogItem) (CatalogItem,error); UpdateCatalogItem(context.Context,string,string,CatalogItem) (CatalogItem,error); DeleteCatalogItem(context.Context,string,string) error
 CreateAd(context.Context,Ad) (Ad,error); UpdateAd(context.Context,string,Ad) (Ad,error); DeleteAd(context.Context,string) error
}