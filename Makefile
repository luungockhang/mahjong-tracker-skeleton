dev:
docker-compose up --build


test:
cd api && go test ./...


migrate-up:
cd api && migrate -path migrations -database $$DB_URL up


migrate-down:
cd api && migrate -path migrations -database $$DB_URL down