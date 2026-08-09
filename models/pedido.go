package models

import "nemo-commerce/interfaces"

// Pedido representa una compra realizada por un cliente y lo contiene como historial.
type Pedido struct {
	id        int
	cliente   *Cliente
	productos []interfaces.Comprable
	total     float64
}

// La función NuevoPedido crea una instancia de Pedido asociandolo con identififcación y cliente.
func NuevoPedido(id int, cliente *Cliente) *Pedido {
	return &Pedido{
		id:        id,
		cliente:   cliente,
		productos: []interfaces.Comprable{},
		total:     0,
	}
}

// GetID devuelve la identificación del pedido.
func (p *Pedido) GetID() int {
	return p.id
}

// GetCliente devuelve el cliente asociado al pedido.
func (p *Pedido) GetCliente() *Cliente {
	return p.cliente
}

// GetProductos devuelve los productos incluidos en el pedido.
func (p *Pedido) GetProductos() []interfaces.Comprable {
	return p.productos
}

// La función AgregarProducto incorpora un producto al pedido con detalles agregados.
func (p *Pedido) AgregarProducto(producto interfaces.Comprable) {
	p.productos = append(p.productos, producto)
	p.total += producto.GetPrecio()
}

// CalcularTotal calcula el valor total del pedido
// sumando el precio de todos los productos incluidos.
func (p *Pedido) GetTotal() float64 {
	total := 0.0

	for _, producto := range p.productos {
		total += producto.GetPrecio()
	}

	return total
}
