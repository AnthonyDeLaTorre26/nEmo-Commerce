package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func Conectar() (*sql.DB, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	usuario := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	puerto := os.Getenv("DB_PORT")
	baseDatos := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		usuario,
		password,
		host,
		puerto,
		baseDatos,
	)

	db, err := sql.Open("mysql", dsn)
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

func ActualizarStock(db *sql.DB, productoID int, cantidad int) error {
	_, err := db.Exec(
		"UPDATE productos SET stock = stock - ? WHERE id = ?",
		cantidad,
		productoID,
	)

	return err
}

func EliminarProducto(db *sql.DB, codigo string) error {
	_, err := db.Exec(
		"DELETE FROM productos WHERE codigo = ?",
		codigo,
	)
	return err
}

func EliminarCliente(db *sql.DB, id int) error {
	_, err := db.Exec(
		"DELETE FROM clientes WHERE id = ?",
		id,
	)
	return err
}

type ClienteDB struct {
	ID     int
	Nombre string
	Correo string
}
