package service

//type Repository interface {
//	Create(g *genreEntity.Genre) (*genreEntity.Genre, error)
//	GetByID(id int64) (*genreEntity.Genre, error)
//	GetAll() ([]*genreEntity.Genre, error)
//	DeleteByID(id int64) error
//}
//
//type Service struct {
//	DB   *sql.DB
//	Log  *slog.Logger
//	Repo Repository
//}
//
//func New(db *sql.DB, log *slog.Logger, repo Repository) *Service {
//	return &Service{
//		DB:   db,
//		Log:  log,
//		Repo: repo,
//	}
//}
//
//func (s *Service) Create(g *genreEntity.Genre) (*genreEntity.Genre, error) {
//	const op = "genre.service.CreateBook"
//
//	genre, err := s.Repo.Create(g)
//	if err != nil {
//		s.Log.Error("failed to create genre", "op", op, "error", err)
//		return nil, err
//	}
//
//	s.Log.Info(
//		"create genre",
//		"id", genre.ID,
//		"name", genre.Name,
//	)
//
//	return genre, nil
//}
