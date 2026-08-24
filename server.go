package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"nemo-commerce/database"
	"net/http"
)

type ProductoWeb struct {
	ID     int
	Codigo string
	Nombre string
	Precio float64
	Stock  int
}

type ClienteWeb struct {
	ID     int
	Nombre string
	Correo string
}

type PedidoWeb struct {
	ID       int
	Cliente  string
	Producto string
	Cantidad int
	Precio   float64
	Total    float64
}

type PedidoPagina struct {
	Clientes  []ClienteWeb
	Productos []ProductoWeb
	Pedidos   []PedidoWeb
}

func iniciarServidorWeb(db *sql.DB) {

	// =========================
	// INICIO
	// =========================
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("web/templates/index.html")
		if err != nil {
			http.Error(w, "Error al cargar la página", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error al mostrar la página", http.StatusInternalServerError)
			return
		}
	})

	// =========================
	// PRODUCTOS
	// =========================
	http.HandleFunc("/productos", func(w http.ResponseWriter, r *http.Request) {

		rows, err := database.ObtenerProductos(db)
		if err != nil {
			http.Error(w, "Error al consultar los productos", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var productos []ProductoWeb

		for rows.Next() {
			var producto ProductoWeb

			err := rows.Scan(
				&producto.ID,
				&producto.Codigo,
				&producto.Nombre,
				&producto.Precio,
				&producto.Stock,
			)

			if err != nil {
				http.Error(w, "Error al leer los productos", http.StatusInternalServerError)
				return
			}

			productos = append(productos, producto)
		}

		tmpl, err := template.ParseFiles("web/templates/productos.html")
		if err != nil {
			http.Error(w, "Error al cargar la página", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, productos)
		if err != nil {
			http.Error(w, "Error al mostrar los productos", http.StatusInternalServerError)
			return
		}
	})

	// =========================
	// AGREGAR PRODUCTO
	// =========================
	http.HandleFunc("/productos/agregar", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		codigo := r.FormValue("codigo")
		nombre := r.FormValue("nombre")
		precioTexto := r.FormValue("precio")
		stockTexto := r.FormValue("stock")

		var precio float64
		var stock int

		_, err := fmt.Sscanf(precioTexto, "%f", &precio)
		if err != nil {
			http.Error(w, "Precio inválido", http.StatusBadRequest)
			return
		}

		_, err = fmt.Sscanf(stockTexto, "%d", &stock)
		if err != nil {
			http.Error(w, "Stock inválido", http.StatusBadRequest)
			return
		}

		_, err = database.InsertarProducto(
			db,
			codigo,
			nombre,
			precio,
			stock,
		)

		if err != nil {
			http.Error(w, "Error al guardar el producto: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/productos", http.StatusSeeOther)
	})

	// =========================
	// CLIENTES
	// =========================
	http.HandleFunc("/clientes", func(w http.ResponseWriter, r *http.Request) {

		rows, err := database.ObtenerClientes(db)
		if err != nil {
			http.Error(w, "Error al consultar los clientes", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var clientes []ClienteWeb

		for rows.Next() {
			var cliente ClienteWeb

			err := rows.Scan(
				&cliente.ID,
				&cliente.Nombre,
				&cliente.Correo,
			)

			if err != nil {
				http.Error(w, "Error al leer los clientes", http.StatusInternalServerError)
				return
			}

			clientes = append(clientes, cliente)
		}

		tmpl, err := template.ParseFiles("web/templates/clientes.html")
		if err != nil {
			http.Error(w, "Error al cargar la página", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, clientes)
		if err != nil {
			http.Error(w, "Error al mostrar los clientes", http.StatusInternalServerError)
			return
		}
	})

	// =========================
	// AGREGAR CLIENTE
	// =========================
	http.HandleFunc("/clientes/agregar", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		_, err := database.InsertarCliente(
			db,
			nombre,
			correo,
		)

		if err != nil {
			http.Error(w, "Error al guardar el cliente: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/clientes", http.StatusSeeOther)
	})

	// =========================
	// MOSTRAR PEDIDOS
	// =========================
	http.HandleFunc("/pedidos", func(w http.ResponseWriter, r *http.Request) {

		var data PedidoPagina

		// Obtener clientes
		rowsClientes, err := database.ObtenerClientes(db)
		if err != nil {
			http.Error(w, "Error al consultar los clientes", http.StatusInternalServerError)
			return
		}

		for rowsClientes.Next() {

			var cliente ClienteWeb

			err := rowsClientes.Scan(
				&cliente.ID,
				&cliente.Nombre,
				&cliente.Correo,
			)

			if err != nil {
				rowsClientes.Close()
				http.Error(w, "Error al leer los clientes", http.StatusInternalServerError)
				return
			}

			data.Clientes = append(data.Clientes, cliente)
		}

		rowsClientes.Close()

		// Obtener productos
		rowsProductos, err := database.ObtenerProductos(db)
		if err != nil {
			http.Error(w, "Error al consultar los productos", http.StatusInternalServerError)
			return
		}

		for rowsProductos.Next() {

			var producto ProductoWeb

			err := rowsProductos.Scan(
				&producto.ID,
				&producto.Codigo,
				&producto.Nombre,
				&producto.Precio,
				&producto.Stock,
			)

			if err != nil {
				rowsProductos.Close()
				http.Error(w, "Error al leer los productos", http.StatusInternalServerError)
				return
			}

			data.Productos = append(data.Productos, producto)
		}

		rowsProductos.Close()

		// Obtener pedidos y detalles
		rowsPedidos, err := db.Query(`
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
			http.Error(w, "Error al consultar los pedidos", http.StatusInternalServerError)
			return
		}

		defer rowsPedidos.Close()

		for rowsPedidos.Next() {

			var pedido PedidoWeb

			err := rowsPedidos.Scan(
				&pedido.ID,
				&pedido.Cliente,
				&pedido.Producto,
				&pedido.Cantidad,
				&pedido.Precio,
				&pedido.Total,
			)

			if err != nil {
				http.Error(w, "Error al leer los pedidos", http.StatusInternalServerError)
				return
			}

			data.Pedidos = append(data.Pedidos, pedido)
		}

		// Cargar plantilla
		tmpl, err := template.ParseFiles("web/templates/pedidos.html")
		if err != nil {
			http.Error(w, "Error al cargar la página", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error al mostrar los pedidos", http.StatusInternalServerError)
			return
		}
	})

	// =========================
	// AGREGAR PEDIDO
	// =========================
	http.HandleFunc("/pedidos/agregar", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		var clienteID int
		var productoID int
		var cantidad int

		_, err := fmt.Sscanf(
			r.FormValue("cliente_id"),
			"%d",
			&clienteID,
		)

		if err != nil {
			http.Error(w, "Cliente inválido", http.StatusBadRequest)
			return
		}

		_, err = fmt.Sscanf(
			r.FormValue("producto_id"),
			"%d",
			&productoID,
		)

		if err != nil {
			http.Error(w, "Producto inválido", http.StatusBadRequest)
			return
		}

		_, err = fmt.Sscanf(
			r.FormValue("cantidad"),
			"%d",
			&cantidad,
		)

		if err != nil || cantidad <= 0 {
			http.Error(w, "Cantidad inválida", http.StatusBadRequest)
			return
		}

		// Obtener precio y stock actual
		var precio float64
		var stock int

		err = db.QueryRow(
			"SELECT precio, stock FROM productos WHERE id = ?",
			productoID,
		).Scan(&precio, &stock)

		if err != nil {
			http.Error(w, "Producto no encontrado", http.StatusBadRequest)
			return
		}

		if stock < cantidad {
			http.Error(w, "Stock insuficiente", http.StatusBadRequest)
			return
		}

		total := precio * float64(cantidad)

		// Insertar pedido
		pedidoID, err := database.InsertarPedido(
			db,
			clienteID,
			total,
		)

		if err != nil {
			http.Error(
				w,
				"Error al guardar el pedido: "+err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		// Insertar detalle
		err = database.InsertarDetallePedido(
			db,
			pedidoID,
			productoID,
			cantidad,
			precio,
		)

		if err != nil {
			http.Error(
				w,
				"Error al guardar el detalle: "+err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		// Actualizar stock
		err = database.ActualizarStock(
			db,
			productoID,
			cantidad,
		)

		if err != nil {
			http.Error(
				w,
				"Error al actualizar el stock: "+err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		// Volver a la página de pedidos
		http.Redirect(
			w,
			r,
			"/pedidos",
			http.StatusSeeOther,
		)
	})

	fmt.Println("Servidor web iniciado en http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
