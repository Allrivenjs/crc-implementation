package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

func textToPolynomial(text string) ([]int, error) {
	re := regexp.MustCompile(`([+-]?\d*x\^\d+|[+-]?\d*x|[+-]?\d+$|[+-]x)`)

	matches := re.FindAllStringSubmatch(text, -1)

	var polynomial []int
	for _, match := range matches {
		// Si hay una coincidencia para la expresión x^n
		if strings.Contains(match[0], "^") {
			// Extraer el exponente del término
			exp, err := strconv.Atoi(strings.Split(match[0], "^")[1])
			if err != nil {
				return nil, fmt.Errorf("error al analizar el exponente: %v", err)
			}
			polynomial = append(polynomial, exp)
			// Si hay una coincidencia para la expresión x
		} else if strings.Contains(match[0], "x") {
			// El coeficiente es 1 si no hay nada delante del "x"
			// o -1 si hay un signo menos delante del "x"
			coeff := 1
			if len(match[0]) > 1 {
				sign := string(match[0][0])
				if sign == "-" {
					coeff = -1
				}
			}
			polynomial = append(polynomial, coeff)
			// Si hay una coincidencia para la expresión constante
		} else if match[0] != "" {
			coeff, err := strconv.Atoi(match[0])
			if err != nil {
				return nil, fmt.Errorf("error al analizar el coeficiente: %v", err)
			}
			polynomial = append(polynomial, coeff)
		} else {
			// Término sin coeficiente (equivalente a coeficiente de 1)
			polynomial = append(polynomial, 1)
		}
	}
	return polynomial, nil
}

func polyToBinary(poly string) ([]int, error) {
	// Convertir el string de polinomio a una lista de coeficientes
	coeffs, err := textToPolynomial(poly)
	logrus.Println(coeffs)
	if err != nil {
		panic(err)
	}
	// Encontrar el valor máximo en el array
	maxVal := 0
	for _, val := range coeffs {
		if val > maxVal {
			maxVal = val
		}
	}

	// Crear un nuevo array de tamaño maxVal + 1 y llenarlo de ceros
	result := make([]int, maxVal+1)

	// Colocar un uno en las posiciones correspondientes a los valores del array original
	for i, val := range coeffs {

		if i == len(coeffs)-1 {
			result[val-1] = 1
		} else {
			result[val] = 1
		}

	}

	return result, nil
}

func polynomialDegree(polynomial string) int {
	return len(polynomial)
}

