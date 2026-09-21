package handlers
import ("context"; "net/http"; "github.com/utopia-development/tonalmaster_backend/internal/repository"; "github.com/utopia-development/tonalmaster_backend/internal/services")
type authContextKey struct{}
func withAuth(c context.Context,u repository.User)context.Context{return context.WithValue(c,authContextKey{},u)}
func currentUser(r *http.Request)(repository.User,bool){u,ok:=r.Context().Value(authContextKey{}).(repository.User);return u,ok}
func RequireAuth(auth *services.AuthService,next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){u,err:=auth.CurrentUser(r.Context(),sessionToken(r));if err!=nil{writeError(w,401,"unauthorized","authentication required");return};next.ServeHTTP(w,r.WithContext(withAuth(r.Context(),u)))})}
