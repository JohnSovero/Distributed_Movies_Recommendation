package fc

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