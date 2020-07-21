Table contato {
  idContato int [pk,increment]
  telefone varchar [not null,unique]
  email varchar [not null,unique]
  facebook varchar [unique]
}

Table pessoa {
idPessoa int [pk, increment]
contato int [not null,unique]
username varchar [not null, unique]
nome varchar [not null]
sobrenome varchar [not null]
tipo varchar [not null]
senha varchar [not null]
}

Ref: pessoa.contato - contato.idContato [delete: cascade, update: cascade]

Table admin {
  idAdmin int [pk,increment]
  idPessoa int [not null, unique]
}

Ref: pessoa.idPessoa - admin.idPessoa

Table medico {
  idMedico int [pk,increment]
CRM varchar [not null, unique]
idPessoa int
}

Ref: pessoa.idPessoa - medico.idPessoa 

Table paciente {
idPaciente int [pk,increment]
idPessoa int
}

Ref: pessoa.idPessoa - paciente.idPessoa

Table medico_paciente {
idMedico int
idPaciente int
Indexes {
  (idMedico,idPaciente) [pk]
}
}

Ref: medico_paciente.idMedico > medico.idMedico
Ref: medico_paciente.idPaciente > paciente.idPaciente

Table agendamento {
idAgendamento int [pk,increment]
idPaciente int 
idMedico int
data TIMESTAMP [not null]
}

Ref: agendamento.idPaciente > paciente.idPaciente
Ref: agendamento.idMedico > medico.idMedico

Table receita {
idReceita int [pk,increment]
idMedico int
idPaciente int
idHistorico int
remedios varchar
}

Ref: receita.idMedico > medico.idMedico
Ref: receita.idPaciente > paciente.idPaciente

Table historico {
idHistorico int [pk,increment]
idMedico int
idPaciente int
idReceita int
data TIMESTAMP [not null]
nota varchar
}

Ref: historico.idMedico > medico.idMedico
Ref: historico.idPaciente > paciente.idPaciente
Ref: historico.idReceita - receita.idReceita

Table tipo {
idPessoa int [pk,not null]
idMedico int
idAdmin int
idPaciente int
Role varchar [not null]
}

Ref: tipo.idPessoa - pessoa.idPessoa
Ref: tipo.idMedico - medico.idMedico
Ref: tipo.idPaciente - paciente.idPaciente
Ref: tipo.idAdmin - admin.idAdmin
