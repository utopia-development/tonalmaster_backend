package repository
import ("context"; "time"; "github.com/jackc/pgx/v5/pgtype")
type Event struct { ID pgtype.UUID; UserID pgtype.UUID; CalendarSystemID string; TargetDate time.Time; Title string; Description *string; CreatedAt time.Time; UpdatedAt time.Time }
type Interpretation struct { ID pgtype.UUID; UserID pgtype.UUID; CalendarSystemID string; TargetDate time.Time; Content string; CreatedAt time.Time; UpdatedAt time.Time }
type ContentRepository interface {
 CreateEvent(context.Context, pgtype.UUID, string, time.Time, string, *string) (Event,error)
 ListEvents(context.Context, pgtype.UUID, time.Time, time.Time) ([]Event,error)
 DeleteEvent(context.Context, pgtype.UUID, pgtype.UUID) error
 CreateInterpretation(context.Context, pgtype.UUID, string, time.Time, string) (Interpretation,error)
 ListInterpretations(context.Context, string, time.Time) ([]Interpretation,error)
}
