# nEmo-Commerce
Sistema de Gestión de Comercio Electrónico desarrollado en Go (Goolang).

## Descripción

nEmo Commerce es un sistema de gestión de comercio electrónico orientado a empresas pequeñas y  medianas. Su objetivo es facilitar la administraciñon de los productos, clientes, inventario y pedidos a través de una aplicación desarrollada en Go (Golang).

## Objetivo

Desarrollar un sistema funcional que nos permita aplicar conceptos de Programación Orientada a Objetos y estructuras de datos para resolver problemas con gestionamientos de un comercio electrónico.

## Tecnolgías

- Go (Golang)
- Git
- GitHub

## Funcionalidades

El sistema se ha configurado para permitirnos:

- Registro de clientes.
- Validar clientes duplicados.
- Mostrar clientes.
- Eliminar clientes.
- Registrar productos.
- Validar productos.
- Mostrar catálogo.
- Controlar stock.
- Realizar pedidos.
- Calcular total de pedidos
- Actualizar inventario.
- Mostrar pedidos registrados.
- Eliminar productos.
- Manejar entradas inválidas.

## Encapsulación

Los atributos principales en las clases "Cliente" y "Producto" se mantendran privados siendo solo accesibles mediante métodos.

## Interfaces

La interfaz "Comprable" va a establecer los comportamientos que debe cumplir un elemento registrado para la realización de un pedido.

## Composición

La clase "Pedido" va a mantener una relacion con un "Cliente" y con los elementos comprables en la compra.

## Manejo de errores

El sistema utiliza errores para poder controlar situaciones tales como:

- Datos inválidos.
- Clientes duplicados.
- Productos duplicados.
- Productos inexistentes.
- Clientes inexistentes.
- Stock insuficiente.

## Estado del Proyecto

Version Funcional

La version actual ha implementado las principales funcionalidades de gestión, aplicando consola y manteniendo los datos en memoria al ejecturase.
