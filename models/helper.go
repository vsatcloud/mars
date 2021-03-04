package models

func Create(value interface{}) error {
	return db.Model(value).Create(value).Error
}

func Save(value interface{}) error {
	return db.Model(value).Save(value).Error
}

func Get(id uint, dest interface{}) error {
	err := db.First(dest, "id = ?", id).Error
	return err
}
