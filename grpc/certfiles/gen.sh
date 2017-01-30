#!/usr/bin/env bash

PASSWORD=`cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1`
HOST=$1
if [ -z $HOST ]; then HOST="localhost"; fi
PREFIX=$2
if [ -z $PREFIX ]; then PREFIX="dev"; fi

openssl genrsa -passout pass:$PASSWORD -des3 -out $PREFIX'server.key' 4096
openssl req -passin pass:$PASSWORD -new -x509 -days 3650 -key $PREFIX'server.key' -out $PREFIX'server.crt' -subj '/C=US/ST=CA/L=Sunnyvale/O=MyApp/CN='$HOST'/emailAddress=foo@bar.com'
openssl rsa -passin pass:$PASSWORD -in $PREFIX'server.key' -out $PREFIX'server.key'