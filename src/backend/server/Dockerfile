# Usa una imagen base de Golang
FROM golang:1.23

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia todo el contenido del directorio actual al directorio de trabajo en el contenedor
COPY . .

# Descarga las dependencias del módulo Go
RUN go mod download

# Limpia el módulo y descarga las dependencias necesarias
RUN go mod tidy

# Compila la aplicación y genera un ejecutable llamado 'server'
RUN go build -o server .

# Asegura que el archivo 'server' tenga permisos de ejecución
RUN chmod +x server

# Expone el puerto 9000 en el que la aplicación escuchará
EXPOSE 9000

# Comando para ejecutar la aplicación
CMD ["./server"]