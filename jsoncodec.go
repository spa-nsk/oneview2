package oneview

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var jsonRawMessageType = reflect.TypeOf(json.RawMessage{})

// DecodeJSON unmarshals OneView JSON into out.
//
// Global Dashboard (API 300) and appliance APIs 3800–8800 disagree on several
// scalar types: eTag as string or number, powerLock as bool or string,
// hostOsType as integer, null or OS name, integers that are null on rack
// servers, and the occasional single object where an array is documented.
// Standard encoding/json rejects those payloads; this decoder coerces them.
func DecodeJSON(data []byte, out any) error {
	return decodeJSON(data, out)
}

func decodeJSON(data []byte, out any) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	err := json.Unmarshal(data, out)
	if err == nil {
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var generic any
	if err := dec.Decode(&generic); err != nil {
		return err
	}
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return err
	}
	if err := assignCoerced(rv.Elem(), generic); err != nil {
		return err
	}
	return nil
}

func assignCoerced(dest reflect.Value, src any) error {
	if !dest.IsValid() {
		return nil
	}
	if src == nil {
		dest.Set(reflect.Zero(dest.Type()))
		return nil
	}
	for dest.Kind() == reflect.Pointer {
		if dest.IsNil() {
			if !dest.CanSet() {
				return nil
			}
			dest.Set(reflect.New(dest.Type().Elem()))
		}
		dest = dest.Elem()
	}
	if dest.CanAddr() {
		if u, ok := dest.Addr().Interface().(json.Unmarshaler); ok {
			b, err := json.Marshal(src)
			if err == nil && u.UnmarshalJSON(b) == nil {
				return nil
			}
		}
	}
	if !dest.CanSet() {
		return nil
	}

	switch dest.Kind() {
	case reflect.Bool:
		dest.SetBool(asBool(src))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dest.SetInt(asInt(src))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dest.SetUint(uint64(asInt(src)))
	case reflect.Float32, reflect.Float64:
		dest.SetFloat(asFloat(src))
	case reflect.String:
		dest.SetString(asString(src))
	case reflect.Slice:
		if dest.Type() == jsonRawMessageType {
			b, err := json.Marshal(src)
			if err != nil {
				return err
			}
			dest.SetBytes(b)
			return nil
		}
		list := asList(src)
		sl := reflect.MakeSlice(dest.Type(), len(list), len(list))
		for i, item := range list {
			if err := assignCoerced(sl.Index(i), item); err != nil {
				return err
			}
		}
		dest.Set(sl)
	case reflect.Map:
		mv, ok := src.(map[string]any)
		if !ok {
			return nil
		}
		if dest.IsNil() {
			dest.Set(reflect.MakeMap(dest.Type()))
		}
		kt := dest.Type().Key()
		et := dest.Type().Elem()
		for k, v := range mv {
			key := reflect.New(kt).Elem()
			if err := assignCoerced(key, k); err != nil {
				return err
			}
			val := reflect.New(et).Elem()
			if err := assignCoerced(val, v); err != nil {
				return err
			}
			dest.SetMapIndex(key, val)
		}
	case reflect.Struct:
		mv, ok := src.(map[string]any)
		if !ok {
			return nil
		}
		return assignStruct(dest, mv)
	case reflect.Interface:
		if dest.NumMethod() != 0 {
			return nil
		}
		dest.Set(reflect.ValueOf(normalizeJSON(src)))
	}
	return nil
}

func assignStruct(dest reflect.Value, m map[string]any) error {
	t := dest.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		fv := dest.Field(i)
		if f.Anonymous && (fv.Kind() == reflect.Struct || (fv.Kind() == reflect.Pointer && fv.Type().Elem().Kind() == reflect.Struct)) {
			if fv.Kind() == reflect.Pointer && fv.IsNil() && fv.CanSet() {
				fv.Set(reflect.New(fv.Type().Elem()))
			}
			target := fv
			if target.Kind() == reflect.Pointer {
				target = target.Elem()
			}
			if err := assignStruct(target, m); err != nil {
				return err
			}
			continue
		}
		name, skip := jsonFieldName(f)
		if skip {
			continue
		}
		val, ok := m[name]
		if !ok {
			continue
		}
		if err := assignCoerced(fv, val); err != nil {
			return err
		}
	}
	return nil
}

func jsonFieldName(f reflect.StructField) (string, bool) {
	tag := f.Tag.Get("json")
	if tag == "-" {
		return "", true
	}
	if tag == "" {
		return f.Name, false
	}
	name, _, _ := strings.Cut(tag, ",")
	if name == "" || name == "-" {
		return "", true
	}
	return name, false
}

func asString(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case json.Number:
		return x.String()
	case bool:
		if x {
			return "true"
		}
		return "false"
	case float64:
		if x == float64(int64(x)) {
			return strconv.FormatInt(int64(x), 10)
		}
		return strconv.FormatFloat(x, 'g', -1, 64)
	default:
		b, err := json.Marshal(x)
		if err != nil {
			return fmt.Sprint(x)
		}
		return string(b)
	}
}

func asInt(v any) int64 {
	switch x := v.(type) {
	case nil:
		return 0
	case json.Number:
		i, err := x.Int64()
		if err != nil {
			f, err2 := x.Float64()
			if err2 != nil {
				return 0
			}
			return int64(f)
		}
		return i
	case float64:
		return int64(x)
	case int:
		return int64(x)
	case int64:
		return x
	case bool:
		if x {
			return 1
		}
		return 0
	case string:
		s := strings.TrimSpace(x)
		if s == "" {
			return 0
		}
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			f, err2 := strconv.ParseFloat(s, 64)
			if err2 != nil {
				return 0
			}
			return int64(f)
		}
		return i
	default:
		return 0
	}
}

func asFloat(v any) float64 {
	switch x := v.(type) {
	case nil:
		return 0
	case json.Number:
		f, err := x.Float64()
		if err != nil {
			return 0
		}
		return f
	case float64:
		return x
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case bool:
		if x {
			return 1
		}
		return 0
	case string:
		s := strings.TrimSpace(x)
		if s == "" {
			return 0
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0
		}
		return f
	default:
		return 0
	}
}

func asBool(v any) bool {
	switch x := v.(type) {
	case nil:
		return false
	case bool:
		return x
	case json.Number:
		i, err := x.Int64()
		if err != nil {
			f, _ := x.Float64()
			return f != 0
		}
		return i != 0
	case float64:
		return x != 0
	case int:
		return x != 0
	case int64:
		return x != 0
	case string:
		switch strings.ToLower(strings.TrimSpace(x)) {
		case "", "false", "0", "no", "off":
			return false
		default:
			return true
		}
	default:
		return false
	}
}

func asList(v any) []any {
	switch x := v.(type) {
	case nil:
		return nil
	case []any:
		return x
	default:
		return []any{x}
	}
}

func normalizeJSON(v any) any {
	switch x := v.(type) {
	case json.Number:
		if i, err := x.Int64(); err == nil {
			return i
		}
		if f, err := x.Float64(); err == nil {
			return f
		}
		return x.String()
	case []any:
		out := make([]any, len(x))
		for i, item := range x {
			out[i] = normalizeJSON(item)
		}
		return out
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, item := range x {
			out[k] = normalizeJSON(item)
		}
		return out
	default:
		return v
	}
}
