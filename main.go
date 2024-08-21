package main

import (
 "bufio"
 "flag"
 "fmt"
 "io/ioutil"
 "os"
 "strings"
)

func main() {
 // Определяем флаги для командной строки
 inputFile := flag.String("input", "", "Входной файл для конвертации")
 inputEncoding := flag.String("input-encoding", "utf-8", "Кодировка входного файла")
 outputEncoding := flag.String("output-encoding", "utf-8", "Кодировка выходного файла")
 outputFile := flag.String("output", "", "Выходной файл после конвертации")
 flag.Parse()

 // Проверяем, что входной и выходной файлы указаны
 if *inputFile == "" || *outputFile == "" {
  fmt.Println("Необходимо указать входной и выходной файлы")
  return
 }

 // Читаем содержимое входного файла
 content, err := ioutil.ReadFile(*inputFile)
 if err != nil {
  fmt.Printf("Ошибка при чтении файла %s: %v\n", *inputFile, err)
  return
 }

 // Конвертируем содержимое файла из одной кодировки в другую
 convertedContent, err := ConvertEncoding(string(content), *inputEncoding, *outputEncoding)
 if err != nil {
  fmt.Printf("Ошибка при конвертации кодировки: %v\n", err)
  return
 }

 // Записываем конвертированное содержимое в выходной файл
 err = ioutil.WriteFile(*outputFile, []byte(convertedContent), 0644)
 if err != nil {
  fmt.Printf("Ошибка при записи файла %s: %v\n", *outputFile, err)
  return
 }

 fmt.Printf("Файл успешно сконвертирован из %s в %s\n", *inputEncoding, *outputEncoding)
}

// ConvertEncoding конвертирует текст из одной кодировки в другую
func ConvertEncoding(text, inputEncoding, outputEncoding string) (string, error) {
 reader := strings.NewReader(text)
 scanner := bufio.NewScanner(reader)
 scanner.Split(bufio.ScanLines)

 var convertedLines []string
 for scanner.Scan() {
  line := scanner.Text()
  convertedLine, err := ConvertLine(line, inputEncoding, outputEncoding)
  if err != nil {
   return "", err
  }
  convertedLines = append(convertedLines, convertedLine)
 }

 if err := scanner.Err(); err != nil {
  return "", err
 }

 return strings.Join(convertedLines, "\n"), nil
}

// ConvertLine конвертирует одну строку из одной кодировки в другую
func ConvertLine(line, inputEncoding, outputEncoding string) (string, error) {
 decoded, err := DecodeString(line, inputEncoding)
 if err != nil {
  return "", err
 }

 encoded, err := EncodeString(decoded, outputEncoding)
 if err != nil {
  return "", err
 }

 return encoded, nil
}

// DecodeString декодирует строку из указанной кодировки
func DecodeString(s, encoding string) (string, error) {
 switch encoding {
 case "utf-8":
  return s, nil
 case "iso-8859-1":
  return DecodeISO88591(s), nil
 default:
  return "", fmt.Errorf("Неподдерживаемая кодировка: %s", encoding)
 }
}

// EncodeString кодирует строку в указанную кодировку
func EncodeString(s, encoding string) (string, error) {
 switch encoding {
 case "utf-8":
  return s, nil
 case "iso-8859-1":
  return EncodeISO88591(s), nil
 default:
  return "", fmt.Errorf("Неподдерживаемая кодировка: %s", encoding)
 }
}

// DecodeISO88591 декодирует строку из ISO-8859-1 кодировки
func DecodeISO88591(s string) string {
 runes := make([]rune, len(s))
 for i, c := range s {
  runes[i] = rune(c)
 }
 return string(runes)
}

// EncodeISO88591 кодирует строку в ISO-8859-1 кодировку
func EncodeISO88591(s string) string {
 runes := []byte(s)
 for i, r := range runes {
  if r > 255 {
   runes[i] = '?'
  }
 }
 return string(runes)
}
