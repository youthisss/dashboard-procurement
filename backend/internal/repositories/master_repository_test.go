package repositories

import (
	"strings"
	"testing"
)

func TestImportCodeFromNameIsSafeAndStable(t *testing.T) {
	name := "ALIEF SUKSES BERDIKARI, PT / BEKASI 17532"

	first := importCodeFromName("V", name)
	second := importCodeFromName("V", name)

	if first != second {
		t.Fatalf("importCodeFromName() unstable: %q != %q", first, second)
	}
	if len(first) > 50 {
		t.Fatalf("code length = %d, want <= 50: %q", len(first), first)
	}
	if strings.ContainsAny(first, " /\\,") {
		t.Fatalf("code contains unsafe separators: %q", first)
	}
	if !strings.HasPrefix(first, "V-") {
		t.Fatalf("code = %q, want V- prefix", first)
	}
}

func TestImportCodeFromNameKeepsDifferentNamesDistinct(t *testing.T) {
	first := importCodeFromName("M", "North Mill")
	second := importCodeFromName("M", "North/Mill")

	if first == second {
		t.Fatalf("codes for different names matched: %q", first)
	}
}
