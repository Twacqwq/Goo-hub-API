package seed

import "gorm.io/gorm"

var (
	seeders []Seeder

	orderdSeederNames []string
)

type SeederFunc func(*gorm.DB)

type Seeder struct {
	Func SeederFunc
	Name string
}

func Add(name string, fn SeederFunc) {
	seeders = append(seeders, Seeder{
		Name: name,
		Func: fn,
	})
}

func SetRunOrder(names []string) {
	orderdSeederNames = names
}
