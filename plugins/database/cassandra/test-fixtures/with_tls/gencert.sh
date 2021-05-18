#!/bin/sh

################################################################
# Usage: ./gencert.sh
# 
# Generates a keystore.jks file that can be used with a
# Cassandra server for TLS connections. This does not update
# a cassandra config file.
################################################################

set -e

KEYFILE="key.pem"
CERTFILE="cert.pem"
PKCSFILE="keystore.p12"
JKSFILE="keystore.jks"

HOST="127.0.0.1"
NAME="cassandra"
ALIAS="cassandra"
PASSWORD="cassandra"

echo "# Generating certificate keypair..."
go run /usr/local/go/src/crypto/tls/generate_cert.go --host=${HOST}

echo "# Creating keystore..."
openssl pkcs12 -export -in ${CERTFILE} -inkey ${KEYFILE} -name ${NAME} -password pass:${PASSWORD} > ${PKCSFILE}

echo "# Creating Java key store"
if [ -e "${JKSFILE}" ]; then
	echo "# Removing old key store"
	rm ${JKSFILE}
fi

set +e
keytool -importkeystore \
	-srckeystore ${PKCSFILE} \
	-srcstoretype PKCS12 \
	-srcstorepass ${PASSWORD} \
	-destkeystore ${JKSFILE} \
	-deststorepass ${PASSWORD} \
	-destkeypass ${PASSWORD} \
	-alias ${ALIAS}

echo "# Removing intermediate files"
rm ${KEYFILE} ${CERTFILE} ${PKCSFILE}
