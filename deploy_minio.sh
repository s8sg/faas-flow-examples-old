#!/bin/bash

docker service rm minio

docker service create --constraint="node.role==manager" \
 --name minio \
 --detach=true --network func_functions \
 --secret s3-access-key \
 --secret s3-secret-key \
 --env MINIO_SECRET_KEY_FILE=s3-secret-key \
 --env MINIO_ACCESS_KEY_FILE=s3-access-key \
minio/minio:RELEASE.2018-11-22T02-51-56Z server /export
