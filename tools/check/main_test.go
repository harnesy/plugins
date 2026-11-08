package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakeHarnsy writes a harnsy that answers `plugin check --json` with the given JSON and `plugin pack -o f` by writing f.
func fakeHarnsy(t *testing.T, answer string) string {
	t.Helper()
	f := filepath.Join(t.TempDir(), "harnsy")
	script := "#!/bin/sh\ncase \"$2\" in\ncheck) cat <<'EOF'\n" + answer + "\nEOF\n;;\npack) while [ $# -gt 0 ]; do [ \"$1\" = -o ] && out=\"$2\"; shift; done; echo x > \"$out\";;\nesac\n"
	if err := os.WriteFile(f, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return f
}

// answer builds `harnsy plugin check --json` for acme/crm with the given summary and description (the card), and the listing docs file.
func answer(summary, description string) string {
	return `{"ok": true, "id": "acme/crm", "version": "1.0.0", "manifest_sha256": "abc123", "needs": ["items.read", "http:api.example.com", "key:crm_token"], "linked": true, ` +
		`"listing": {"summary": "` + summary + `", "description": "` + description + `", "icon": {"file": "listing/icon.png", "size": 900, "w": 256, "h": 256}, ` +
		`"docs": [{"slug": "usage", "title": "Usage", "file": "listing/docs/usage.md", "size": 20}], "screens": [{"file": "listing/screens/a.png", "alt": "The board", "w": 800, "h": 500}]}}`
}

var goodAnswer = answer("Deals and clients on a page.", "Keeps deals. Reads the project's items and sends them to api.example.com with the key crm_token.")

func write(t *testing.T, dir, rel, content string) {
	t.Helper()
	f := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(f), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// catalog writes a catalog with one acme/crm 1.0.0; edit lets a test change it.
func catalog(t *testing.T, edit func(write func(rel, content string))) string {
	t.Helper()
	dir := t.TempDir()
	w := func(rel, content string) { write(t, dir, rel, content) }
	w("publishers/acme.json", `{"name": "Acme", "description": "Test publisher.", "github": ["acme"], "github_ids": [42]}`)
	w("plugins/acme/crm/1.0.0/plugin.yaml", "name: crm\nversion: 1.0.0\ntitle: Mini CRM\nneeds: [items.read]\n")
	w("plugins/acme/crm/1.0.0/lua/main.lua", "function run() end\n")
	w("plugins/acme/crm/1.0.0/LICENSE", "MIT License\n\nCopyright (c) 2026 Acme\n")
	w("plugins/acme/crm/1.0.0/listing/docs/usage.md", "# Usage\n\nOpen the Deals page.\n")
	if edit != nil {
		edit(w)
	}
	return dir
}

func run(t *testing.T, dir, harnsy string) []string {
	t.Helper()
	var p problems
	var notes []string
	checkTree(&p, &notes, dir, options{harnsy: harnsy}, nil, account{"acme", 42})
	return p.list
}

func want(t *testing.T, got []string, sub string) {
	t.Helper()
	for _, g := range got {
		if strings.Contains(g, sub) {
			return
		}
	}
	t.Fatalf("no problem mentioning %q in:\n%s", sub, strings.Join(got, "\n"))
}

func TestGoodCatalogPasses(t *testing.T) {
	if got := run(t, catalog(t, nil), fakeHarnsy(t, goodAnswer)); len(got) != 0 {
		t.Fatalf("unexpected problems:\n%s", strings.Join(got, "\n"))
	}
}

func TestNoHarnsyOnlyReads(t *testing.T) {
	var p problems
	var notes []string
	checkTree(&p, &notes, catalog(t, nil), options{}, nil, account{})
	if len(p.list) != 0 || len(notes) == 0 {
		t.Fatalf("problems %v notes %v", p.list, notes)
	}
}

func TestFileRules(t *testing.T) {
	h := fakeHarnsy(t, goodAnswer)
	cases := map[string]struct {
		edit func(func(string, string))
		want string
	}{
		"bundle": {func(w func(string, string)) { w("plugins/acme/crm/1.0.0/x.harnsy", "zip") }, ".harnsy file"},
		"hidden": {func(w func(string, string)) { w("plugins/acme/crm/1.0.0/.env", "A=1\n") }, "hidden"},
		"secret": {func(w func(string, string)) {
			w("plugins/acme/crm/1.0.0/lua/k.lua", "k = 'ghp_"+strings.Repeat("a", 36)+"'\n")
		}, "key or token"},
		"http":     {func(w func(string, string)) { w("plugins/acme/crm/1.0.0/lua/h.lua", "u = 'http://example.com/x'\n") }, "plain http"},
		"minified": {func(w func(string, string)) { w("plugins/acme/crm/1.0.0/web/p.js", strings.Repeat("a=1;", 400)+"\n") }, "minified"},
		"binary":   {func(w func(string, string)) { w("plugins/acme/crm/1.0.0/b.bin", "a\x00b") }, "not UTF-8"},
		"item #":   {func(w func(string, string)) { w("plugins/acme/crm/1.0.0/lua/i.lua", "-- see item #123\n") }, "item numbers"},
		"licence":  {func(w func(string, string)) { w("plugins/acme/crm/1.0.0/LICENSE", "All rights reserved\n") }, "not MIT"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) { want(t, run(t, catalog(t, c.edit), h), c.want) })
	}
	t.Run("missing LICENSE", func(t *testing.T) {
		dir := catalog(t, nil)
		os.Remove(filepath.Join(dir, "plugins/acme/crm/1.0.0/LICENSE"))
		want(t, run(t, dir, h), "LICENSE")
	})
	t.Run("symlink", func(t *testing.T) {
		dir := catalog(t, nil)
		os.Symlink("/etc/passwd", filepath.Join(dir, "plugins/acme/crm/1.0.0/lua/l.lua"))
		want(t, run(t, dir, h), "symbolic link")
	})
}

func TestCardRules(t *testing.T) {
	dir := catalog(t, nil)
	cases := map[string]struct{ answer, want string }{
		"host not named":  {answer("Deals.", "It reads items and sends them out. Key crm_token."), "api.example.com"},
		"key not named":   {answer("Deals.", "It reads items and sends them to api.example.com."), "crm_token"},
		"linked sentence": {answer("Deals.", "Talks to api.example.com with crm_token."), "sends"},
		"boast":           {answer("A safe CRM.", "Reads and sends to api.example.com with crm_token."), "safe"},
		"long summary":    {answer(strings.Repeat("a", 161), "x api.example.com crm_token sends"), "summary"},
		"non-english":     {answer("Сделки", "api.example.com crm_token sends"), "English"},
		"big icon":        {strings.Replace(goodAnswer, `"w": 256, "h": 256`, `"w": 2048, "h": 2048`, 1), "icon"},
		"no listing":      {`{"ok": true, "id": "acme/crm", "version": "1.0.0", "manifest_sha256": "abc", "needs": [], "linked": false}`, "no listing"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) { want(t, run(t, dir, fakeHarnsy(t, c.answer)), c.want) })
	}
	t.Run("names found in the docs page", func(t *testing.T) {
		d := catalog(t, func(w func(string, string)) {
			w("plugins/acme/crm/1.0.0/listing/docs/usage.md", "# Usage\n\nUses api.example.com and the key crm_token; it sends the items.\n")
		})
		if got := run(t, d, fakeHarnsy(t, answer("Deals.", "Keeps deals."))); len(got) != 0 {
			t.Fatalf("a card whose docs name the rights: %v", got)
		}
	})
	t.Run("non-english docs", func(t *testing.T) {
		d := catalog(t, func(w func(string, string)) {
			w("plugins/acme/crm/1.0.0/listing/docs/usage.md", "Привет api.example.com crm_token sends\n")
		})
		want(t, run(t, d, fakeHarnsy(t, answer("Deals.", "Keeps deals."))), "English")
	})
}

