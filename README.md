# go-tcp-server
A dummy tcp server written Go to handle TPS and return 1800 as a successful Health Check

Run these commands
- go mod tidy      # updates your go.sum and records the replace
- go run ./cmd/tcpserver # this will start the go server

Once the server is up, in an adjacent terminal, run this command:
- echo -n "$(head -c16 /dev/urandom | base64 | tr -dc 'A-Za-z0-9' | head -c16)" | nc localhost 8080 | xxd
Generates a random 16‑char alphanumeric string, pipes it to your server on localhost:8080 via nc, and uses xxd to hex‑dump the 4‑byte reply.
