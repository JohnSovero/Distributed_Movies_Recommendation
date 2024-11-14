# Usa una imagen base de Golang
FROM golang:1.23

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos del proyecto al directorio de trabajo
COPY . .

# Descarga las dependencias del módulo Go
RUN go mod download

# Actualiza todas las dependencias del módulo Go a la última versión
RUN go get -u ./...

# Limpia el módulo y descarga las dependencias necesarias
RUN go mod tidy

# Compila la aplicación y genera un ejecutable llamado 'api'
RUN go build -o api .

# Comando para ejecutar la aplicación
CMD ["./api"]