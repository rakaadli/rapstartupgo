package user

import "gorm.io/gorm"

type Repository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
	FindByID(ID int) (User, error)
	Update(user User) (User, error)
	FindAll() ([]User, error)
}

type repository struct {
	db *gorm.DB
}

func (r *repository) FindByEmail(email string) (User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) FindByID(ID int) (User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Update(user User) (User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) FindAll() ([]User, error) {
	//TODO implement me
	panic("implement me")
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
