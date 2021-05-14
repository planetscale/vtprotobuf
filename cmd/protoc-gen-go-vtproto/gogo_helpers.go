package main

import (
	"path"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"

	"google.golang.org/protobuf/compiler/protogen"
)

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func CamelCaseSlice(elem []string) string { return CamelCase(strings.Join(elem, "_")) }

func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

func FileName(file *protogen.File) string {
	fname := path.Base(file.Desc.Path())
	fname = strings.Replace(fname, ".proto", "", -1)
	fname = strings.Replace(fname, "-", "_", -1)
	fname = strings.Replace(fname, ".", "_", -1)
	return CamelCase(fname)
}

// GoMapValueTypes returns the map value Go type and the alias map value Go type (for casting), taking into
// account whether the map is nullable or the value is a message.
func GoMapValueTypes(mapField, valueField protoreflect.FieldDescriptor, goValueType, goValueAliasType string) (nullable bool, outGoType string, outGoAliasType string) {
	nullable = valueField.Kind() == protoreflect.MessageKind
	if nullable {
		// ensure the non-aliased Go value type is a pointer for consistency
		if strings.HasPrefix(goValueType, "*") {
			outGoType = goValueType
		} else {
			outGoType = "*" + goValueType
		}
		outGoAliasType = goValueAliasType
	} else {
		outGoType = strings.Replace(goValueType, "*", "", 1)
		outGoAliasType = strings.Replace(goValueAliasType, "*", "", 1)
	}
	return
}
