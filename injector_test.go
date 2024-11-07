package injector3_test

import (
	"testing"

	"github.com/minhajuddin/injector3"
)

type (
	DBConnectionPoolImpl struct{}
	DBConnectionPool     interface{}
)

// Test the injector
func TestInjector(t *testing.T) {
	i := injector3.NewInjector()
	singleton := &DBConnectionPoolImpl{}
	injector3.Register(i, func(i *injector3.Injector) DBConnectionPool {
		return singleton
	})

	db := injector3.Resolve[DBConnectionPool](i)
	if db != singleton {
		t.Errorf("Expected %p, got %p", singleton, db)
	}
}
