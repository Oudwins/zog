package zog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEq(t *testing.T) {
	validator := Int().EQ(5)
	val, errs, ok := validator.Parse(5)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	val, errs, ok = validator.Parse(4)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}

	validator2 := Float().EQ(5.0)
	val, errs, ok = validator2.Parse(5.0)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	val, errs, ok = validator2.Parse(4.0)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, 4.0, val)
}

func TestGt(t *testing.T) {
	validator := Int().GT(5)
	val, errs, ok := validator.Parse(6)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	val, errs, ok = validator.Parse(5)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	val, errs, ok = validator.Parse(4)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}

	validator2 := Float().GT(5.0)
	val, errs, ok = validator2.Parse(6.0)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	val, errs, ok = validator2.Parse(5.0)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	val, errs, ok = validator2.Parse(4.0)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}

	assert.Equal(t, 4.0, val)
}

func TestGte(t *testing.T) {
	validator := Int().GTE(5)
	val, errs, ok := validator.Parse(6)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	val, errs, ok = validator.Parse(5)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	val, errs, ok = validator.Parse(4)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}

	validator2 := Float().GTE(5.0)
	val, errs, ok = validator2.Parse(6.0)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	val, errs, ok = validator2.Parse(5.0)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	val, errs, ok = validator2.Parse(4.0)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}

	assert.Equal(t, 4.0, val)
}

func TestLt(t *testing.T) {
	validator := Int().LT(5)
	val, errs, ok := validator.Parse(4)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	val, errs, ok = validator.Parse(5)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	val, errs, ok = validator.Parse(6)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}

	validator2 := Float().LT(5.0)
	val, errs, ok = validator2.Parse(4.0)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	val, errs, ok = validator2.Parse(5.0)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	val, errs, ok = validator2.Parse(6.0)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, 6.0, val)
}

func TestLte(t *testing.T) {
	validator := Int().LTE(5)
	val, errs, ok := validator.Parse(4)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	val, errs, ok = validator.Parse(5)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	val, errs, ok = validator.Parse(6)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}

	validator2 := Float().LTE(5.0)
	val, errs, ok = validator2.Parse(4.0)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	val, errs, ok = validator2.Parse(5.0)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	val, errs, ok = validator2.Parse(6.0)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, 6.0, val)
}
