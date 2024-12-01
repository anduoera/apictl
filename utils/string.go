package utils

import (
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

func InArray(val interface{}, array interface{}) (exists bool) {
	exists = false
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				exists = true
				return
			}
		}
	}
	return
}

func ContainsOnlyLetters(s string) error {
	letterRegex := `^[a-zA-Z]+$`
	matched, err := regexp.MatchString(letterRegex, s)
	if err != nil {
		return fmt.Errorf("error occurred: %s", err.Error())
	}

	if !matched || !unicode.IsUpper(rune(s[0])) || len(s) < 2 {
		return fmt.Errorf("error only input [a-zA-Z] and first letter is capitalized")
	}

	return nil
}

func CamelToSnake(name string) string {
	var builder strings.Builder
	for i, char := range name {
		if unicode.IsUpper(char) {
			if i != 0 {
				builder.WriteRune('_')
			}
			builder.WriteRune(unicode.ToLower(char))
		} else {
			builder.WriteRune(char)
		}
	}
	return builder.String()
}

func NameToFileNamePath(path, name string) string {
	return filepath.Join(path, CamelToSnake(name)+".go")
}
