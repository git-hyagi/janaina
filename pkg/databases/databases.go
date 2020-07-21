package databases

import (
	"database/sql"
	"github.com/git-hyagi/janaina/pkg/types"
	_ "github.com/go-sql-driver/mysql"
)

type DbConnection struct {
	Conn *sql.DB
}

// database connection function
func Connect(database, user, password, address string) (db *sql.DB, err error) {
	dsn := user + ":" + password + "@tcp(" + address + ")/" + database
	db, err = sql.Open("mysql", dsn)

	return db, err
}

// function to find a person by username in the database
func (db *DbConnection) FindPersonByUsername(username string) (types.Person, error) {
	rows, err := db.Conn.Query("SELECT * from pessoa WHERE username like '%" + username + "%'")
	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	rows.Next()
	var user types.Person
	err = rows.Scan(&user.IdPerson, &user.Contact, &user.Username, &user.Name, &user.Surname, &user.Role, &user.Password)
	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}
	rows.Close()
	return user, nil
}

// function to remove a user from database
func (db *DbConnection) RemoveUserByName(name, table string) (err error) {
	_, err = db.Conn.Query("DELETE FROM " + table + " WHERE name='" + name + "'")
	return err
}

// function to update an address from an user
func (db *DbConnection) UpdateUserAddress(name, address, table string) (err error) {
	_, err = db.Conn.Query("UPDATE " + table + " SET address='" + address + "' WHERE name='" + name + "'")
	return err
}

// function to update an crm from an user
func (db *DbConnection) UpdateUserCrm(name, crm, table string) (err error) {
	_, err = db.Conn.Query("UPDATE " + table + " SET crm='" + crm + "' WHERE name='" + name + "'")
	return err
}

// function to update a name from an user
func (db *DbConnection) UpdateUserName(name, new_name, table string) (err error) {
	_, err = db.Conn.Query("UPDATE " + table + " SET name='" + new_name + "' WHERE name='" + name + "'")
	return err
}

// add a new contact to contato table
func (db *DbConnection) AddContact(telephone, email, facebook string) (err error) {
	_, err = db.Conn.Query("INSERT INTO contato VALUES(NULL,'" + telephone + "','" + email + "','" + facebook + "' )")
	return err
}

// add a new person to pessoa table
func (db *DbConnection) AddPerson(idContact, username, name, surname, userType, password string) (err error) {
	_, err = db.Conn.Query("INSERT INTO pessoa VALUES(NULL,'" + idContact + "','" + username + "','" + name + "','" + surname + "','" + userType + "','" + password + "' )")
	return err
}

// add a new person to tipo table
func (db *DbConnection) AddType(idPerson, idMedic, idAdmin, idPatient string) (err error) {
	if idMedic != "NULL" {
		_, err = db.Conn.Query("INSERT INTO tipo VALUES('" + idPerson + "','" + idMedic + "',NULL,NULL,'medico')")
	} else if idAdmin != "NULL" {
		_, err = db.Conn.Query("INSERT INTO tipo VALUES('" + idPerson + "',NULL,'" + idAdmin + "',NULL,'admin')")
	} else if idPatient != "NULL" {
		_, err = db.Conn.Query("INSERT INTO tipo VALUES('" + idPerson + "',NULL,NULL,'" + idPatient + "','paciente')")
	}
	return err
}

// add a new person to admin table
func (db *DbConnection) AddAdmin(idPerson string) (err error) {
	_, err = db.Conn.Query("INSERT INTO admin VALUES(NULL,'" + idPerson + "')")
	return err
}

// add a new person to medico table
func (db *DbConnection) AddMedic(CRM, idPerson string) (err error) {
	_, err = db.Conn.Query("INSERT INTO medico VALUES(NULL,'" + CRM + "','" + idPerson + "' )")
	return err
}

// add a new patient to paciente table
func (db *DbConnection) AddPatient(idPerson string) (err error) {
	_, err = db.Conn.Query("INSERT INTO paciente VALUES(NULL,'" + idPerson + "' )")
	return err
}

