package models

import "errors"

// Cliente representa un usuario que sea registrado y aplica encapsulación con atributos privados.
type Cliente struct {
	id     int
	nombre string
	correo string
}

// La función NuevoCliente crea un cliente validando la datos básicos antes de registrar.
func NuevoCliente(id int, nombre, correo string) (*Cliente, error) {
	if id <= 0 {
		return nil, errors.New("El ID debe ser mayor que cero")
	}

	if nombre == "" {
		return nil, errors.New("El nombre no puede estar vacío")
	}

	if correo == "" {
		return nil, errors.New("El correo no puede estar vacío")
	}

	return &Cliente{
		id:     id,
		nombre: nombre,
		correo: correo,
	}, nil
}

// GetID devuelve la identificación del cliente.
func (c *Cliente) GetID() int {
	return c.id
}

// GetNombre devuelve el nombre del cliente.
func (c *Cliente) GetNombre() string {
	return c.nombre
}

// GetCorreo devuelve el correo del cliente.
func (c *Cliente) GetCorreo() string {
	return c.correo
}

// SetCorreo permite modificar el correo del cliente validando que el nuevo correo no este vacío.
func (c *Cliente) SetCorreo(correo string) error {
	if correo == "" {
		return errors.New("El correo no puede estar vacío")
	}

	c.correo = correo
	return nil
}
