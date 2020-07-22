# ABOUT THIS PROJECT
This is just a personal project of a basic telemedicine app.


# FEATURES
- [x] `GRPC` unnary connections through a secure communication using `TLS` and `JWT` authentication
- [ ] [WIP] RBAC
- [x] CRUD user/medic/patient  
   **note**: even though this is the first working version, there are a lot of issues and pending work to be done
- [x] Chat (`Bi-directional Streaming GRPC` server (`golang`) and client (`flutter/dart`))  
   **note**: as of CRUD, there are a lot of issues and pending work (for now, it is already sending and receiving messages with some disconnections and other bugs)
- [ ] schedule doctor appointment
- [ ] prescription
  - [ ] digital signature
- [ ] video streaming


# The database
- [ ] [WIP] First draft running as a `MySQL` `Kubernetes/OpenShift` pods using script [tables.sh](deploy/tables.sh) to create the database schema. 
![alt text](https://github.com/git-hyagi/janaina/blob/master/MER.png?raw=true)

# Usage

To run the server as a `Kubernetes`/`OpenShift` pod, create a secret called `config.yaml` and mount it in the `/opt/app-root` (default for `OpenShift s2i`) folder.

## config file
- should be named config.yaml
- format
~~~
database:
  name: "<db name>"
  user: "<db user name>"
  password:  "<db user password>"
  address: "<db host>:<db port>"

table:
  name: "medic"
~~~

#### DRAFT openshift
* Database settings [**DEPRECATED!** (database schema changed)]
~~~
pod=$(oc -n test3 get pods -l name=mariadb -o custom-columns=:.metadata.name)
oc exec $pod -- mysql -u root -e "create table medic (name VARCHAR(100) CHARACTER SET utf8, address VARCHAR(100) CHARACTER SET utf8, CRM VARCHAR(100) CHARACTER SET utf8);" telemedicine
oc exec $pod -- mysql -u root -e "select * from medic\G" telemedicine
oc exec $pod -- mysql -u root -e "select * from medic" --table telemedicine
~~~

* Server app settings
~~~
- ssh key to pull the repo
oc create secret generic --type=kubernetes.io/ssh-auth github --from-file=ssh-privatekey=<private key>
~~~

* deploy/config.yaml file with database connection info
~~~
oc create secret generic my-config  --from-file=deploy/config.yaml
~~~

* create app
~~~
oc new-app golang~git@github.com/git-hyagi/janaina.git --source-secret=github --name telemedicine
~~~

* create a svc nodePort to expose the app
~~~
oc create svc nodeport my-api --tcp=9000:9000 --node-port=30900
~~~

* mount the secret
~~~
oc set volume deployment/telemedicine --add -m /opt/app-root/config.yaml --sub-path=config.yaml -t secret --secret-name=my-config
~~~

# REFERENCES
Big thanks to:  
[TECH SCHOOL](https://gitlab.com/techschool/pcbook)  
[TENSOR PROGRAMMING](https://github.com/tensor-programming/docker_grpc_chat_tutorial)  
