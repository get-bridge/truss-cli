#!/usr/bin/env bash

cat <<END >"$1"
secrets:
  api:
    db_password: "$DB_PASSWORD"
END
