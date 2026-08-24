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

func iniciarServidorWeb(db *sql.DB) {

	// Página principal
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

	// Mostrar productos
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

	// Agregar producto
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

	// Mostrar clientes
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

	// Agregar cliente
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

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
