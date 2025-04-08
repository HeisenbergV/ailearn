package tools

import (
	"fmt"
	"testing"
)

func TestRunNmap(t *testing.T) {
	result, err := RunNmap("127.0.0.1")
	if err != nil {
		t.Fatalf("RunNmap failed: %v", err)
	}
	fmt.Println(result)
}

func TestRunNikto(t *testing.T) {
	result, err := RunNikto("127.0.0.1")
	if err != nil {
		t.Fatalf("RunNikto failed: %v", err)
	}
	fmt.Println(result)
}

func TestRunSqlmap(t *testing.T) {
	result, err := RunSqlmap("127.0.0.1")
	if err != nil {
		t.Fatalf("RunSqlmap failed: %v", err)
	}
	fmt.Println(result)
}
