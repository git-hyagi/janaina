#!/bin/sh

#############################################################################
#   Just a little script to create/recreate the database inside mysql pod   #
#############################################################################


# pod='mariadb-1-nnp5z'
project='test3'
pod=$(oc -n $project get pods -l name=mariadb -o custom-columns=:.metadata.name)
oc project $project
oc -n $project exec $pod -- mysql -uroot -e 'drop database telemedicine';
oc -n $project exec $pod -- mysql -uroot -e 'create database telemedicine';
#oc -n $project exec $pod -- mysql -uroot -e 'drop table tipo' telemedicine;
#oc -n $project exec $pod -- mysql -uroot -e 'drop table historico' telemedicine;
#oc -n $project exec $pod -- mysql -uroot -e 'drop table medico_paciente' telemedicine;
#oc -n $project exec $pod -- mysql -uroot -e 'drop table agendamento' telemedicine;
#oc -n $project exec $pod -- mysql -uroot -e 'drop table receita' telemedicine;
#oc -n $project exec $pod -- mysql -uroot -e 'drop table admin' telemedicine;
#oc -n $project exec $pod -- mysql -uroot -e 'drop table medico' telemedicine;
#oc -n $project exec $pod -- mysql -uroot -e 'drop table paciente' telemedicine;
#oc -n $project exec $pod -- mysql -uroot -e 'drop table pessoa' telemedicine;
#oc -n $project exec $pod -- mysql -uroot -e 'drop table contato' telemedicine;

oc exec $pod -- mysql -uroot -e 'create table contato ( \
idContato int NOT NULL UNIQUE AUTO_INCREMENT, \
telefone VARCHAR(255) NOT NULL UNIQUE, \
email VARCHAR(255) NOT NULL UNIQUE, \
facebook VARCHAR(255) UNIQUE, \
PRIMARY KEY (idContato) \
)' telemedicine;

oc exec $pod -- mysql -uroot -e 'create table pessoa ( \
idPessoa int NOT NULL UNIQUE AUTO_INCREMENT, \
contato int NOT NULL UNIQUE, \
username VARCHAR(255) NOT NULL UNIQUE, \
nome VARCHAR(255) NOT NULL, \
sobrenome VARCHAR(255) NOT NULL, \
tipo VARCHAR(255) NOT NULL, \
senha VARCHAR(255) NOT NULL, \
PRIMARY KEY (idPessoa), \
FOREIGN KEY (contato) REFERENCES contato(idContato) ON DELETE CASCADE ON UPDATE CASCADE \
)' telemedicine;

oc exec $pod -- mysql -uroot -e 'create table admin ( \
idAdmin int NOT NULL UNIQUE AUTO_INCREMENT, \
idPessoa int NOT NULL UNIQUE, \
PRIMARY KEY (idAdmin), \
FOREIGN KEY (idPessoa) REFERENCES pessoa (idPessoa) \
)' telemedicine;

oc exec $pod -- mysql -uroot -e 'create table medico ( \
idMedico int NOT NULL UNIQUE AUTO_INCREMENT, \
CRM VARCHAR(255) NOT NULL UNIQUE, \
idPessoa int, \
PRIMARY KEY (idMedico), \
FOREIGN KEY (idPessoa) REFERENCES pessoa(idPessoa) \
)' telemedicine;

oc exec $pod -- mysql -uroot -e 'create table paciente ( \
idPaciente int NOT NULL UNIQUE AUTO_INCREMENT, \
idPessoa int, \
PRIMARY KEY (idPaciente), \
FOREIGN KEY (idPessoa) REFERENCES pessoa(idPessoa) \
)' telemedicine;

oc exec $pod -- mysql -uroot -e 'create table medico_paciente ( \
idMedico int, \
idPaciente int, \
PRIMARY KEY(idMedico,idPaciente), \
FOREIGN KEY (idMedico) REFERENCES medico(idMedico), \
FOREIGN KEY (idPaciente) REFERENCES paciente(idPaciente) \
)' telemedicine;

