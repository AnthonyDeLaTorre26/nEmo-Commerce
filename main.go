package main

import (
	"fmt"
	"nemo-commerce/models"
	"nemo-commerce/utils"
)

func mostrarMenu() {
	fmt.Println("\n===== nEmo Commerce =====")
	fmt.Println("1. Registrar cliente")
	fmt.Println("2. Agregar producto")
	fmt.Println("3. Mostrar productos")
	fmt.Println("4. Realizar pedido")
	fmt.Println("5. Mostrar pedidos")
	fmt.Println("6. Mostrar clientes")
	fmt.Println("7. Eliminar producto")
	fmt.Println("8. Eliminar cliente")
	fmt.Println("9. Salir")
}

func main() {
	tienda := models.NuevaTienda()

	for {
		mostrarMenu()
		opcion := utils.LeerEntero("Seleccione una opción: ")

		switch opcion {

		case 1:
			id := utils.LeerEntero("ID del cliente: ")
			nombre := utils.LeerTexto("Nombre: ")
			correo := utils.LeerTexto("Correo: ")

			cliente, err := models.NuevoCliente(id, nombre, correo)
			if err != nil {
				fmt.Println("Error:", err)
				break
			}

			err = tienda.RegistrarCliente(cliente)
			if err != nil {
				fmt.Println("Error:", err)
				break
			}

			fmt.Println("Cliente registrado correctamente.")

		case 2:
			codigo := utils.LeerTexto("Código: ")
			nombre := utils.LeerTexto("Nombre del producto: ")
			precio := utils.LeerDecimal("Precio: ")
			stock := utils.LeerEntero("Stock: ")

			producto, err := models.NuevoProducto(codigo, nombre, precio, stock)
			if err != nil {
				fmt.Println("Error:", err)
				break
			}

			err = tienda.AgregarProducto(producto)
			if err != nil {
				fmt.Println("Error:", err)
				break
			}

			fmt.Println("Producto agregado correctamente.")

		case 3:
			productos := tienda.GetProductos()

			if len(productos) == 0 {
				fmt.Println("No hay productos registrados.")
				break
			}

			fmt.Println("\n===== Catálogo de Productos =====")

			for _, p := range productos {
				fmt.Printf("Código: %s | Nombre: %s | Precio: $%.2f | Stock: %d\n",
					p.GetCodigo(),
					p.GetNombre(),
					p.GetPrecio(),
					p.GetStock())
			}

		case 4:
			idCliente := utils.LeerEntero("ID del cliente: ")
			cliente := tienda.BuscarClientePorID(idCliente)

			if cliente == nil {
				fmt.Println("Cliente no encontrado.")
				break
			}

			codigo := utils.LeerTexto("Código del producto: ")
			producto := tienda.BuscarProductoPorCodigo(codigo)

			if producto == nil {
				fmt.Println("Producto no encontrado.")
				break
			}

			cantidad := utils.LeerEntero("Cantidad: ")

			if producto.GetStock() < cantidad {
				fmt.Println("Stock insuficiente.")
				break
			}

			err := producto.ReducirStock(cantidad)
			if err != nil {
				fmt.Println("Error:", err)
				break
			}

			pedido := models.NuevoPedido(len(tienda.GetPedidos())+1, cliente)

			for i := 0; i < cantidad; i++ {
				pedido.AgregarProducto(producto)
			}

			tienda.AgregarPedido(pedido)

			fmt.Printf("Pedido registrado correctamente. Total: $%.2f\n", pedido.GetTotal())

		case 5:
			pedidos := tienda.GetPedidos()

			if len(pedidos) == 0 {
				fmt.Println("No existen pedidos registrados.")
				break
			}

			fmt.Println("\n===== Pedidos Registrados =====")

			for _, p := range pedidos {
				fmt.Printf("\nPedido #%d\n", p.GetID())
				fmt.Printf("Cliente: %s\n", p.GetCliente().GetNombre())
				fmt.Println("Productos:")

				for _, prod := range p.GetProductos() {
					fmt.Printf(" - %s ($%.2f)\n", prod.GetNombre(), prod.GetPrecio())
				}

				fmt.Printf("Total: $%.2f\n", p.GetTotal())
			}

		case 6:
			clientes := tienda.GetClientes()

			if len(clientes) == 0 {
				fmt.Println("No hay clientes registrados.")
				break
			}

			fmt.Println("\n===== Clientes Registrados =====")

			for _, c := range clientes {
				fmt.Printf("ID: %d | Nombre: %s | Correo: %s\n",
					c.GetID(),
					c.GetNombre(),
					c.GetCorreo())
			}

		case 7:
			codigo := utils.LeerTexto("Código del producto a eliminar: ")
			err := tienda.EliminarProducto(codigo)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Producto eliminado correctamente.")
			}

		case 8:
			id := utils.LeerEntero("ID del cliente a eliminar: ")
			err := tienda.EliminarCliente(id)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Cliente eliminado correctamente.")
			}

		case 9:
			fmt.Println("Saliendo del sistema...")
			return

		default:
			fmt.Println("Opción inválida")
		}
	}
}
