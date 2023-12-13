package unpack

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const (
	Space     = 0x30
	Backslash = 0x5C
)

type Token struct {
	char  string // Символ
	count int    // Количество его повторений в итоговой строке
}

// Проверка корректности экранируемого символа.
func isCorrectEscapedChar(char rune) bool {
	return unicode.IsDigit(char) || char == Backslash
}

// Вычисление необходимого количества повторений.
// Используется вместо strconv.Atoi, чтобы избежать многословности в обработке ошибок
// Базируется на особенности реализации таблицы ASCII, в которой цифры идут подряд, начиная с 0, который имеет код 0x30
// Например, если string(char) = "5", то его ASCII-код = 0x35, а значит количество повторений равно 5 (0x35 - 0x30 = 5).
func calcCounterValue(char rune) int {
	return int(char) - Space
}

// Разбивка строки на токены для дальнейшей распаковки.
func splitToTokens(data []rune) ([]Token, error) {
	result := make([]Token, 0)
	for i := 0; i < len(data); i++ {
		if unicode.IsDigit(data[i]) {
			continue
		}
		count := 1
		char := string(data[i])
		// Если символ экранирования - обратаываем его блок
		if data[i] == Backslash {
			i++
			if !isCorrectEscapedChar(data[i]) {
				return nil, ErrInvalidString
			}
			char = string(data[i])
			// Если за экранируемым символом есть цифра, запоминаем ее в счетчике повторений
			if i+1 < len(data) && unicode.IsDigit(data[i+1]) {
				count = calcCounterValue(data[i+1])
				i++
			}
		} else if i < len(data)-1 && unicode.IsDigit(data[i+1]) {
			count = calcCounterValue(data[i+1])
			i++
		}
		result = append(result, Token{char: char, count: count})
	}
	return result, nil
}

// Проверка, что входная строка соответствует требуемому формату:
//   - не представляет собой число
//   - не начинается с цифры
//   - не содержит идущие подряд 2 и более числа, если перед первым нет \
//   - не заканчивается на символ экранирования
//
// Проверка корректности экранируемых блоков осуществляется при разбивке в функции splitToTokens.
func isCorrectString(data string) bool {
	_, err := strconv.Atoi(data)
	incorrect := regexp.MustCompile(`^\d|[^\\]\d{2,}`).FindString(data)
	return err != nil && incorrect == "" && !strings.HasSuffix(data, "\\")
}

// Unpack Распаковка строки с поддержкой экранировки
// Допускается экранирование цифр и символа "\".
func Unpack(data string) (string, error) {
	if !isCorrectString(data) {
		return "", ErrInvalidString
	}

	tokens, err := splitToTokens([]rune(data))
	if err != nil {
		return "", ErrInvalidString
	}

	var builder strings.Builder
	for _, t := range tokens {
		builder.WriteString(strings.Repeat(t.char, t.count))
	}
	return builder.String(), nil
}
