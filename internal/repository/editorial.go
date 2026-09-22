package repository

import "context"

type EditorialRepository interface {
 CreateArticle(context.Context, Article) (Article,error)
 UpdateArticle(context.Context,string, Article) (Article,error)
 DeleteArticle(context.Context,string) error
 SetArticleBibliography(context.Context,string,[]string) error
 CreateBibliography(context.Context,Bibliography) (Bibliography,error)
 UpdateBibliography(context.Context,string,Bibliography) (Bibliography,error)
 DeleteBibliography(context.Context,string) error
 CreateCatalog(context.Context,Catalog) (Catalog,error)
 UpdateCatalog(context.Context,string,Catalog) (Catalog,error)
 DeleteCatalog(context.Context,string) error
 CreateCatalogItem(context.Context,CatalogItem) (CatalogItem,error)
 UpdateCatalogItem(context.Context,string,string,CatalogItem) (CatalogItem,error)
 DeleteCatalogItem(context.Context,string,string) error
 CreateAd(context.Context,Ad) (Ad,error)
 UpdateAd(context.Context,string,Ad) (Ad,error)
 DeleteAd(context.Context,string) error
}