package models

type FooModel struct {
	ID uint64 `gorm:"primary_key" json:"id"`
	V  string `json:"v"`
}

func (*FooModel) TableName() string {
	return "foo"
}