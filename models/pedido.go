package models

import "nemo-commerce/interfaces"

// Pedido representa una compra realizada por un cliente.
type Pedido struct {
	id        int
	cliente   *Cliente
	productos []interfaces.Comprable
	total     float64
}

// La función NuevoPedido crea un nuevo pedido.
func NuevoPedido(id int, cliente *Cliente) *Pedido {
	return &Pedido{
		id:        id,
		cliente:   cliente,
		productos: []interfaces.Comprable{},
		total:     0,
	}
}

// La función AgregarProducto incorpora un producto al pedido con detalles agregados.
func (p *Pedido) AgregarProducto(producto interfaces.Comprable) {
	p.productos = append(p.productos, producto)
	p.total += producto.GetPrecio()
}

func (p *Pedido) GetID() int {
	return p.id
}

func (p *Pedido) GetCliente() *Cliente {
	return p.cliente
}

func (p *Pedido) GetProductos() []interfaces.Comprable {
	return p.productos
}

func (p *Pedido) GetTotal() float64 {
	return p.total
}
