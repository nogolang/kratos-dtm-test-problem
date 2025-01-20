package gormUtils

import "gorm.io/gorm"

func Pagination(page int, pageSize int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		//这里要减1，因为第1页数据从0开始
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
