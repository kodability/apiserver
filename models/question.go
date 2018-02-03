package models

import "github.com/jinzhu/gorm"

// Question defines a table that describes each tryout question.
type Question struct {
	gorm.Model
	Level         int
	EstimatedTime int
	Desctiptions  []QuestionDesc `gorm:"ForeignKey:QuestionID"`
	Codes         []QuestionCode `gorm:"ForeignKey:QuestionID"`
	Tags          string         `gorm:"type:varchar(50)"`
	Demo          bool
}

// QuestionDesc defines title and description for each locale
type QuestionDesc struct {
	ID          int
	QuestionID  int    `gorm:"index;unique"`
	LocaleID    string `gorm:"type:varchar(10);not null;unique"`
	Title       string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:varchar(4000)"`
}

// QuestionCode contains initially provided source code and Test code.
type QuestionCode struct {
	ID         int
	QuestionID int    `gorm:"index"`
	Lang       string `gorm:"type:varchar(10);not null;index"`
	InitCode   string `gorm:"type:varchar(255)"`
	TestCode   string `gorm:"type:varchar(30000)"`
}

// QuestionTag defines tag
type QuestionTag struct {
	ID       int
	Tag      string `gorm:"type:varchar(20);not null;unique"`
	LocaleID string `gorm:"type:varchar(10);not null;unique"`
	Name     string `gorm:"type:varchar(30)"`
}
