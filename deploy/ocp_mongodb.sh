#!/bin/bash

oc project test3
oc new-app mongodb-ephemeral

oc extract secret/mongodb --to=- --keys=database-admin-password
oc create svc nodeport mongodb-nodeport --tcp 27017:27017 --node-port=30901
oc set selector svc mongodb-nodeport deploymentconfig=mongodb
