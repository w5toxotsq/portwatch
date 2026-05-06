package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
)

func TestNew_ValidIncludeExclude(t *testing.T) {
	f, err := filter.New("22,80,8000-9000", "8080")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNew_InvalidInclude(t *testing.T) {
	_, err := filter.New("abc", "")
	if err == nil {
		t.Fatal("expected error for invalid include")
	}
}

func TestNew_InvalidRange(t *testing.T) {
	_, err := filter.New("9000-8000", "")
	if err == nil {
		t.Fatal("expected error for reversed range")
	}
}

func TestNew_InvalidExclude(t *testing.T) {
	_, err := filter.New("", "not-a-port")
	if err == nil {
		t.Fatal("expected error for invalid exclude")
	}
}

func TestAllow_NoRules(t *testing.T) {
	f, _ := filter.New("", "")
	for _, port := range []int{22, 80, 443, 8080} {
		if !f.Allow(port) {
			t.Errorf("expected port %d to be allowed with no rules", port)
		}
	}
}

func TestAllow_IncludeRange(t *testing.T) {
	f, _ := filter.New("8000-8100", "")
	if !f.Allow(8080) {
		t.Error("expected 8080 to be allowed")
	}
	if f.Allow(80) {
		t.Error("expected 80 to be excluded by include rule")
	}
}

func TestAllow_ExcludePort(t *testing.T) {
	f, _ := filter.New("", "22")
	if f.Allow(22) {
		t.Error("expected 22 to be excluded")
	}
	if !f.Allow(80) {
		t.Error("expected 80 to be allowed")
	}
}

func TestAllow_IncludeAndExclude(t *testing.T) {
	f, _ := filter.New("80,8000-9000", "8080")
	if !f.Allow(80) {
		t.Error("expected 80 to be allowed")
	}
	if !f.Allow(8000) {
		t.Error("expected 8000 to be allowed")
	}
	if f.Allow(8080) {
		t.Error("expected 8080 to be excluded")
	}
	if f.Allow(443) {
		t.Error("expected 443 to be excluded by include rule")
	}
}

func TestApply_FiltersSlice(t *testing.T) {
	f, _ := filter.New("80,443", "")
	input := []int{22, 80, 443, 8080}
	result := f.Apply(input)
	if len(result) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(result))
	}
	if result[0] != 80 || result[1] != 443 {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	f, _ := filter.New("80", "")
	result := f.Apply([]int{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}
