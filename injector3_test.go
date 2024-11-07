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

	controller := injector3.Resolve[HomeController](i)
	assert.Equal(t, controller.Master, singleton)
	assert.Equal(t, controller.ReadReplica, singleton)
}
