package litter

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type dumpState struct {
	w                 io.Writer
	depth             int
	pointers          ptrmap
	visitedPointers   ptrmap
	parentPointers    ptrmap
	currentPointer    *ptrinfo
	homePackageRegexp *regexp.Regexp
}

func (s *dumpState) write(b []byte) {
	if _, err := s.w.Write(b); err != nil {
		panic(err)
	}
}

func (s *dumpState) writeString(str string) {
	s.write([]byte(str))
}

func (s *dumpState) indent() {
	s.write(bytes.Repeat([]byte("  "), s.depth))
}

func (s *dumpState) newlineWithPointerNameComment() {
	if ptr := s.currentPointer; ptr != nil {
		s.write([]byte(fmt.Sprintf(" // %s\n", ptr.label())))
		s.currentPointer = nil
		return
	}
	s.write([]byte("\n"))
}

func (s *dumpState) dumpType(v reflect.Value) {
	typeName := v.Type().String()
	s.write([]byte(typeName))
}

func (s *dumpState) dumpSlice(v reflect.Value) {
	s.dumpType(v)
	numEntries := v.Len()
	if numEntries == 0 {
		s.write([]byte("{}"))
		return
	}
	s.write([]byte("{"))
	s.newlineWithPointerNameComment()
	s.depth++
	for i := 0; i < numEntries; i++ {
		s.indent()
		s.dumpVal(v.Index(i))
		s.write([]byte(","))
		s.newlineWithPointerNameComment()
	}
	s.depth--
	s.indent()
	s.write([]byte("}"))
}

func (s *dumpState) dumpStruct(v reflect.Value) {
	dumpPreamble := func() {
		s.dumpType(v)
		s.write([]byte("{"))
		s.newlineWithPointerNameComment()
		s.depth++
	}
	preambleDumped := false
	vt := v.Type()
	numFields := v.NumField()
	for i := 0; i < numFields; i++ {
		vtf := vt.Field(i)
		if !preambleDumped {
			dumpPreamble()
			preambleDumped = true
		}
		s.indent()
		s.write([]byte(vtf.Name))
		s.write([]byte(": "))
		s.dumpVal(v.Field(i))
		s.write([]byte(","))
		s.newlineWithPointerNameComment()
	}
	if preambleDumped {
		s.depth--
		s.indent()
		s.write([]byte("}"))
	} else {
		// There were no fields dumped
		s.dumpType(v)
		s.write([]byte("{}"))
	}
}

func (s *dumpState) dumpMap(v reflect.Value) {
	if v.IsNil() {
		s.dumpType(v)
		s.writeString("(nil)")
		return
	}

	s.dumpType(v)

	keys := v.MapKeys()
	if len(keys) == 0 {
		s.write([]byte("{}"))
		return
	}

	s.write([]byte("{"))
	s.newlineWithPointerNameComment()
	s.depth++
	sort.Sort(mapKeySorter{keys: keys})
	for _, key := range keys {
		s.indent()
		s.dumpVal(key)
		s.write([]byte(": "))
		s.dumpVal(v.MapIndex(key))
		s.write([]byte(","))
		s.newlineWithPointerNameComment()
	}
	s.depth--
	s.indent()
	s.write([]byte("}"))
}

func (s *dumpState) dumpFunc(v reflect.Value) {
	parts := strings.Split(runtime.FuncForPC(v.Pointer()).Name(), "/")
	name := parts[len(parts)-1]

	// Anonymous function
	if strings.Count(name, ".") > 1 {
		s.dumpType(v)
	} else {
		s.write([]byte(name))
	}
}

func (s *dumpState) dumpChan(v reflect.Value) {
	vType := v.Type()
	res := []byte(vType.String())
	s.write(res)
}

func (s *dumpState) dump(value interface{}) {
	if value == nil {
		printNil(s.w)
		return
	}
	v := reflect.ValueOf(value)
	s.dumpVal(v)
}

func (s *dumpState) descendIntoPossiblePointer(value reflect.Value, f func()) {
	canonicalize := true
	if isPointerValue(value) {
		// If elision disabled, and this is not a circular reference, don't canonicalize
		if s.parentPointers.add(value) {
			canonicalize = false
		}

		// Add to stack of pointers we're recursively descending into
		s.parentPointers.add(value)
		defer s.parentPointers.remove(value)
	}

	if !canonicalize {
		ptr, _ := s.pointerFor(value)
		s.currentPointer = ptr
		f()
		return
	}

	ptr, firstVisit := s.pointerFor(value)
	if ptr == nil {
		f()
		return
	}
	if firstVisit {
		s.currentPointer = ptr
		f()
		return
	}
	s.write([]byte(ptr.label()))
}

