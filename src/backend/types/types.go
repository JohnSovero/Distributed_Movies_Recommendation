package types

// Movie representa una película con su ID, título y géneros
type Movie struct {
	MovieID  int      `json:"id"`
	Title    string   `json:"title"`
	Year     int      `json:"year"`
	Genres   []string `json:"genres"`
	IMDBLink string   `json:"imdb_link"`
	TMDBLink string   `json:"tmdb_link"`
}

// Rating representa una calificación de un usuario a un ítem
type Rating struct {
	UserID int
	ItemID int
	Score  float64
}

// Item representa un ítem con sus calificaciones
type Item struct {
	ID      int
	Ratings map[int]float64
}

// Estructura para representar a un usuario
type User struct {
	ID      int
	Ratings map[int]float64
}

// Estructura para recibir datos del cliente
type FromClientData struct {
	Similarity float64 `json:"similarity"`
	UserID     string  `json:"userID"`
}

// Estructura para enviar datos al cliente
type ToClientData struct {
	User1 map[int]float64 `json:"user1"`
	User2 map[int]User    `json:"user2"`
}

// Estructura para ordenar las similitudes
type Kv struct {
	Key   int
	Value float64
}
