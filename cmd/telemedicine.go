/*
  This is just a program to test the methods to work with the database.
  [TO-DO] This code should be removed and a test file should be used instead.
*/
package main

import (
	"fmt"
	"github.com/git-hyagi/janaina/pkg/databases"
	config "github.com/git-hyagi/janaina/pkg/loadconfigs"
	"github.com/git-hyagi/janaina/pkg/types"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type contact struct {
	tel      string
	email    string
	facebook string
}

type person struct {
	idContact string
	username  string
	name      string
	surname   string
	userType  string
	password  string
}

type medic struct {
	crm      string
	idPerson string
}

type patient struct {
	idPerson string
}

func main() {

	var database, db_user, db_pass, db_addr, _ = config.LoadConfig()
	connection, err := databases.Connect(database, db_user, db_pass, db_addr)
	db := &databases.DbConnection{Conn: connection}
	defer db.Conn.Close()
	if err != nil {
		panic(err.Error())
	}

	/* BEGIN NEW ADMIN USER */
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	contact1 := contact{"11 1234-5678", "john_doe@email.com", "John Doe"}
	person1 := person{"1", "johnDoe", "John", "Doe", "admin", string(hashedPassword)}
	err = db.AddContact(contact1.tel, contact1.email, contact1.facebook)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Contact " + contact1.email + " added successfuly!")
	}

	err = db.AddPerson(person1.idContact, person1.username, person1.name, person1.surname, person1.userType, person1.password)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Person " + person1.username + " added successfuly")
	}

	err = db.AddAdmin("1")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Person " + person1.username + " added as admin successfuly")
	}

	err = db.AddType("1", "NULL", "1", "NULL")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Type for person " + person1.username + " configured successfuly")
	}
	/* END NEW ADMIN USER */

	/* BEGIN NEW MEDIC USER */
	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	contact2 := contact{"11 2345-6789", "gregory_house@email.com", "Gregory House"}
	person2 := person{"2", "house", "Gregory", "House", "medico", string(hashedPassword)}
	medic := medic{"CRM-123456", "2"}
	err = db.AddContact(contact2.tel, contact2.email, contact2.facebook)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Contact " + contact2.email + " added successfuly!")
	}

	err = db.AddPerson(person2.idContact, person2.username, person2.name, person2.surname, person2.userType, person2.password)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Person " + person2.username + " added successfuly")
	}

	err = db.AddMedic(medic.crm, medic.idPerson)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Person " + person2.username + " added as medic successfuly")
	}

	err = db.AddType("2", "1", "NULL", "NULL")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Type for person " + person2.username + " configured successfuly")
	}
	/* END NEW MEDIC */

	/* BEGIN NEW PATIENT USER */
	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	contact3 := contact{"11 3456-7890", "patch_adams@email.com", "Patch Adams"}
	person3 := person{"3", "padams", "Patch", "Adams", "paciente", string(hashedPassword)}
	patient := patient{"3"}
	err = db.AddContact(contact3.tel, contact3.email, contact3.facebook)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Contact " + contact3.email + " added successfuly!")
	}

	err = db.AddPerson(person3.idContact, person3.username, person3.name, person3.surname, person3.userType, person3.password)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Person " + person3.username + " added successfuly")
	}

	err = db.AddPatient(patient.idPerson)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Person " + person3.username + " added as patient successfuly")
	}

	err = db.AddType("3", "NULL", "NULL", "1")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Type for person " + person3.username + " configured successfuly")
	}
	/* END NEW PATIENT */

	/* FIND MEDIC */
	var medic2 types.Medic
	medic2, err = db.FindMedicByCRM("CRM-123456")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(medic2)

	var medic3 types.Medic
	medic3, err = db.FindMedicByName("Gregory", "House")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(medic3)
	/* END FIND MEDIC */

	/* FIND ADMIN */
	var admin types.Admin
	admin, err = db.FindAdminByName("John", "Doe")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(admin)
	/* END FIND ADMIN */

	/* FIND PATIENT */
	var patient2 types.Patient
	patient2, err = db.FindPatientByName("Patch", "Adams")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(patient2)
	/* END FIND PATIENT */

	var contact types.Contact
	var person4 types.Person
	person4, contact, err = db.FindPatientById("1")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(person4, contact)

	var contact4 types.Contact
	var person5 types.Person
	person5, contact4, err = db.FindMedicById("1")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(person5, contact4)

	var contact5 types.Contact
	var person6 types.Person
	person6, contact5, err = db.FindAdminById("1")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(person6, contact5)

	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	contactN := types.Contact{"NULL", "12 1234-5678", "john_doe_2@email.com", "John Doe 2"}
	personN := types.Person{"NULL", "NULL", "johnDoe2", "John", "Doe 2", "admin", string(hashedPassword)}
	err = db.CreateAdmin(personN, contactN)
	if err != nil {
		fmt.Println(err.Error())
	}

	contactN = types.Contact{"NULL", "13 1234-5678", "john_doe_3@email.com", "John Doe 3"}
	personN = types.Person{"NULL", "NULL", "johnDoe3", "John", "Doe 3", "admin", string(hashedPassword)} //role is defined as admin just as a test
	err = db.CreatePatient(personN, contactN)
	if err != nil {
		fmt.Println(err.Error())
	}

	contactN = types.Contact{"NULL", "14 1234-5678", "john_doe_4@email.com", "John Doe 4"}
	personN = types.Person{"NULL", "NULL", "johnDoe4", "John", "Doe 4", "admin", string(hashedPassword)} //role is defined as admin just as a test
	crm := "CRM-11123445"
	err = db.CreateMedic(personN, contactN, crm)
	if err != nil {
		fmt.Println(err.Error())
	}

	db.RemovePersonByUsername("johnDoe4")
	db.RemovePersonByUsername("johnDoe3")
	db.UpdateEmail("johnDoe2", "johnDoe@domain.com")
	db.UpdatePhone("johnDoe2", "98 7654-4321")
	db.UpdateFacebook("johnDoe2", "jDoe")
	db.UpdateName("johnDoe2", "JOHN")
	db.UpdateSurname("johnDoe2", "DOE 2")

	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("password2"), bcrypt.DefaultCost)
	db.UpdatePassword("johnDoe2", string(hashedPassword))

}
