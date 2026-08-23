package models

import "errors"

// Producto representará un artículo del catálogo.
// Cumple con princio de encapsulación al tener sus atributos privados.
type Producto struct {
	id     int
	codigo string
	nombre string
	precio float64
	stock  int
}

// La función NuevoProducto crea un Producto validando los datos proporcionados.
func NuevoProducto(codigo, nombre string, precio float64, stock int) (*Producto, error) {
	if codigo == "" {
		return nil, errors.New("El código no puede estar vacío")
	}

	if nombre == "" {
		return nil, errors.New("El nombre no puede estar vacío")
	}

	if precio <= 0 {
		return nil, errors.New("El precio debe ser mayor que cero")
	}

	if stock < 0 {
		return nil, errors.New("El stock no puede ser negativo")
	}

	return &Producto{
		codigo: codigo,
		nombre: nombre,
		precio: precio,
		stock:  stock,
	}, nil
}

// GetID devuelve la identificación del producto.
func (p *Producto) GetID() int {
	return p.id
}

// SetID asigna el identificador generado por la base de datos.
func (p *Producto) SetID(id int) {
	p.id = id
}

// GetCodigo devuelve el código del producto.
func (p *Producto) GetCodigo() string {
	return p.codigo
}

// GetNombre devuelve el nombre del producto.
func (p *Producto) GetNombre() string {
	return p.nombre
}

// GetPrecio devuelve el precio del producto.
func (p *Producto) GetPrecio() float64 {
	return p.precio
}

// GetStock devuelve la cantidad disponible del producto.
func (p *Producto) GetStock() int {
	return p.stock
}

// SetPrecio permite modificar el precio del producto evitando que sea mayor que cero.
func (p *Producto) SetPrecio(precio float64) error {
	if precio <= 0 {
		return errors.New("El precio debe ser mayor que cero")
	}

	p.precio = precio
	return nil
}

// ReducirStock disminuye la cantidad disponible del producto evitando reducir una cantidad mayo al stock.
func (p *Producto) ReducirStock(cantidad int) error {
	if cantidad <= 0 {
		return errors.New("La cantidad debe ser positiva")
	}

	if cantidad > p.stock {
		return errors.New("Stock insuficiente")
	}

	p.stock -= cantidad
	return nil
}

// AumentarStock aumenta la cantidad disponible del producto sin límite.
func (p *Producto) AumentarStock(cantidad int) error {
	if cantidad <= 0 {
		return errors.New("La cantidad debe ser positiva")
	}

	p.stock += cantidad
	return nil
}
