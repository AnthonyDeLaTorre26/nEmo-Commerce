package main

import (
	"fmt"
	"log"

	"nemo-commerce/database"
	"nemo-commerce/models"
	"nemo-commerce/utils"
)

// mostrarMenu presenta las opciones disponibles parala interfaz del usuario.
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

	// Conexión con la base de datos MySQL.
	db, err := database.Conectar()
	if err != nil {
		log.Fatal("Error al conectar con MySQL:", err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("No se pudo establecer conexión con MySQL:", err)
	}

	fmt.Println("Conexión exitosa con MySQL.")

	// Se crea una instancia de Tienda
	tienda := models.NuevaTienda()

	// Cargar los clientes existentes desde MySQL.
	rows, err := database.ObtenerClientes(db)
	if err != nil {
		log.Fatal("Error al cargar los clientes:", err)
	}

	for rows.Next() {
		var (
			id     int
			nombre string
			correo string
		)

		err := rows.Scan(&id, &nombre, &correo)
		if err != nil {
			log.Fatal("Error al leer un cliente:", err)
		}

		cliente, err := models.NuevoCliente(id, nombre, correo)
		if err != nil {
			log.Fatal("Error al crear el cliente:", err)
		}

		err = tienda.RegistrarCliente(cliente)
		if err != nil {
			log.Fatal("Error al registrar el cliente en la tienda:", err)
		}
	}

	rows.Close()

	// Iniciar el servidor web en segundo plano.
	go iniciarServidorWeb(db)

	fmt.Println("Servidor web iniciado en http://localhost:8080")

	// Se crea el ciclo principal manteniendo el sistema activo hasta el cierre de sesión.
	for {
		mostrarMenu()
		opcion := utils.LeerEntero("Seleccione una opción: ")

		switch opcion {

		// Registrar un nuevo cliente.
		case 1:
			nombre := utils.LeerTexto("Nombre: ")
			correo := utils.LeerTexto("Correo: ")

			// Guardar el cliente en MySQL.
			id, err := database.InsertarCliente(db, nombre, correo)
			if err != nil {
				fmt.Println("Error al guardar el cliente en MySQL:", err)
				break
			}

			// Crear el objeto Cliente con el ID generado por MySQL.
			cliente, err := models.NuevoCliente(int(id), nombre, correo)
			if err != nil {
				fmt.Println("Error:", err)
				break
			}

			// Agregar el cliente a la tienda en memoria.
			err = tienda.RegistrarCliente(cliente)
			if err != nil {
				fmt.Println("Error:", err)
				break
			}

			fmt.Printf("Cliente registrado correctamente. ID asignado: %d\n", id)

		// Agrega un nuevo producto al catálogo.
		case 2:
			codigo := utils.LeerTexto("Código: ")
			nombre := utils.LeerTexto("Nombre del producto: ")
			precio := utils.LeerDecimal("Precio: ")
			stock := utils.LeerEntero("Stock: ")

			// Crear el objeto Producto utilizando el modelo.
			producto, err := models.NuevoProducto(codigo, nombre, precio, stock)
			if err != nil {
				fmt.Println("Error:", err)
				break
			}

			// Guardar el producto en MySQL.
			id, err := database.InsertarProducto(db, codigo, nombre, precio, stock)
			if err != nil {
				fmt.Println("Error al guardar el producto en MySQL:", err)
				break
			}

			// Asignar al objeto el ID generado por MySQL.
			producto.SetID(int(id))

			// Agregar el producto a la tienda.
			err = tienda.AgregarProducto(producto)
			if err != nil {
				fmt.Println("Error:", err)
				break
			}

			fmt.Printf("Producto agregado correctamente. ID asignado: %d\n", id)

		// Muestra los productos disponibles.
		case 3:
			rows, err := database.ObtenerProductos(db)
			if err != nil {
				fmt.Println("Error al consultar los productos:", err)
				break
			}

			defer rows.Close()

			fmt.Println("\n===== Catálogo de Productos =====")

			hayProductos := false

			for rows.Next() {
				var (
					id     int
					codigo string
					nombre string
					precio float64
					stock  int
				)

				err := rows.Scan(&id, &codigo, &nombre, &precio, &stock)
				if err != nil {
					fmt.Println("Error al leer el producto:", err)
					break
				}

				hayProductos = true

				fmt.Printf(
					"ID: %d | Código: %s | Nombre: %s | Precio: $%.2f | Stock: %d\n",
					id,
					codigo,
					nombre,
					precio,
					stock,
				)
			}

			if !hayProductos {
				fmt.Println("No hay productos registrados.")
			}

		// Realización de un nuevo pedido.
		case 4:
			// Realización de un nuevo pedido.
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

			if cantidad <= 0 {
				fmt.Println("La cantidad debe ser mayor que cero.")
				break
			}

			if producto.GetStock() < cantidad {
				fmt.Println("Stock insuficiente.")
				break
			}

			total := producto.GetPrecio() * float64(cantidad)

			pedidoID, err := database.InsertarPedido(
				db,
				cliente.GetID(),
				total,
			)

			if err != nil {
				fmt.Println("Error al guardar el pedido:", err)
				break
			}

			err = database.InsertarDetallePedido(
				db,
				pedidoID,
				producto.GetID(),
				cantidad,
				producto.GetPrecio(),
			)

			if err != nil {
				fmt.Println("Error al guardar el detalle del pedido:", err)
				break
			}

			// Actualizar el stock directamente en MySQL.
			err = database.ActualizarStock(db, producto.GetID(), cantidad)
			if err != nil {
				fmt.Println("Error al actualizar el stock en MySQL:", err)
				break
			}

			// Actualizar también el stock del objeto en memoria.
			err = producto.ReducirStock(cantidad)
			if err != nil {
				fmt.Println("Error al actualizar el stock:", err)
				break
			}

			// Crear el objeto Pedido.
			pedido := models.NuevoPedido(int(pedidoID), cliente)

			// Agregar al pedido la cantidad de productos comprados.
			for i := 0; i < cantidad; i++ {
				pedido.AgregarProducto(producto)
			}

			// Guardar el pedido también en la tienda.
			tienda.AgregarPedido(pedido)

			fmt.Printf(
				"Pedido registrado correctamente. ID: %d | Total: $%.2f\n",
				pedidoID,
				total,
			)

		//Muestra todos loss pedidos registrados.
		case 5:

			rows, err := db.Query(`
		SELECT
			p.id,
			c.nombre,
			pr.nombre,
			dp.cantidad,
			dp.precio_unitario,
			(dp.cantidad * dp.precio_unitario) AS total
		FROM pedidos p
		INNER JOIN clientes c ON p.cliente_id = c.id
		INNER JOIN detalle_pedidos dp ON p.id = dp.pedido_id
		INNER JOIN productos pr ON dp.producto_id = pr.id
		ORDER BY p.id DESC
	`)

			if err != nil {
				fmt.Println("Error al consultar los pedidos:", err)
				break
			}

			fmt.Println("\n===== Pedidos Registrados =====")

			hayPedidos := false

			for rows.Next() {

				var (
					id       int
					cliente  string
					producto string
					cantidad int
					precio   float64
					total    float64
				)

				err := rows.Scan(
					&id,
					&cliente,
					&producto,
					&cantidad,
					&precio,
					&total,
				)

				if err != nil {
					fmt.Println("Error al leer el pedido:", err)
					break
				}

				hayPedidos = true

				fmt.Printf(
					"Pedido #%d | Cliente: %s | Producto: %s | Cantidad: %d | Precio: $%.2f | Total: $%.2f\n",
					id,
					cliente,
					producto,
					cantidad,
					precio,
					total,
				)
			}

			rows.Close()

			if !hayPedidos {
				fmt.Println("No existen pedidos registrados.")
			}

		// Muestra todos los clientes registrados.
		case 6:
			rows, err := database.ObtenerClientes(db)
			if err != nil {
				fmt.Println("Error al consultar los clientes:", err)
				break
			}

			defer rows.Close()

			fmt.Println("\n===== Clientes Registrados =====")

			hayClientes := false

			for rows.Next() {
				var (
					id     int
					nombre string
					correo string
				)

				err := rows.Scan(&id, &nombre, &correo)
				if err != nil {
					fmt.Println("Error al leer el cliente:", err)
					break
				}

				hayClientes = true

				fmt.Printf(
					"ID: %d | Nombre: %s | Correo: %s\n",
					id,
					nombre,
					correo,
				)
			}

			if !hayClientes {
				fmt.Println("No hay clientes registrados.")
			}

		// Elimina productos del catálogo.
		case 7:
			codigo := utils.LeerTexto("Código del producto a eliminar: ")

			err := database.EliminarProducto(db, codigo)

			if err != nil {
				fmt.Println("Error al eliminar el producto de MySQL:", err)
				break
			}

			err = tienda.EliminarProducto(codigo)

			if err != nil {
				fmt.Println("Error al eliminar el producto de la tienda:", err)
				break
			}

			fmt.Println("Producto eliminado correctamente.")

		// Elimina clientes registrados.
		case 8:
			id := utils.LeerEntero("ID del cliente a eliminar: ")

			err := database.EliminarCliente(db, id)

			if err != nil {
				fmt.Println("Error al eliminar el cliente de MySQL:", err)
				break
			}

			err = tienda.EliminarCliente(id)

			if err != nil {
				fmt.Println("Error al eliminar el cliente de la tienda:", err)
				break
			}

			fmt.Println("Cliente eliminado correctamente.")

		// Finaliza la sesión cerrando el programa.
		case 9:
			fmt.Println("Saliendo del sistema...")
			return

		// Evita opciones que no existan en el menú.
		default:
			fmt.Println("Opción inválida")
		}
	}
}
