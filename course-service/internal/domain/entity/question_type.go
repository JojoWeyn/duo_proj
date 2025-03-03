package entity

type QuestionType struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
}
