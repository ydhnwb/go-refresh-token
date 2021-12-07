package repo

import (
	"log"

	"github.com/ydhnwb/go-refresh-token-example/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	CreateAccount(user entity.User) (*entity.User, error)
	FindUserByEmail(email string) (*entity.User, error)
	FindUserByID(id uint) (*entity.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepoInterface {
	return &userRepo{
		db: db,
	}
}

func (c *userRepo) CreateAccount(user entity.User) (*entity.User, error) {
	user.Password = hashAndSalt([]byte(user.Password))
	err := c.db.Save(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *userRepo) FindUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	result := c.db.Where("email = ?", email).Take(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (c *userRepo) FindUserByID(id uint) (*entity.User, error) {
	var user entity.User
	result := c.db.Where("id = ?", id).Take(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}
	return string(hash)
}
