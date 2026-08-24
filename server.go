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

	// Página de productos
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

	// Agregar un producto desde la página web.
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

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
