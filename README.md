# Movies App Recomender with Golang using Distributed and Concurrent Programming
---
## Dataset
[MovieLens Dataset](https://grouplens.org/datasets/movielens/latest/)

---
# How to Use

- To execute the code, download the ZIP of 25 million data and extract them in the folder DATASET
- Then execute the `main.go` with the `Go Run Main.go` command and customers with `Go Run Client.go`, placing on the corresponding route for each one.
- The datasets rise to the repository are of an example since Github does not allow 660 MB CSV files
- Do not forget to rename the `Movies.csv` file as `movies25.csv` and `ratings.csv` as `ratings25.csv`

If it is the first time you run the application, you must use the `Docker Compose Up -Build`.

For the proper functioning of the border, [create a developer account in TMDB] (https://developer.themoviedb.org/reference/Inter/geting-started) is required to get access to an API key and a reading token.

Once with these credentials, generate within the route `src/frontend/environments` a file `environment.ts` with content:
<!-- CODE -->
```
export const environment = {
    production: false,
    tmdbApiKey: <API-KEY>,
    tmdbApiReadToken: <API-READ-TOKEN>
};
```
---
# Members:
- André Dario Pilco Chiuyare
- John Davids Sovero Cubillas
