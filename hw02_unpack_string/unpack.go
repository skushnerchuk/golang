package unpack

import (
	"errors"
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

// Разбивка строки на токены для дальнейшей распаковки с учетом экранировки.
// Если за символом нет цифры - то добавляется токен с 1 повторением
// Если за символом есть цифра - то добавляется токен с повторением, равном этой цифре.
func splitToTokens(data []rune) ([]Token, error) {
	result := make([]Token, 0)
	for i := 0; i < len(data); i++ {
		// Если нам попалась цифра - значит строка неверная
		if unicode.IsDigit(data[i]) {
			return nil, ErrInvalidString
		}
		count := 1
		char := string(data[i])
		// Если символ экранирования - обратаываем его блок
		if data[i] == Backslash {
			// Перескакиваем на экранированный символ
			i++
			char = string(data[i])
			if !isCorrectEscapedChar(data[i]) {
				return nil, ErrInvalidString
			}
			// Если за экранируемым символом есть цифра, запоминаем ее в счетчике повторений и перескакиваем эту цифру
			if len(data) > i+1 && unicode.IsDigit(data[i+1]) {
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

// Unpack Распаковка строки с поддержкой экранировки
// Допускается экранирование цифр и символа "\".
func Unpack(data string) (string, error) {
	if strings.HasSuffix(data, "\\") && data != `\\` {
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
