package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	maxFile      = 256 << 10 // any text file in a plugin folder (harnsy's own script limit)
	maxVersion   = 8 << 20   // one version folder in all (harnsy's bundle limit)
	maxFiles     = 254       // files of a bundle
	maxMinLine   = 1200      // a line longer than this in code looks minified
	maxSummary   = 160
	maxIconPx    = 1024
	maxIconBytes = 100 << 10
)

var (
	titleRE    = regexp.MustCompile(`(?m)^title:\s*(.+)$`)
	codeExtRE  = regexp.MustCompile(`\.(lua|js|mjs)$`)
	boastRE    = regexp.MustCompile(`(?i)\b(secure|securely|safe|safely|guarantee|guaranteed)\b`)
	sendsRE    = regexp.MustCompile(`(?i)\bsends?\b`)
	licenceIDs = map[string]string{"MIT": "MIT License", "Apache-2.0": "Apache License", "CC0-1.0": "CC0"}
	// plain rights have no name in the card to look for: only the rights that carry a value are matched by that value
	valued = map[string]bool{"http": true, "key": true, "mcp": true, "chat": true, "feed": true, "sub": true}
)

// checkResult is what `harnsy plugin check --json` says about one folder (the fields the card is matched against).
type checkResult struct {
	OK    bool `json:"ok"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	ID             string   `json:"id"`
	Version        string   `json:"version"`
	ManifestSHA256 string   `json:"manifest_sha256"`
	Needs          []string `json:"needs"`
	Linked         bool     `json:"linked"`
	Listing        *listing `json:"listing"`
}

// listing is the card as `harnsy plugin check --json` reports it (plugin.yaml `listing` + the files under listing/); the node's own check has
// already enforced the format limits (summary 160, description 500, icon, docs, screens), so this tool only adds the editorial rules.
type listing struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Icon        *struct {
		File string `json:"file"`
		Size int    `json:"size"`
		W    int    `json:"w"`
		H    int    `json:"h"`
	} `json:"icon"`
	Docs []struct {
		Slug  string `json:"slug"`
		Title string `json:"title"`
		File  string `json:"file"`
		Size  int    `json:"size"`
	} `json:"docs"`
	Screens []struct {
		File string `json:"file"`
		Alt  string `json:"alt"`
		W    int    `json:"w"`
		H    int    `json:"h"`
	} `json:"screens"`
}

// checkVersion checks one plugins/<publisher>/<name>/<version>/ folder, its card, and (run=true, with a harnsy) what the node's own check says.
// It reports whether harnsy ran.
func checkVersion(p *problems, notes *[]string, dir, v string, opt options, run bool, author account) bool {
	rel, _ := filepath.Rel(dir, v)
	rel = filepath.ToSlash(rel)
	parts := strings.Split(rel, "/")
	pub, name, ver := parts[1], parts[2], parts[3]
	if !partRE.MatchString(pub) || !partRE.MatchString(name) || !semverRE.MatchString(ver) || name == "card" || name == "screens" {
		p.add(rel, "want plugins/<publisher>/<name>/<MAJOR.MINOR.PATCH> in lowercase letters, digits and dashes (a plugin cannot be named card or screens)")
		return false
	}
	if _, err := os.Stat(filepath.Join(dir, "publishers", pub+".json")); err != nil {
		p.add(rel, "no publishers/%s.json — add it with your name, description and GitHub login", pub)
	}
	// Files: text only, no links, no bundles, no hidden files, no minified code.
	var total int64
	nFiles := 0
	hasLicence := false
	var licenceText string
	filepath.WalkDir(v, func(f string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		fr, _ := filepath.Rel(dir, f)
		fr = filepath.ToSlash(fr)
		in := filepath.ToSlash(strings.TrimPrefix(f, v+string(filepath.Separator)))
		info, _ := d.Info()
		total += info.Size()
		nFiles++
		if d.Type()&fs.ModeSymlink != 0 {
			p.add(fr, "a symbolic link — commit the file itself")
			return nil
		}
		if strings.HasSuffix(in, ".harnsy") {
			p.add(fr, "a .harnsy file — a pull request carries sources; the catalog packs them itself")
			return nil
		}
		for _, seg := range strings.Split(in, "/") {
			if strings.HasPrefix(seg, ".") {
				p.add(fr, "a hidden file or folder — leave it out (harnsy refuses it in a catalog build)")
				return nil
			}
		}
		b, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if len(b) > maxFile {
			p.add(fr, "%d bytes, at most %d", len(b), maxFile)
		}
		if !utf8.Valid(b) || bytes.IndexByte(b, 0) >= 0 {
			p.add(fr, "not UTF-8 text — no binaries in a plugin; scripts are source code only")
			return nil
		}
		s := string(b)
		scanText(p, fr, s)
		if codeExtRE.MatchString(in) {
			for _, line := range strings.Split(s, "\n") {
				if len(line) > maxMinLine {
					p.add(fr, "a line of %d characters — obfuscated or minified code is not accepted; send the readable source", len(line))
					break
				}
			}
		}
		if in == "LICENSE" {
			hasLicence, licenceText = true, s
		}
		return nil
	})
	if total > maxVersion {
		p.add(rel, "%d bytes in all, at most %d", total, maxVersion)
	}
	if nFiles > maxFiles {
		p.add(rel, "%d files, at most %d", nFiles, maxFiles)
	}
	licenceID := ""
	switch {
	case !hasLicence:
		p.add(rel+"/LICENSE", "missing — put the plugin's LICENSE (MIT, Apache-2.0 or CC0-1.0) in its folder on your node before «Prepare for catalog»")
	default:
		for id, marker := range licenceIDs {
			if strings.Contains(licenceText, marker) {
				licenceID = id
			}
		}
		if licenceID == "" {
			p.add(rel+"/LICENSE", "not MIT, Apache-2.0 or CC0-1.0 text — only those licences are accepted (TERMS.md)")
		}
	}
	var res checkResult
	ran := false
	if run && opt.harnsy != "" {
		abs, _ := filepath.Abs(v)
		res, ran = runHarnsyCheck(p, rel, abs, pub, name, ver, opt.harnsy)
	} else if run {
		*notes = append(*notes, rel+": harnsy not given (--harnsy), so `plugin check` and `plugin pack` were not run")
	}
	checkCard(p, v, rel, res, ran)
	if run {
		checkProof(p, notes, rel, res, ran, opt, author)
	}
	return ran
}

// scanText: secrets, plain http links and internal item numbers in any text file of the pull request.
func scanText(p *problems, fr, s string) {
	if secretRE.MatchString(s) {
		p.add(fr, "looks like a key or token — remove it; never put secrets in a plugin")
	}
	for _, m := range httpRE.FindAllString(s, -1) {
		if !localHTTPRE.MatchString(m) {
			p.add(fr, "plain http link %s — use https:// (only http://localhost is allowed)", m)
			break
		}
	}
	if m := itemRefRE.FindString(s); m != "" {
		p.add(fr, "internal reference %q — retell it without item numbers", m)
	}
}

// runHarnsyCheck: `harnsy plugin check <abs folder> --json --publisher <pub>`, then `harnsy plugin pack` to be sure it packs.
func runHarnsyCheck(p *problems, rel, abs, pub, name, ver, harnsy string) (checkResult, bool) {
	var res checkResult
	out, err := exec.Command(harnsy, "plugin", "check", abs, "--json", "--publisher", pub).Output()
	if jerr := json.Unmarshal(out, &res); jerr != nil {
		msg := "no JSON answer"
		if err != nil {
			msg = err.Error()
		}
		p.add(rel, "`harnsy plugin check` gave no answer (%s) — is the pinned harnsy installed?", msg)
		return res, false
	}
	if !res.OK {
		p.add(rel, "harnsy plugin check refuses it: %s %s — fix it and run the check on your node", res.Error.Code, oneLineMsg(res.Error.Message))
		return res, false
	}
	if res.ID != pub+"/"+name {
		p.add(rel, "the folder says %s but the path says %s/%s — the folder, the id in plugin.yaml and the publisher must agree", res.ID, pub, name)
	}
	if res.Version != ver {
		p.add(rel, "plugin.yaml says version %s but the folder is %s", res.Version, ver)
	}
	tmp, _ := os.MkdirTemp("", "plugcheck")
	defer os.RemoveAll(tmp)
	if o, err := exec.Command(harnsy, "plugin", "pack", abs, "--publisher", pub, "-o", filepath.Join(tmp, "x.harnsy")).CombinedOutput(); err != nil {
		p.add(rel, "harnsy plugin pack refuses it: %s", oneLineMsg(string(o)))
		return res, false
	}
	return res, true
}

func oneLineMsg(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > 300 {
		s = s[:300] + "…"
	}
	return s
}

// checkCard: the card is the plugin's `listing` (plugin.yaml + listing/), read from `harnsy plugin check --json`. The node's check enforces the format;
// here the rules of CONTRIBUTING.md that need a reader: English, no boasting, every right that carries a value named in the card text, the «sends» sentence
// for a linked plugin. The curator reads the card in full; this reports the misses a machine can see.
func checkCard(p *problems, v, rel string, res checkResult, ran bool) {
	if !ran {
		return
	}
	if res.Listing == nil {
		p.add(rel+"/plugin.yaml", "no listing — add the `listing:` block (summary, description, icon, docs, screens) and the files under listing/; see CONTRIBUTING.md, «The card»")
		return
	}
	l := res.Listing
	title := ""
	if b, err := os.ReadFile(filepath.Join(v, "plugin.yaml")); err == nil {
		if m := titleRE.FindSubmatch(b); m != nil {
			title = strings.Trim(string(m[1]), `"' `)
		}
	}
	if err := oneLine("summary", l.Summary, maxSummary*4); err != nil || utf8.RuneCountInString(l.Summary) > maxSummary {
		p.add(rel+"/plugin.yaml", "listing.summary: one line, 1..%d characters", maxSummary)
	}
	if l.Icon != nil {
		if l.Icon.W > maxIconPx || l.Icon.H > maxIconPx || l.Icon.Size > maxIconBytes {
			p.add(rel+"/"+l.Icon.File, "the icon is %dx%d, %d bytes: at most %dx%d and %d bytes", l.Icon.W, l.Icon.H, l.Icon.Size, maxIconPx, maxIconPx, maxIconBytes)
		}
	}
	text := l.Summary + "\n" + l.Description + "\n" + title
	for _, d := range l.Docs {
		b, err := os.ReadFile(filepath.Join(v, filepath.FromSlash(d.File)))
		if err != nil {
			continue
		}
		text += "\n" + string(b)
		latin(p, rel+"/"+d.File, d.Title+" "+string(b))
	}
	for _, sc := range l.Screens {
		latin(p, rel+"/"+sc.File, sc.Alt)
	}
	latin(p, rel+"/plugin.yaml", l.Summary+" "+l.Description+" "+title)
	if m := boastRE.FindString(text); m != "" {
		p.add(rel+"/plugin.yaml", "«%s» in the card: no «secure», «safe» or «guaranteed» — say what the plugin does and does not do", m)
	}
	lower := strings.ToLower(text)
	for _, need := range res.Needs {
		kind, val, hasVal := strings.Cut(need, ":")
		if !hasVal || !valued[kind] {
			continue
		}
		if !strings.Contains(lower, strings.ToLower(val)) {
			p.add(rel+"/plugin.yaml", "the plugin asks for %q but the card (summary, description, docs) never names %q — say what it is used for", need, val)
		}
	}
	if res.Linked && !sendsRE.MatchString(text) {
		p.add(rel+"/plugin.yaml", "the plugin reads the project and also sends data out (it is «linked»): say so in one sentence with the word «sends», e.g. «Reads the project's items and sends them to api.example.com»")
	}
}