oc exec $pod -- mysql -uroot -e 'create table agendamento ( \
idAgendamento int NOT NULL UNIQUE AUTO_INCREMENT, \
idPaciente int NOT NULL UNIQUE, \
idMedico int NOT NULL UNIQUE, \
data TIMESTAMP NOT NULL, \
PRIMARY KEY (idAgendamento), \
FOREIGN KEY (idPaciente) REFERENCES paciente(idPaciente), \
FOREIGN KEY (idMedico) REFERENCES medico(idMedico) \
)' telemedicine;

oc exec $pod -- mysql -uroot -e 'create table receita ( \
idReceita int NOT NULL UNIQUE AUTO_INCREMENT, \
idMedico int, \
idPaciente int, \
idHistorico int, \
remedios VARCHAR(255), \
PRIMARY KEY(idReceita), \
FOREIGN KEY (idMedico) REFERENCES medico (idMedico), \
FOREIGN KEY (idPaciente) REFERENCES paciente (idPaciente) \
)' telemedicine;

oc exec $pod -- mysql -uroot -e 'create table historico ( \
idHistorico int NOT NULL UNIQUE AUTO_INCREMENT, \
idMedico int, \
idPaciente int, \
idReceita int, \
data TIMESTAMP NOT NULL, \
nota VARCHAR(255), \
PRIMARY KEY (idHistorico), \
FOREIGN KEY (idMedico) REFERENCES medico (idMedico), \
FOREIGN KEY (idPaciente) REFERENCES paciente (idPaciente), \
FOREIGN KEY (idReceita) REFERENCES receita (idReceita) \
)' telemedicine;

oc exec $pod -- mysql -uroot -e 'create table tipo ( \
idPessoa int,
idMedico int,
idAdmin int,
idPaciente int,
Role VARCHAR(255) NOT NULL,
PRIMARY KEY (idPessoa),
FOREIGN KEY (idPessoa) REFERENCES pessoa (idPessoa), \
FOREIGN KEY (idMedico) REFERENCES medico (idMedico), \
FOREIGN KEY (idAdmin) REFERENCES admin (idAdmin), \
FOREIGN KEY (idPaciente) REFERENCES paciente (idPaciente) \
)' telemedicine;


echo "Pessoa";          oc exec $pod -- mysql -uroot -e 'describe pessoa' --table telemedicine
echo "Contato";         oc exec $pod -- mysql -uroot -e 'describe contato' --table telemedicine
echo "Admin";           oc exec $pod -- mysql -uroot -e 'describe admin' --table telemedicine
echo "Medico";          oc exec $pod -- mysql -uroot -e 'describe medico' --table telemedicine
echo "Paciente";        oc exec $pod -- mysql -uroot -e 'describe paciente' --table telemedicine
echo "Medico-Paciente"; oc exec $pod -- mysql -uroot -e 'describe medico_paciente' --table telemedicine
echo "Agendamento";     oc exec $pod -- mysql -uroot -e 'describe agendamento' --table telemedicine
echo "Receita";         oc exec $pod -- mysql -uroot -e 'describe receita' --table telemedicine
echo "Historico";       oc exec $pod -- mysql -uroot -e 'describe historico' --table telemedicine
echo "Tipo";            oc exec $pod -- mysql -uroot -e 'describe tipo' --table telemedicine


