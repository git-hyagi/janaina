package person

import (
	"log"
	"strconv"

	"github.com/git-hyagi/janaina/pkg/databases"
	config "github.com/git-hyagi/janaina/pkg/loadconfigs"
	"github.com/git-hyagi/janaina/pkg/types"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type Server struct {
}

var database, db_user, db_pass, db_addr, table = config.LoadConfig()
var dsn = db_user + ":" + db_pass + "@tcp(" + db_addr + "/" + database
var connection, _ = databases.Connect(database, db_user, db_pass, db_addr)
var db = &databases.DbConnection{Conn: connection}

// Add a new medic to the database
func (s *Server) CreateMedic(ctx context.Context, medic *Medic) (*Message, error) {
	log.Printf("Adding a new medic to the database: %v", medic)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(medic.Person.GetPassword()), bcrypt.DefaultCost)

	contact := types.Contact{IdContact: "NULL", Phone: medic.Contact.GetPhone(), Email: medic.Contact.GetEmail(), Facebook: medic.Contact.GetFacebook()}
	person := types.Person{IdPerson: "NULL", Contact: "NULL", Username: medic.Person.GetUsername(), Name: medic.Person.GetName(), Surname: medic.Person.GetSurname(), Role: "NULL", Password: string(hashedPassword)}
	crm := medic.GetCrm()

	err := db.CreateMedic(person, contact, crm)
	if err != nil {
		return &Message{Body: "Error creating medic " + medic.Person.GetUsername()}, err
	}
	return &Message{Body: "User " + medic.Person.GetUsername() + " added to the database!"}, nil
}

// Add a new admin to the database
func (s *Server) CreateAdmin(ctx context.Context, admin *Admin) (*Message, error) {
	log.Printf("Adding a new admin to the database: %v", admin)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(admin.Person.GetPassword()), bcrypt.DefaultCost)

	contact := types.Contact{IdContact: "NULL", Phone: admin.Contact.GetPhone(), Email: admin.Contact.GetEmail(), Facebook: admin.Contact.GetFacebook()}
	person := types.Person{IdPerson: "NULL", Contact: "NULL", Username: admin.Person.GetUsername(), Name: admin.Person.GetName(), Surname: admin.Person.GetSurname(), Role: "NULL", Password: string(hashedPassword)}

	err := db.CreateAdmin(person, contact)
	if err != nil {
		return &Message{Body: "Error creating admin " + admin.Person.GetUsername()}, err
	}
	return &Message{Body: "User " + admin.Person.GetUsername() + " added to the database!"}, nil
}

// Add a new patient to the database
func (s *Server) CreatePatient(ctx context.Context, patient *Patient) (*Message, error) {
	log.Printf("Adding a new patient to the database: %v", patient)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(patient.Person.GetPassword()), bcrypt.DefaultCost)

	contact := types.Contact{IdContact: "NULL", Phone: patient.Contact.GetPhone(), Email: patient.Contact.GetEmail(), Facebook: patient.Contact.GetFacebook()}
	person := types.Person{IdPerson: "NULL", Contact: "NULL", Username: patient.Person.GetUsername(), Name: patient.Person.GetName(), Surname: patient.Person.GetSurname(), Role: "NULL", Password: string(hashedPassword)}

	err := db.CreatePatient(person, contact)
	if err != nil {
		return &Message{Body: "Error creating admin " + patient.Person.GetUsername()}, err
	}
	return &Message{Body: "User " + patient.Person.GetUsername() + " added to the database!"}, nil
}

// Remove person by username
func (s *Server) RemovePerson(ctx context.Context, username *PrimaryKey) (*Message, error) {
	log.Printf("Removing a person from the database: %v", username)

	err := db.RemovePersonByUsername(username.GetUsername())
	if err != nil {
		return &Message{Body: "Error removing user " + username.GetUsername()}, err
	}
	return &Message{Body: "User " + username.GetUsername() + " removed from the database!"}, nil
}

// Find person by username
func (s *Server) FindPerson(ctx context.Context, username *PrimaryKey) (*Person, error) {
	log.Printf("Looking for %v", username)

	found, err := db.FindPersonByUsername(username.GetUsername())

	idPerson, _ := strconv.Atoi(found.IdPerson)
	contact, _ := strconv.Atoi(found.Contact)
	person := &Person{IdPerson: uint32(idPerson), Contact: uint32(contact), Username: found.Username, Name: found.Name, Surname: found.Surname, Type: found.Role, Password: found.Password}
	if err != nil {
		return nil, err
	}
	return person, nil
}

// Update email
func (s *Server) UpdateEmail(ctx context.Context, newMail *UpdateField) (*Message, error) {
	log.Printf("Updating mail from %v", newMail.GetUsername())
	err := db.UpdateEmail(newMail.GetUsername(), newMail.GetNewContent())
	if err != nil {
		return nil, err
	}
	return &Message{Body: "User " + newMail.GetUsername() + " mail (" + newMail.GetNewContent() + ") updated successfuly!"}, nil
}

// Update phone
func (s *Server) UpdatePhone(ctx context.Context, newMail *UpdateField) (*Message, error) {
	log.Printf("Updating phone from %v", newMail.GetUsername())
	err := db.UpdatePhone(newMail.GetUsername(), newMail.GetNewContent())
	if err != nil {
		return nil, err
	}
	return &Message{Body: "User " + newMail.GetUsername() + " phone (" + newMail.GetNewContent() + ") updated successfuly!"}, nil
}

// Update facebook
func (s *Server) UpdateFacebook(ctx context.Context, newMail *UpdateField) (*Message, error) {
	log.Printf("Updating facebook from %v", newMail.GetUsername())
	err := db.UpdateFacebook(newMail.GetUsername(), newMail.GetNewContent())
	if err != nil {
		return nil, err
	}
	return &Message{Body: "User " + newMail.GetUsername() + " facebook (" + newMail.GetNewContent() + ") updated successfuly!"}, nil
}

// Update Name
func (s *Server) UpdateName(ctx context.Context, newMail *UpdateField) (*Message, error) {
	log.Printf("Updating name from %v", newMail.GetUsername())
	err := db.UpdateName(newMail.GetUsername(), newMail.GetNewContent())
	if err != nil {
		return nil, err
	}
	return &Message{Body: "User " + newMail.GetUsername() + " name (" + newMail.GetNewContent() + ") updated successfuly!"}, nil
}

// Update surname
func (s *Server) UpdateSurname(ctx context.Context, newMail *UpdateField) (*Message, error) {
	log.Printf("Updating surname from %v", newMail.GetUsername())
	err := db.UpdateSurname(newMail.GetUsername(), newMail.GetNewContent())
	if err != nil {
		return nil, err
	}
	return &Message{Body: "User " + newMail.GetUsername() + " surname (" + newMail.GetNewContent() + ") updated successfuly!"}, nil
}

// Update Password
func (s *Server) UpdatePassword(ctx context.Context, newMail *UpdateField) (*Message, error) {
	log.Printf("Updating password from %v", newMail.GetUsername())
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newMail.GetNewContent()), bcrypt.DefaultCost)
	err := db.UpdatePassword(newMail.GetUsername(), string(hashedPassword))
	if err != nil {
		return nil, err
	}
	return &Message{Body: "User " + newMail.GetUsername() + " password (****************) updated successfuly!"}, nil
}