// create a new admin
func (db *DbConnection) CreateAdmin(person types.Person, contact types.Contact) error {

	err := db.AddContact(contact.Phone, contact.Email, contact.Facebook)
	if err != nil {
		return err
	}

	// this is not very efficient
	// we are creating a contact and, just after that, looking for it just to retrieve its id
	// maybe would be better if the db.AddContact could return the idContato just after its creation
	idContato, err := db.findContactByEmail(contact.Email)
	if err != nil {
		return err
	}

	err = db.AddPerson(idContato, person.Username, person.Name, person.Surname, "admin", person.Password)
	if err != nil {
		return err
	}

	// the same issue as with contato
	idPessoa, err := db.findPersonByContact(idContato)
	if err != nil {
		return err
	}

	err = db.AddAdmin(idPessoa)
	if err != nil {
		return err
	}

	// the same issue as with contato and pessoa
	idAdmin, err := db.findAdmin(idPessoa)
	if err != nil {
		return err
	}

	err = db.AddType(idPessoa, "NULL", idAdmin, "NULL")
	if err != nil {
		return err
	}

	return nil
}

// create a new patient
func (db *DbConnection) CreatePatient(person types.Person, contact types.Contact) error {

	err := db.AddContact(contact.Phone, contact.Email, contact.Facebook)
	if err != nil {
		return err
	}

	// this is not very efficient
	// we are creating a contact and, just after that, looking for it just to retrieve its id
	// maybe would be better if the db.AddContact could return the idContato just after its creation
	idContato, err := db.findContactByEmail(contact.Email)
	if err != nil {
		return err
	}

	err = db.AddPerson(idContato, person.Username, person.Name, person.Surname, "paciente", person.Password)
	if err != nil {
		return err
	}

	// the same issue as with contato
	idPessoa, err := db.findPersonByContact(idContato)
	if err != nil {
		return err
	}

	err = db.AddPatient(idPessoa)
	if err != nil {
		return err
	}

	// the same issue as with contato and pessoa
	idPatient, err := db.findPatient(idPessoa)
	if err != nil {
		return err
	}

	err = db.AddType(idPessoa, "NULL", "NULL", idPatient)
	if err != nil {
		return err
	}

	return nil
}

// create a new medic
func (db *DbConnection) CreateMedic(person types.Person, contact types.Contact, CRM string) error {

	err := db.AddContact(contact.Phone, contact.Email, contact.Facebook)
	if err != nil {
		return err
	}

	// this is not very efficient
	// we are creating a contact and, just after that, looking for it just to retrieve its id
	// maybe would be better if the db.AddContact could return the idContato just after its creation
	idContato, err := db.findContactByEmail(contact.Email)
	if err != nil {
		return err
	}

	err = db.AddPerson(idContato, person.Username, person.Name, person.Surname, "medico", person.Password)
	if err != nil {
		return err
	}

	// the same issue as with contato
	idPessoa, err := db.findPersonByContact(idContato)
	if err != nil {
		return err
	}

	err = db.AddMedic(CRM, idPessoa)
	if err != nil {
		return err
	}

	// the same issue as with contato and pessoa
	idMedic, err := db.findMedic(idPessoa)
	if err != nil {
		return err
	}

	err = db.AddType(idPessoa, idMedic, "NULL", "NULL")
	if err != nil {
		return err
	}

	return nil
}

// find a contact by email
// returns idContato
func (db *DbConnection) findContactByEmail(email string) (string, error) {
	rows, err := db.Conn.Query("select idContato from contato where email = '" + email + "'")

	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	var idContato string
	rows.Next()
	err = rows.Scan(&idContato)
	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	return idContato, nil
}

// find a person by contato
// returns the idPessoa
func (db *DbConnection) findPersonByContact(idContact string) (string, error) {
	rows, err := db.Conn.Query("select idPessoa from pessoa where contato = '" + idContact + "'")

	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	var idPessoa string
	rows.Next()
	err = rows.Scan(&idPessoa)
	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	return idPessoa, nil
}

// find an admin by idPessoa
// returns idAdmin
func (db *DbConnection) findAdmin(idPerson string) (string, error) {
	rows, err := db.Conn.Query("select idAdmin from admin where idPessoa = '" + idPerson + "'")

	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	var idAdmin string
	rows.Next()
	err = rows.Scan(&idAdmin)
	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	return idAdmin, nil
}

// find medic by idPessoa
// returns idMedico
func (db *DbConnection) findMedic(idPerson string) (string, error) {
	rows, err := db.Conn.Query("select idMedico from medico where idPessoa = '" + idPerson + "'")

	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	var idMedic string
	rows.Next()
	err = rows.Scan(&idMedic)
	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	return idMedic, nil
}

