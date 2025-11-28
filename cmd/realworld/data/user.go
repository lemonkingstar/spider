package data

import "gorm.io/gorm"

type User struct {
	gorm.Model
}

type UserStorage struct {
	Data
}

func (p *UserStorage) Get(id int) (*User, error) {
	result := new(User)
	if err := p.db().Where("id = ?", id).First(result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
