package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/utopia-development/tonalmaster_backend/internal/services"
)

type AuthHandler struct {
	auth *services.AuthService
	secureCookie bool
}

func NewAuthHandler(auth *services.AuthService, secureCookie bool) *AuthHandler {
	return &AuthHandler{auth: auth, secureCookie: secureCookie}
}

type authRequest struct {
	Email string `json:"email"`
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}

type userDTO struct {
	ID string `json:"id"`
	Email string `json:"email"`
	Username string `json:"username"`
	Role string `json:"role"`
}

func (h *AuthHandler) Register(w http.ResponseWriter,r *http.Request){
	var req authRequest
	if err:=json.NewDecoder(r.Body).Decode(&req);err!=nil||req.Username==""{writeError(w,400,"invalid_request","email, username and password are required");return}
	u,t,err:=h.auth.Register(r.Context(),req.Email,req.Username,req.Password);if err!=nil{writeError(w,400,"registration_failed","could not create account");return}
	h.setSessionCookie(w,t);writeJSON(w,201,userDTO{ID:u.ID.String(),Email:u.Email,Username:u.Username,Role:u.Role})
}
func (h *AuthHandler) Login(w http.ResponseWriter,r *http.Request){
	var req authRequest
	if err:=json.NewDecoder(r.Body).Decode(&req);err!=nil{writeError(w,400,"invalid_request","email and password are required");return}
	u,t,err:=h.auth.Login(r.Context(),req.Email,req.Password);if err!=nil{writeError(w,401,"invalid_credentials","invalid credentials");return}
	h.setSessionCookie(w,t);writeJSON(w,200,userDTO{ID:u.ID.String(),Email:u.Email,Username:u.Username,Role:u.Role})
}
func (h *AuthHandler) Me(w http.ResponseWriter,r *http.Request){
	u,err:=h.auth.CurrentUser(r.Context(),sessionToken(r));if err!=nil{writeError(w,401,"unauthorized","authentication required");return}
	writeJSON(w,200,userDTO{ID:u.ID.String(),Email:u.Email,Username:u.Username,Role:u.Role})
}
func (h *AuthHandler) Logout(w http.ResponseWriter,r *http.Request){
	if err:=h.auth.Logout(r.Context(),sessionToken(r));err!=nil{writeError(w,500,"logout_failed","could not revoke session");return}
	http.SetCookie(w,&http.Cookie{Name:"tonalmaster_session",Value:"",Path:"/",MaxAge:-1,HttpOnly:true,Secure:h.secureCookie,SameSite:http.SameSiteLaxMode})
	w.WriteHeader(http.StatusNoContent)
}
func (h *AuthHandler) setSessionCookie(w http.ResponseWriter,token string){
	http.SetCookie(w,&http.Cookie{Name:"tonalmaster_session",Value:token,Path:"/",HttpOnly:true,Secure:h.secureCookie,SameSite:http.SameSiteLaxMode,MaxAge:30*24*60*60})
}
func sessionToken(r *http.Request)string{
	if c,err:=r.Cookie("tonalmaster_session");err==nil{return c.Value}
	v:=strings.TrimSpace(r.Header.Get("Authorization"));if strings.HasPrefix(v,"Bearer "){return strings.TrimSpace(strings.TrimPrefix(v,"Bearer "))}
	return ""
}
