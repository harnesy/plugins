// check reads the plugin catalog and says what the plugins.harnsy.dev signer and curator would refuse. It never runs a plugin:
// `harnsy plugin check` only compiles the Lua and reads the page scripts.
//
//	check [--dir <catalog checkout>] [--harnsy harnsy]                every publisher and plugin version
//	check --dir <pr checkout> --base <base checkout> --author <login> --author-id <id> [--body <file>] [--proofs-url <url>]
//	                                                                  plus the pull-request rules (and the template's boxes)
//
// Each problem is one line «path: what is wrong — how to fix»; in GitHub Actions also an ::error annotation. Exit 1 on any.
// The harnsy release that runs `plugin check/pack` is the one pinned in the BASE branch (tools/harnsy-version).
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	maxPubName   = 80
	maxPubDesc   = 300
	maxPublisher = 4 << 10
)

var (
	// publisher ids that are never a contributor's: harnsy is ours alone (the base branch's publishers/harnsy.json lists who may send it).
	reserved = map[string]bool{"local": true, "own": true}
	partRE   = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)
	semverRE = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)
	githubRE = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9-]{0,38})$`)
	secretRE = regexp.MustCompile(`AKIA[0-9A-Z]{16}|-----BEGIN [A-Z ]*PRIVATE KEY-----|gh[pousr]_[A-Za-z0-9]{36}|github_pat_[A-Za-z0-9_]{40,}|` +
		`sk-ant-[A-Za-z0-9_-]{20,}|sk-[A-Za-z0-9]{32,}|xox[abprs]-[A-Za-z0-9-]{10,}|AIza[0-9A-Za-z_-]{35}`)
	httpRE      = regexp.MustCompile(`http://[^\s"'<>)\]]+`)
	localHTTPRE = regexp.MustCompile(`^http://(localhost|127\.0\.0\.1)(:[0-9]{1,5})?([/?#]|$)`)
	itemRefRE   = regexp.MustCompile(`\b[Ii]tems? #[0-9]+`)
	dayRE       = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`)
	termsRE     = regexp.MustCompile("(?m)^Version `([0-9]{4}-[0-9]{2}-[0-9]{2}\\.[0-9]+)`")
	boxOffRE    = regexp.MustCompile(`(?m)^\s*[-*] \[ \]`)
	boxOnRE     = regexp.MustCompile(`(?m)^\s*[-*] \[[xX]\]`)
)

// minBoxes: the confirmation boxes of .github/PULL_REQUEST_TEMPLATE.md.
const minBoxes = 9

// options: what the command line gave.
type options struct {
	harnsy    string // the harnsy binary for `plugin check/pack`; empty = those two steps are skipped, and the output says so
	proofsURL string // the catalog's proof endpoint; empty = «proof not checked yet»
}

type problems struct{ list []string }

func (p *problems) add(file, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	p.list = append(p.list, file+": "+msg)
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		fmt.Printf("::error file=%s::%s\n", file, strings.ReplaceAll(msg, "\n", " "))
	}
}

func main() {
	dir := flag.String("dir", ".", "the catalog checkout to check")
	base := flag.String("base", "", "the base branch checkout: turns on the pull-request rules")
	author := flag.String("author", "", "the pull request's GitHub login (with --base)")
	authorID := flag.Int64("author-id", 0, "the pull request author's numeric GitHub id (with --base)")
	body := flag.String("body", "", "a file with the pull request's description: the template's boxes and terms version (with --base)")
	harnsy := flag.String("harnsy", "", "the harnsy binary (pinned in tools/harnsy-version) for `plugin check` and `plugin pack`")
	proofs := flag.String("proofs-url", "", "the catalog's proof endpoint (GET <url>/<manifest_sha256>?github=<login>); empty: not checked yet")
	flag.Parse()
	opt := options{harnsy: *harnsy, proofsURL: strings.TrimRight(*proofs, "/")}
	var p problems
	var notes []string
	// With a base checkout only the versions the pull request adds are run through harnsy; a whole-tree run checks them all.
	only := map[string]bool(nil)
	if *base != "" {
		only = changedVersions(*dir, *base)
	}
	versions := checkTree(&p, &notes, *dir, opt, only, account{*author, *authorID})
	if *base != "" {
		checkPullRequest(&p, *dir, *base, account{*author, *authorID})
		if *body != "" {
			checkBody(&p, *base, *body)
		}
	}
	for _, n := range notes {
		fmt.Println("note: " + n)
	}
	for _, l := range p.list {
		fmt.Println(l)
	}
	if len(p.list) > 0 {
		fmt.Printf("%d problem(s)\n", len(p.list))
		os.Exit(1)
	}
	fmt.Printf("ok: %d plugin version(s)\n", versions)
}

// checkTree checks every publisher file and every plugin version under dir. `only` (the pull request's new versions) limits the
// harnsy runs and the proof lookup; nil = all.
func checkTree(p *problems, notes *[]string, dir string, opt options, only map[string]bool, author account) int {
	pubs, _ := filepath.Glob(filepath.Join(dir, "publishers", "*.json"))
	for _, f := range pubs {
		checkPublisher(p, dir, f)
	}
	n := 0
	vers, _ := filepath.Glob(filepath.Join(dir, "plugins", "*", "*", "*"))
	sort.Strings(vers)
	for _, v := range vers {
		st, err := os.Stat(v)
		if err != nil || !st.IsDir() {
			continue
		}
		if filepath.Base(v) == "card" || filepath.Base(v) == "screens" || filepath.Base(v) == "icon.png" {
			continue // a plugin's own side folders: checked with the plugin
		}
		rel, _ := filepath.Rel(dir, v)
		rel = filepath.ToSlash(rel)
		run := only == nil || only[rel]
		checkVersion(p, notes, dir, v, opt, run, author)
		n++
	}
	checkSidecars(p, dir)
	if b, err := os.ReadFile(filepath.Join(dir, "revoked.json")); err == nil {
		var r []string
		if err := json.Unmarshal(b, &r); err != nil {
			p.add("revoked.json", "not a JSON list of strings — write [\"<publisher>/<name>@<version>\", …]")
		}
		for _, x := range r {
			id, ver, ok := strings.Cut(x, "@")
			pub, name, ok2 := strings.Cut(id, "/")
			if !ok || !ok2 || !partRE.MatchString(pub) || !partRE.MatchString(name) || !semverRE.MatchString(ver) {
				p.add("revoked.json", "%q — want <publisher>/<name>@<MAJOR.MINOR.PATCH>", x)
			}
		}
	}
	return n
}

type publisher struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	GitHub      []string `json:"github"`
	// GitHubIDs: the numeric ids of the accounts in "github". A login can be renamed and taken by someone else; the id
	// stays, so the id decides who may send.
	GitHubIDs []int64 `json:"github_ids"`
	// Banned: the day the publisher was blocked, YYYY-MM-DD, and nothing else (the reason goes in the pull request).
	Banned string `json:"banned,omitempty"`
}

// account is a pull request's author.
type account struct {
	login string
	id    int64
}

// owns: the publisher lists the account by its numeric id; by login only for a publisher file without ids.
func (pb publisher) owns(a account) bool {
	if len(pb.GitHubIDs) == 0 {
		return contains(pb.GitHub, a.login)
	}
	for _, id := range pb.GitHubIDs {
		if id == a.id && a.id != 0 {
			return true
		}
	}
	return false
}

func readPublisher(file string) (publisher, error) {
	var pb publisher
	b, err := os.ReadFile(file)
	if err != nil {
		return pb, err
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	return pb, dec.Decode(&pb)
}

func checkPublisher(p *problems, dir, file string) {
	rel, _ := filepath.Rel(dir, file)
	id := strings.TrimSuffix(filepath.Base(file), ".json")
	if !partRE.MatchString(id) || reserved[id] {
		p.add(rel, "publisher id %q is not allowed — use lowercase letters, digits and dashes; `local` and `own` are taken by every node", id)
	}
	if st, err := os.Stat(file); err == nil && st.Size() > maxPublisher {
		p.add(rel, "larger than %d bytes", maxPublisher)
		return
	}
	pb, err := readPublisher(file)
	if err != nil {
		p.add(rel, "%v — want {\"name\", \"description\", \"github\": [\"<login>\"], \"github_ids\": [<id>]}", err)
		return
	}
	if err := oneLine("name", pb.Name, maxPubName*4); err != nil || utf8.RuneCountInString(pb.Name) > maxPubName {
		p.add(rel, "name: one line, 1..%d characters", maxPubName)
	}
	if utf8.RuneCountInString(pb.Description) > maxPubDesc || strings.ContainsAny(pb.Description, "\n\r") {
		p.add(rel, "description: one paragraph, at most %d characters", maxPubDesc)
	}
	if len(pb.GitHub) == 0 {
		p.add(rel, "github: list the GitHub login(s) allowed to send this publisher's plugins")
	}
	for _, g := range pb.GitHub {
		if !githubRE.MatchString(g) {
			p.add(rel, "github: %q is not a GitHub login", g)
		}
	}
	if len(pb.GitHubIDs) == 0 {
		p.add(rel, "github_ids: list the numeric id of each account in \"github\" (\"id\" in https://api.github.com/users/<login>)")
	}
	for _, id := range pb.GitHubIDs {
		if id <= 0 {
			p.add(rel, "github_ids: %d is not a GitHub account id", id)
		}
	}
	if pb.Banned != "" && !dayRE.MatchString(pb.Banned) {
		p.add(rel, "banned: the day of the block, YYYY-MM-DD, and nothing else — the reason goes in the pull request")
	}
	latin(p, rel, pb.Name+" "+pb.Description)
}

func oneLine(field, v string, max int) error {
	if v == "" || len(v) > max || !utf8.ValidString(v) {
		return fmt.Errorf("%s: one line, 1..%d bytes", field, max)
	}
	for _, r := range v {
		if !unicode.IsPrint(r) {
			return fmt.Errorf("%s: one line, printable characters only", field)
		}
	}
	return nil
}

// latin refuses letters of other scripts: catalog texts are English.
func latin(p *problems, file, s string) {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.Is(unicode.Latin, r) {
			p.add(file, "letter %q: catalog texts are English (Latin letters only)", r)
			return
		}
	}
}

func contains(xs []string, s string) bool {
	for _, x := range xs {
		if strings.EqualFold(x, s) && s != "" {
			return true
		}
	}
	return false
}

// changedVersions: the plugin version folders whose files differ between the pull request and its base («plugins/<pub>/<name>/<ver>»).
func changedVersions(dir, base string) map[string]bool {
	out := map[string]bool{}
	for f := range changedFiles(dir, base) {
		parts := strings.Split(f, "/")
		if len(parts) >= 5 && parts[0] == "plugins" {
			out[strings.Join(parts[:4], "/")] = true
		}
	}
	return out
}

func changedFiles(dir, base string) map[string]bool {
	changed := map[string]bool{}
	for _, root := range []string{dir, base} {
		filepath.WalkDir(root, func(f string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.Name() == ".git" { // a checkout's .git: a folder, or a file in a git worktree
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if d.IsDir() {
				return nil
			}
			rel, _ := filepath.Rel(root, f)
			rel = filepath.ToSlash(rel)
			a, errA := os.ReadFile(filepath.Join(dir, rel))
			b, errB := os.ReadFile(filepath.Join(base, rel))
			if (errA == nil) != (errB == nil) || !bytes.Equal(a, b) {
				changed[rel] = true
			}
			return nil
		})
	}
	return changed
}

// checkPullRequest: one publisher, at most one new plugin version, published versions untouched, the publisher's own GitHub account.
func checkPullRequest(p *problems, dir, base string, author account) {
	changed := changedFiles(dir, base)
	pubs, versions := map[string]bool{}, map[string]bool{}
	for f := range changed {
		parts := strings.Split(f, "/")
		switch {
		case parts[0] == "plugins" && len(parts) >= 5 && semverRE.MatchString(parts[3]):
			pubs[parts[1]] = true
			v := strings.Join(parts[:4], "/")
			versions[v] = true
			if st, err := os.Stat(filepath.Join(base, filepath.FromSlash(v))); err == nil && st.IsDir() {
				p.add(f, "%s is published and never changes — send a new version folder instead", v)
			}
		case parts[0] == "plugins" && len(parts) >= 4: // card/, screens/, icon.png: the old card format
			p.add(f, "not used any more — the card is the `listing` of plugin.yaml and the files under listing/ inside the version folder")
		case parts[0] == "publishers" && len(parts) == 2:
			pubs[strings.TrimSuffix(parts[1], ".json")] = true
		case f == "revoked.json":
			checkRevokedChange(p, dir, base, author)
		default:
			p.add(f, "outside plugins/ and publishers/ — only maintainers change this; a contribution touches its own folders")
		}
	}
	if len(pubs) > 1 {
		p.add("pull request", "touches %d publishers — one pull request, one publisher", len(pubs))
	}
	if len(versions) > 1 {
		p.add("pull request", "touches %d plugin versions — one pull request, one plugin version", len(versions))
	}
	for pub := range pubs {
		owner := filepath.Join(base, "publishers", pub+".json")
		if pb, err := readPublisher(owner); err == nil {
			if pb.Banned != "" {
				p.add("publishers/"+pub+".json", "this publisher is blocked since %s", pb.Banned)
			} else if !pb.owns(author) {
				p.add("publishers/"+pub+".json", "publisher %q belongs to %s — send under your own publisher id", pub, strings.Join(pb.GitHub, ", "))
			}
		} else if pb, err := readPublisher(filepath.Join(dir, "publishers", pub+".json")); err != nil || !contains(pb.GitHub, author.login) || len(pb.GitHubIDs) == 0 || !pb.owns(author) {
			p.add("publishers/"+pub+".json", "a new publisher: add publishers/%s.json with your GitHub login %q in \"github\" and your id %d in \"github_ids\"", pub, author.login, author.id)
		} else if pub == "harnsy" {
			p.add("publishers/harnsy.json", "`harnsy` is ours alone — only the accounts listed in the base branch's publishers/harnsy.json send under it")
		}
	}
}

// checkRevokedChange: entries are only added, and only for the author's own publishers.
func checkRevokedChange(p *problems, dir, base string, author account) {
	read := func(root string) map[string]bool {
		var r []string
		b, _ := os.ReadFile(filepath.Join(root, "revoked.json"))
		json.Unmarshal(b, &r)
		m := map[string]bool{}
		for _, x := range r {
			m[x] = true
		}
		return m
	}
	before, after := read(base), read(dir)
	for x := range before {
		if !after[x] {
			p.add("revoked.json", "%s was removed — a revocation is never undone by a pull request", x)
		}
	}
	for x := range after {
		if before[x] {
			continue
		}
		pub, _, _ := strings.Cut(x, "/")
		if pb, err := readPublisher(filepath.Join(base, "publishers", pub+".json")); err != nil || !pb.owns(author) {
			p.add("revoked.json", "%s: only its publisher withdraws a version", x)
		}
	}
}

// checkBody: the pull request's description keeps the template's confirmations, every box ticked, and names the terms
// version of the base branch's TERMS.md. No TERMS.md on the base: nothing to accept yet.
func checkBody(p *problems, base, file string) {
	terms, err := os.ReadFile(filepath.Join(base, "TERMS.md"))
	if err != nil {
		return
	}
	m := termsRE.FindSubmatch(terms)
	if m == nil {
		p.add("TERMS.md", "no «Version `YYYY-MM-DD.N`» line — maintainers: put it back")
		return
	}
	b, _ := os.ReadFile(file)
	body := strings.ReplaceAll(string(b), "\r\n", "\n")
	switch {
	case !strings.Contains(body, "version "+string(m[1])):
		p.add("pull request", "the description does not name the contributor terms version %s — use the pull request template (Confirmations)", m[1])
	case boxOffRE.MatchString(body):
		p.add("pull request", "a confirmation box is not ticked — tick every box under Confirmations, or say in a comment what stops you")
	case len(boxOnRE.FindAllString(body, -1)) < minBoxes:
		p.add("pull request", "the description has fewer than %d ticked confirmations — use the pull request template", minBoxes)
	}
}