func TestHarnsyAnswers(t *testing.T) {
	dir := catalog(t, nil)
	want(t, run(t, dir, fakeHarnsy(t, `{"ok": false, "error": {"code": "plugin_lua", "message": "lua/main.lua:1:3 syntax error"}}`)), "plugin_lua")
	want(t, run(t, dir, fakeHarnsy(t, strings.Replace(goodAnswer, "acme/crm", "acme/other", 1))), "must agree")
	want(t, run(t, dir, fakeHarnsy(t, strings.Replace(goodAnswer, `"1.0.0"`, `"1.0.1"`, 1))), "plugin.yaml says version")
}

func TestPublisherRules(t *testing.T) {
	h := fakeHarnsy(t, goodAnswer)
	want(t, run(t, catalog(t, func(w func(string, string)) {
		w("publishers/local.json", `{"name":"L","description":"d","github":["a"],"github_ids":[1]}`)
	}), h), "not allowed")
	want(t, run(t, catalog(t, func(w func(string, string)) {
		w("publishers/acme.json", `{"name":"A","description":"d","github":["acme"]}`)
	}), h), "github_ids")
	want(t, run(t, catalog(t, func(w func(string, string)) {
		w("publishers/acme.json", `{"name":"A","description":"d","github":["acme"],"github_ids":[1],"extra":1}`)
	}), h), "extra")
}

