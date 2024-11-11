package injector3

import (
	"fmt"
	"reflect"
	"sync"
)

// Public APIs

// Injector is a simple dependency injection container
type Injector struct {
	ctors map[reflect.Type]any
	mu    sync.RWMutex
}

// NewInjector creates a new instance of the Injector
func NewInjector() *Injector {
	return &Injector{
		ctors: make(map[reflect.Type]any),
	}
}

// Register registers a constructor for a type
func Register[T any](i *Injector, ctor func(*Injector) T) *Injector {
	ctorType := reflect.TypeOf(ctor).Out(0)
	i.mu.Lock()
	defer i.mu.Unlock()
	i.ctors[ctorType] = ctor
	return i
}

// Resolve resolves a type from the injector
func Resolve[T any](injector *Injector) T {
	t := reflect.TypeOf([0]T{}).Elem()
	ctor, ok := getConstructor(injector, t)

	// We have a constructor for the type, so this is a simple case
	if ok {
		return ctor.(func(i *Injector) T)(injector)
	}

	// if t is a struct, we need to resolve its fields and build the struct
	if t.Kind() == reflect.Struct {
		v := resolveStruct(injector, t)
		return v.Interface().(T)
	}

	// if t is a pointer to a struct, we need to resolve its fields and build the struct
	if t.Kind() == reflect.Pointer {
		tt := t.Elem()
		if tt.Kind() == reflect.Struct {
			v := resolveStruct(injector, tt)
			return v.Addr().Interface().(T)
		}
	}

	// NOTE: We don't handle other types in our injector
	panic(fmt.Sprintf("could not resolve type %s", t))
}

func getConstructor(i *Injector, t reflect.Type) (any, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	c, ok := i.ctors[t]
	return c, ok
}

func resolve(injector *Injector, t reflect.Type) reflect.Value {
	ctor, ok := getConstructor(injector, t)
	if ok {
		ctorFn := reflect.ValueOf(ctor)
		return ctorFn.Call([]reflect.Value{reflect.ValueOf(injector)})[0]
	}

	// if t is a struct, we need to resolve its fields and build the struct
	if t.Kind() == reflect.Struct {
		v := resolveStruct(injector, t)
		return v
	}

	// if t is a pointer to a struct, we need to resolve its fields and build the struct
	if t.Kind() == reflect.Pointer {
		tt := t.Elem()
		if tt.Kind() == reflect.Struct {
			return resolveStruct(injector, tt).Addr()
		}
	}

	panic(fmt.Sprintf("could not resolve inner type %s", t))
}

func resolveStruct(injector *Injector, t reflect.Type) reflect.Value {
	// create a new instance of the struct
	v := reflect.New(t).Elem()
	// get the number of fields in the struct
	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		fieldType := field.Type
		// resolve the field
		fieldValue := resolve(injector, fieldType)
		// set the field value
		v.Field(i).Set(fieldValue)
	}
	return v
}
