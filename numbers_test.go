package zog

import "testing"

func TestEq(t *testing.T) {
	validator := Int().EQ(5)
	errs, ok := validator.Validate(5)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs, ok = validator.Validate(4)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}

	validator2 := Float().EQ(5.0)
	errs, ok = validator2.Validate(5.0)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs, ok = validator2.Validate(4.0)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
}

func TestGt(t *testing.T) {
	validator := Int().GT(5)
	errs, ok := validator.Validate(6)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs, ok = validator.Validate(5)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	errs, ok = validator.Validate(4)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}

	validator2 := Float().GT(5.0)
	errs, ok = validator2.Validate(6.0)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs, ok = validator2.Validate(5.0)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	errs, ok = validator2.Validate(4.0)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
}

func TestGte(t *testing.T) {
	validator := Int().GTE(5)
	errs, ok := validator.Validate(6)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs, ok = validator.Validate(5)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs, ok = validator.Validate(4)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}

	validator2 := Float().GTE(5.0)
	errs, ok = validator2.Validate(6.0)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs, ok = validator2.Validate(5.0)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs, ok = validator2.Validate(4.0)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
}

func TestLt(t *testing.T) {
	validator := Int().LT(5)
	errs, ok := validator.Validate(4)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs, ok = validator.Validate(5)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	errs, ok = validator.Validate(6)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}

	validator2 := Float().LT(5.0)
	errs, ok = validator2.Validate(4.0)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs, ok = validator2.Validate(5.0)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	errs, ok = validator2.Validate(6.0)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
}

func TestLte(t *testing.T) {
	validator := Int().LTE(5)
	errs, ok := validator.Validate(4)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs, ok = validator.Validate(5)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs, ok = validator.Validate(6)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}

	validator2 := Float().LTE(5.0)
	errs, ok = validator2.Validate(4.0)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs, ok = validator2.Validate(5.0)
	if !ok || len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs, ok = validator2.Validate(6.0)
	if ok || len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
}
