package main

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strconv"
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
		crc = (crc >> 8) ^ table[byte(crc)^byte(data&0xFF)]
		data >>= 8
	}

	checksum := ^crc
	logrus.Infof("Checksum CRC-32 calculado: %08X", checksum)

	return checksum, nil
}
func isBinary(s string) bool {
	for _, r := range s {
		if r != '0' && r != '1' {
			return false
		}
	}
	return true
}
func textToBinary(s string) (string, error) {
	if isBinary(s) {
		return s, nil
	}

	var binaryString string
	for _, r := range s {
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

func crc32Binary(data, poly string) string {
	logrus.Info("Iniciando el cálculo del CRC-32 en binario")
	degree := polynomialDegree(poly)
	// Añadir ceros al final de los datos igual al grado del polinomio
	data += generateZerosWithGradPoly(degree)
	logrus.Infof("Datos con ceros añadidos: %d", degree)

	divisor := poly
	remainder := data[:degree+1]

	for i := degree; i < len(data); i++ {
		if remainder[0] == '1' {
			logrus.Infof("Iteración %d: Dividiendo %s por %s", i-32, remainder, divisor)
			remainder = xor(remainder, divisor)
		}
		remainder = remainder[1:] + string(data[i])
		logrus.Infof("Iteración %d: Residuo actual: %s", i-32, remainder)
	}

	if remainder[0] == '1' {
		logrus.Infof("Última iteración: Dividiendo %s por %s", remainder, divisor)
		remainder = xor(remainder, divisor)
	}
	remainder = remainder[1:]

	return remainder
}

func binPolynomial(polynomial string) string {
	if isBinary(polynomial) {
		return polynomial
	}
	_polynomial, err := textToPolynomial(polynomial)
	if err != nil {
		logrus.Fatal("Ocurrio un error al convertir el polinomio", err)
		panic(err)
	}
	return polynomialToBinary(_polynomial)
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el polinomio a codificar: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		logrus.Fatal("Ocurrio un error al leer el polinomio", err)
		panic(err)
	}

	text = text[:len(text)-1]
	logrus.Println("El polinomio es: ", text)

	_binPolynomial := binPolynomial(text)
	logrus.Println("El polinomio en binario es: ", _binPolynomial)

	grad := polynomialDegree(_binPolynomial)
	logrus.Println("El grado del polinomio es: ", grad)

	//polynomial = binaryToPolynomial(binPolynomial)
	//logrus.Println("El polinomio en binario es: ", binPolynomial)
	//poly := uint32(0xEDB88320)
	fmt.Print("Ingrese la trama: ")
	trama, err := reader.ReadString('\n')
	if err != nil {
		logrus.Fatal("Ocurrio un error al leer el polinomio", err)
		panic(err)
	}

	trama = trama[:len(trama)-1]
	logrus.Println("El trama es: ", trama)

	fmt.Print("Desea corromper los datos? (Y: yes, N: no): ")
	c, err := reader.ReadString('\n')
	if err != nil {
		logrus.Fatal("Ocurrio un error al leer el polinomio", err)
		panic(err)
	}

	c = c[:len(c)-1]
	binTrama, err := textToBinary(trama)
	originalCRC := crc32Binary(binTrama, _binPolynomial)
	logrus.Infof("Checksum CRC-32: %s\n", originalCRC)
	if c == "Y" {
		logrus.Println("Corrompiendo datos...")
		binTrama = binTrama[:len(binTrama)-1] + "1"
		// Simular la corrupción de datos cambiando un bit
		corruptedData := binTrama[:len(binTrama)-2] + "11011011" + binTrama[len(binTrama):]
		fmt.Printf("Datos corrompidos: %s\n", corruptedData)
		// Calcular el CRC-32 de los datos corrompidos
		corruptedCRC := crc32Binary(corruptedData, _binPolynomial)
		fmt.Printf("CRC-32 corrompido: %s\n", corruptedCRC)
		// Comprobar si los valores CRC coinciden
		if originalCRC == corruptedCRC {
			fmt.Println("Los datos no están corrompidos.")
		} else {
			fmt.Println("Los datos están corrompidos.")
		}
	} else {
		logrus.Println("No se corrompieron los datos")
	}

	//x^6 + x^5 + x^4 + x^3 + x^2 + x
}
