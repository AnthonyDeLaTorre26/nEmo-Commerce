package interfaces

// Comprable va a definir el limite mínimo de un elemento para poder ser comprado.
type Comprable interface {
	GetNombre() string
	GetPrecio() float64
}
