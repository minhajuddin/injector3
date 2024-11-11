package injector3_test

import (
	"testing"

	"github.com/minhajuddin/injector3"
	"github.com/stretchr/testify/assert"
)

type (
	DBConnectionPoolImpl struct{}
	DBConnectionPool     interface{}
	HomeController       struct {
		ReadReplica DBConnectionPool
		Master      DBConnectionPool
	}

	HomeControllerPrivate struct {
		ReadReplica DBConnectionPool
		awesome     AwesomeString
		Master      DBConnectionPool
	}

	AwesomeString string

	A struct {
		B *B
	}
	B struct {
		C *C
	}
	C struct {
		D D
	}
	D struct {
		E E
	}
	E struct {
		F F
	}
	F struct {
		G *G
	}
	G struct {
		Name AwesomeString
	}
)

// Test the injector
func TestInjector(t *testing.T) {
	i := injector3.NewInjector()
	singleton := &DBConnectionPoolImpl{}
	injector3.Register(i, func(i *injector3.Injector) DBConnectionPool {
		return singleton
	})

	db := injector3.Resolve[DBConnectionPool](i)
	assert.Equal(t, db, singleton)
}

func TestDefaultResolver(t *testing.T) {
	i := injector3.NewInjector()
	singleton := &DBConnectionPoolImpl{}
	injector3.Register(i, func(i *injector3.Injector) DBConnectionPool {
		return singleton
	})

	injector3.Register(i, func(i *injector3.Injector) AwesomeString {
		return "awesome"
	})

	t.Run("Resolve a struct", func(t *testing.T) {
		controller := injector3.Resolve[HomeController](i)
		assert.Equal(t, singleton, controller.Master)
		assert.Equal(t, singleton, controller.ReadReplica)
	})

	t.Run("Resolve a struct ptr", func(t *testing.T) {
		controller := injector3.Resolve[*HomeController](i)
		assert.Equal(t, singleton, controller.Master)
		assert.Equal(t, singleton, controller.ReadReplica)
	})

	t.Run("Deep nested struct", func(t *testing.T) {
		a := injector3.Resolve[A](i)
		assert.Equal(t, AwesomeString("awesome"), a.B.C.D.E.F.G.Name)
	})

	t.Run("Does not set private fields", func(t *testing.T) {
		a := injector3.Resolve[HomeControllerPrivate](i)
		assert.Equal(t, a.awesome, AwesomeString(""))
		assert.Equal(t, a.Master, singleton)
		assert.Equal(t, a.ReadReplica, singleton)
	})
}
