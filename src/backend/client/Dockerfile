# Usa una imagen base de Golang
FROM golang:1.23

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos del proyecto al directorio de trabajo
COPY . .

# Descarga las dependencias del módulo Go
RUN go mod download

# Compila la aplicación y genera un ejecutable llamado 'api'
RUN go build -o client .

# Comando para ejecutar la aplicación
CMD ["./client"]