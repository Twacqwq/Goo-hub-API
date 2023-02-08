package seed

import (
	"thub/pkg/console"
	"thub/pkg/database"

	"gorm.io/gorm"
)

var (
	seeders []Seeder

	orderdSeederNames []string
)

type SeederFunc func(*gorm.DB)

type Seeder struct {
	Func SeederFunc
	Name string
}

// Add 添加 Seeder
func Add(name string, fn SeederFunc) {
	seeders = append(seeders, Seeder{
		Name: name,
		Func: fn,
	})
}

func SetRunOrder(names []string) {
	orderdSeederNames = names
}

// GetSeeder 通过名称获取 Seeder 对象
func GetSeeder(name string) Seeder {
	for _, sdr := range seeders {
		if name == sdr.Name {
			return sdr
		}
	}

	return Seeder{}
}

// RunAll 运行所有 Seeder
func RunAll() {
	executed := make(map[string]string)
	for _, name := range orderdSeederNames {
		sdr := GetSeeder(name)
		if len(sdr.Name) > 0 {
			console.Warning("Running Orderd Seeder:" + sdr.Name)
			sdr.Func(database.DB)
			executed[name] = name
		}
	}

	for _, sdr := range seeders {
		if _, ok := executed[sdr.Name]; !ok {
			console.Warning("Running Seeder: " + sdr.Name)
			sdr.Func(database.DB)
		}
	}
}

// RunSeeder 运行单个 Seeder
func RunSeeder(name string) {
	for _, sdr := range seeders {
		if name == sdr.Name {
			sdr.Func(database.DB)
			break
		}
	}
}
