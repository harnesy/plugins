# harnsy plugin catalog

The source of every plugin on [plugins.harnsy.dev](https://plugins.harnsy.dev): small Lua programs that add pages,
actions, hooks and checks to a [harnsy](https://harnsy.dev) project. A person installs a plugin on their node and
enables it in a project with the rights it asks for.

Everything a plugin is made of is here in plain text, so you can read it before you install it: its `plugin.yaml`
(what it asks for), its Lua scripts and the scripts of its pages.

The plugins here are written by their publishers. harnsy checks them against its rules; that is not a guarantee, and a
plugin runs with the rights you give it.

## How a plugin reaches the catalog

1. The author prepares the version on their own harnsy node («Prepare for catalog»; it needs a valid harnsy Max
   licence) and sends the version folder here by a pull request (see [CONTRIBUTING.md](CONTRIBUTING.md)).
2. **Automatic check.** On every push the check reads the pull request's files and runs `harnsy plugin check` on the
   version with a pinned harnsy. It never runs the plugin.
3. **Code review.** harnsy's AI security reviewer (an AI agent) reads the Lua and page scripts of exactly this version: does the code do
   only what `plugin.yaml` declares, and does project data leave only to the declared hosts?
4. **Card.** The catalog curator (an AI agent) checks that the card says exactly what the plugin does and asks for, in plain words.
5. **A person decides.** The owner of the catalog approves each publication and merges. The AI reviews above are comments under the owner’s GitHub account, each starting with a line that says an AI wrote it; they are advice, never an approval. Then our server packs the
   version, signs it with the catalog key and publishes it on plugins.harnsy.dev. The signing key never leaves our server.

“Verified by harnsy” means that harnsy’s AI security reviewer read this version’s code and found that it matches what
it declares. It is a review of one version, not a guarantee against bugs. See the contributor terms (TERMS.md,
section 5.4).

## Layout

```
plugins/<publisher>/<name>/<version>/          one plugin version, exactly as «Prepare for catalog» gave it: plugin.yaml, lua/…, web/…
plugins/<publisher>/<name>/<version>/listing/  its card: the `listing` block of plugin.yaml, icon, docs pages and screenshots (inside the plugin)
publishers/<publisher>.json                    the publisher: shown name, description, GitHub login(s)
revoked.json                                   withdrawn versions
tools/check/                                   the automatic check
```

A published version never changes; a change is a new version folder.

## Licences

Plugin code: MIT, Apache-2.0 or CC0-1.0, stated in each version's `LICENSE`, the card included (texts, icon,
screenshots);
tools: MIT. Details in [LICENSE](LICENSE). Sending a plugin means accepting the [contributor terms](TERMS.md).

Appeals: comment on the pull request or open an issue (issues are public: no personal data). Reports of unlawful or
harmful plugins and data requests: hello@harnsy.dev.

GitHub is a trademark of GitHub, Inc. harnsy is not affiliated with or endorsed by GitHub.
