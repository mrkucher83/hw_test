package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var re = regexp.MustCompile(`[\n\t]`)

type Word struct {
	Value string
	Count int
}

func Top10(str string) []string {
	// загоняем текст в срез, заменяя "\n" и "\t" на пробелы
	arrWord := strings.Split(string(re.ReplaceAll([]byte(str), []byte(" "))), " ")

	// создаем хэш-таблицу для подсчета слов
	top := make(map[string]int)

	for _, word := range arrWord {
		top[word]++
	}

	// удаляем пробелы, которые попали, как слова (2 и более пробелов подряд)
	delete(top, "")

	// создаем слайс для сортировки значений
	arrSort := make([]Word, 0, len(top))

	for key, val := range top {
		arrSort = append(arrSort, Word{key, val})
	}

	// сортируем слова по количеству и лексикографически (если кол-во одинаковое)
	sort.Slice(arrSort, func(i, j int) bool {
		if arrSort[i].Count == arrSort[j].Count {
			return arrSort[i].Value < arrSort[j].Value
		}

		return arrSort[i].Count > arrSort[j].Count
	})

	if len(arrSort) > 10 {
		arrSort = arrSort[:10]
	}

	// создаем итоговый срез отсортированных слов
	result := make([]string, 0, len(arrSort))

	for i := 0; i < len(arrSort); i++ {
		result = append(result, arrSort[i].Value)
	}

	// возвращаем итоговый срез из 10 или менее слов
	return result
}
