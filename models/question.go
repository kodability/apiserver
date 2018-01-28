package models

import "github.com/jinzhu/gorm"

// Question defines a table that describes each tryout question.
type Question struct {
	gorm.Model
	Level         int
	EstimatedTime int
	Desctiptions  []QuestionDesc
	Codes         []QuestionCode
	Tags          []QuestionTag
	Demo          bool
}

// QuestionDesc defines title and description for each locale
type QuestionDesc struct {
	ID          int
	LocaleID    string `gorm:"type:varchar(10);not null;unique"`
	Title       string `gorm:"type:varchar(100);not null;unique"`
	Description string
}

// QuestionCode contains initially provided source code and Test code.
type QuestionCode struct {
	ID       int
	Lang     string `gorm:"type:varchar(10);not null;index"`
	InitCode string
	TestCode string
}

// QuestionTag defines tag for each locale
type QuestionTag struct {
	ID       int
	Tag      string `gorm:"type:varchar(20);not null;unique"`
	LocaleID string `gorm:"type:varchar(10);not null;unique"`
	Name     string
}
