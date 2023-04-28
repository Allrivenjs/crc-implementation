package main

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func textToPolynomial(text string) ([]int, error) {
	re := regexp.MustCompile(`x\^(\d+)`)
	//get x alone
	reX := regexp.MustCompile(`x`)

	matches := re.FindAllStringSubmatch(text, -1)

	var polynomial []int
	for _, match := range matches {
		exp, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, fmt.Errorf("error al analizar el exponente: %v", err)
		}
		polynomial = append(polynomial, exp)
	}
	if reX != nil {
		polynomial = append(polynomial, 1)
	}

	return polynomial, nil
}

func polynomialDegree(polynomial string) int {
	return len(polynomial) - 1
}

func polynomialToBinary(polynomial []int) string {
	maxExp := 0
	for _, exp := range polynomial {
		if exp > maxExp {
			maxExp = exp
		}
	}

	bitLength := maxExp + 1
	binaryRepresentation := make([]byte, bitLength)

	for i := 0; i < bitLength; i++ {
		binaryRepresentation[i] = '0'
	}

	for _, exp := range polynomial {
		binaryRepresentation[bitLength-exp-1] = '1'
	}

	return string(binaryRepresentation)
}
func binaryToPolynomial(binaryRepresentation string) []int {
	var polynomial []int
	bitLength := len(binaryRepresentation)

	for i := 0; i < bitLength; i++ {
		if binaryRepresentation[i] == '1' {
			exp := bitLength - i - 1
			polynomial = append(polynomial, exp)
		}
	}

	return polynomial
}
func generateCRCTable(poly uint32) [256]uint32 {
	var table [256]uint32
	for i := uint32(0); i < 256; i++ {
		c := i
		for j := 0; j < 8; j++ {
			if c&1 == 1 {
				c = (c >> 1) ^ poly
			} else {
				c >>= 1
			}
		}
		table[i] = c
	}
	return table
}
func crc32(binData string, poly uint32) (uint32, error) {
	data, err := strconv.ParseUint(binData, 2, 64)
	if err != nil {
		return 0, fmt.Errorf("error al convertir la cadena binaria a número: %v", err)
	}

	table := generateCRCTable(poly)
	crc := uint32(0xFFFFFFFF)

	logrus.Infof("Procesando datos: %s", binData)
	for i := 0; i < len(binData); i++ {
		logrus.Infof("Iteración %d, CRC actual: %08X", i, crc)
		//crc = (crc >> 8) ^ table[byte(crc)^byte(data&0xFF)]: Actualiza el valor de crc realizando una operación XOR
		//entre el valor actual de crc desplazado 8 bits a la derecha y el valor en la tabla de CRC correspondiente al byte menos significativo de data.
		crc = (crc >> 8) ^ table[byte(crc)^byte(data&0xFF)]
		//data >>= 8: Desplaza los bits de data 8 posiciones a la derecha, descartando el byte menos significativo que se acaba de procesar.
		data >>= 8
	}

	checksum := ^crc
	logrus.Infof("Checksum CRC-32 calculado: %08X", checksum)

	return checksum, nil
}
func isBinary(s string) bool {
	// un string es binario si lo unico que contiene son '0' y '1'
	for _, r := range s {
		if r != '0' && r != '1' {
			return false
		}
	}
	return true
}

func textToBinary(s string) (string, error) {
	// primero, verificamos si la trama ya se encuentra en binario. Si lo esta, se devuelve tal cual
	if isBinary(s) {
		return s, nil
	}

	var binaryString string
	for _, r := range s {
		// se convierte cada valor de r en una cadena de 8 bits (un byte) en formato binario usando la
		// verbosidad %08b, que indica que la salida debe tener 8 caracteres y se debe completar con
		// ceros a la izquierda si es necesario
		binaryString += fmt.Sprintf("%08b", r)
	}

	return binaryString, nil
}

func xor(a, b string) string {
	result := ""
	length := len(a)
	if len(b) < length {
		length = len(b)
	}

	for i := 0; i < length; i++ {
		if a[i] == b[i] {
			result += "0"
		} else {
			result += "1"
		}
	}
	return result
}

func generateZerosWithGradPoly(i int) string {
	var result string
	for _i := 0; _i < i; _i++ {
		result += fmt.Sprintf("%d", 0)
	}
	return result
}

//fmt.Sprintf("Tamaño: %d, Cadena: %s, Divisor: %s, Residuo: %s", len(remainder), oldremainder, divisor, remainder)
//1001
//10011110

