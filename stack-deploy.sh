#!/bin/sh

if ! [ -x "$(command -v docker)" ]; then
  echo 'Unable to find docker command, please install Docker (https://www.docker.com/) and retry' >&2
  exit 1
fi

export BASIC_AUTH="true"
export AUTH_URL="http://basic-auth-plugin:8080/validate"

# Secrets should be created even if basic-auth is disabled.
echo "Attempting to create credentials for gateway.."
echo "admin" | docker secret create basic-auth-user -
# For local environment setup
# We will use admin as password
secret="admin"
echo "${secret}" | docker secret create basic-auth-password -
if [ $? = 0 ];
then
  echo "[Credentials]\n username: admin \n password: ${secret}\n echo -n "${secret}" | faas-cli login --username=admin --password-stdin"
else
  echo "[Credentials]\n already exist, not creating"
fi

if [ ${BASIC_AUTH} = "true" ];
then
  echo ""
  echo "Enabling basic authentication for gateway.."
  echo ""
else
  echo ""
  echo "Disabling basic authentication for gateway.."
  echo ""
fi

docker stack deploy scaling --compose-file docker-compose.yml
