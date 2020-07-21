package loadconfigs

import (
	"github.com/olebedev/config"
	"io/ioutil"
)

// load database config
func LoadConfig() (database, db_user, db_pass, db_addr, table string) {

	file, err := ioutil.ReadFile("/opt/app-root/config.yaml")
	if err != nil {
		panic(err)
	}
	yamlString := string(file)

	cfg, err := config.ParseYaml(yamlString)
	if err != nil {
		panic(err)
	}

	database, err = cfg.String("database.name")
	if err != nil {
		panic(err)
	}

	db_user, err = cfg.String("database.user")
	if err != nil {
		panic(err)
	}

	db_pass, err = cfg.String("database.password")
	if err != nil {
		panic(err)
	}
	db_addr, err = cfg.String("database.address")
	if err != nil {
		panic(err)
	}

	table, err = cfg.String("table.name")
	if err != nil {
		panic(err)
	}

	return database, db_user, db_pass, db_addr, table
}
