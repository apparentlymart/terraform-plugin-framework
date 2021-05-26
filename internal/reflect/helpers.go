package reflect

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// trueReflectValue returns the reflect.Value for `in` after derefencing all
// the pointers and unwrapping all the interfaces. It's the concrete value
// beneath it all.
func trueReflectValue(in interface{}) reflect.Value {
	val := reflect.ValueOf(in)
	kind := val.Type().Kind()
	for kind == reflect.Interface || kind == reflect.Ptr {
		val = val.Elem()
		kind = val.Type().Kind()
	}
	return val
}

// commaSeparatedString returns an English joining of the strings in `in`,
// using "and" and commas as appropriate.
func commaSeparatedString(in []string) string {
	switch len(in) {
	case 0:
		return ""
	case 1:
		return in[0]
	case 2:
		return strings.Join(in, " and ")
	default:
		in[len(in)-1] = "and " + in[len(in)-1]
		return strings.Join(in, ", ")
	}
}

// getStructTags returns a map of Terraform field names to their position in
// the tags of the struct `in`. `in` must be a struct.
func getStructTags(ctx context.Context, in interface{}, path *tftypes.AttributePath) (map[string]int, error) {
	tags := map[string]int{}
	typ := trueReflectValue(in).Type()
	if typ.Kind() != reflect.Struct {
		return nil, path.NewErrorf("can't get struct tags of %T, is not a struct", in)
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			// skip unexported fields
			continue
		}
		tag := field.Tag.Get(`tfsdk`)
		if tag == "-" {
			// skip explicitly excluded fields
			continue
		}
		if tag == "" {
			return nil, path.NewErrorf(`need a struct tag for "tfsdk" on %s`, field.Name)
		}
		path := path.WithAttributeName(tag)
		if !isValidFieldName(tag) {
			return nil, path.NewError(errors.New("invalid field name, must only use lowercase letters, underscores, and numbers, and must start with a letter"))
		}
		if other, ok := tags[tag]; ok {
			return nil, path.NewErrorf("can't use field name for both %s and %s", typ.Field(other).Name, field.Name)
		}
		tags[tag] = i
	}
	return tags, nil
}

// isValidFieldName returns true if `name` can be used as a field name in a
// Terraform resource or data source.
func isValidFieldName(name string) bool {
	re := regexp.MustCompile("^[a-z][a-z0-9_]*$")
	return re.MatchString(name)
}
