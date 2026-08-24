# nEmo-Commerce
Sistema de Gestión de Comercio Electrónico desarrollado en Go (Goolang).

## Descripción

nEmo Commerce es un sistema de gestión de comercio electrónico orientado a empresas pequeñas y  medianas. Su objetivo es facilitar la administraciñon de los productos, clientes, inventario y pedidos a través de una aplicación desarrollada en Go (Golang).
El sistema cuenta con una aplicación de consola y una interfaz web conectada a una base de datos MySQL.

## Objetivo

Desarrollar un sistema funcional que nos permita aplicar conceptos de Programación Orientada a Objetos y estructuras de datos para resolver problemas con gestionamientos de un comercio electrónico.

## Tecnolgías

- Go (Golang) - lenguaje principal de programación.
- Git 
- GitHub - Control de versiones y almacenamiento del proyecto.
- MySQL - Sistema de gestión de base de datos.
- HTML - Estructura de las páginas web.
- Bootstrap - Diseño y estilo de la interfaz web.

## Funcionalidades
Productos

Permite:
- Registrar nuevos productos.
- Consultar el catálogo.
- Mostrar código, nombre, precio y stock.
- Actualizar el stock al realizar un pedido.

Clientes

Permite:
- Registrar clientes.
- Consultar los clientes registrados.
- Almacenar nombre y correo electrónico.

Pedidos

Permite:
- Seleccionar un cliente.
- Seleccionar un producto.
- Indicar la cantidad.
- Verificar el stock disponible.
- Registrar el pedido.
- Registrar el detalle del pedido.
- Calcular el total.
- Descontar automáticamente la cantidad solicitada del stock.
## Encapsulación

Los atributos principales en las clases "Cliente" y "Producto" se mantendran privados siendo solo accesibles mediante métodos.

## Interfaces

La interfaz "Comprable" va a establecer los comportamientos que debe cumplir un elemento registrado para la realización de un pedido.

## Interfaz web

El proyecto incorpora una interfaz web accesible desde:

http://localhost:8080

La interfaz contiene las siguientes secciones:

- Inicio
- Productos
- Clientes
- Pedidos

La página utiliza Bootstrap para mejorar la presentación visual y facilitar la navegación entre las diferentes funcionalidades del sistema.

## Composición

La clase "Pedido" va a mantener una relacion con un "Cliente" y con los elementos comprables en la compra.

## Funcionamiento

Al iniciar el programa se muestra el menú principal de la aplicación:

===== nEmo Commerce =====
1. Registrar cliente
2. Agregar producto
3. Mostrar productos
4. Realizar pedido
5. Mostrar pedidos
6. Mostrar clientes
7. Eliminar producto
8. Eliminar cliente
9. Salir

Además, el sistema permite acceder a las mismas funcionalidades principales mediante la interfaz web.

## Manejo de errores

El sistema utiliza errores para poder controlar situaciones tales como:

- Datos inválidos.
- Clientes duplicados.
- Productos duplicados.
- Productos inexistentes.
- Clientes inexistentes.
- Stock insuficiente.

## Gestión de stock

Cuando se realiza un pedido, el sistema verifica que exista suficiente stock del producto.

Si la cantidad solicitada supera el stock disponible, el pedido no se registra y se muestra un mensaje indicando:

- Stock insuficiente

Cuando el pedido se registra correctamente, el stock del producto se actualiza automáticamente.

## Estado del Proyecto

Version Funcional

La version actual ha implementado las principales funcionalidades de gestión, aplicando consola y manteniendo los datos en la base de datos MySQL, actualizando automáticamente la tienda completa incluyendo calculos y validaciones.
El proyecto esta sincronizado con GitHub y con sus funcionalidades proncipales probadas.