package postgres

import (
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	masterEngine            *gorm.DB
	slaveEngine             *gorm.DB
	ErrFailedToConnectToSQL = "Failed to connect to postgres %v\n"
)

// postgres struct
type PGMaster struct {
	Hostname     string
	Username     string
	Password     string
	MaxOpenConns int
	MaxIdleConns int
	Schema       string
}

type PGSlave struct {
	Hostname     string
	Username     string
	Password     string
	MaxOpenConns int
	MaxIdleConns int
	Schema       string
}

func InitialisePostgres(masterConfig *PGMaster, slaveConfig *PGSlave) {
	postgresMasterConfigData := masterConfig
	masterDsn := fmt.Sprintf("postgresql://%s:%s@%s/%s",
		postgresMasterConfigData.Username,
		postgresMasterConfigData.Password,
		postgresMasterConfigData.Hostname,
		postgresMasterConfigData.Schema)

	var err error

	masterEngine, err = gorm.Open(postgres.Open(masterDsn), &gorm.Config{})
	if err != nil {
		err = fmt.Errorf(ErrFailedToConnectToSQL, err)
		panic(err.Error())
	}

	if err := masterEngine.Use(otelgorm.NewPlugin()); err != nil {
		panic(err)
	}

	if err != nil {
		err = fmt.Errorf(ErrFailedToConnectToSQL, err)
		panic(err.Error())
	}

	postgresSlaveConfigData := slaveConfig
	slaveDsn := fmt.Sprintf("postgresql://%s:%s@%s/%s",
		postgresSlaveConfigData.Username,
		postgresSlaveConfigData.Password,
		postgresSlaveConfigData.Hostname,
		postgresSlaveConfigData.Schema)

	slaveEngine, err = gorm.Open(postgres.Open(slaveDsn), &gorm.Config{})
	if err != nil {
		err = fmt.Errorf(ErrFailedToConnectToSQL, err)
		panic(err.Error())
	}

	if err := slaveEngine.Use(otelgorm.NewPlugin()); err != nil {
		panic(err)
	}

}

func GetMasterEngine() *gorm.DB {
	return masterEngine
}

func GetSlaveEngine() *gorm.DB {
	return slaveEngine
}