func TestOldCardFormatIsRefused(t *testing.T) {
	h := fakeHarnsy(t, goodAnswer)
	for _, f := range []string{"plugins/acme/crm/card/1.0.0/card.md", "plugins/acme/crm/screens/1.0.0/1.png", "plugins/acme/crm/icon.png"} {
		t.Run(f, func(t *testing.T) {
			want(t, run(t, catalog(t, func(w func(string, string)) { w(f, "x") }), h), "not used any more")
		})
	}
}

func TestPullRequest(t *testing.T) {
	base := catalog(t, nil)
	pr := func(edit func(func(string, string))) []string {
		dir := catalog(t, edit)
		var p problems
		checkPullRequest(&p, dir, base, account{"acme", 42})
		return p.list
	}
	if got := pr(nil); len(got) != 0 {
		t.Fatalf("an unchanged pull request: %v", got)
	}
	want(t, pr(func(w func(string, string)) { w("plugins/acme/crm/1.0.0/lua/main.lua", "function run() x() end\n") }), "never changes")
	want(t, pr(func(w func(string, string)) {
		w("plugins/acme/crm/1.1.0/plugin.yaml", "name: crm\n")
		w("plugins/acme/crm/1.2.0/plugin.yaml", "name: crm\n")
	}), "plugin versions")
	want(t, pr(func(w func(string, string)) { w("README.md", "changed") }), "outside plugins/")
	want(t, pr(func(w func(string, string)) { w("plugins/other/x/1.0.0/plugin.yaml", "name: x\n") }), "a new publisher")
	var p problems
	checkPullRequest(&p, catalog(t, func(w func(string, string)) { w("plugins/acme/crm/1.1.0/plugin.yaml", "name: crm\n") }), base, account{"mallory", 7})
	want(t, p.list, "belongs to acme")
	want(t, pr(func(w func(string, string)) { w("plugins/acme/crm/card/1.1.0/card.md", "x") }), "not used any more")
	if got := pr(func(w func(string, string)) { w("plugins/acme/crm/1.1.0/plugin.yaml", "name: crm\n") }); len(got) != 0 {
		t.Fatalf("a new version alone: %v", got)
	}
}

func TestRevoked(t *testing.T) {
	base := catalog(t, func(w func(string, string)) { w("revoked.json", `["acme/crm@0.9.0"]`) })
	var p problems
	checkPullRequest(&p, catalog(t, func(w func(string, string)) { w("revoked.json", `[]`) }), base, account{"acme", 42})
	want(t, p.list, "never undone")
	p = problems{}
	checkPullRequest(&p, catalog(t, func(w func(string, string)) { w("revoked.json", `["acme/crm@0.9.0","other/x@1.0.0"]`) }), base, account{"acme", 42})
	want(t, p.list, "only its publisher")
}

func TestBody(t *testing.T) {
	base := t.TempDir()
	write(t, base, "TERMS.md", "# Terms\n\nVersion `2026-11-04.1` · effective\n")
	body := func(s string) []string {
		f := filepath.Join(t.TempDir(), "body.md")
		os.WriteFile(f, []byte(s), 0o644)
		var p problems
		checkBody(&p, base, f)
		return p.list
	}
	boxes := strings.Repeat("- [x] ok\n", minBoxes)
	if got := body("version 2026-11-04.1\n" + boxes); len(got) != 0 {
		t.Fatalf("a good body: %v", got)
	}
	want(t, body("version 2026-01-01.1\n"+boxes), "terms version")
	want(t, body("version 2026-11-04.1\n"+boxes+"- [ ] no\n"), "not ticked")
	want(t, body("version 2026-11-04.1\n- [x] one\n"), fmt.Sprintf("fewer than %d", minBoxes))
}
