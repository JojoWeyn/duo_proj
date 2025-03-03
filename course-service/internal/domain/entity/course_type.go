package entity

type CourseType struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
}
