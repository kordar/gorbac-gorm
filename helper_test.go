package gorbac_gorm

import (
	"testing"
	"time"

	"github.com/kordar/gorbac"
)

func TestToAuthRule_ExecuteName(t *testing.T) {
	rule := gorbac.Rule{
		Name:        "r1",
		ExecuteName: "exec",
		CreateTime:  time.Unix(1, 0),
		UpdateTime:  time.Unix(2, 0),
	}
	auth := ToAuthRule(rule)
	if auth.Name != "r1" {
		t.Fatalf("expected name r1, got=%s", auth.Name)
	}
	if auth.ExecuteName != "exec" {
		t.Fatalf("expected execute_name exec, got=%s", auth.ExecuteName)
	}
	if !auth.CreateTime.Equal(rule.CreateTime) || !auth.UpdateTime.Equal(rule.UpdateTime) {
		t.Fatalf("expected times preserved")
	}
}

