FROM golang:1.24-alpine

WORKDIR /app

# Primero copiamos solo los archivos de requisitos.
# Si estos archivos no cambian, Docker no vuelve a descargar nada (ahorra tiempo).
COPY go.mod go.sum ./

# Descargar las librerías (como el driver de Postgres)
RUN go mod download

# Copiar el código fuente (main.go, etc)
COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]