package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// Estructura para representar un producto
type Producto struct {
	Nombre   string
	Cantidad int
	Precio   float64
}

func main() {
	// Creamos la variable entradaConEspacios usando la libreria "bufio" y "os" para ingresar el nombre del cliente y del producto con espacios
	entradaConEspacios := bufio.NewReader(os.Stdin)

	// Solicitamos el nombre del cliente
	fmt.Println("Ingrese el nombre del cliente:")
	nombreCliente, _ := entradaConEspacios.ReadString('\n')
	nombreCliente = strings.TrimSpace(nombreCliente)

	// Obtenemos la fecha actual para la factura
	fechaActual := time.Now().Format("02/01/2006")

	// Creamos una lista (slice) para guardar los productos
	listaProductos := []Producto{}

	// Ciclo para pedir información de los productos
	for {
		// Pedimos el nombre del producto
		fmt.Println("Ingrese el nombre del producto (o 'f' para terminar):")
		nombreProducto, _ := entradaConEspacios.ReadString('\n')
		nombreProducto = strings.TrimSpace(nombreProducto)
		if nombreProducto == "f" {
			break
		}

		// Pedimos la cantidad del producto
		var cantidadProducto int
		fmt.Println("Ingrese la cantidad:")
		fmt.Scanln(&cantidadProducto)

		// Pedimos el precio del producto
		var precioProducto float64
		fmt.Println("Ingrese el precio unitario:")
		fmt.Scanln(&precioProducto)

		// Agregamos el producto a la lista
		listaProductos = append(listaProductos, Producto{Nombre: nombreProducto, Cantidad: cantidadProducto, Precio: precioProducto})
	}
	// Calculamos el valor total y mostramos el resumen de la factura
	var totalFactura float64
	fmt.Println("\nResumen de la factura:")
	fmt.Printf("Cliente: %s\n", nombreCliente)
	fmt.Printf("Fecha: %s\n", fechaActual)
	fmt.Println("Productos:")

	// Recorremos la lista de productos y calculamos el total de cada uno
	for _, producto := range listaProductos {
		totalProducto := float64(producto.Cantidad) * producto.Precio
		totalFactura += totalProducto
		fmt.Printf("- %s: %d x %.2f = %.2f\n", producto.Nombre, producto.Cantidad, producto.Precio, totalProducto)
	}

	// Mostramos el valor total
	fmt.Printf("Valor Total: %.2f\n", totalFactura)

	// Preguntamos si el usuario desea generar la factura en PDF
	var respuesta string
	fmt.Println("¿Quiere imprimir la factura? (si/no):")
	fmt.Scanln(&respuesta)

	// Si el usuario responde 'si', generamos la factura en PDF
	if strings.ToLower(respuesta) == "si" {
		// Manejamos errores al generar el pdf
		error := FacturaPDF(nombreCliente, fechaActual, listaProductos, totalFactura)
		if error != nil {
			fmt.Println("Error al generar el PDF: ", error)
		} else {
			fmt.Println("Ver.pdf")
		}
	} else {
		fmt.Println("No se ha generado la factura en PDF.")
	}
}

// Función para generar la factura en formato PDF
func FacturaPDF(cliente, fechaActual string, productos []Producto, total float64) error {
	// Creamos un nuevo archivo PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Times", "B", 14)

	// Añadimos una campo con el título "Factura" y el tamaño
	pdf.Cell(200, 10, "Resumen de la Factura")
	// salto de 12 unidades
	pdf.Ln(12)

	// Mostramos el nombre del cliente y la fecha
	pdf.SetFont("Times", "", 12)
	pdf.CellFormat(95, 10, "Cliente: "+cliente, "", 0, "L", false, 0, "")
	pdf.CellFormat(95, 10, "Fecha: "+fechaActual, "", 1, "R", false, 0, "")
	pdf.Ln(5)

	// Creamos una tabla con los productos (los campos del encabezado)
	pdf.SetFont("Times", "B", 12)
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(65, 10, "Producto", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 10, "Cantidad", "1", 0, "C", true, 0, "")
	pdf.CellFormat(45, 10, "Precio Unitario", "1", 0, "C", true, 0, "")
	pdf.CellFormat(45, 10, "Total", "1", 1, "C", true, 0, "")

	// Llenamos la tabla con los datos de los productos
	pdf.SetFont("Times", "", 12)
	// itera sobe los productos, calculamos el total por cada producto y los creamos una fila para cada producto
	for _, producto := range productos {
		totalProducto := float64(producto.Cantidad) * producto.Precio
		pdf.CellFormat(65, 10, producto.Nombre, "1", 0, "L", false, 0, "")
		pdf.CellFormat(35, 10, fmt.Sprintf("%d", producto.Cantidad), "1", 0, "C", false, 0, "")
		pdf.CellFormat(45, 10, fmt.Sprintf("%.2f", producto.Precio), "1", 0, "R", false, 0, "")
		pdf.CellFormat(45, 10, fmt.Sprintf("%.2f", totalProducto), "1", 1, "R", false, 0, "")
	}

	// Mostramos el valor total de la factura
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(145, 10, "Valor Total:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(45, 10, fmt.Sprintf("%.2f", total), "1", 1, "R", false, 0, "")

	// Guardamos el PDF en un archivo
	return pdf.OutputFileAndClose("Ver.pdf")
}
