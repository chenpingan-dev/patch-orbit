package report

import (
	"bytes"
	"github.com/chenpingan-dev/patch-orbit/internal/impact"
	"testing"
)

func TestHTMLIsSelfContainedAndEscapesPayload(t *testing.T) {
	value := impact.Report{Title: "</script><script>alert(1)</script>"}
	html, err := HTML(value)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(html, []byte(value.Title)) {
		t.Fatal("unescaped script payload")
	}
	if bytes.Contains(html, []byte("<script src=")) {
		t.Fatal("external script")
	}
	if !bytes.Contains(html, []byte(`\u003c/script\u003e`)) {
		t.Fatal("missing JSON escape")
	}
}
