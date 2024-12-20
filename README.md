# Distributed_Movies_Recommendation
PC4 del curso de Programacion Concurrente y Distribuida

Dataset: https://grouplens.org/datasets/movielens/25m/


- Para ejecutar el código, descargar el zip de 25 millones de datos y extraerlos en la carpeta dataset/
- Luego ejecutar el main.go con el comando "go run main.go" y los clientes con "go run client.go", situándose en la ruta correspondiente para cada uno.
- Los datasets subidos al repositorio son de ejemplo ya que github no permite archivos csv de 660 MB
- No olvidar renombrar el archivo movies.csv como movies25.csv y ratings.csv como ratings25.csv

Si es la primera vez que ejecuta la aplicación, debe utilizar el comando `docker compose up --build`.

Para el correcto funcionamiento del frontend, se requiere [crear una cuenta de desarrollador en TMDB](https://developer.themoviedb.org/reference/intro/getting-started) para obtener acceso a una llave API y a un token de lectura.

Una vez con estas credenciales, generar dentro de la ruta `src/frontend/src/environments` un archivo `environment.ts` con contenido:
<!-- codigo -->
```
export const environment = {
    production: false,
    tmdbApiKey: <API-KEY>,
    tmdbApiReadToken: <API-READ-TOKEN>
};
```

Integrantes:
- André Dario Pilco Chiuyare
- John Davids Sovero Cubillas
