# Chirpy
boot.dev http server example


To build and run:
`go build -o out && ./out`

Link to server: [https://localhost:8080](https://localhost:8080)


To start the postgres server: `sudo service postgresql start`

Enter shell: `sudo -u postgres psql`

`psql "postgres://postgres:postgres@localhost:5432/chirpy"`

goose postgres "postgres://postgres:postgres@localhost:5432/chirpy" up