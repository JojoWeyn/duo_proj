package entity

type Rank struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}
