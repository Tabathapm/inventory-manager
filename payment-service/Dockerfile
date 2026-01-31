# Uasr la imagen oficial de Go
FROM golang:1.23

# Crear la carpeta de trabajo adentro del contenedor
WORKDIR /app

# Copiar todo el código
COPY . .

# Descargar librerías (si hubiera) y construir el ejecutable
RUN go mod download
RUN go build -o main .

# IMPORTANTE: Exponer el puerto 8081 (El de pagos)
EXPOSE 8081

# Ejecutar
CMD ["./main"]