package resolvers

import (
	"gorm.io/gorm"
	"weather-subscriptions/internal/db/models"
)

// TODO: add insert/update functional

type Resolver interface {
	UserByID(id string) (*models.User, error)
	UserByEmail(email string) (*models.User, error)
	Token(token string) (*models.Token, error)
	UnsubToken(userID string) (*models.Token, error)
	UserToken(userID, tokenType string) (*models.Token, error)
	Subscription(userID string) (*models.Subscription, error)
	Subscriptions(subscriptionType models.SubscriptionType) ([]*models.Subscription, error)
	City(name string) (*models.City, error)
	CityByID(id string) (*models.City, error)
	Weather(CityID string) (*models.Weather, error)
	WeatherByCityID(cityID string) (*models.Weather, error)
	Save(model any) error
	Remove(model any) error
}

type DBResolver struct {
	db *gorm.DB
}

func New(db *gorm.DB) Resolver {
	return &DBResolver{
		db: db,
	}
}

func (r *DBResolver) UserByID(id string) (user *models.User, err error) {
	return user, r.db.First(&user, "id = ?", id).Error
}

func (r *DBResolver) UserByEmail(email string) (user *models.User, err error) {
	return user, r.db.First(&user, "email = ?", email).Error
}

func (r *DBResolver) Token(token string) (t *models.Token, err error) {
	return t, r.db.First(&t, "token = ?", token).Error
}

func (r *DBResolver) UnsubToken(userID string) (t *models.Token, err error) {
	return t, r.db.First(&t, "user_id = ? AND type = ?", userID, models.Unsub).Error
}

func (r *DBResolver) UserToken(userID, tokenType string) (token *models.Token, err error) {
	return token, r.db.First(&token, "user_id = ? AND token_type = ?", userID, tokenType).Error
}

func (r *DBResolver) Subscription(userID string) (subscription *models.Subscription, err error) {
	return subscription, r.db.First(&subscription, "user_id = ?", userID).Error
}

func (r *DBResolver) Subscriptions(subscriptionType models.SubscriptionType) (subscriptions []*models.Subscription, err error) {
	return subscriptions, r.db.Preload("User").Where("frequency = ?", subscriptionType).Find(&subscriptions).Error
}

func (r *DBResolver) CityByID(id string) (city *models.City, err error) {
	return city, r.db.First(&city, "id = ?", id).Error
}

func (r *DBResolver) City(name string) (city *models.City, err error) {
	return city, r.db.First(&city, "name ILIKE ?", name).Error
}

func (r *DBResolver) Weather(CityID string) (weather *models.Weather, err error) {
	return weather, r.db.
		Order("time DESC").
		First(&weather, "city_id = ?", CityID).
		Error
}

func (r *DBResolver) WeatherByCityID(cityID string) (weather *models.Weather, err error) {
	return weather, r.db.Order("time desc").First(&weather, "city_id = ?", cityID).Error
}

func (r *DBResolver) Save(model any) error {
	return r.db.Save(model).Error
}

func (r *DBResolver) Remove(model any) error { return r.db.Delete(model).Error }
