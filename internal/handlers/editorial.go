package handlers

import (
 "encoding/json"
 "net/http"
 "time"
 "github.com/utopia-development/tonalmaster_backend/internal/repository"
 "github.com/utopia-development/tonalmaster_backend/internal/services"
)

type EditorialHandler struct{s *services.EditorialService}
func NewEditorialHandler(s *services.EditorialService)*EditorialHandler{return &EditorialHandler{s:s}}
func(h *EditorialHandler)auth(w http.ResponseWriter,r *http.Request)(repository.User,bool){u,ok:=currentUser(r);if !ok||!h.s.IsAuthorized(u.Role){writeError(w,http.StatusForbidden,"forbidden","editorial permission required");return repository.User{},false};return u,true}
type articleRequest struct {
 ID string `json:"id"`
 Titulo string `json:"titulo"`
 Autor *string `json:"autor"`
 Fecha string `json:"fecha"`
 Resumen string `json:"resumen"`
 ContenidoHTML string `json:"contenido_html"`
 ImagenDestacada *string `json:"imagen_destacada"`
 ImagenAlt *string `json:"imagen_alt"`
 Categoria *string `json:"categoria"`
 Etiquetas []string `json:"etiquetas"`
 Visible bool `json:"visible"`
 BibliographyIDs []string `json:"bibliografía_relacionada"`
}
func articleFrom(q articleRequest)(repository.Article,error){d,e:=time.Parse("2006-01-02",q.Fecha);return repository.Article{ID:q.ID,Titulo:q.Titulo,Autor:q.Autor,Fecha:d,Resumen:q.Resumen,ContenidoHTML:q.ContenidoHTML,ImagenDestacada:q.ImagenDestacada,ImagenAlt:q.ImagenAlt,Categoria:q.Categoria,Etiquetas:q.Etiquetas,Visible:q.Visible,BibliographyIDs:q.BibliographyIDs},e}
func(h *EditorialHandler)CreateArticle(w http.ResponseWriter,r *http.Request){u,ok:=h.auth(w,r);_ = u;if !ok{return};var q articleRequest;if json.NewDecoder(r.Body).Decode(&q)!=nil||q.ID==""||q.Titulo==""||q.Autor==nil||*q.Autor==""||q.Fecha==""||q.Resumen==""||q.ContenidoHTML==""{writeError(w,400,"invalid_request","id, titulo, autor, fecha, resumen and contenido_html are required");return};x,e:=articleFrom(q);if e!=nil{writeError(w,400,"invalid_date","fecha must use YYYY-MM-DD");return};x,e=h.s.CreateArticle(r.Context(),u.Role,x);if e!=nil{writeError(w,400,"article_create_failed","could not create article");return};if e=h.s.SetArticleBibliography(r.Context(),u.Role,x.ID,q.BibliographyIDs);e!=nil{writeError(w,400,"article_bibliography_update_failed","could not update article bibliography");return};writeJSON(w,201,makeArticleResponse(x))}
func(h *EditorialHandler)UpdateArticle(w http.ResponseWriter,r *http.Request){u,ok:=h.auth(w,r);if !ok{return};var q articleRequest;if json.NewDecoder(r.Body).Decode(&q)!=nil{writeError(w,400,"invalid_request","invalid JSON");return};q.ID=r.PathValue("id");x,e:=articleFrom(q);if e!=nil{writeError(w,400,"invalid_date","fecha must use YYYY-MM-DD");return};x,e=h.s.UpdateArticle(r.Context(),u.Role,q.ID,x);if e==repository.ErrNotFound{writeError(w,404,"article_not_found","article not found");return};if e!=nil{writeError(w,400,"article_update_failed","could not update article");return};if e=h.s.SetArticleBibliography(r.Context(),u.Role,q.ID,q.BibliographyIDs);e!=nil{writeError(w,400,"article_bibliography_update_failed","could not update article bibliography");return};writeJSON(w,200,makeArticleResponse(x))}
func(h *EditorialHandler)DeleteArticle(w http.ResponseWriter,r *http.Request){u,ok:=h.auth(w,r);if !ok{return};e:=h.s.DeleteArticle(r.Context(),u.Role,r.PathValue("id"));if e==repository.ErrNotFound{writeError(w,404,"article_not_found","article not found");return};if e!=nil{writeError(w,500,"article_delete_failed","could not delete article");return};w.WriteHeader(204)}
type bibliographyRequest struct {
 ID string `json:"id"`
 Titulo string `json:"titulo"`
 Autores []string `json:"autores"`
 Anio *int `json:"año"`
 Tipo string `json:"tipo"`
 Editorial *string `json:"editorial"`
 Resumen *string `json:"resumen"`
 URL *string `json:"url"`
 Visible bool `json:"visible"`
 ArticleIDs []string `json:"articulos_relacionados"`
}
func(h *EditorialHandler)CreateBibliography(w http.ResponseWriter,r *http.Request){u,ok:=h.auth(w,r);if !ok{return};var q bibliographyRequest;if json.NewDecoder(r.Body).Decode(&q)!=nil||q.ID==""||q.Titulo==""||q.Tipo==""{writeError(w,400,"invalid_request","id, titulo and tipo are required");return};x:=repository.Bibliography{ID:q.ID,Titulo:q.Titulo,Autores:q.Autores,Anio:q.Anio,Tipo:q.Tipo,Editorial:q.Editorial,Resumen:q.Resumen,URL:q.URL,Visible:q.Visible,ArticleIDs:q.ArticleIDs};x,e:=h.s.CreateBibliography(r.Context(),u.Role,x);if e!=nil{writeError(w,400,"bibliography_create_failed","could not create bibliography");return};writeJSON(w,201,x)}
func(h *EditorialHandler)UpdateBibliography(w http.ResponseWriter,r *http.Request){u,ok:=h.auth(w,r);if !ok{return};var q bibliographyRequest;if json.NewDecoder(r.Body).Decode(&q)!=nil{writeError(w,400,"invalid_request","invalid JSON");return};x:=repository.Bibliography{ID:r.PathValue("id"),Titulo:q.Titulo,Autores:q.Autores,Anio:q.Anio,Tipo:q.Tipo,Editorial:q.Editorial,Resumen:q.Resumen,URL:q.URL,Visible:q.Visible};x,e:=h.s.UpdateBibliography(r.Context(),u.Role,x.ID,x);if e==repository.ErrNotFound{writeError(w,404,"bibliography_not_found","not found");return};if e!=nil{writeError(w,400,"bibliography_update_failed","could not update");return};writeJSON(w,200,x)}
func(h *EditorialHandler)DeleteBibliography(w http.ResponseWriter,r *http.Request){u,ok:=h.auth(w,r);if !ok{return};e:=h.s.DeleteBibliography(r.Context(),u.Role,r.PathValue("id"));if e==repository.ErrNotFound{writeError(w,404,"bibliography_not_found","not found");return};if e!=nil{writeError(w,500,"bibliography_delete_failed","could not delete");return};w.WriteHeader(204)}
type catalogRequest struct {
 ID string `json:"id"`
 Titulo string `json:"titulo"`
 Descripcion *string `json:"descripcion"`
 ImagenPortada *string `json:"imagen_portada"`
 Detalles json.RawMessage `json:"detalles"`
 Visible bool `json:"visible"`
}
func(h *EditorialHandler)CreateCatalog(w http.ResponseWriter,r *http.Request){u,ok:=h.auth(w,r);if !ok{return};var q catalogRequest;if json.NewDecoder(r.Body).Decode(&q)!=nil||q.ID==""||q.Titulo==""{writeError(w,400,"invalid_request","id and titulo are required");return};if len(q.Detalles)==0{q.Detalles=[]byte("{}")};x,e:=h.s.CreateCatalog(r.Context(),u.Role,repository.Catalog{ID:q.ID,Titulo:q.Titulo,Descripcion:q.Descripcion,ImagenPortada:q.ImagenPortada,Detalles:q.Detalles,Visible:q.Visible});if e!=nil{writeError(w,400,"catalog_create_failed","could not create catalog");return};writeJSON(w,201,makeCatalogResponse(x,nil))}
func(h *EditorialHandler)UpdateCatalog(w http.ResponseWriter,r *http.Request){u,ok:=h.auth(w,r);if !ok{return};var q catalogRequest;if json.NewDecoder(r.Body).Decode(&q)!=nil{writeError(w,400,"invalid_request","invalid JSON");return};if len(q.Detalles)==0{q.Detalles=[]byte("{}")};x,e:=h.s.UpdateCatalog(r.Context(),u.Role,r.PathValue("id"),repository.Catalog{ID:r.PathValue("id"),Titulo:q.Titulo,Descripcion:q.Descripcion,ImagenPortada:q.ImagenPortada,Detalles:q.Detalles,Visible:q.Visible});if e==repository.ErrNotFound{writeError(w,404,"catalog_not_found","not found");return};if e!=nil{writeError(w,400,"catalog_update_failed","could not update");return};writeJSON(w,200,makeCatalogResponse(x,nil))}
func(h *EditorialHandler)DeleteCatalog(w http.ResponseWriter,r *http.Request){u,ok:=h.auth(w,r);if !ok{return};e:=h.s.DeleteCatalog(r.Context(),u.Role,r.PathValue("id"));if e==repository.ErrNotFound{writeError(w,404,"catalog_not_found","not found");return};if e!=nil{writeError(w,500,"catalog_delete_failed","could not delete");return};w.WriteHeader(204)}
type itemRequest struct {
 ID string `json:"id"`
 Titulo string `json:"titulo"`
 Imagen string `json:"imagen"`
 Detalles json.RawMessage `json:"detalles"`
}
func(h *EditorialHandler)CreateItem(w http.ResponseWriter,r *http.Request){u,ok:=h.auth(w,r);if !ok{return};var q itemRequest;if json.NewDecoder(r.Body).Decode(&q)!=nil||q.ID==""||q.Titulo==""||q.Imagen==""{writeError(w,400,"invalid_request","id, titulo and imagen are required");return};details,_:=json.Marshal(map[string]any{"imagen_alt":q.ImagenAlt,"descripcion":q.Descripcion,"categorias":q.Categorias,"año_descubrimiento":q.AnioDescubrimiento,"ubicacion":q.Ubicacion});x,e:=h.s.CreateCatalogItem(r.Context(),u.Role,repository.CatalogItem{ID:q.ID,CatalogID:r.PathValue("id"),Titulo:q.Titulo,Imagen:q.Imagen,Detalles:details});if e!=nil{writeError(w,400,"item_create_failed","could not create item");return};writeJSON(w,201,makeItemResponse(x))}
func(h *EditorialHandler)UpdateItem(w http.ResponseWriter,r *http.Request){u,ok:=h.auth(w,r);if !ok{return};var q itemRequest;if json.NewDecoder(r.Body).Decode(&q)!=nil{writeError(w,400,"invalid_request","invalid JSON");return};if len(q.Detalles)==0{q.Detalles=[]byte("{}")};x,e:=h.s.UpdateCatalogItem(r.Context(),u.Role,r.PathValue("id"),r.PathValue("item_id"),repository.CatalogItem{ID:r.PathValue("item_id"),CatalogID:r.PathValue("id"),Titulo:q.Titulo,Imagen:q.Imagen,Detalles:q.Detalles});if e==repository.ErrNotFound{writeError(w,404,"item_not_found","not found");return};if e!=nil{writeError(w,400,"item_update_failed","could not update");return};writeJSON(w,200,makeItemResponse(x))}
func(h *EditorialHandler)DeleteItem(w http.ResponseWriter,r *http.Request){u,ok:=h.auth(w,r);if !ok{return};e:=h.s.DeleteCatalogItem(r.Context(),u.Role,r.PathValue("id"),r.PathValue("item_id"));if e==repository.ErrNotFound{writeError(w,404,"item_not_found","not found");return};if e!=nil{writeError(w,500,"item_delete_failed","could not delete");return};w.WriteHeader(204)}
type adRequest struct {
 ID string `json:"id"`
 Imagen *string `json:"imagen"`
 ImagenAlt *string `json:"imagen_alt"`
 Contacto *string `json:"contacto"`
 Slogan *string `json:"slogan"`
 Descripcion *string `json:"descripcion"`
 Inicio *string `json:"vigencia_inicio"`
 Fin *string `json:"vigencia_fin"`
 Activo bool `json:"activo"`
 Peso int `json:"peso"`
 Enlace *string `json:"enlace"`
 Tipo string `json:"tipo"`
 Paginas []string `json:"paginas"`
 PrioridadSlot *string `json:"prioridad_slot"`
}
func(h *EditorialHandler)CreateAd(w http.ResponseWriter,r *http.Request){u,ok:=h.auth(w,r);if !ok{return};var q adRequest;if json.NewDecoder(r.Body).Decode(&q)!=nil{writeError(w,400,"invalid_request","invalid JSON");return};if q.ID==""{writeError(w,400,"invalid_request","id is required");return};x,e:=parseAd(q,q.ID);if e!=nil{writeError(w,400,"invalid_request",e.Error());return};x,e=h.s.CreateAd(r.Context(),u.Role,x);if e!=nil{writeError(w,400,"ad_create_failed","could not create ad");return};writeJSON(w,201,x)}
func parseAd(q adRequest,id string)(repository.Ad,error){var a repository.Ad;a.ID=id;a.Imagen=q.Imagen;a.ImagenAlt=q.ImagenAlt;a.Contacto=q.Contacto;a.Slogan=q.Slogan;a.Descripcion=q.Descripcion;a.Activo=q.Activo;a.Peso=q.Peso;a.Enlace=q.Enlace;a.Tipo=q.Tipo;a.Paginas=q.Paginas;a.PrioridadSlot=q.PrioridadSlot;if q.Peso==0{a.Peso=1};if q.Inicio!=nil && *q.Inicio!=""{d,err:=time.Parse("2006-01-02",*q.Inicio);if err!=nil{return a,err};a.Inicio=&d};if q.Fin!=nil && *q.Fin!=""{d,err:=time.Parse("2006-01-02",*q.Fin);if err!=nil{return a,err};a.Fin=&d};return a,nil}
func(h *EditorialHandler)UpdateAd(w http.ResponseWriter,r *http.Request){u,ok:=h.auth(w,r);if !ok{return};var q adRequest;if json.NewDecoder(r.Body).Decode(&q)!=nil{writeError(w,400,"invalid_request","invalid JSON");return};x,e:=parseAd(q,r.PathValue("id"));if e!=nil{writeError(w,400,"invalid_request",e.Error());return};x,e=h.s.UpdateAd(r.Context(),u.Role,x.ID,x);if e==repository.ErrNotFound{writeError(w,404,"ad_not_found","not found");return};if e!=nil{writeError(w,400,"ad_update_failed","could not update");return};writeJSON(w,200,x)}
func(h *EditorialHandler)DeleteAd(w http.ResponseWriter,r *http.Request){u,ok:=h.auth(w,r);if !ok{return};e:=h.s.DeleteAd(r.Context(),u.Role,r.PathValue("id"));if e==repository.ErrNotFound{writeError(w,404,"ad_not_found","not found");return};if e!=nil{writeError(w,500,"ad_delete_failed","could not delete");return};w.WriteHeader(204)}