// find an patient by idPessoa
// returns idPaciente
func (db *DbConnection) findPatient(idPerson string) (string, error) {
	rows, err := db.Conn.Query("select idPaciente from paciente where idPessoa = '" + idPerson + "'")

	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	var idPessoa string
	rows.Next()
	err = rows.Scan(&idPessoa)
	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	return idPessoa, nil
}

// find a medic by CRM
// returns a Medic{idMedic, idCRM, idPerson}
func (db *DbConnection) FindMedicByCRM(CRM string) (types.Medic, error) {
	rows, err := db.Conn.Query("select * from medico where CRM = '" + CRM + "'")

	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	var medic types.Medic
	rows.Next()
	err = rows.Scan(&medic.IdMedic, &medic.CRM, &medic.IdPerson)
	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	return medic, nil
}

// find a medic by Name+Surname
// returns a Medic{idMedic, idCRM, idPerson}
func (db *DbConnection) FindMedicByName(name, surname string) (types.Medic, error) {
	rows, err := db.Conn.Query("select medico.idMedico,medico.CRM,medico.idPessoa from medico join pessoa on medico.idPessoa = pessoa.idPessoa where pessoa.nome = '" + name + "' and pessoa.sobrenome = '" + surname + "'")

	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	var medic types.Medic
	rows.Next()
	err = rows.Scan(&medic.IdMedic, &medic.CRM, &medic.IdPerson)
	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	return medic, nil
}

// find an admin by Name+Surname
// returns an Admin{idAdmin,idPerson}
func (db *DbConnection) FindAdminByName(name, surname string) (types.Admin, error) {
	rows, err := db.Conn.Query("select admin.idAdmin,admin.idPessoa from admin join pessoa on admin.idPessoa = pessoa.idPessoa where pessoa.nome = '" + name + "' and pessoa.sobrenome = '" + surname + "'")

	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	var admin types.Admin
	rows.Next()
	err = rows.Scan(&admin.IdAdmin, &admin.IdPerson)
	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	return admin, nil
}

// find an patient by Name+Surname
// returns an Patient{idPatient,idPerson}
func (db *DbConnection) FindPatientByName(name, surname string) (types.Patient, error) {
	rows, err := db.Conn.Query("select paciente.idPaciente,paciente.idPessoa from paciente join pessoa on paciente.idPessoa = pessoa.idPessoa where pessoa.nome = '" + name + "' and pessoa.sobrenome = '" + surname + "'")

	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	var patient types.Patient
	rows.Next()
	err = rows.Scan(&patient.IdPatient, &patient.IdPerson)
	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	return patient, nil
}

// find person by idMedic
func (db *DbConnection) FindMedicById(id string) (types.Person, types.Contact, error) {
	rows, err := db.Conn.Query("select pessoa.idPessoa,pessoa.contato,pessoa.username,pessoa.nome,pessoa.sobrenome,pessoa.tipo,pessoa.senha,contato.idContato,contato.telefone,contato.email,contato.facebook from medico join pessoa join contato on medico.idPessoa = pessoa.idPessoa and pessoa.contato = contato.idContato where idMedico = '" + id + "'")

	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	var person types.Person
	var contact types.Contact
	rows.Next()
	err = rows.Scan(&person.IdPerson, &person.Contact, &person.Username, &person.Name, &person.Surname, &person.Role, &person.Password, &contact.IdContact, &contact.Phone, &contact.Email, &contact.Facebook)
	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	return person, contact, nil
}

// find person by idPatient
func (db *DbConnection) FindPatientById(id string) (types.Person, types.Contact, error) {
	rows, err := db.Conn.Query("select pessoa.idPessoa,pessoa.contato,pessoa.username,pessoa.nome,pessoa.sobrenome,pessoa.tipo,pessoa.senha,contato.idContato,contato.telefone,contato.email,contato.facebook from paciente join pessoa join contato on paciente.idPessoa = pessoa.idPessoa and pessoa.contato = contato.idContato where idPaciente = '" + id + "'")

	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	var person types.Person
	var contact types.Contact
	rows.Next()
	err = rows.Scan(&person.IdPerson, &person.Contact, &person.Username, &person.Name, &person.Surname, &person.Role, &person.Password, &contact.IdContact, &contact.Phone, &contact.Email, &contact.Facebook)
	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	return person, contact, nil
}

