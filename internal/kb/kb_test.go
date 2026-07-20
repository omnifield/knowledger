package kb

import "testing"

func TestValidRefKind(t *testing.T) {
	t.Parallel()
	cases := map[string]bool{
		RefLink:       true,
		RefTransclude: true,
		RefRelates:    true,
		RefBlocks:     true,
		RefDependsOn:  true,
		RefDuplicate:  true,
		"":            false,
		"bogus":       false,
		"Link":        false,
		"tag":         false, // tags are a separate m2m, not a ref kind
	}
	for k, want := range cases {
		if got := ValidRefKind(k); got != want {
			t.Errorf("ValidRefKind(%q) = %v, want %v", k, got, want)
		}
	}
}

func TestLooksLikeKey(t *testing.T) {
	t.Parallel()
	keys := []string{"KNOW-1", "KNOW-42", "DEVOPSER-128", "A0-9"}
	for _, k := range keys {
		if !LooksLikeKey(k) {
			t.Errorf("LooksLikeKey(%q) = false, want true", k)
		}
	}
	nonKeys := []string{
		"", "know-1", "KNOW", "KNOW-", "-1",
		"550e8400-e29b-41d4-a716-446655440000", // a UUID, not a key
	}
	for _, s := range nonKeys {
		if LooksLikeKey(s) {
			t.Errorf("LooksLikeKey(%q) = true, want false", s)
		}
	}
}
