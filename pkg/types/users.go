package types

type Person struct {
	IdPerson string `json:"idPerson"`
	Contact  string `json:"contact"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Role     string `json:"type"`
	Password string `json:"password"`
}

type Medic struct {
	IdMedic  string `json:"idMedic"`
	CRM      string `json:"crm"`
	IdPerson string `json:"idPerson"`
}

type Patient struct {
	IdPatient string `json:"idPatient"`
	IdPerson  string `json:"idPerson"`
}

type Admin struct {
	IdAdmin  string `json:"idAdmin"`
	IdPerson string `json:"idPerson"`
}

type Contact struct {
	IdContact string `json:"idContact"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Facebook  string `json:"facebook"`
}
