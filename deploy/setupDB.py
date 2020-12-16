#!/bin/bash
import os, argparse

parser = argparse.ArgumentParser()

parser.add_argument("-D", "--databasename", help="Database name")
parser.add_argument("-U", "--username", help="User name")
parser.add_argument("-P", "--password", help="Password")

args = parser.parse_args()

if(args.databasename and args.username and args.password):

    userCreate = os.popen('sudo -u postgres -H -- psql -c "CREATE USER {} WITH PASSWORD \'{}\';"'.format(args.username, args.password))
    databaseCreate = os.popen('sudo -u postgres -H -- psql -c "CREATE DATABASE {} OWNER {};"'.format(args.databasename, args.username))

    print(userCreate.read())
    print(databaseCreate.read())

else:
    print("Specify database name, password and username.")
