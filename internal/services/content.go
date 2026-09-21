package services
import ("context"; "time"; "github.com/jackc/pgx/v5/pgtype"; "github.com/utopia-development/tonalmaster_backend/internal/repository")
type ContentService struct{repo repository.ContentRepository}
func NewContentService(r repository.ContentRepository)*ContentService{return &ContentService{repo:r}}
func(s *ContentService)CreateEvent(c context.Context,u pgtype.UUID,sys string,d time.Time,t string,desc *string)(repository.Event,error){return s.repo.CreateEvent(c,u,sys,d,t,desc)}
func(s *ContentService)ListEvents(c context.Context,u pgtype.UUID,f,to time.Time)([]repository.Event,error){return s.repo.ListEvents(c,u,f,to)}
func(s *ContentService)DeleteEvent(c context.Context,u,id pgtype.UUID)error{return s.repo.DeleteEvent(c,u,id)}
func(s *ContentService)CreateInterpretation(c context.Context,u pgtype.UUID,sys string,d time.Time,v string)(repository.Interpretation,error){return s.repo.CreateInterpretation(c,u,sys,d,v)}
func(s *ContentService)ListInterpretations(c context.Context,sys string,d time.Time)([]repository.Interpretation,error){return s.repo.ListInterpretations(c,sys,d)}
