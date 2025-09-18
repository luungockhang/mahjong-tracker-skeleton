#!/bin/sh
set -e
psql $DB_URL -f migrations/seed.sql