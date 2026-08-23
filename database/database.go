package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Conectar() (*sql.DB, error) {
	db, err := sql.Open(
		"mysql",
		"nemo_app:123456789@@tcp(127.0.0.1:3306)/nemo_commerce",
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func InsertarCliente(db *sql.DB, nombre string, correo string) (int64, error) {
	resultado, err := db.Exec(
		"INSERT INTO clientes (nombre, correo) VALUES (?, ?)",
		nombre,
		correo,
	)

	if err != nil {
		return 0, err
	}

	id, err := resultado.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func InsertarProducto(db *sql.DB, codigo string, nombre string, precio float64, stock int) (int64, error) {
	resultado, err := db.Exec(
		"INSERT INTO productos (codigo, nombre, precio, stock) VALUES (?, ?, ?, ?)",
		codigo,
		nombre,
		precio,
		stock,
	)

	if err != nil {
		return 0, err
	}

	id, err := resultado.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func ObtenerProductos(db *sql.DB) (*sql.Rows, error) {
	return db.Query(
		"SELECT id, codigo, nombre, precio, stock FROM productos",
	)
}

func ObtenerClientes(db *sql.DB) (*sql.Rows, error) {
	return db.Query(
		"SELECT id, nombre, correo FROM clientes",
	)
}

func InsertarPedido(db *sql.DB, clienteID int, total float64) (int64, error) {
	resultado, err := db.Exec(
		"INSERT INTO pedidos (cliente_id, total) VALUES (?, ?)",
		clienteID,
		total,
	)

	if err != nil {
		return 0, err
	}

	id, err := resultado.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func InsertarDetallePedido(
	db *sql.DB,
	pedidoID int64,
	productoID int,
	cantidad int,
	precioUnitario float64,
) error {

	_, err := db.Exec(
		`INSERT INTO detalle_pedidos
		(pedido_id, producto_id, cantidad, precio_unitario)
		VALUES (?, ?, ?, ?)`,
		pedidoID,
		productoID,
		cantidad,
		precioUnitario,
	)

	return err
}
