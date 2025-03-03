package utils

import (
	"reflect"
	"unicode"
)

func toSnakeCase(input string) string {
	var result []rune

	for i, char := range input {
		if input == "ID" {
			return string(input)
		}
		// Check if the character is uppercase
		if unicode.IsUpper(char) {
			// If it's not the first character, add an underscore before the uppercase letter
			if i > 0 {
				result = append(result, '_')
			}
			// Convert the uppercase letter to lowercase
			result = append(result, unicode.ToLower(char))
		} else {
			// Add the lowercase character as is
			result = append(result, char)
		}
	}

	return string(result)
}

func EditObject(req interface{}) map[string]interface{} {
	updateData := map[string]interface{}{}

	val := reflect.ValueOf(req)
	typ := reflect.TypeOf(req)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := toSnakeCase(typ.Field(i).Name)
		if field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface() {
			if fieldName != "ID" {
				updateData[fieldName] = field.Interface()
			}
		}
	}
	return updateData
}
