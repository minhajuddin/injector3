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

func Resolve[T any](i *Injector) T {
	var t0 [0]T
	t := reflect.TypeOf(t0).Elem()
	ctor, ok := i.ctors[t]
	if !ok {
		panic(fmt.Sprintf("No such type registered %s", t))
	}
	return ctor.(func(i *Injector) T)(i)
}

//
// type Constructor[T any] struct {
// 	Fn func(*Injector) T
// }

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
