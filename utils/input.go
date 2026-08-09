package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

// La función LeerTexto solicita y devuelve un texto ingresado por el usuario.
func LeerTexto(mensaje string) string {
	fmt.Print(mensaje)
	texto, _ := reader.ReadString('\n')
	return strings.TrimSpace(texto)
}

// La función LeerEntero solicita un número entero y valida la entrada.
func LeerEntero(mensaje string) int {
	for {
		texto := LeerTexto(mensaje)
		valor, err := strconv.Atoi(texto)

		if err == nil {
			return valor
		}

		fmt.Println("Ingrese un número válido.")
	}
}

// La función LeerDecimal solicita un número decimal y valida la entrada.
func LeerDecimal(mensaje string) float64 {
	for {
		texto := LeerTexto(mensaje)
		valor, err := strconv.ParseFloat(texto, 64)

		if err == nil {
			return valor
		}

		fmt.Println("Ingrese un valor decimal válido.")
	}
}
