package protocol

import (
	"fmt"
	"reflect"
	"strconv"
)

func appendTail(buf []byte) []byte {
	return append(buf, '\r', '\n')
}

func appendInt(buf []byte, n int64) []byte {
	buf = append(buf, byte('$'))
	buf = strconv.AppendInt(buf, numberLen(int64(n)), 10)
	buf = appendTail(buf)
	buf = strconv.AppendInt(buf, int64(n), 10)
	return appendTail(buf)
}

func appendBytes(buf []byte, b []byte) []byte {
	buf = append(buf, byte('$'))
	buf = strconv.AppendInt(buf, int64(len(b)), 10)
	buf = appendTail(buf)
	buf = append(buf, b...)
	return appendTail(buf)
}

func appendFloat(buf []byte, f float64) []byte {
	return appendBytes(buf, []byte(strconv.FormatFloat(f, 'f', -1, 64)))
}

func numberLen(number int64) int64 {
	var count int64 = 1
	if number < 0 {
		number = -number
		count = 2
	}
	for number > 9 {
		number /= 10
		count++
	}
	return count
}

func Pack(command string, args ...interface{}) ([]byte, error) {
	n := len(args)
	buf := make([]byte, 0, 10*n)

	buf = append(buf, byte('*'))
	buf = strconv.AppendInt(buf, int64(n+1), 10)
	buf = appendTail(buf)

	buf = appendBytes(buf, []byte(command))
	for _, arg := range args {
		switch v := arg.(type) {

		case nil:
			debug(v, "nil------")
			buf = appendBytes(buf, []byte{})

		case bool:
			debug(v, "bool------")
			if v {
				buf = appendBytes(buf, []byte{'1'})
			} else {
				buf = appendBytes(buf, []byte{'0'})
			}

		case []byte:
			debug(v, "byte[]------")
			buf = appendBytes(buf, v)
		case string:
			debug(v, "string------")
			buf = appendBytes(buf, []byte(v))

		case int:
			debug(v, "int------")
			buf = appendInt(buf, int64(v))
		case int8:
			debug(v, "int8------")
			buf = appendInt(buf, int64(v))
		case int16:
			debug(v, "int16------")
			buf = appendInt(buf, int64(v))
		case int32:
			debug(v, "int32------")
			buf = appendInt(buf, int64(v))
		case int64:
			debug(v, "int64------")
			buf = appendInt(buf, v)
		case uint:
			debug(v, "uint------")
			buf = appendInt(buf, int64(v))
		case uint8:
			debug(v, "uint8------")
			buf = appendInt(buf, int64(v))
		case uint16:
			debug(v, "uint16------")

			buf = appendInt(buf, int64(v))
		case uint32:
			debug(v, "uint32------")

			buf = appendInt(buf, int64(v))
		case uint64:
			debug(v, "uint64------")
			buf = appendInt(buf, int64(v))

		case float32:
			debug(v, "float32------")
			buf = appendFloat(buf, float64(v))
		case float64:
			debug(v, "fl64------")
			buf = appendFloat(buf, v)

		default:
			debug(v, "default------")
			buf = appendBytes(buf, []byte(fmt.Sprint(arg)))
		}
	}
	return buf, nil
}

func NormalizeArgs(args ...interface{}) []interface{} {
	normal := make([]interface{}, 0, 64)
	for _, arg := range args {
		v := reflect.ValueOf(arg)
		switch v.Kind() {
		case reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				normal = append(normal, v.Index(i).Interface())
			}
		case reflect.Map:
			for _, mk := range v.MapKeys() {
				normal = append(normal, mk.Interface(), v.MapIndex(mk).Interface())
			}
		default:
			normal = append(normal, v.Interface())
		}
	}

	return normal
}
