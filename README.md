# Código de ejemplo para calcular CRC-32

Este código en Go proporciona una implementación para calcular el CRC-32 de una trama de bits. El programa toma como entrada una trama de bits, la polinomio generador y el bit de relleno y devuelve el valor CRC-32 correspondiente.

El algoritmo para calcular el CRC-32 funciona de la siguiente manera:

1. Se toma la trama de bits y se convierte a una secuencia de bits en formato binario.
2. Se aplica el bit de relleno a la secuencia de bits.
3. Se genera una tabla de CRC utilizando el polinomio generador.
4. Se calcula el CRC-32 iterando sobre cada bit de la secuencia de bits y actualizando el valor del CRC-32 en cada iteración utilizando la tabla de CRC generada en el paso anterior.

## Requisitos

Para ejecutar el código, necesitas tener instalado Go (versión 1.16 o superior).

Además, debes tener instaladas las siguientes bibliotecas de Go:

- logrus
- regexp

## Uso

Para utilizar el programa, ejecuta el siguiente comando en la línea de comandos:

```
go run main.go
```

## Explicación del Código

### Función `textToPolynomial(text string) ([]int, error)`

Esta función convierte una cadena de texto que representa un polinomio en su representación de lista de exponentes de las variables. Por ejemplo, la cadena `x^4 + x^2 + 1` se convertiría en una lista con los valores `[4, 2, 0]`.

### Función `polynomialDegree(polynomial string) int`

Esta función devuelve el grado de un polinomio representado como una cadena de texto. El grado de un polinomio es el exponente más alto de la variable en el polinomio. Por ejemplo, el grado del polinomio `x^4 + x^2 + 1` es 4.

### Función `polynomialToBinary(polynomial []int) string`

Esta función convierte una lista de exponentes de variables en una cadena de texto que representa el polinomio en notación binaria. Los coeficientes de la cadena de texto representan los coeficientes binarios del polinomio. Por ejemplo, la lista `[4, 2, 0]` se convertiría en la cadena `10101`.

### Función `binaryToPolynomial(binaryRepresentation string) []int`

Esta función convierte una cadena de texto que representa un polinomio en notación binaria en su representación de lista de exponentes de variables. Por ejemplo, la cadena `10101` se convertiría en la lista `[4, 2, 0]`.

### Función `generateCRCTable(poly uint32) [256]uint32`

Esta función genera una tabla de valores CRC-32 para el polinomio especificado. El resultado es una matriz de `256` elementos, cada uno de los cuales es un valor CRC-32 para un byte de datos de entrada.

### Función `crc32(binData string, poly uint32) (uint32, error)`

Esta función calcula el valor de verificación CRC-32 para una cadena de texto binaria de entrada y un polinomio especificado. La función utiliza la tabla CRC-32 generada por `generateCRCTable()` para realizar el cálculo.

### Función `isBinary(s string) bool`

Esta función verifica si una cadena de texto contiene solo caracteres binarios `(0 y 1)`. Devuelve true si es así, false en caso contrario.

### Función `textToBinary(s string) (string, error)`

Esta función convierte una cadena de texto en una cadena de texto binaria. Si la cadena de entrada ya está en formato binario, la función simplemente devuelve la cadena original. De lo contrario, la función convierte cada carácter en su representación binaria y concatena los resultados para formar una cadena binaria completa.

### Función `xor(a, b string) string`

Esta función realiza una operación XOR a nivel de bit entre dos cadenas de texto binarias de igual longitud. La función devuelve una cadena binaria que representa el resultado de la operación.

### Función `generateZerosWithGradPoly(i int) string`

Esta función genera una cadena de texto que representa un número binario con un cierto número de ceros, correspondiente al grado del polinomio más alto, más un número adicional de ceros especificado por el argumento i. Por ejemplo, si el grado del polinomio más alto es 3 y el argumento i es 2, la función devolverá
