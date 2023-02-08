package seeders

import "thub/pkg/seed"

func Initialize() {
	seed.SetRunOrder([]string{
		"SeedUsersTable",
	})
}