// find person by idAdmin
func (db *DbConnection) FindAdminById(id string) (types.Person, types.Contact, error) {
	rows, err := db.Conn.Query("select pessoa.idPessoa,pessoa.contato,pessoa.username,pessoa.nome,pessoa.sobrenome,pessoa.tipo,pessoa.senha,contato.idContato,contato.telefone,contato.email,contato.facebook from admin join pessoa join contato on admin.idPessoa = pessoa.idPessoa and pessoa.contato = contato.idContato where idAdmin = '" + id + "'")

	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	var person types.Person
	var contact types.Contact
	rows.Next()
	err = rows.Scan(&person.IdPerson, &person.Contact, &person.Username, &person.Name, &person.Surname, &person.Role, &person.Password, &contact.IdContact, &contact.Phone, &contact.Email, &contact.Facebook)
	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	return person, contact, nil
}

// find a person by username
// returns idPessoa,tipo, contato
func (db *DbConnection) findPersonByUsername(username string) (idPerson, tipo, contato string, err error) {
	rows, err := db.Conn.Query("select pessoa.idPessoa,pessoa.tipo,pessoa.contato from tipo join pessoa on pessoa.idPessoa = tipo.idPessoa where pessoa.username='" + username + "'")

	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	rows.Next()
	err = rows.Scan(&idPerson, &tipo, &contato)
	if err != nil {
		rows.Close()
		panic(err.Error())
	}

	rows.Close()
	return idPerson, tipo, contato, nil
}

// find a contact by idPessoa
// returns idContato
func (db *DbConnection) findContactByIdPerson(idPerson string) (idContact string, err error) {
	rows, err := db.Conn.Query("select contato from pessoa where idPessoa = '" + idPerson + "'")

	if err != nil {
		db.Conn.Close()
		panic(err.Error())
	}

	rows.Next()
	err = rows.Scan(&idContact)
	if err != nil {
		rows.Close()
		panic(err.Error())
	}
	rows.Close()
	return idContact, nil

}

// remove a person by its username
func (db *DbConnection) RemovePersonByUsername(username string) error {

	// before removing it, we will need to get the person id
	idPerson, role, _, err := db.findPersonByUsername(username)
	if err != nil {
		panic(err.Error())
	}

	// and also the idContact
	idContact, err := db.findContactByIdPerson(idPerson)
	if err != nil {
		panic(err.Error())
	}

	tx, _ := db.Conn.Begin()
	_, err = tx.Exec("delete from tipo where idPessoa = '" + idPerson + "'")
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("delete from " + role + " where idPessoa = '" + idPerson + "'")
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("delete from pessoa where idPessoa = '" + idPerson + "'")
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("delete from contato where idContato = '" + idContact + "'")
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// update email (table contato) by username
func (db *DbConnection) UpdateEmail(username, newEmail string) error {

	_, _, idContact, err := db.findPersonByUsername(username)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec("update contato set email = '" + newEmail + "' where idContato = '" + idContact + "'")
	if err != nil {
		return err
	}

	return nil
}

// update telefone (table contato) by username
func (db *DbConnection) UpdatePhone(username, newPhone string) error {

	_, _, idContact, err := db.findPersonByUsername(username)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec("update contato set telefone = '" + newPhone + "' where idContato = '" + idContact + "'")
	if err != nil {
		return err
	}

	return nil
}

// update facebook (table contato) by username
func (db *DbConnection) UpdateFacebook(username, newFacebook string) error {

	_, _, idContact, err := db.findPersonByUsername(username)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec("update contato set facebook = '" + newFacebook + "' where idContato = '" + idContact + "'")
	if err != nil {
		return err
	}

	return nil
}

// update name from table pessoa
func (db *DbConnection) UpdateName(username, newName string) error {

	idPerson, _, _, err := db.findPersonByUsername(username)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec("update pessoa set nome = '" + newName + "' where idPessoa = '" + idPerson + "'")
	if err != nil {
		return err
	}

	return nil
}

// update surname from table pessoa
func (db *DbConnection) UpdateSurname(username, newName string) error {

	idPerson, _, _, err := db.findPersonByUsername(username)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec("update pessoa set sobrenome = '" + newName + "' where idPessoa = '" + idPerson + "'")
	if err != nil {
		return err
	}

	return nil
}

// update password from table pessoa
func (db *DbConnection) UpdatePassword(username, newPassword string) error {

	idPerson, _, _, err := db.findPersonByUsername(username)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec("update pessoa set senha = '" + newPassword + "' where idPessoa = '" + idPerson + "'")
	if err != nil {
		return err
	}

	return nil
}
