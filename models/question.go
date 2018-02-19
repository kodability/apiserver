package models

import "github.com/jinzhu/gorm"

// Question defines a table that describes each tryout question.
type Question struct {
	gorm.Model
	Level         int
	EstimatedTime int
	Desctiptions  []QuestionDescription `gorm:"ForeignKey:QuestionID"`
	Codes         []QuestionCode        `gorm:"ForeignKey:QuestionID"`
	Tags          string                `gorm:"type:varchar(50)"`
	Demo          bool
}

// QuestionDescription defines title and description for each locale
type QuestionDescription struct {
	gorm.Model
	QuestionID  uint   `gorm:"index;unique_index:question_desc_uq"`
	LocaleID    string `gorm:"type:varchar(10);not null;unique_index:question_desc_uq"`
	Title       string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:varchar(4000)"`
}

// QuestionCode contains initially provided source code and Test code.
type QuestionCode struct {
	gorm.Model
	QuestionID uint   `gorm:"index;unique_index:question_code_uq"`
	Lang       string `gorm:"type:varchar(10);not null;index;unique_index:question_code_uq"`
	InitCode   string `gorm:"type:varchar(255)"`
	TestCode   string `gorm:"type:varchar(30000)"`
}

// QuestionTag defines tag
type QuestionTag struct {
	gorm.Model
	Tag      string `gorm:"type:varchar(20);not null;unique_index:question_tag_uq"`
	LocaleID string `gorm:"type:varchar(10);not null;unique_index:question_tag_uq"`
	Name     string `gorm:"type:varchar(30)"`
}
