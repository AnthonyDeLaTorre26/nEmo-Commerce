package models

import "errors"

// Producto representará un artículo del catálogo.
type Producto struct {
	codigo string
	nombre string
	precio float64
	stock  int
}

// La función NuevoProducto crea un producto validando los datos.
func NuevoProducto(codigo, nombre string, precio float64, stock int) (*Producto, error) {
	if codigo == "" {
		return nil, errors.New("el código no puede estar vacío")
	}

	if nombre == "" {
		return nil, errors.New("el nombre no puede estar vacío")
	}

	if precio <= 0 {
		return nil, errors.New("el precio debe ser mayor que cero")
	}

	if stock < 0 {
		return nil, errors.New("el stock no puede ser negativo")
	}

	return &Producto{
		codigo: codigo,
		nombre: nombre,
		precio: precio,
		stock:  stock,
	}, nil
}

func (p *Producto) GetCodigo() string {
	return p.codigo
}

func (p *Producto) GetNombre() string {
	return p.nombre
}

func (p *Producto) GetPrecio() float64 {
	return p.precio
}

func (p *Producto) GetStock() int {
	return p.stock
}

func (p *Producto) SetPrecio(precio float64) error {
	if precio <= 0 {
		return errors.New("el precio debe ser mayor que cero")
	}

	p.precio = precio
	return nil
}

func (p *Producto) ReducirStock(cantidad int) error {
	if cantidad <= 0 {
		return errors.New("la cantidad debe ser positiva")
	}

	if cantidad > p.stock {
		return errors.New("stock insuficiente")
	}

	p.stock -= cantidad
	return nil
}

func (p *Producto) AumentarStock(cantidad int) error {
	if cantidad <= 0 {
		return errors.New("la cantidad debe ser positiva")
	}

	p.stock += cantidad
	return nil
}
