package criteria

import (
	"net/http"
	"regexp"
	"strconv"
)

func convertValueType(value string) interface{} {
	// Try integer conversion
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}

	// Try float conversion
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}

	// Try boolean conversion
	if value == "true" {
		return true
	} else if value == "false" {
		return false
	}

	// Default to string
	return value
}

func QueryParamsToFilterParser(r *http.Request) []*FilterPrimitive {
	var filters []*FilterPrimitive
	re := regexp.MustCompile(`filters\[(\d+)\]\[(\w+)\]`)

	for key, values := range r.URL.Query() {
		matches := re.FindStringSubmatch(key)
		if len(matches) == 3 {
			index, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}
			field := matches[2]

			for len(filters) <= index {
				filters = append(filters, &FilterPrimitive{})
			}

			if len(values) > 0 {
				switch field {
				case "field":
					filters[index].Field = values[0]
				case "operator":
					filters[index].Operator = values[0]
				case "value":
					filters[index].Value = convertValueType(values[0])
				}
			}
		}
	}

	return filters
}
