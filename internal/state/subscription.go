package state

import "weather-subscriptions/internal/db/models"

type SubscriptionsState interface {
	GetSubscription(userID string) (*models.Subscription, error)
	GetSubscriptions(subscriptionType models.SubscriptionType) ([]*models.Subscription, error)
	SaveSubscription(subscription *models.Subscription) error
	RemoveSubscription(subscription *models.Subscription) error
}

func (s *State) GetSubscription(userID string) (*models.Subscription, error) {
	subscription, ok := s.subscriptions[userID]
	if !ok {
		foundSubscription, err := s.resolver.Subscription(userID)
		if err != nil {
			return nil, err
		}
		subscription = foundSubscription
	}
	s.subscriptions[userID] = subscription

	return subscription, nil
}

func (s *State) GetSubscriptions(subscriptionType models.SubscriptionType) ([]*models.Subscription, error) {
	subscriptions, err := s.resolver.Subscriptions(subscriptionType)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (s *State) SaveSubscription(subscription *models.Subscription) error {
	err := s.resolver.Save(subscription)
	if err != nil {
		return err
	}
	s.subscriptions[subscription.UserID] = subscription

	return nil
}

func (s *State) RemoveSubscription(subscription *models.Subscription) error {
	err := s.resolver.Remove(subscription)
	if err != nil {
		return err
	}

	delete(s.subscriptions, subscription.UserID)
	delete(s.user, subscription.UserID)

	return nil
}
