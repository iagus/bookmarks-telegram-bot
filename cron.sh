#!/usr/bin/bash

NPM_TESTS="Tests"
FORMATTER="Format"

# # START

# Updating repository
cd /opt/bookmarks-telegram-bot
git fetch origin main

# Checking for updates
LOCAL_HASH=$(git rev-parse HEAD)
REMOTE_HASH=$(git rev-parse "origin/main")
if [ "$LOCAL_HASH" = "$REMOTE_HASH" ]; then
  printf "[telegram-bookmarks-bot] No updates detected."
  exit 0
else
  git reset --soft origin/main
fi

# Cd'ing to Node server
cd server

npm install

if npm run --silent test; then
  echo "${NPM_TESTS}: OK"

  if npm run --silent check; then
    echo "${FORMATTER}: OK"

    cd ..

    printf "\n"
    printf "[telegram-bookmarks-bot] Compiling Go script"
    /usr/local/go/bin/go build main.go

    printf "\n\n"
    echo "[telegram-bookmarks-bot] All systems operational!"
  else
    echo "[telegram-bookmarks-bot] ${FORMATTER}: Failure"
  fi
else
  echo "[telegram-bookmarks-bot] ${NPM_TESTS}: Failure"
fi

