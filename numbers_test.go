package zog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEq(t *testing.T) {
	validator := Int().EQ(5)
	dest := 0
	errs := validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(4, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}

	dest2 := 0.0
	validator2 := Float().EQ(5.0)
	errs = validator2.Parse(5.0, &dest2)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator2.Parse(4.0, &dest2)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, 4.0, dest2)
}

func TestGt(t *testing.T) {
	validator := Int().GT(5)
	dint := 0
	errs := validator.Parse(6, &dint)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(5, &dint)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	errs = validator.Parse(4, &dint)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, 4, dint)

	dfl := 0.0
	validator2 := Float().GT(5.0)
	errs = validator2.Parse(6.0, &dfl)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator2.Parse(5.0, &dfl)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	errs = validator2.Parse(4.0, &dfl)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}

	assert.Equal(t, 4.0, dfl)
}

func TestGte(t *testing.T) {
	dint := 0
	validator := Int().GTE(5)
	errs := validator.Parse(6, &dint)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(5, &dint)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(4, &dint)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	dfl := 0.0
	validator2 := Float().GTE(5.0)
	errs = validator2.Parse(6.0, &dfl)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator2.Parse(5.0, &dfl)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator2.Parse(4.0, &dfl)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}

	assert.Equal(t, 4.0, dfl)
}

func TestLt(t *testing.T) {
	dint := 0
	validator := Int().LT(5)
	errs := validator.Parse(4, &dint)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(5, &dint)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	errs = validator.Parse(6, &dint)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	dfl := 0.0
	validator2 := Float().LT(5.0)
	errs = validator2.Parse(4.0, &dfl)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator2.Parse(5.0, &dfl)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	errs = validator2.Parse(6.0, &dfl)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, 6.0, dfl)
}

func TestLte(t *testing.T) {
	dint := 0
	validator := Int().LTE(5)
	errs := validator.Parse(4, &dint)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(5, &dint)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(6, &dint)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	dfl := 0.0
	validator2 := Float().LTE(5.0)
	errs = validator2.Parse(4.0, &dfl)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator2.Parse(5.0, &dfl)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator2.Parse(6.0, &dfl)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, 6.0, dfl)
}
