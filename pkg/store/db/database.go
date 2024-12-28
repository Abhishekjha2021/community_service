package db

import (
	"fmt"

	"github.com/Abhishekjha321/community_service/storage/db/postgres"
	"github.com/Abhishekjha321/community_service/pkg/config"
	"github.com/Abhishekjha321/community_service/pkg/store/db/model"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

const (
	SortDirectionASC  = "ASC"
	SortDirectionDESC = "DESC"
)

type Store struct {
	MasterDB *gorm.DB
	SlaveDB  *gorm.DB
}

func NewPostgresStorage() (*Store, error) {
	postgres.InitialisePostgres(config.Config.PostgresMaster, config.Config.PostgresSlave)
	s := new(Store)
	s.MasterDB = postgres.GetMasterEngine()
	s.SlaveDB = postgres.GetSlaveEngine()
	err := s.MasterDB.AutoMigrate(
		model.Post{},
		model.UserDetails{},
		model.UserActions{},
		model.UserStatus{},
		model.MasterReport{},
		model.Forum{},
		model.ForumEventLink{},
		model.Reports{},
	)
	if err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	fmt.Println("Connected to db")
	return s, nil
}
