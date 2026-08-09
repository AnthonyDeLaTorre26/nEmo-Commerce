package models

import "errors"

// Tienda va a administrar clientes, productos y pedidos.
type Tienda struct {
	productos []*Producto
	clientes  []*Cliente
	pedidos   []*Pedido
}

// La función NuevaTienda crea una tienda vacía.
func NuevaTienda() *Tienda {
	return &Tienda{
		productos: []*Producto{},
		clientes:  []*Cliente{},
		pedidos:   []*Pedido{},
	}
}

// La función AgregarProducto registra un producto evitando códigos duplicados.
func (t *Tienda) AgregarProducto(p *Producto) error {
	if t.BuscarProductoPorCodigo(p.GetCodigo()) != nil {
		return errors.New("Ya existe un producto con ese código")
	}

	t.productos = append(t.productos, p)
	return nil
}

// La función RegistrarCliente registra un cliente evitando IDs duplicados.
func (t *Tienda) RegistrarCliente(c *Cliente) error {
	if t.BuscarClientePorID(c.GetID()) != nil {
		return errors.New("Ya existe un cliente con ese ID")
	}

	t.clientes = append(t.clientes, c)
	return nil
}

// La función AgregarPedido almacena un pedido.
func (t *Tienda) AgregarPedido(p *Pedido) {
	t.pedidos = append(t.pedidos, p)
}

// La función GetProductos devuelve todos los productos.
func (t *Tienda) GetProductos() []*Producto {
	return t.productos
}

// La función GetClientes devuelve todos los clientes.
func (t *Tienda) GetClientes() []*Cliente {
	return t.clientes
}

// La función GetPedidos devuelve todos los pedidos.
func (t *Tienda) GetPedidos() []*Pedido {
	return t.pedidos
}

// La función BuscarClientePorID busca un cliente por su ID.
func (t *Tienda) BuscarClientePorID(id int) *Cliente {
	for _, c := range t.clientes {
		if c.GetID() == id {
			return c
		}
	}
	return nil
}

// La función BuscarProductoPorCodigo busca un producto por su código.
func (t *Tienda) BuscarProductoPorCodigo(codigo string) *Producto {
	for _, p := range t.productos {
		if p.GetCodigo() == codigo {
			return p
		}
	}
	return nil
}

// La función EliminarProducto elimina un producto por su código.
func (t *Tienda) EliminarProducto(codigo string) error {
	for i, p := range t.productos {
		if p.GetCodigo() == codigo {
			t.productos = append(t.productos[:i], t.productos[i+1:]...)
			return nil
		}
	}
	return errors.New("Producto no encontrado")
}

// La función EliminarCliente elimina un cliente por el ID.
func (t *Tienda) EliminarCliente(id int) error {
	for i, c := range t.clientes {
		if c.GetID() == id {
			t.clientes = append(t.clientes[:i], t.clientes[i+1:]...)
			return nil
		}
	}
	return errors.New("Cliente no encontrado")
}
