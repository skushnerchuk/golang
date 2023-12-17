package topwords

import (
	"sort"
	"strings"
	"unicode"
)

func trimWord(word string) string {
	return strings.TrimFunc(word, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r) && !unicode.IsMark(r) && r != '-'
	})
}

func Top10(data string) []string {
	// Разбиваем текст на слова и готовим частотную карту
	words := strings.Fields(strings.ToLower(data))
	pairs := make(map[string]int)
	for _, s := range words {
		word := trimWord(s)
		pairs[word]++
	}
	// Удаляем все что не считаем словами
	delete(pairs, "")
	delete(pairs, "-")
	words = make([]string, 0, len(pairs))
	for k := range pairs {
		words = append(words, k)
	}
	// Сначала сортируем слова по алфавиту
	sort.Strings(words)
	// а далее - по их количеству в построенной карте
	sort.SliceStable(words, func(i, j int) bool {
		return pairs[words[i]] > pairs[words[j]]
	})
	var a []string
	for i := 0; i < len(words) && i < 10; i++ {
		a = append(a, words[i])
	}
	return a
}
