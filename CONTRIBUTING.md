# Sending a plugin

You need a GitHub account, must be at least 18, and need a valid **harnsy Max** licence: only a version proven from a
licensed node is signed. If you send for a company, you may bind it. The public repository keeps its whole history: what you send here, and your GitHub name, can be
read by everyone for good.

How a plugin works, what it may ask for and how it is packed: the harnsy documentation that ships with harnsy
(`/llms/plugins.txt` and `/docs/plugin-v1.md` on your node's dashboard).

## 1. Your publisher id

Pick an id of lowercase letters, digits and dashes (for example `anna-petrova`). The first time, add
`publishers/<id>.json`:

```json
{
  "name": "Anna Petrova",
  "description": "One or two sentences about you or your team, at most 300 characters.",
  "github": ["your-github-login"],
  "github_ids": [1234567]
}
```

`github_ids` holds the numeric id of each account in `github` (the `id` at `https://api.github.com/users/<login>`): only
those accounts can send plugins under this id. Set the same id as `plugins.publisher` on your node. `harnsy`, `local` and
`own` are taken. The first proven version ties the id to your licence; later versions must be proven with the same one.

## 2. Prepare the version on your node

1. Build and test the plugin in a project of your own node (`<project>/.harnsy/plugins/<name>/`).
2. Run `harnsy plugin check <folder> --json` and fix what it reports.
3. On the dashboard, Plugins › your plugin › **Prepare for catalog**, with your GitHub login. The node packs the version,
   proves it with your licence directly to plugins.harnsy.dev, and gives you the version folder. Your licence key and
   the proof never go into the pull request.

Any change to the code afterwards changes the version's identity: prepare it again.

## 3. The pull request

One plugin version per pull request:

- `plugins/<id>/<name>/<version>/` — the folder exactly as «Prepare for catalog» gave it. It must carry a `LICENSE`
  file (MIT, Apache-2.0 or CC0-1.0): put it in the plugin's folder on your node before you prepare. Do not edit the
  folder by hand: the catalog packs it again and the result must match what your node proved.
- **The card** is part of the plugin: the `listing` block of `plugin.yaml` and the files under `listing/`, inside the
  version folder (the node's own check reads and limits it; `harnsy plugin check <folder> --json` shows it under `listing`):

  ```yaml
  listing:
    summary: "Deals and clients on a page of the project."     # one line, at most 160 characters
    description: "What it does, for whom, an example."          # plain text, at most 500 characters
    icon: listing/icon.png                                      # square PNG or WebP, at most 100 KB
    docs:                                                       # up to 8 pages of Markdown, 24 KB each
      - {slug: usage, title: How to use it, file: listing/docs/usage.md}
    screens:                                                    # up to 6 pictures with a short description
      - {file: listing/screens/board.png, alt: The board with five stages}
  ```

  The card is English, and it says **what the plugin asks for and why**, each right in plain words (the description
  or a docs page): every host, key, MCP server and chat it asks for is named by its name. If the plugin both reads
  the project and sends data out, say so in one sentence («Reads the project's items and sends them to api.github.com»).
  It also says what the plugin does not do. The icon and pictures: you must have the right to use them (they go out under the
  plugin's licence), with no personal data or secrets on them. The title is the `title` of `plugin.yaml`.

Rules the check and the review apply:

- The plugin asks only for what its code uses. Every right, host, key, MCP server and chat it asks for is named in
  the card.
- No malware and no hidden data flows: project data, messages or files leave the computer only to a host the plugin
  declares and the card names.
- The plugin does not bypass harnsy’s licence checks, editions or limits.
- If it reads other people’s messages, mail or forms, or sends personal data anywhere, say so in the card; you are
  responsible for that under the law that applies to you.
- Names, titles and icons do not use another’s trademark or suggest that harnsy made or endorses the plugin.
- No superlatives or promises: not «secure», «safe» or «guaranteed».
- No obfuscated or minified code; no code that downloads and runs other code; nothing that hides what it sends.
- No secrets, licence keys, proofs, tokens or personal data in any file.
- Prices: not yet. A plugin with a price is not listed.
- A published version never changes: send a new version folder (`1.0.1`, `1.1.0` …).

## 4. Review

- **Automatic check** on every push: fix what it reports. It also says whether your version was proven.
- **Code review** by harnsy’s AI security reviewer (an AI agent), on the pull request. Fix the findings in the same pull request; a
  code change means preparing the version again on your node.
- **Card review** by the catalog curator (an AI agent of harnsy): it matches every line of your card to what
  `harnsy plugin check` reports and lists what to change.
- **A person decides.** The AI reviews are posted as comments under the owner’s GitHub account; the first line of each says an AI wrote it, and they are advice, never an approval. The owner of the catalog reads the reviews and merges, asks for changes or closes. Merged
  plugins appear on plugins.harnsy.dev within an hour. Nobody but that person promises publication.

## 5. Licence and terms

By opening a pull request with the confirmation boxes ticked you accept the [contributor terms](TERMS.md) in the
version named in the template. In short: you have the right to publish what you send; the code goes out under its MIT,
Apache-2.0 or CC0 licence, the card included; harnsy may check, pack, sign, list and display
the plugin with your shown name and GitHub link; the history and forks stay; our signature and «Verified by harnsy» mean
«reviewed against our rules», not «safe», «approved» or «free of bugs».

## 6. Withdrawing, revoking, appeals

- To withdraw your version, send a pull request adding `"<id>/<name>@<version>"` to `revoked.json`. A node turns the
  version off when it next reads the list.
- A harmful plugin, impersonation, false licences or a flood of junk gets the version revoked and the publisher
  blocked. The owner of the catalog decides every revoke and block, also for a harmful plugin; we prepare it and tell
  him at once, with no promise of a time, and we say why on the pull request.
- **Appeals:** comment on the pull request or open an Issue (Issues are public: no personal data). Reports of harmful
  or unlawful plugins: hello@harnsy.dev, not in an Issue.

## 7. Your data

Your GitHub name, shown name and description are public, and so are the names and e-mails in your commits: use GitHub's
private «noreply» address. plugins.harnsy.dev keeps which licence (as a coded ID) your publisher id belongs to and, per
version, its identity, your GitHub login and the date, besides a server log kept 7 days. Details: TERMS.md section 10. Deleting or correcting: hello@harnsy.dev; we cannot erase git history or
forks. Full terms: [TERMS.md](TERMS.md).
