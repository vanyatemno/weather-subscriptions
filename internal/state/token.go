package state

import "weather-subscriptions/internal/db/models"

type TokensState interface {
	GetToken(tokens string) (*models.Token, error)
	GetTokenByUserIDAndType(userID string, tokenType string) (*models.Token, error)
	SaveToken(token *models.Token) error
	RemoveToken(token *models.Token) error
}

func (s *State) GetToken(token string) (*models.Token, error) {
	userToken, ok := s.tokens[token]
	if !ok {
		foundToken, err := s.resolver.Token(token)
		if err != nil {
			return nil, err
		}
		userToken = foundToken
	}
	s.tokens[token] = userToken

	return userToken, nil
}

func (s *State) GetTokenByUserIDAndType(userID string, tokenType models.TokenType) (*models.Token, error) {
	token, err := s.resolver.UserToken(userID, tokenType)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *State) SaveToken(token *models.Token) error {
	err := s.resolver.Save(token)
	if err != nil {
		return err
	}
	s.tokens[token.Token] = token

	return nil
}

func (s *State) RemoveToken(token *models.Token) error {
	err := s.resolver.Remove(token)
	if err != nil {
		return err
	}

	delete(s.tokens, token.Token)

	return nil
}
