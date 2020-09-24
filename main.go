package main

// #include <string.h>
// #include <stdbool.h>
// #include <mysql.h>
// #cgo CFLAGS: -O3 -I/usr/include/mysql -fno-omit-frame-pointer
import "C"
import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"
)

// main function is needed even for generating shared object files
func main() {}

var l = log.New(os.Stderr, "csv-to-json: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

func msg(message *C.char, s string) {
	m := C.CString(s)
	defer C.free(unsafe.Pointer(m))

	C.strcpy(message, m)
}

//export csv_to_json_init
func csv_to_json_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) C.bool {
	if args.arg_count != 1 {
		msg(message, "`csv_to_json` requires 1 parameter: the CSV file string to be encoded")
		return C.bool(true)
	}

	argsTypes := (*[2]uint32)(unsafe.Pointer(args.arg_type))

	argsTypes[0] = C.STRING_RESULT
	initid.maybe_null = 1

	return C.bool(false)
}

//export csv_to_json
func csv_to_json(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *uint64, isNull *C.char, message *C.char) *C.char {
	c := 1
	argsArgs := (*[1 << 30]*C.char)(unsafe.Pointer(args.args))[:c:c]
	argsLengths := (*[1 << 30]uint64)(unsafe.Pointer(args.lengths))[:c:c]

	*length = 0
	*isNull = 1
	if argsArgs[0] == nil {
		return nil
	}

	a := make([]string, c, c)
	for i, argsArg := range argsArgs {
		a[i] = C.GoStringN(argsArg, C.int(argsLengths[i]))
	}

	stringSlices, err := csv.NewReader(strings.NewReader(a[0])).ReadAll()
	if err != nil {
		l.Println(err)
		return nil
	}
	var b []byte

	if len(stringSlices) > 2 {
		objects := make([]map[string]string, 0, len(stringSlices)-1)

		h := stringSlices[0]
		uniqueColumns := make(map[string]struct{}, len(h))
		columnIndexes := make(map[int]string, len(h))

		stringSlices = stringSlices[1:]
		for _, slice := range stringSlices {
			newObject := make(map[string]string, len(h))

			for i, s := range slice {
				var column string
				var ok bool
				if column, ok = columnIndexes[i]; !ok {
					var newColumn string
					if i < len(h) {
						newColumn = strings.TrimSpace(h[i])
					}
					originalName := newColumn
					j := 2
					for {
						if _, ok := uniqueColumns[newColumn]; !ok {
							break
						}

						newColumn = strings.TrimSpace(fmt.Sprintf("%s %d", originalName, j))
						j++
					}

					column = newColumn
					columnIndexes[i] = newColumn
					uniqueColumns[newColumn] = struct{}{}
				}

				newObject[column] = s
			}

			objects = append(objects, newObject)
		}

		b, err = json.Marshal(objects)
		if err != nil {
			l.Println(err)
			return nil
		}
	}

	if b == nil {
		b = []byte("[]")
	}

	*length = uint64(len(b))
	*isNull = 0
	return C.CString(string(b))
}

//export csv_to_json_no_headers_init
func csv_to_json_no_headers_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) C.bool {
	if args.arg_count != 1 {
		msg(message, "`csv_to_json_no_headers` requires 1 parameter: the CSV file string to be encoded")
		return C.bool(true)
	}

	argsTypes := (*[2]uint32)(unsafe.Pointer(args.arg_type))

	argsTypes[0] = C.STRING_RESULT
	initid.maybe_null = 1

	return C.bool(false)
}

//export csv_to_json_no_headers
func csv_to_json_no_headers(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *uint64, isNull *C.char, message *C.char) *C.char {
	c := 1
	argsArgs := (*[1 << 30]*C.char)(unsafe.Pointer(args.args))[:c:c]
	argsLengths := (*[1 << 30]uint64)(unsafe.Pointer(args.lengths))[:c:c]

	*length = 0
	*isNull = 1
	if argsArgs[0] == nil {
		return nil
	}

	a := make([]string, c, c)
	for i, argsArg := range argsArgs {
		a[i] = C.GoStringN(argsArg, C.int(argsLengths[i]))
	}

	stringSlices, err := csv.NewReader(strings.NewReader(a[0])).ReadAll()
	if err != nil {
		l.Println(err)
		return nil
	}
	if stringSlices == nil {
		stringSlices = make([][]string, 0)
	}
	b, err := json.Marshal(stringSlices)
	if err != nil {
		l.Println(err)
		return nil
	}

	*length = uint64(len(b))
	*isNull = 0
	return C.CString(string(b))
}
