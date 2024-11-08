package fc

// Movie representa una película con su ID, título y géneros
type Movie struct {
    MovieID int
    Title   string
    Genres  []string
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
    User2 map[int]User `json:"user2"`
}
// Estructura para ordenar las similitudes
type kv struct {
	Key   int
	Value float64
}