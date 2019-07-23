#!/bin/sh

#Check if server is up

if [[ -z "${CONDUCTOR_HOST}" ]]; then
  echo ERRO: CONDUCTOR_HOST não informado
  exit 1
fi

if [[ -z "${CONDUCTOR_PORT}" ]]; then
  echo ERRO: CONDUCTOR__PORT não informado
  exit 1
fi

while [ "$(nc -z $CONDUCTOR_HOST $CONDUCTOR_PORT </dev/null; echo $?)" !=  "0" ];
do sleep 5;
echo "Waiting for CONDUCTOR SERVER is UP and RESPONDING";
done

sleep 5;

BACKY2LS=$(backy2 ls)
if [[ $BACKY2LS = *"Please run initdb first"* ]]
then
    echo "== Initializing Backy DB =="
    backy2 initdb
else
    echo "== Backy2 DB is inited! =="
fi

cat /backy.cfg.template | envsubst > /etc/backy.cfg
cat /etc/backy.cfg

#redirect backy2 logs to stdout and remove internal log file to avoid increasing endlessly
tail -f /var/log/backy.log&
while true; do if [ -f /var/log/backy.log ]; then rm /var/log/backy.log; fi; sleep 86400; done&

/bin/backy-scraper


