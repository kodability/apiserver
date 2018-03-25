package models

import "time"

// Question defines a table that describes each tryout question.
type Question struct {
	SurrogateModel
	AuditModel
	Level         int
	EstimatedTime int
	Desctiptions  []QuestionDescription `gorm:"foreignKey:QuestionID"`
	Codes         []QuestionCode        `gorm:"foreignKey:QuestionID"`
	Tags          string                `gorm:"type:varchar(50)"`
	Demo          bool
}

// QuestionDescription defines title and description for each locale
type QuestionDescription struct {
	AuditModel
	QuestionID  uint   `gorm:"index;unique_index:question_desc_uq"`
	LocaleID    string `gorm:"type:varchar(10);not null;unique_index:question_desc_uq"`
	Title       string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:varchar(4000)"`
}

// QuestionCode contains initially provided source code and Test code.
type QuestionCode struct {
	AuditModel
	QuestionID uint   `gorm:"index;unique_index:question_code_uq"`
	Lang       string `gorm:"type:varchar(10);not null;index;unique_index:question_code_uq"`
	InitCode   string `gorm:"type:varchar(255)"`
	TestCode   string `gorm:"type:varchar(30000)"`
}

// QuestionTag defines tag
type QuestionTag struct {
	AuditModel
	Tag      string `gorm:"type:varchar(20);not null;unique_index:question_tag_uq"`
	LocaleID string `gorm:"type:varchar(10);not null;unique_index:question_tag_uq"`
	Name     string `gorm:"type:varchar(30)"`
}

type Tryout struct {
	SurrogateModel
	Question   Question `gorm:"foreignKey:QuestionID"`
	QuestionID uint
	Lang       string
	Code       string
	CreatedAt  time.Time
}

type TryoutResult struct {
	Tryout       Tryout `gorm:"foreignKey:TryoutID"`
	TryoutID     uint
	TestCount    int
	ErrorCount   int
	ErrorNames   string `gorm:"type:varchar(1000)"`
	FailureCount int
	FailureNames string `gorm:"type:varchar(1000)"`
	ElapsedTime  float64
	ErrorMsg     string `gorm:"type:varchar(255)"`
	CreatedAt    time.Time
}
