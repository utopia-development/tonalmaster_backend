package handlers
import("net/http";"github.com/utopia-development/tonalmaster_backend/internal/repository";"github.com/utopia-development/tonalmaster_backend/internal/services")
type ULEHandler struct{s *services.ULEService}
func NewULEHandler(s *services.ULEService)*ULEHandler{return &ULEHandler{s:s}}
func(h *ULEHandler)Articles(w http.ResponseWriter,r *http.Request){v,e:=h.s.Articles(r.Context());if e!=nil{writeError(w,500,"articles_failed","could not load articles");return};writeJSON(w,200,v)}
func(h *ULEHandler)Article(w http.ResponseWriter,r *http.Request){v,e:=h.s.Article(r.Context(),r.PathValue("id"));if e==repository.ErrNotFound{writeError(w,404,"article_not_found","article not found");return};if e!=nil{writeError(w,500,"article_failed","could not load article");return};writeJSON(w,200,v)}
func(h *ULEHandler)Bibliography(w http.ResponseWriter,r *http.Request){v,e:=h.s.Bibliography(r.Context());if e!=nil{writeError(w,500,"bibliography_failed","could not load bibliography");return};writeJSON(w,200,v)}
func(h *ULEHandler)Biblio(w http.ResponseWriter,r *http.Request){v,e:=h.s.Biblio(r.Context(),r.PathValue("id"));if e==repository.ErrNotFound{writeError(w,404,"bibliography_not_found","bibliography entry not found");return};if e!=nil{writeError(w,500,"bibliography_failed","could not load bibliography entry");return};writeJSON(w,200,v)}
func(h *ULEHandler)Catalogs(w http.ResponseWriter,r *http.Request){v,e:=h.s.Catalogs(r.Context());if e!=nil{writeError(w,500,"catalogs_failed","could not load catalogs");return};writeJSON(w,200,v)}
func(h *ULEHandler)Catalog(w http.ResponseWriter,r *http.Request){v,e:=h.s.Catalog(r.Context(),r.PathValue("id"));if e==repository.ErrNotFound{writeError(w,404,"catalog_not_found","catalog not found");return};if e!=nil{writeError(w,500,"catalog_failed","could not load catalog");return};writeJSON(w,200,v)}
func(h *ULEHandler)Ads(w http.ResponseWriter,r *http.Request){v,e:=h.s.Ads(r.Context());if e!=nil{writeError(w,500,"ads_failed","could not load ads");return};writeJSON(w,200,v)}