// checkSidecars: card/, screens/ and icon.png beside a version are the old card format; the listing inside the version replaces them.
func checkSidecars(p *problems, dir string) {
	plugs, _ := filepath.Glob(filepath.Join(dir, "plugins", "*", "*"))
	for _, pd := range plugs {
		st, err := os.Stat(pd)
		if err != nil || !st.IsDir() {
			continue
		}
		prel, _ := filepath.Rel(dir, pd)
		prel = filepath.ToSlash(prel)
		for _, old := range []string{"card", "screens", "icon.png"} {
			if _, err := os.Stat(filepath.Join(pd, old)); err == nil {
				p.add(prel+"/"+old, "not used any more — the card is the `listing` block of plugin.yaml and the files under listing/ inside the version folder")
			}
		}
	}
}

// checkProof: the catalog recorded this manifest_sha256 as proven for this GitHub author. The endpoint is not built yet:
// without --proofs-url the check says so and passes (plain note, no finding).
func checkProof(p *problems, notes *[]string, rel string, res checkResult, ran bool, opt options, author account) {
	if !ran || res.ManifestSHA256 == "" {
		return
	}
	if opt.proofsURL == "" {
		*notes = append(*notes, rel+": proof not checked yet (no --proofs-url) — manifest_sha256 "+res.ManifestSHA256)
		return
	}
	u := opt.proofsURL + "/" + url.PathEscape(res.ManifestSHA256) + "?github=" + url.QueryEscape(author.login)
	cl := http.Client{Timeout: 15 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	r, err := cl.Get(u)
	if err != nil {
		p.add(rel, "the proof service did not answer (%v) — try again; this is not a verdict on your plugin", err)
		return
	}
	defer r.Body.Close()
	var ans struct {
		Proven bool `json:"proven"`
	}
	json.NewDecoder(io.LimitReader(r.Body, 4<<10)).Decode(&ans)
	if r.StatusCode != 200 || !ans.Proven {
		p.add(rel, "no proof for this version (manifest %s) from the author %s: on your node run «Prepare for catalog» again, and do not edit the folder afterwards", res.ManifestSHA256, author.login)
	}
}
