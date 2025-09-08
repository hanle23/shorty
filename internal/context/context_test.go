package context_test

import (
	"testing"

	"github.com/hanle23/shorty/internal/context"
)

func TestGetContext(t *testing.T) {
	context1 := context.GetContext()
	if context1 == nil {
		t.Error("Context can never be nil")
	}
	context2 := context.GetContext()
	if context1 != context2 {
		t.Error("Get context should be the same instance")
	}
}

func TestSetContext(t *testing.T) {
	c := context.GetContext()
	if c.Debug != false {
		t.Error("Initial debug state should be False")
	}
	c.SetContext(true)
	if c.Debug == false {
		t.Error("Debug should be true after set")
	}
	c1 := context.GetContext()
	if c1.Debug == false {
		t.Error("Debug should not be false after set")
	}
}
