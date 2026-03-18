package components

import (
	"strings"
	"testing"
)

func TestRenderASCII_ProducesCorrectLines(t *testing.T) {
	result := RenderASCII("AB")
	lines := strings.Split(result, "\n")
	if len(lines) != 6 {
		t.Fatalf("expected 6 lines, got %d", len(lines))
	}
}

func TestRenderASCII_HandlesSpaces(t *testing.T) {
	result := RenderASCII("A B")
	lines := strings.Split(result, "\n")
	if len(lines) != 6 {
		t.Fatalf("expected 6 lines, got %d", len(lines))
	}
	// Each line should contain the space glyph between A and B
	for _, line := range lines {
		if len(line) == 0 {
			t.Fatal("unexpected empty line")
		}
	}
}

func TestRenderASCII_CaseInsensitive(t *testing.T) {
	upper := RenderASCII("HELLO")
	lower := RenderASCII("hello")
	if upper != lower {
		t.Fatal("expected case-insensitive rendering")
	}
}

func TestRenderASCII_UnknownChar(t *testing.T) {
	// Should not panic on unknown characters
	result := RenderASCII("A!B")
	lines := strings.Split(result, "\n")
	if len(lines) != 6 {
		t.Fatalf("expected 6 lines, got %d", len(lines))
	}
}

func TestASCIIWidth(t *testing.T) {
	w := ASCIIWidth("A")
	if w != 6 {
		t.Fatalf("expected width 6 for 'A', got %d", w)
	}

	w2 := ASCIIWidth("AB")
	// A=6 + spacing=1 + B=6 = 13
	if w2 != 13 {
		t.Fatalf("expected width 13 for 'AB', got %d", w2)
	}
}

func TestRenderASCII_FullName(t *testing.T) {
	result := RenderASCII("Mayur Athavale")
	lines := strings.Split(result, "\n")
	if len(lines) != 6 {
		t.Fatalf("expected 6 lines, got %d", len(lines))
	}

	// Print for visual inspection
	t.Logf("Rendered 'Mayur Athavale':\n%s", result)
}