# # NEW ADMIN
# oc exec $pod -- mysql -uroot -e 'insert into contato values(NULL,"11 1234-5678","john_doe@email.com","John Doe")' telemedicine
# oc exec $pod -- mysql -uroot -e 'insert into pessoa values(NULL,(select idContato from contato where email = "john_doe@email.com"), "johndoe", "John", "Doe", "admin", "password" )' telemedicine 
# oc exec $pod -- mysql -uroot -e 'insert into admin values(NULL,(select idPessoa from pessoa where username = "johndoe"))' telemedicine
# oc exec $pod -- mysql -uroot -e 'insert into tipo values((select idPessoa from pessoa where username = "johndoe"), NULL,(select admin.idAdmin from admin join pessoa on admin.idPessoa = pessoa.idPessoa),NULL,"admin" )' telemedicine
#
#
# # NEW MEDIC
# oc exec $pod -- mysql -uroot -e 'insert into contato values(NULL, "11 2345-6789", "gregory_house@email.com","Gregory House")' telemedicine
# oc exec $pod -- mysql -uroot -e 'insert into pessoa values(NULL, (select idContato from contato where email="gregory_house@email.com"),"house","Gregory", "House", "medic","password")' telemedicine
# oc exec $pod -- mysql -uroot -e 'insert into medico values (NULL,"CRM-12345",NULL,NULL,NULL,NULL,(select idPessoa from pessoa where username="house"))' telemedicine
# oc exec $pod -- mysql -uroot -e 'insert into tipo values((select idPessoa from pessoa where username = "house"),(select medico.idMedico from medico join pessoa on medico.idPessoa = pessoa.idPessoa),NULL,NULL,"medico" )' telemedicine
#
# # NEW PATIENT
# oc exec $pod -- mysql -uroot -e 'insert into contato values(NULL,"11 3456-7890","patch_adams@email.com","Patch Adams")' telemedicine
# oc exec $pod -- mysql -uroot -e 'insert into pessoa values(NULL,(select idContato from contato where email="patch_adams@email.com"),"padams", "Patch","Adams", "paciente", "password")' telemedicine
# oc exec $pod -- mysql -uroot -e 'insert into paciente values(NULL,(select medico.idMedico from medico join pessoa on pessoa.username="house"),NULL,NULL,(select idPessoa from pessoa where username="padams"),NULL )' telemedicine
# oc exec $pod -- mysql -uroot -e 'insert into tipo values((select idPessoa from pessoa where username="padams"),NULL,NULL,(select paciente.idPaciente from pessoa join paciente on pessoa.idPessoa = paciente.idPaciente),"paciente")' telemedicine
#
# # NEW APPOINTMENT
# oc exec $pod -- mysql -uroot -e 'insert into agendamento values(NULL, (select paciente.idPaciente from paciente join pessoa on pessoa.username = "padams"), (select medico.idMedico from medico join pessoa on pessoa.username = "house"), "2020-05-13 14:30:00" )' telemedicine
# oc exec $pod -- mysql -uroot -e 'update paciente join agendamento set paciente.idAgendamento = agendamento.idAgendamento on agendamento.idPaciente = paciente.idPaciente' telemedicine
#
# # NEW Historico
# oc exec $pod -- mysql -uroot telemedicine -e 'insert into historico values(NULL,(select idMedico from agendamento where idAgendamento = 1),(select idPaciente from agendamento where idAgendamento = 1), NULL, (select data from agendamento where idAgendamento = 1), "Foi constatdo que o paciente está com deficiência de vitaminas. Melhorar as refeições e evitar refrigerantes.")'
#
# # NEW PRESCRIPTION
# oc exec $pod -- mysql -uroot telemedicine -e 'insert into receita values(NULL, (select idMedico from agendamento where idAgendamento = 1),(select idPaciente from paciente where idAgendamento = 1), NULL, "- Comer de 3 em 3 horas.\n - Realizar refeições balanceadas.\n")'


# PRINTING
oc exec $pod -- mysql -uroot -e 'select * from pessoa' --table telemedicine
oc exec $pod -- mysql -uroot -e 'select * from contato' --table telemedicine
oc exec $pod -- mysql -uroot -e 'select * from admin' --table telemedicine
oc exec $pod -- mysql -uroot -e 'select * from medico' --table telemedicine
oc exec $pod -- mysql -uroot -e 'select * from paciente' --table telemedicine
oc exec $pod -- mysql -uroot -e 'select * from tipo' --table telemedicine
oc exec $pod -- mysql -uroot -e 'select * from receita' --table telemedicine
oc exec $pod -- mysql -uroot -e 'select * from historico' --table telemedicine


# print all person type medic with medic information
oc exec $pod -- mysql -uroot --table telemedicine -e 'select * from pessoa join tipo join medico on pessoa.idPessoa = tipo.idPessoa and pessoa.idPessoa = medico.idPessoa where tipo.idMedico is NOT NULL'
oc exec $pod -- mysql -uroot --table telemedicine -e 'select * from paciente join pessoa on paciente.idPessoa = pessoa.idPessoa'


# print all patient
oc exec $pod -- mysql -uroot --table telemedicine -e 'select paciente.idPaciente,pessoa.idPessoa,pessoa.contato,pessoa.username,pessoa.nome,pessoa.sobrenome,pessoa.tipo,pessoa.senha from paciente join pessoa on paciente.idPessoa = pessoa.idPessoa'
