package models

type Person struct {
	ID        uint   `gorm:"primaryKey; not null" json:"id,omitempty" yaml:"id,omitempty"`
	Name      string `gorm:"type:varchar(200); not null" json:"name" yaml:"name"`
	Surname   string `gorm:"type:varchar(200); not null" json:"surname" yaml:"surname"`
	Age       int    `gorm:"age; not null" json:"age" yaml:"age"`
	Email     string `gorm:"type:varchar(200); uniqueIndex; not null" json:"email" yaml:"email"`
	Telephone string `gorm:"type:varchar(200); not null" json:"telephone" yaml:"telephone"`
}
