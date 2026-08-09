package models

import "errors"

// Cliente representa un usuario que sea registrado.
type Cliente struct {
	id     int
	nombre string
	correo string
}

// La función NuevoCliente crea un cliente validando la información.
func NuevoCliente(id int, nombre, correo string) (*Cliente, error) {
	if id <= 0 {
		return nil, errors.New("el ID debe ser mayor que cero")
	}

	if nombre == "" {
		return nil, errors.New("el nombre no puede estar vacío")
	}

	if correo == "" {
		return nil, errors.New("el correo no puede estar vacío")
	}

	return &Cliente{
		id:     id,
		nombre: nombre,
		correo: correo,
	}, nil
}

func (c *Cliente) GetID() int {
	return c.id
}

func (c *Cliente) GetNombre() string {
	return c.nombre
}

func (c *Cliente) GetCorreo() string {
	return c.correo
}

func (c *Cliente) SetCorreo(correo string) error {
	if correo == "" {
		return errors.New("el correo no puede estar vacío")
	}

	c.correo = correo
	return nil
}
