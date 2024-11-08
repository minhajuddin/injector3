package injector3

import (
	"fmt"
	"reflect"
)

type Injector struct {
	ctors map[reflect.Type]any
}

func (i *Injector) Inspect() {
	for k, v := range i.ctors {
		fmt.Printf("%s %T\n", k, v)
	}
}

func Register[T any](i *Injector, ctor func(*Injector) T) *Injector {
	ctorType := reflect.TypeOf(ctor).Out(0)
	i.ctors[ctorType] = ctor
	return i
}

func resolve(injector *Injector, t reflect.Type) reflect.Value {
	ctor, ok := injector.ctors[t]
	if ok {
		ctorFn := reflect.ValueOf(ctor)
		return ctorFn.Call([]reflect.Value{reflect.ValueOf(injector)})[0]
	}
	panic("could not resolve inner type")
}

func Resolve[T any](injector *Injector) T {
	t := reflect.TypeOf([0]T{}).Elem()
	ctor, ok := injector.ctors[t]
	if ok {
		return ctor.(func(i *Injector) T)(injector)
	}

	// if t is a struct, we need to resolve its fields and build the struct
	if t.Kind() == reflect.Struct {
		// create a new instance of the struct
		v := reflect.New(t).Elem()
		// get the number of fields in the struct
		numFields := t.NumField()
		for i := 0; i < numFields; i++ {
			field := t.Field(i)
			fieldType := field.Type
			// resolve the field
			fieldValue := resolve(injector, fieldType)
			// set the field value
			v.Field(i).Set(fieldValue)
		}
		return v.Interface().(T)
	}
	panic("could not resolve type")
}

func NewInjector() *Injector {
	return &Injector{
		ctors: make(map[reflect.Type]any),
	}
}

type DBPoolImpl struct{}

type DBPool interface{}

type HomeController struct {
	db DBPool
}

func main() {
	i := NewInjector()

	fmt.Println("--------------------")
	Register(i, func(i *Injector) DBPool {
		return &DBPoolImpl{}
	})

	fmt.Println("--------------------")
	Register(i, func(i *Injector) *HomeController {
		return &HomeController{
			db: Resolve[DBPool](i),
		}
	})

	fmt.Println("--------------------")
	fmt.Printf("%#v\n", Resolve[*HomeController](i))
	fmt.Println("--------------------")
}
