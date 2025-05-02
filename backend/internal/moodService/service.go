package moodservice

type MoodService interface {
	GetMoods(user_id string) ([]Mood, error)
	CreateMood(user_id string, score int, emotions string, description *string) (Mood, error)
	UpdateMood(m Mood) error
	DeleteMood(id string) error
}

type moodService struct {
	repo MoodRepository
}

func (s *moodService) CreateMood(user_id string, score int, emotions string, description *string) (Mood, error) {

	newMood := Mood{
		Score:       score,
		Emotions:    emotions,
		Description: description,
	}

	if err := s.repo.CreateMood(&newMood); err != nil {
		return Mood{}, err
	}
	return newMood, nil
}

func (s *moodService) DeleteMood(id string) error {
	return s.repo.DeleteMood(id)
}

func (s *moodService) GetMoods(user_id string) ([]Mood, error) {
	return s.repo.GetMoods(user_id)
}

func (s *moodService) UpdateMood(m Mood) error {
	return s.repo.UpdateMood(&m)
}

func NewService(repo MoodRepository) MoodService {
	return &moodService{repo: repo}
}