func (s *dumpState) dumpVal(value reflect.Value) {
	if value.Kind() == reflect.Ptr && value.IsNil() {
		s.write([]byte("nil"))
		return
	}

	v := deInterface(value)
	kind := v.Kind()

	// Check if the type implements the String method
	if m := v.MethodByName("String"); m.IsValid() && m.CanInterface() {
		result := m.Call([]reflect.Value{})[0]
		s.writeString(result.String())
		return
	}

	switch kind {
	case reflect.Invalid:
		// Do nothing.  We should never get here since invalid has already
		// been handled above.
		s.write([]byte("<invalid>"))

	case reflect.Bool:
		printBool(s.w, v.Bool())

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		printInt(s.w, v.Int(), 10)

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		printUint(s.w, v.Uint(), 10)

	case reflect.Float32:
		printFloat(s.w, v.Float(), 32)

	case reflect.Float64:
		printFloat(s.w, v.Float(), 64)

	case reflect.Complex64:
		printComplex(s.w, v.Complex(), 32)

	case reflect.Complex128:
		printComplex(s.w, v.Complex(), 64)

	case reflect.String:
		s.write([]byte(strconv.Quote(v.String())))

	case reflect.Slice:
		if v.IsNil() {
			printNil(s.w)
			break
		}
		fallthrough

	case reflect.Array:
		s.descendIntoPossiblePointer(v, func() {
			s.dumpSlice(v)
		})

	case reflect.Interface:
		// The only time we should get here is for nil interfaces due to
		// unpackValue calls.
		if v.IsNil() {
			printNil(s.w)
		}

	case reflect.Ptr:
		s.descendIntoPossiblePointer(v, func() {
			s.writeString("&")
			s.dumpVal(v.Elem())
		})

	case reflect.Map:
		s.descendIntoPossiblePointer(v, func() {
			s.dumpMap(v)
		})

	case reflect.Struct:
		s.dumpStruct(v)

	case reflect.Func:
		s.dumpFunc(v)

	case reflect.Chan:
		s.dumpChan(v)

	default:
		if v.CanInterface() {
			s.writeString(fmt.Sprintf("%v", v.Interface()))
		} else {
			s.writeString(fmt.Sprintf("%v", v.String()))
		}
	}
}

// registers that the value has been visited and checks to see if it is one of the
// pointers we will see multiple times. If it is, it returns a temporary name for this
// pointer. It also returns a boolean value indicating whether this is the first time
// this name is returned so the caller can decide whether the contents of the pointer
// has been dumped before or not.
func (s *dumpState) pointerFor(v reflect.Value) (*ptrinfo, bool) {
	if isPointerValue(v) {
		if info, ok := s.pointers.get(v); ok {
			firstVisit := s.visitedPointers.add(v)
			return info, firstVisit
		}
	}
	return nil, false
}

// prepares a new state object for dumping the provided value
func newDumpState(value reflect.Value, writer io.Writer) *dumpState {
	result := &dumpState{
		pointers: mapReusedPointers(value),
		w:        writer,
	}

	return result
}

// Sdump dumps a value to a string according to the options
func Sdump(values ...interface{}) string {
	buf := new(bytes.Buffer)
	for i, value := range values {
		if i > 0 {
			_, _ = buf.Write([]byte(" "))
		}
		state := newDumpState(reflect.ValueOf(value), buf)
		state.dump(value)
	}
	return buf.String()
}

type mapKeySorter struct {
	keys []reflect.Value
}

func (s mapKeySorter) Len() int {
	return len(s.keys)
}

func (s mapKeySorter) Swap(i, j int) {
	s.keys[i], s.keys[j] = s.keys[j], s.keys[i]
}

func (s mapKeySorter) Less(i, j int) bool {
	ibuf := new(bytes.Buffer)
	jbuf := new(bytes.Buffer)
	newDumpState(s.keys[i], ibuf).dumpVal(s.keys[i])
	newDumpState(s.keys[j], jbuf).dumpVal(s.keys[j])
	return ibuf.String() < jbuf.String()
}