func crc32Binary(data, poly string, check ...string) string {
	logrus.Info("Iniciando el cálculo del CRC-32 en binario")
	// Hallamos nuevamente el grado del polinomio
	degree := polynomialDegree(poly)
	// Añadir ceros al final de los datos igual al grado del polinomio
	if len(check) > 0 {
		data += check[0]
	} else {
		data += generateZerosWithGradPoly(degree)
	}

	logrus.Infof("Datos con ceros añadidos: %s", data)

	divisor := poly
	remainder := data[:degree+1]

	// Crear un array para almacenar el proceso de cada iteración
	process := []string{fmt.Sprintf("Tamaño: %d, Cadena: %s, Divisor: %s, Residuo: %s", len(remainder), data, divisor, remainder)}

	i := degree
	for i < len(data) {
		oldremainder := remainder
		if remainder[0] == '1' {
			logrus.Infof("Iteración %d: Dividiendo %s por %s", i-degree, remainder, divisor)
			remainder = xor(remainder, divisor)
			process[len(process)-1] = fmt.Sprintf("Tamaño: %d, Cadena: %s, Divisor: %s, Residuo: %s", len(remainder), oldremainder, divisor, remainder)
		}

		if i < len(data)-1 {
			remainder = remainder[1:] + string(data[i+1])
			i++
		} else {
			break
		}

		// Añadir el tamaño, divisor y residuo actual al array process
		//process = append(process, fmt.Sprintf("Tamaño: %d, Cadena: %s, Divisor: %s, Residuo: %s", len(remainder), oldremainder, divisor, remainder))
	}

	// Eliminar la última entrada del array process, ya que contiene el residuo final incorrecto
	process = process[:len(process)-1]

	// Imprimir el array process con todos los pasos del proceso
	logrus.Info("Proceso completo:")
	for i, step := range process {
		logrus.Infof("Paso %d: %s", i, step)
	}

	return remainder
}

func binPolynomial(polynomial string) string {
	// si el polinomio ya esta expresado en binario, lo retornamos tal cual
	if isBinary(polynomial) {
		return polynomial
	}
	// transformamos el polinomio en string a un array de enteros, presentando el polinomio binario
	_polynomial, err := textToPolynomial(polynomial)
	if err != nil {
		logrus.Fatal("Ocurrio un error al convertir el polinomio", err)
		panic(err)
	}
	return polynomialToBinary(_polynomial)
}

func main() {
	// Creamos una instancia de reader
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el polinomio a codificar: ")

	// ejecutamos la lectura de string del reader
	text, err := reader.ReadString('\n')
	if err != nil {
		logrus.Fatal("Ocurrio un error al leer el polinomio", err)
		panic(err)
	}

	// eliminamos el ultimo caracter del reader (\n)
	text = text[:len(text)-1]
	text = strings.ReplaceAll(text, "\r", "")
	logrus.Println("El polinomio es: ", text)

	// extraemos el número binario del polinomio
	_binPolynomial := binPolynomial(text)

	logrus.Println("El polinomio en binario es: ", _binPolynomial)

	// hallamos el grado del polinomio
	_binPolynomial = strings.ReplaceAll(_binPolynomial, "\r", "")
	grad := polynomialDegree(_binPolynomial)
	logrus.Println("El grado del polinomio es: ", grad)

	//polynomial = binaryToPolynomial(binPolynomial)
	//logrus.Println("El polinomio en binario es: ", binPolynomial)
	//poly := uint32(0xEDB88320)
	fmt.Print("Ingrese la trama: ")

	// usamos el reader para pedir la trama de datos
	trama, err := reader.ReadString('\n')
	if err != nil {
		logrus.Fatal("Ocurrio un error al leer el polinomio", err)
		panic(err)
	}

	// eliminamos el ultimo caracter del reader (\n)
	trama = trama[:len(trama)-1]
	trama = strings.ReplaceAll(trama, "\r", "")
	logrus.Println("El trama es: ", trama)

	fmt.Print("Desea corromper los datos? (Y: yes, N: no): ")
	c, err := reader.ReadString('\n')
	if err != nil {
		logrus.Fatal("Ocurrio un error al leer el polinomio", err)
		panic(err)
	}

	c = c[:len(c)-1]

	// convertimos la trama de datos a números binarios
	binTrama, err := textToBinary(trama)

	// finalmente, una vez obtenidos el polinomio generador en binario y la trama en binario,
	// empezamos el algoritmo de redundancia ciclica
	originalCRC := crc32Binary(binTrama, _binPolynomial)
	logrus.Infof("Checksum CRC-32: %s\n", originalCRC)
	c = strings.TrimSpace(c)
	if c == "Y" || c == "y" {
		logrus.Println("Corrompiendo datos...")
		binTrama = binTrama[:len(binTrama)-1] + "1"
		// Simular la corrupción de datos cambiando un bit
		corruptedData := binTrama[:len(binTrama)-2] + "1111111" + binTrama[len(binTrama):]
		logrus.Printf("Datos corrompidos: %s", corruptedData)
		// Calcular el CRC-32 de los datos corrompidos
		corruptedCRC := crc32Binary(corruptedData, _binPolynomial)
		check := crc32Binary(binTrama, _binPolynomial, corruptedCRC)
		logrus.Printf("CRC-32 corrompido: %s", corruptedCRC)
		// Comprobar si los valores CRC coinciden
		if check == "0000" {
			logrus.Println("Los datos no están corruptos.")
		} else {
			logrus.Println("Los datos están corruptos.")
		}
	} else {
		check := crc32Binary(binTrama, _binPolynomial, originalCRC)
		logrus.Printf("Check %v", check)
		if check == "0000" {
			logrus.Println("Los datos no están corruptos.")
		} else {
			logrus.Println("Los datos están corruptos.")
		}
	}

	//x^6 + x^5 + x^4 + x^3 + x^2 + x

	//example for error check
	//data := "010010000110111101101100011011000110000101101100001011110110001101110010" // Datos en binario
	//poly := "10000100110001011101101111011011"

}
