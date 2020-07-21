package databases

import (
	config "github.com/git-hyagi/janaina/pkg/loadconfigs"
	"github.com/git-hyagi/janaina/pkg/types"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

var user = types.User{Name: "Jose Silva", Address: "Rua 1, 123", CRM: "CRM-3294932"}
var database, db_user, db_pass, db_addr, table = config.LoadConfig()
var dsn = db_user + ":" + db_pass + "@tcp(" + db_addr + "/" + database
var connection, _ = Connect(database, db_user, db_pass, db_addr)
var db = &DbConnection{Conn: connection}

var testUser = "John Doe"

// Test connection function
func TestConnect(t *testing.T) {
	_, err := Connect("test", testUser, "", db_addr)
	if err != nil {
		t.Errorf("Failed to connect to the database!!")
	}
}

// Test create user function
func TestAdduser(t *testing.T) {
	err := db.AddUser(&user, table)

	if err != nil {
		t.Errorf("Failed to add user to the database!!")
	}
}

// Test find user
func TestFindUser(t *testing.T) {
	_, err := db.FindUserByName(testUser, table)

	if err != nil {
		t.Errorf("Failed to add user to the database!!")
	}
}

// Test remove user
func TestRemoveUser(t *testing.T) {
	err := db.RemoveUserByName(testUser, table)

	if err != nil {
		t.Errorf("Failed to remove user from database!!")
	}
}
