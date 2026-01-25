package state

import "weather-subscriptions/internal/db/models"

type UsersState interface {
	GetUser(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	SaveUser(user *models.User) error
	RemoveUser(user *models.User) error
}

func (s *State) GetUser(id string) (*models.User, error) {
	user, ok := s.user[id]
	if !ok {
		foundUser, err := s.resolver.UserByID(id)
		if err != nil {
			return nil, err
		}
		user = foundUser
	}
	s.user[id] = user

	return user, nil
}

func (s *State) GetUserByEmail(email string) (*models.User, error) {
	user, ok := s.user[email]
	if !ok {
		foundUser, err := s.resolver.UserByEmail(email)
		if err != nil {
			return nil, err
		}
		user = foundUser
	}

	s.user[email] = user
	return user, nil
}

func (s *State) SaveUser(user *models.User) error {
	err := s.resolver.Save(user)
	if err != nil {
		return err
	}
	s.user[user.Email] = user
	s.user[user.ID] = user

	return nil
}

func (s *State) RemoveUser(user *models.User) error {
	err := s.resolver.Remove(user)
	if err != nil {
		return err
	}

	delete(s.user, user.ID)
	delete(s.subscriptions, user.ID)

	return nil
}
