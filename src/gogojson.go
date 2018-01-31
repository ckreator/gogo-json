package gogojson

/*
var numCheck = regexp.MustCompile("[0-9]")

// var numCheck = regexp.MustCompile("[0-9]")

func GogoJson(source string) map[string]interface{} {
	json := MakeIterator(source)
	out := make(map[string]interface{})
	skipWhitespace(json)

	if json.Peek() == "{" {
		json.Next()
	}

	for json.Peek() != "}" {
		skipWhitespace(json)
		key, value := parseKeyValue(json)
		skipWhitespace(json)
		skipComma(json)

		out[key] = value
	}

	return out
}
*/