type articleDTOResponse struct {
 ID string `json:"id"`; Titulo string `json:"titulo"`; Autor *string `json:"autor"`; Fecha string `json:"fecha"`; Resumen string `json:"resumen"`; ContenidoHTML string `json:"contenido_html"`; ImagenDestacada *string `json:"imagen_destacada"`; ImagenAlt *string `json:"imagen_alt"`; Categoria *string `json:"categoria"`; Etiquetas []string `json:"etiquetas"`; BibliographyIDs []string `json:"bibliografía_relacionada"`; Visible bool `json:"visible""
}
func makeArticleResponse(x repository.Article) articleDTOResponse { return articleDTOResponse{ID:x.ID,Titulo:x.Titulo,Autor:x.Autor,Fecha:x.Fecha.Format("2006-01-02"),Resumen:x.Resumen,ContenidoHTML:x.ContenidoHTML,ImagenDestacada:x.ImagenDestacada,ImagenAlt:x.ImagenAlt,Categoria:x.Categoria,Etiquetas:x.Etiquetas,BibliographyIDs:x.BibliographyIDs,Visible:x.Visible} }

type bibliographyDTOResponse struct { ID string `json:"id"`; Titulo string `json:"titulo"`; Autores []string `json:"autores"`; Anio *int `json:"año"`; Tipo string `json:"tipo"`; Editorial *string `json:"editorial"`; Resumen *string `json:"resumen"`; URL *string `json:"url"`; ArticleIDs []string `json:"articulos_relacionados"`; Visible bool `json:"visible"` }
func makeBibliographyResponse(x repository.Bibliography) bibliographyDTOResponse { return bibliographyDTOResponse{ID:x.ID,Titulo:x.Titulo,Autores:x.Autores,Anio:x.Anio,Tipo:x.Tipo,Editorial:x.Editorial,Resumen:x.Resumen,URL:x.URL,ArticleIDs:x.ArticleIDs,Visible:x.Visible} }

type catalogDTOResponse struct { ID string `json:"id"`; Titulo string `json:"titulo"`; Descripcion *string `json:"descripcion"`; ImagenPortada *string `json:"imagen_portada"`; CategoriasDisponibles map[string]any `json:"categorias_disponibles"`; Elementos []itemDTOResponse `json:"elementos"`; Visible bool `json:"visible"` }
func makeCatalogResponse(x repository.Catalog, items []repository.CatalogItem) catalogDTOResponse { var d map[string]any; _=json.Unmarshal(x.Detalles,&d); cats:=map[string]any{};if v,ok:=d["categorias_disponibles"].(map[string]any);ok{cats=v};out:=make([]itemDTOResponse,0,len(items));for _,i:=range items{out=append(out,makeItemResponse(i))};return catalogDTOResponse{ID:x.ID,Titulo:x.Titulo,Descripcion:x.Descripcion,ImagenPortada:x.ImagenPortada,CategoriasDisponibles:cats,Elementos:out,Visible:x.Visible} }

type itemDTOResponse struct { ID string `json:"id"`; Titulo string `json:"titulo"`; Imagen string `json:"imagen"`; ImagenAlt *string `json:"imagen_alt"`; Descripcion *string `json:"descripcion"`; Categorias map[string]any `json:"categorias"`; AnioDescubrimiento *int `json:"año_descubrimiento"`; Ubicacion *string `json:"ubicacion"` }
func makeItemResponse(x repository.CatalogItem) itemDTOResponse { var d map[string]any;_=json.Unmarshal(x.Detalles,&d); out:=itemResponse{ID:x.ID,Titulo:x.Titulo,Imagen:x.Imagen}; if v,ok:=d["imagen_alt"].(string);ok{out.ImagenAlt=&v};if v,ok:=d["descripcion"].(string);ok{out.Descripcion=&v};if v,ok:=d["categorias"].(map[string]any);ok{out.Categorias=v};if v,ok:=d["año_descubrimiento"].(float64);ok{n:=int(v);out.AnioDescubrimiento=&n};if v,ok:=d["ubicacion"].(string);ok{out.Ubicacion=&v};return out }

type adResponse struct { ID string `json:"id"`; Imagen *string `json:"imagen"`; ImagenAlt *string `json:"imagen_alt"`; Contacto *string `json:"contacto"`; Slogan *string `json:"slogan"`; Descripcion *string `json:"descripcion"`; Inicio *string `json:"vigencia_inicio"`; Fin *string `json:"vigencia_fin"`; Activo bool `json:"activo"`; Peso int `json:"peso"`; Enlace *string `json:"enlace"`; Tipo string `json:"tipo"`; Paginas []string `json:"paginas"`; PrioridadSlot *string `json:"prioridad_slot"` }
func adResponse(x repository.Ad) adResponse { f:=func(t *time.Time)*string{if t==nil{return nil};v:=t.Format("2006-01-02");return &v};return adResponse{ID:x.ID,Imagen:x.Imagen,ImagenAlt:x.ImagenAlt,Contacto:x.Contacto,Slogan:x.Slogan,Descripcion:x.Descripcion,Inicio:f(x.Inicio),Fin:f(x.Fin),Activo:x.Activo,Peso:x.Peso,Enlace:x.Enlace,Tipo:x.Tipo,Paginas:x.Paginas,PrioridadSlot:x.PrioridadSlot} }
