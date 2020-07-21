/*
  This is just a program to test the methods to work with the person grpc methods.
  [TO-DO] This code should be removed and a test file should be used instead.
*/

package main

import (
	"github.com/git-hyagi/janaina/pkg/authentication"
	"github.com/git-hyagi/janaina/pkg/person"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	username        = "John Doe"
	password        = "password"
	refreshDuration = 30 * time.Second
	caFile          = "deploy/ca.crt"
	servername      = "golab"
)

func authMethods() map[string]bool {
	return map[string]bool{
		"/person.ManagePerson/CreateMedic":    true,
		"/person.ManagePerson/CreateAdmin":    true,
		"/person.ManagePerson/CreatePatient":  true,
		"/person.ManagePerson/RemovePerson":   true,
		"/person.ManagePerson/FindPerson":     true,
		"/person.ManagePerson/UpdateEmail":    true,
		"/person.ManagePerson/UpdatePhone":    true,
		"/person.ManagePerson/UpdateFacebook": true,
		"/person.ManagePerson/UpdateName":     true,
		"/person.ManagePerson/UpdateSurname":  true,
		"/person.ManagePerson/UpdatePassword": true,
	}
}

func main() {

	var conn *grpc.ClientConn

	// Load CA file to verify server authencity
	serverCA, err := authentication.LoadServerCA(caFile)
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	// create a grpc connection to server
	conn, err = grpc.Dial(servername+":9000", grpc.WithTransportCredentials(serverCA))
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}
	defer conn.Close()

	clientAuth := authentication.NewClientAuth(conn, "johnDoe2", "password2")
	interceptor, err := authentication.NewClientAuthInterceptor(clientAuth, authMethods(), refreshDuration)
	if err != nil {
		log.Fatalf("cannot create auth interceptor: %v", err)
	}

	authConn, err := grpc.Dial(servername+":9000", grpc.WithTransportCredentials(serverCA), grpc.WithUnaryInterceptor(interceptor.Unary()))
	if err != nil {
		log.Fatalf("cannot dial server: %v", err)
	}

	e := person.NewManagePersonClient(authConn)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	medic := &person.Medic{IdMedic: 123, Crm: "CRM-12345X", IdPerson: 123, Person: &person.Person{IdPerson: 123, Contact: 123, Username: "jsilva2", Name: "Joao", Surname: "da Silva", Type: "medico", Password: string(hashedPassword)}, Contact: &person.Contact{IdContact: 123, Phone: "234-56789", Email: "jsilva2@email.com", Facebook: "Joao da Silva2"}}
	msg, err := e.CreateMedic(context.Background(), medic)
	if err != nil {
		log.Fatalf("Error when calling AddMedic: %s", err)
	}
	log.Printf("Resonse from server: %s", msg)

	admin := &person.Admin{IdAdmin: 123, IdPerson: 123, Person: &person.Person{IdPerson: 123, Contact: 123, Username: "jsilva3", Name: "Joao", Surname: "da Silva", Type: "admin", Password: string(hashedPassword)}, Contact: &person.Contact{IdContact: 123, Phone: "234-567890", Email: "jsilva3@email.com", Facebook: "Joao da Silva3"}}
	msg, err = e.CreateAdmin(context.Background(), admin)
	if err != nil {
		log.Fatalf("Error when calling AddAdmin: %s", err)
	}
	log.Printf("Resonse from server: %s", msg)

	patient := &person.Patient{IdPatient: 123, IdPerson: 123, Person: &person.Person{IdPerson: 123, Contact: 123, Username: "jsilva4", Name: "Joao", Surname: "da Silva", Type: "admin", Password: string(hashedPassword)}, Contact: &person.Contact{IdContact: 123, Phone: "234-5678901", Email: "jsilva4@email.com", Facebook: "Joao da Silva4"}}
	msg, err = e.CreatePatient(context.Background(), patient)
	if err != nil {
		log.Fatalf("Error when calling AddPatient: %s", err)
	}
	log.Printf("Resonse from server: %s", msg)

	msg, err = e.RemovePerson(context.Background(), &person.PrimaryKey{Username: "jsilva2"})
	if err != nil {
		log.Fatalf("Error when calling AddMedic: %s", err)
	}
	log.Printf("Resonse from server: %s", msg)
	msg, err = e.RemovePerson(context.Background(), &person.PrimaryKey{Username: "jsilva3"})
	if err != nil {
		log.Fatalf("Error when calling AddMedic: %s", err)
	}
	log.Printf("Resonse from server: %s", msg)
	msg, err = e.RemovePerson(context.Background(), &person.PrimaryKey{Username: "jsilva4"})
	if err != nil {
		log.Fatalf("Error when calling AddMedic: %s", err)
	}
	log.Printf("Resonse from server: %s", msg)

	found, err := e.FindPerson(context.Background(), &person.PrimaryKey{Username: "johnDoe"})
	if err != nil {
		log.Fatalf("Error trying to find person: %s", err)
	}
	log.Printf("Resonse from server: %v", found)

	msg, err = e.UpdateEmail(context.Background(), &person.UpdateField{Username: "johnDoe", NewContent: "JD@email.com"})
	if err != nil {
		log.Fatalf("Error trying to update person: %s", err)
	}
	log.Printf("Resonse from server: %v", msg)

	msg, err = e.UpdatePhone(context.Background(), &person.UpdateField{Username: "johnDoe", NewContent: "234-5678"})
	if err != nil {
		log.Fatalf("Error trying to update person: %s", err)
	}
	log.Printf("Resonse from server: %v", msg)

	msg, err = e.UpdateFacebook(context.Background(), &person.UpdateField{Username: "johnDoe", NewContent: "NIL"})
	if err != nil {
		log.Fatalf("Error trying to update person: %s", err)
	}
	log.Printf("Resonse from server: %v", msg)

	msg, err = e.UpdateName(context.Background(), &person.UpdateField{Username: "johnDoe", NewContent: "JOHNN"})
	if err != nil {
		log.Fatalf("Error trying to update person: %s", err)
	}
	log.Printf("Resonse from server: %v", msg)

	msg, err = e.UpdateSurname(context.Background(), &person.UpdateField{Username: "johnDoe", NewContent: "Due"})
	if err != nil {
		log.Fatalf("Error trying to update person: %s", err)
	}
	log.Printf("Resonse from server: %v", msg)

	msg, err = e.UpdatePassword(context.Background(), &person.UpdateField{Username: "johnDoe", NewContent: "password"})
	if err != nil {
		log.Fatalf("Error trying to update person: %s", err)
	}
	log.Printf("Resonse from server: %v", msg)
}
