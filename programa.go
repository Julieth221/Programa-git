package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
	// Usamos un scanner para leer la entrada del usuario
	lector := bufio.NewScanner(os.Stdin)

	// Solicitamos el nombre del cliente
	fmt.Println("Ingrese el nombre del cliente:")
	lector.Scan()
	nombreCliente := lector.Text()

	// Obtenemos la fecha actual para la factura
	fechaActual := time.Now().Format("02/01/2006")

	// Creamos una lista (slice) para guardar los productos
	var listaProductos []Producto

	// Ciclo para pedir información de los productos
	for {
		// Pedimos el nombre del producto
		fmt.Println("Ingrese el nombre del producto (o 'f' para terminar):")
		lector.Scan()
		nombreProducto := lector.Text()

		// Si el nombre del producto es 'f', salimos del ciclo
		if nombreProducto == "f" {
			break
		}

		// Pedimos la cantidad del producto
		fmt.Println("Ingrese la cantidad:")
		lector.Scan()
		cantidadProducto, _ := strconv.Atoi(lector.Text())

		// Pedimos el precio del producto
		fmt.Println("Ingrese el precio unitario:")
		lector.Scan()
		precioProducto, _ := strconv.ParseFloat(lector.Text(), 64)

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
	fmt.Println("¿Quiere imprimir la factura? (si/no):")
	lector.Scan()
	decision := strings.ToLower(strings.TrimSpace(lector.Text()))

	// Si el usuario responde 'si', generamos la factura en PDF
	if decision == "si" {
		// Manejamos errores al generar el pdf
		error := generarFacturaPDF(nombreCliente, fechaActual, listaProductos, totalFactura)
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
func generarFacturaPDF(cliente, fecha string, productos []Producto, total float64) error {
	// Creamos un nuevo archivo PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)

	// Añadimos una campo con el título "Factura" y el tamaño
	pdf.Cell(200, 10, "Factura")
	// salto de 15 unidades
	pdf.Ln(15)

	// Mostramos el nombre del cliente y la fecha
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(200, 10, "Cliente: "+cliente)
	pdf.Ln(10)
	pdf.Cell(200, 10, "Fecha: "+fecha)
	pdf.Ln(15)

	// Creamos una tabla con los productos (los campos del encabezado)
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(60, 10, "Producto", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Cantidad", "1", 0, "C", true, 0, "")
	pdf.CellFormat(45, 10, "Precio Unitario", "1", 0, "C", true, 0, "")
	pdf.CellFormat(45, 10, "Total", "1", 1, "C", true, 0, "")

	// Llenamos la tabla con los datos de los productos
	pdf.SetFont("Arial", "", 12)
	// itera sobe los productos, calculamos el total por producto y los creamos una fila para cada producto
	for _, producto := range productos {
		totalProducto := float64(producto.Cantidad) * producto.Precio
		pdf.CellFormat(60, 10, producto.Nombre, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 10, strconv.Itoa(producto.Cantidad), "1", 0, "C", false, 0, "")
		pdf.CellFormat(45, 10, fmt.Sprintf("%.2f", producto.Precio), "1", 0, "R", false, 0, "")
		pdf.CellFormat(45, 10, fmt.Sprintf("%.2f", totalProducto), "1", 1, "R", false, 0, "")
	}

	// Mostramos el valor total
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(135, 10, "Valor Total:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(45, 10, fmt.Sprintf("%.2f", total), "1", 1, "R", false, 0, "")

	// Guardamos el PDF en un archivo
	return pdf.OutputFileAndClose("Ver.pdf")
}
