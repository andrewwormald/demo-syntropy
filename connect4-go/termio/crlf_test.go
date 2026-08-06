package termio

import (
	"bytes"
	"testing"
)

func TestCRLFWriterTranslatesBareNewlines(t *testing.T) {
	var buf bytes.Buffer
	w := NewCRLFWriter(&buf)

	n, err := w.Write([]byte("line one\nline two\nline three"))
	if err != nil {
		t.Fatalf("Write() returned error %v, want nil", err)
	}
	if want := len("line one\nline two\nline three"); n != want {
		t.Errorf("Write() n = %d, want %d", n, want)
	}

	want := "line one\r\nline two\r\nline three"
	if got := buf.String(); got != want {
		t.Errorf("Write() output = %q, want %q", got, want)
	}
}

func TestCRLFWriterLeavesExistingCRLFAlone(t *testing.T) {
	var buf bytes.Buffer
	w := NewCRLFWriter(&buf)

	if _, err := w.Write([]byte("already\r\nfine\r\n")); err != nil {
		t.Fatalf("Write() returned error %v, want nil", err)
	}

	want := "already\r\nfine\r\n"
	if got := buf.String(); got != want {
		t.Errorf("Write() output = %q, want %q", got, want)
	}
}

func TestCRLFWriterEmptyInput(t *testing.T) {
	var buf bytes.Buffer
	w := NewCRLFWriter(&buf)

	n, err := w.Write(nil)
	if err != nil {
		t.Fatalf("Write() returned error %v, want nil", err)
	}
	if n != 0 {
		t.Errorf("Write() n = %d, want 0", n)
	}
	if buf.Len() != 0 {
		t.Errorf("Write() wrote %q, want nothing", buf.String())
	}
}
