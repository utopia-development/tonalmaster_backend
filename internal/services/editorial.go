package services

import (
 "context"
 "errors"
 "github.com/utopia-development/tonalmaster_backend/internal/repository"
)

var ErrForbidden = errors.New("forbidden")

type EditorialService struct { repo repository.EditorialRepository }
func NewEditorialService(r repository.EditorialRepository)*EditorialService{return &EditorialService{repo:r}}

func(s *EditorialService) authorize(role string)error{if role!="contributor"&&role!="admin"{return ErrForbidden};return nil}
func(s *EditorialService) CreateArticle(c context.Context,role string,x repository.Article)(repository.Article,error){if e:=s.authorize(role);e!=nil{return repository.Article{},e};return s.repo.CreateArticle(c,x)}
func(s *EditorialService) UpdateArticle(c context.Context,role,id string,x repository.Article)(repository.Article,error){if e:=s.authorize(role);e!=nil{return repository.Article{},e};return s.repo.UpdateArticle(c,id,x)}
func(s *EditorialService) DeleteArticle(c context.Context,role,id string)error{if e:=s.authorize(role);e!=nil{return e};return s.repo.DeleteArticle(c,id)}
func(s *EditorialService) CreateBibliography(c context.Context,role string,x repository.Bibliography)(repository.Bibliography,error){if e:=s.authorize(role);e!=nil{return repository.Bibliography{},e};return s.repo.CreateBibliography(c,x)}
func(s *EditorialService) UpdateBibliography(c context.Context,role,id string,x repository.Bibliography)(repository.Bibliography,error){if e:=s.authorize(role);e!=nil{return repository.Bibliography{},e};return s.repo.UpdateBibliography(c,id,x)}
func(s *EditorialService) DeleteBibliography(c context.Context,role,id string)error{if e:=s.authorize(role);e!=nil{return e};return s.repo.DeleteBibliography(c,id)}
func(s *EditorialService) SetArticleBibliography(c context.Context,role,id string,ids []string)error{if e:=s.authorize(role);e!=nil{return e};return s.repo.SetArticleBibliography(c,id,ids)}
func(s *EditorialService) CreateCatalog(c context.Context,role string,x repository.Catalog)(repository.Catalog,error){if e:=s.authorize(role);e!=nil{return repository.Catalog{},e};return s.repo.CreateCatalog(c,x)}
func(s *EditorialService) UpdateCatalog(c context.Context,role,id string,x repository.Catalog)(repository.Catalog,error){if e:=s.authorize(role);e!=nil{return repository.Catalog{},e};return s.repo.UpdateCatalog(c,id,x)}
func(s *EditorialService) DeleteCatalog(c context.Context,role,id string)error{if e:=s.authorize(role);e!=nil{return e};return s.repo.DeleteCatalog(c,id)}
func(s *EditorialService) CreateCatalogItem(c context.Context,role string,x repository.CatalogItem)(repository.CatalogItem,error){if e:=s.authorize(role);e!=nil{return repository.CatalogItem{},e};return s.repo.CreateCatalogItem(c,x)}
func(s *EditorialService) UpdateCatalogItem(c context.Context,role,cid,id string,x repository.CatalogItem)(repository.CatalogItem,error){if e:=s.authorize(role);e!=nil{return repository.CatalogItem{},e};return s.repo.UpdateCatalogItem(c,cid,id,x)}
func(s *EditorialService) DeleteCatalogItem(c context.Context,role,cid,id string)error{if e:=s.authorize(role);e!=nil{return e};return s.repo.DeleteCatalogItem(c,cid,id)}
func(s *EditorialService) CreateAd(c context.Context,role string,x repository.Ad)(repository.Ad,error){if e:=s.authorize(role);e!=nil{return repository.Ad{},e};return s.repo.CreateAd(c,x)}
func(s *EditorialService) UpdateAd(c context.Context,role,id string,x repository.Ad)(repository.Ad,error){if e:=s.authorize(role);e!=nil{return repository.Ad{},e};return s.repo.UpdateAd(c,id,x)}
func(s *EditorialService) DeleteAd(c context.Context,role,id string)error{if e:=s.authorize(role);e!=nil{return e};return s.repo.DeleteAd(c,id)}
