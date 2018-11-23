#!/bin/bash

echo "Attempting to create credentials for minio.."

SECRET_KEY=$(head -c 12 /dev/urandom | shasum| cut -d' ' -f1)
echo -n "$SECRET_KEY" | docker secret create s3-secret-key -
if [ $? = 0 ];
then
  echo "[Credentials] s3-secret-key : $SECRET_KEY"
else
  echo "[Credentials]\n s3-secret-key already exist, not creating"
fi

ACCESS_KEY=$(head -c 12 /dev/urandom | shasum| cut -d' ' -f1)
echo -n "$ACCESS_KEY" | docker secret create s3-access-key -
if [ $? = 0 ];
then
  echo "[Credentials] s3-access-key : $SECRET_KEY"
else
  echo "[Credentials]\n s3-access-key already exist, not creating"
fi

faas-cli deploy -f stack.yml


