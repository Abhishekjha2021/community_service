package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	masterEngine            = &sql.DB{}
	slaveEngine             = &sql.DB{}
	ErrFailedToConnectToSQL = "Failed to connect to mysql %v\n"
)

// MySQL struct
type MySQLMaster struct {
	Hostname     string
	Username     string
	Password     string
	MaxOpenConns int
	MaxIdleConns int
	Schema       string
}

type MySQLSlave struct {
	Hostname     string
	Username     string
	Password     string
	MaxOpenConns int
	MaxIdleConns int
	Schema       string
}

func InitialiseMysql(masterConfig *MySQLMaster, slaveConfig *MySQLSlave) {
	mysqlMasterConfigData := masterConfig
	masterDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlMasterConfigData.Username, mysqlMasterConfigData.Password, mysqlMasterConfigData.Hostname, mysqlMasterConfigData.Schema)
	var err error
	masterdb, err := gorm.Open(mysql.Open(masterDsn), &gorm.Config{})
	if err != nil {
		err = fmt.Errorf(ErrFailedToConnectToSQL, err)
		panic(err.Error())
	}
	masterEngine, err = masterdb.DB()
	if err != nil {
		err = fmt.Errorf(ErrFailedToConnectToSQL, err)
		panic(err.Error())
	}
	masterEngine.SetMaxOpenConns(mysqlMasterConfigData.MaxOpenConns)
	masterEngine.SetMaxIdleConns(mysqlMasterConfigData.MaxOpenConns)
	masterEngine.SetConnMaxLifetime(-1)

	mysqlSlaveConfigData := slaveConfig
	slaveDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlSlaveConfigData.Username, mysqlSlaveConfigData.Password, mysqlSlaveConfigData.Hostname, mysqlSlaveConfigData.Schema)
	slavedb, err := gorm.Open(mysql.Open(slaveDsn), &gorm.Config{})
	if err != nil {
		err = fmt.Errorf(ErrFailedToConnectToSQL, err)
		panic(err.Error())
	}
	slaveEngine, err = slavedb.DB()
	if err != nil {
		err = fmt.Errorf(ErrFailedToConnectToSQL, err)
		panic(err.Error())
	}
	slaveEngine.SetMaxOpenConns(mysqlMasterConfigData.MaxOpenConns)
	slaveEngine.SetMaxIdleConns(mysqlMasterConfigData.MaxOpenConns)
	slaveEngine.SetConnMaxLifetime(-1)
}

func GetMasterEngine() *sql.DB {
	return masterEngine
}

func GetSlaveEngine() *sql.DB {
	return slaveEngine
}
