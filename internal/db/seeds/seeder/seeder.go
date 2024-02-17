package seeder

import (
	"context"
	"errors"
	"fmt"
	"go_gin/internal/domain/model"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type Seeder struct {
	DB            *gorm.DB
	CustomSeeders []model.CustomSeeds
}

func NewSeeder(DB *gorm.DB, customSeeders []model.CustomSeeds) *Seeder {
	return &Seeder{DB: DB, CustomSeeders: customSeeders}
}

func (s *Seeder) Seeds(ctx context.Context, seedName string, fill int) error {
	if len(s.CustomSeeders) == 0 {
		return errors.New("seeds is nil must be register")
	}
	for _, customSeeder := range s.CustomSeeders {
		if t := reflect.TypeOf(customSeeder); t.Kind() == reflect.Ptr {
			if t.Elem().Name() == seedName {
				customSeeder.Run(ctx, s.DB, fill)
			} else {
				return fmt.Errorf("not found seeders '%s' must be register", seedName)
			}
			return nil
		}
		return errors.New("seeds must be pointer")
	}
	return nil
}

func Execute(seedName string, fill int, seeder *Seeder) error {
	var seedMethods []string
	ctx := context.Background()
	if strings.Contains(seedName, ",") {
		seedMethods = strings.Split(seedName, ",")
	} else if seedName == "" {
		for _, customSeeder := range seeder.CustomSeeders {
			seedType := reflect.TypeOf(customSeeder)
			seedMethods = append(seedMethods, seedType.Elem().Name())
		}
	} else {
		seedMethods = append(seedMethods, seedName)
	}

	for _, method := range seedMethods {
		err := seeder.Seeds(ctx, method, fill)
		if err != nil {
			return err
		}
	}
	return nil
}
