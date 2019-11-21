package utils

import (
	"bytes"
	"testing"
)

func TestCalculateHash(t *testing.T) {
	dest := "d9CFre4DiBDJMtlBcQWiYSIIl581qqJlbhMM0QTI0qA="
	reader := bytes.NewReader([]byte("this is test2"))
	out := CalculateHash(reader)
	t.Log(out)
	if out != dest {
		t.Error("Calculathash error")
	}
}
