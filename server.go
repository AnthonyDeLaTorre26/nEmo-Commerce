package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"nemo-commerce/database"
	"net/http"
	"strconv"
	"sync"
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

	// INICIO

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

	// PRODUCTOS

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

	// AGREGAR PRODUCTO

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
			http.Error(
				w,
				"Error al guardar el producto: "+err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		http.Redirect(w, r, "/productos", http.StatusSeeOther)
	})

	// CLIENTES

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

	// AGREGAR CLIENTE

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
			http.Error(
				w,
				"Error al guardar el cliente: "+err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		http.Redirect(w, r, "/clientes", http.StatusSeeOther)
	})

	// MOSTRAR PEDIDO

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

		// Obtener pedidos
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

	// AGREGAR PEDIDO

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

		http.Redirect(
			w,
			r,
			"/pedidos",
			http.StatusSeeOther,
		)
	})

	// SERVICIOS WEB - JSON

	var pedidoMutex sync.Mutex

	// 1. OBTENER PRODUCTOS
	http.HandleFunc("/api/productos", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(productos)
	})

	// 2. AGREGAR PRODUCTO
	http.HandleFunc("/api/productos/agregar", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		var producto ProductoWeb

		err := json.NewDecoder(r.Body).Decode(&producto)
		if err != nil {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}

		_, err = database.InsertarProducto(
			db,
			producto.Codigo,
			producto.Nombre,
			producto.Precio,
			producto.Stock,
		)

		if err != nil {
			http.Error(
				w,
				"Error al guardar el producto",
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		respuesta := map[string]interface{}{
			"mensaje":  "Producto registrado correctamente",
			"producto": producto,
		}

		json.NewEncoder(w).Encode(respuesta)
	})

	// 3. OBTENER CLIENTES
	http.HandleFunc("/api/clientes", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clientes)
	})

	// 4. AGREGAR CLIENTE
	http.HandleFunc("/api/clientes/agregar", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		var cliente ClienteWeb

		err := json.NewDecoder(r.Body).Decode(&cliente)
		if err != nil {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}

		_, err = database.InsertarCliente(
			db,
			cliente.Nombre,
			cliente.Correo,
		)

		if err != nil {
			http.Error(
				w,
				"Error al guardar el cliente",
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		respuesta := map[string]interface{}{
			"mensaje": "Cliente registrado correctamente",
			"cliente": cliente,
		}

		json.NewEncoder(w).Encode(respuesta)
	})

	// 5. OBTENER PEDIDOS
	http.HandleFunc("/api/pedidos", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

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
			http.Error(w, "Error al consultar los pedidos", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var pedidos []PedidoWeb

		for rows.Next() {

			var pedido PedidoWeb

			err := rows.Scan(
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

			pedidos = append(pedidos, pedido)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pedidos)
	})

	// 6. CREAR PEDIDO
	http.HandleFunc("/api/pedidos/agregar", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		var pedido struct {
			ClienteID  int `json:"cliente_id"`
			ProductoID int `json:"producto_id"`
			Cantidad   int `json:"cantidad"`
		}

		err := json.NewDecoder(r.Body).Decode(&pedido)
		if err != nil {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}

		if pedido.ClienteID <= 0 ||
			pedido.ProductoID <= 0 ||
			pedido.Cantidad <= 0 {

			http.Error(
				w,
				"Datos del pedido inválidos",
				http.StatusBadRequest,
			)
			return
		}

		// Control de concurrencia
		pedidoMutex.Lock()
		defer pedidoMutex.Unlock()

		var precio float64
		var stock int

		err = db.QueryRow(
			"SELECT precio, stock FROM productos WHERE id = ?",
			pedido.ProductoID,
		).Scan(&precio, &stock)

		if err != nil {
			http.Error(
				w,
				"Producto no encontrado",
				http.StatusBadRequest,
			)
			return
		}

		if stock < pedido.Cantidad {
			http.Error(
				w,
				"Stock insuficiente",
				http.StatusBadRequest,
			)
			return
		}

		total := precio * float64(pedido.Cantidad)

		pedidoID, err := database.InsertarPedido(
			db,
			pedido.ClienteID,
			total,
		)

		if err != nil {
			http.Error(
				w,
				"Error al guardar el pedido",
				http.StatusInternalServerError,
			)
			return
		}

		err = database.InsertarDetallePedido(
			db,
			pedidoID,
			pedido.ProductoID,
			pedido.Cantidad,
			precio,
		)

		if err != nil {
			http.Error(
				w,
				"Error al guardar el detalle",
				http.StatusInternalServerError,
			)
			return
		}

		err = database.ActualizarStock(
			db,
			pedido.ProductoID,
			pedido.Cantidad,
		)

		if err != nil {
			http.Error(
				w,
				"Error al actualizar el stock",
				http.StatusInternalServerError,
			)
			return
		}

		respuesta := map[string]interface{}{
			"mensaje":     "Pedido registrado correctamente",
			"pedido_id":   pedidoID,
			"cliente_id":  pedido.ClienteID,
			"producto_id": pedido.ProductoID,
			"cantidad":    pedido.Cantidad,
			"precio":      precio,
			"total":       total,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(respuesta)
	})

	// 7. ELIMINAR PRODUCTO
	http.HandleFunc("/api/productos/eliminar/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodDelete {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		idTexto := r.URL.Path[len("/api/productos/eliminar/"):]

		id, err := strconv.Atoi(idTexto)
		if err != nil {
			http.Error(
				w,
				"ID de producto inválido",
				http.StatusBadRequest,
			)
			return
		}

		resultado, err := db.Exec(
			"DELETE FROM productos WHERE id = ?",
			id,
		)

		if err != nil {
			http.Error(
				w,
				"Error al eliminar el producto",
				http.StatusInternalServerError,
			)
			return
		}

		filas, err := resultado.RowsAffected()
		if err != nil {
			http.Error(
				w,
				"Error al comprobar la eliminación",
				http.StatusInternalServerError,
			)
			return
		}

		if filas == 0 {
			http.Error(
				w,
				"Producto no encontrado",
				http.StatusNotFound,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		respuesta := map[string]interface{}{
			"mensaje":     "Producto eliminado correctamente",
			"producto_id": id,
		}

		json.NewEncoder(w).Encode(respuesta)
	})

	// 8. ELIMINAR CLIENTE
	http.HandleFunc("/api/clientes/eliminar/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodDelete {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		idTexto := r.URL.Path[len("/api/clientes/eliminar/"):]

		id, err := strconv.Atoi(idTexto)
		if err != nil {
			http.Error(
				w,
				"ID de cliente inválido",
				http.StatusBadRequest,
			)
			return
		}

		resultado, err := db.Exec(
			"DELETE FROM clientes WHERE id = ?",
			id,
		)

		if err != nil {
			http.Error(
				w,
				"Error al eliminar el cliente",
				http.StatusInternalServerError,
			)
			return
		}

		filas, err := resultado.RowsAffected()
		if err != nil {
			http.Error(
				w,
				"Error al comprobar la eliminación",
				http.StatusInternalServerError,
			)
			return
		}

		if filas == 0 {
			http.Error(
				w,
				"Cliente no encontrado",
				http.StatusNotFound,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		respuesta := map[string]interface{}{
			"mensaje":    "Cliente eliminado correctamente",
			"cliente_id": id,
		}

		json.NewEncoder(w).Encode(respuesta)
	})

	// INICIAR SERVIDOR

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