func polynomialToBinary(polynomial []int) string {
	var binaryString strings.Builder

	for i := len(polynomial) - 1; i >= 0; i-- {
		if polynomial[i] == 1 {
			binaryString.WriteString("1")
		} else {
			binaryString.WriteString("0")
		}
	}

	return binaryString.String()
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
	//logrus.Info("Iniciando el cálculo del CRC-32 en binario")
	response = append(response, fmt.Sprintf("Los datos no están corruptos."))
	// Hallamos nuevamente el grado del polinomio
	degree := polynomialDegree(poly)
	// Añadir ceros al final de los datos igual al grado del polinomio
	if len(check) > 0 {
		data += check[0]
	} else {
		data += generateZerosWithGradPoly(degree)
	}

	//logrus.Infof("Datos con ceros añadidos: %s", data)
	response = append(response, fmt.Sprintf("Datos con ceros añadidos: %s", data))

	divisor := poly
	remainder := data[:degree+1]

	// Crear un array para almacenar el proceso de cada iteración
	process := []string{fmt.Sprintf("Tamaño: %d, Cadena: %s, Divisor: %s, Residuo: %s", len(remainder), data, divisor, remainder)}

	i := degree
	for i < len(data) {
		oldremainder := remainder
		if remainder[0] == '1' {
			//logrus.Infof("Iteración %d: Dividiendo %s por %s", i-degree, remainder, divisor)
			response = append(response, fmt.Sprintf("Iteración %d: Dividiendo %s por %s", i-degree, remainder, divisor))
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
	//logrus.Info("Proceso completo:")
	response = append(response, fmt.Sprintf("Proceso completo:"))
	for i, step := range process {
		//logrus.Infof("Paso %d: %s", i, step)
		response = append(response, fmt.Sprintf("Paso %d: %s", i, step))
	}

	return remainder
}

func binPolynomial(polynomial string) string {
	// si el polinomio ya esta expresado en binario, lo retornamos tal cual
	if isBinary(polynomial) {
		return polynomial
	}
	// transformamos el polinomio en string a un array de enteros, presentando el polinomio binario
	_polynomial, err := polyToBinary(polynomial)
	if err != nil {
		logrus.Fatal("Ocurrio un error al convertir el polinomio", err)
		panic(err)
	}
	return polynomialToBinary(_polynomial)
}

var response []string

func eject(poly, trama string, c bool) {
	// Creamos una instancia de reader
	//reader := bufio.NewReader(os.Stdin)
	//fmt.Print("Ingrese el polinomio a codificar: ")

	// ejecutamos la lectura de string del reader
	//text, err := reader.ReadString('\n')
	//if err != nil {
	//	logrus.Fatal("Ocurrio un error al leer el polinomio", err)
	//	panic(err)
	//}

	// eliminamos el ultimo caracter del reader (\n)
	//text = text[:len(text)-1]
	//text = strings.ReplaceAll(text, "\r", "")
	response = append(response, fmt.Sprintf("El polinomio rectificador es: %v", poly))

	// extraemos el número binario del polinomio
	_binPolynomial := binPolynomial(poly)

	response = append(response, fmt.Sprintf("El polinomio rectificador en binario es: %v", _binPolynomial))

	// hallamos el grado del polinomio
	_binPolynomial = strings.ReplaceAll(_binPolynomial, "\r", "")
	grad := polynomialDegree(_binPolynomial)
	response = append(response, fmt.Sprintf("El grado del polinomio rectificador es: %v", grad))

	//polynomial = binaryToPolynomial(binPolynomial)
	//logrus.Println("El polinomio en binario es: ", binPolynomial)
	//poly := uint32(0xEDB88320)
	//fmt.Print("Ingrese la trama: ")

	// usamos el reader para pedir la trama de datos
	//trama, err := reader.ReadString('\n')
	//if err != nil {
	//	logrus.Fatal("Ocurrio un error al leer el polinomio", err)
	//	panic(err)
	//}

	// eliminamos el ultimo caracter del reader (\n)
	//trama = trama[:len(trama)-1]
	//trama = strings.ReplaceAll(trama, "\r", "")
	response = append(response, fmt.Sprintf("El trama es: %v", trama))

	//fmt.Print("Desea corromper los datos? (Y: yes, N: no): ")
	//c, err := reader.ReadString('\n')
	//if err != nil {
	//	logrus.Fatal("Ocurrio un error al leer el polinomio", err)
	//	panic(err)
	//}

	//c = c[:len(c)-1]

	// convertimos la trama de datos a números binarios

	binTrama := binPolynomial(trama)

	// finalmente, una vez obtenidos el polinomio generador en binario y la trama en binario,
	// empezamos el algoritmo de redundancia ciclica
	originalCRC := crc32Binary(binTrama, _binPolynomial)
	response = append(response, fmt.Sprintf("Checksum CRC-32: %s\n", originalCRC))

	//c = strings.TrimSpace(c)
	if c {
		response = append(response, fmt.Sprintf("Corrompiendo datos..."))
		binTrama = binTrama[:len(binTrama)-1] + "1"
		// Simular la corrupción de datos cambiando un bit
		corruptedData := binTrama[:len(binTrama)-2] + "1111111" + binTrama[len(binTrama):]
		response = append(response, fmt.Sprintf("Datos corrompidos: %s", corruptedData))
		// Calcular el CRC-32 de los datos corrompidos
		corruptedCRC := crc32Binary(corruptedData, _binPolynomial)
		check := crc32Binary(binTrama, _binPolynomial, corruptedCRC)
		response = append(response, fmt.Sprintf("CRC-32 corrompido: %s", corruptedCRC))
		// Comprobar si los valores CRC coinciden
		if check == generateZerosWithGradPoly(grad) {
			response = append(response, fmt.Sprintf("Los datos no están corruptos."))
		} else {
			response = append(response, fmt.Sprintf("Los datos están corruptos."))
		}
	} else {
		check := crc32Binary(binTrama, _binPolynomial, originalCRC)
		logrus.Printf("Check %v", check)
		//logrus.Println(generateZerosWithGradPoly(grad + 1))
		if check == generateZerosWithGradPoly(grad) {
			response = append(response, fmt.Sprintf("Los datos no están corruptos."))
		} else {
			response = append(response, fmt.Sprintf("Los datos están corruptos."))
		}
	}

	//x^6 + x^5 + x^4 + x^3 + x^2 + x

	//example for error check
	//1001
	//10011110
}

type request struct {
	Poly      string `json:"poly"`
	Trama     string `json:"trama"`
	Corructed bool   `json:"corructed"`
}

type Response struct {
	Resultado []string `json:"resultado"`
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Post("/crc", func(c *fiber.Ctx) error {
		var body request
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		eject(body.Poly, body.Trama, body.Corructed)
		//logrus.Println(response)
		response_ := Response{Resultado: response}
		jsonResponse, err := json.Marshal(response_)
		if err != nil {
			// Manejo de error
		}

		err = c.JSON(string(jsonResponse))

		if err != nil {
			return err
		}
		return nil
	})
	logrus.Fatal(app.Listen(":3000"))
}
