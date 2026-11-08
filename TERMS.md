# harnsy plugin catalog contributor terms

Version `2026-11-04.1` · effective on the day the catalog goes live

## In short

- You send a plugin version you built as a **pull request** to a public repository. Anyone can read it, including its history, and forks and copies cannot be called back.
- You keep it. Everyone may use its code under **MIT, Apache-2.0 or CC0-1.0** (the licence in the version's `LICENSE`), the card included.
- Only an author with a valid **harnsy Max** licence can publish. Publishing itself is free; there are no prices yet.
- Your plugin runs on other people's computers, with the rights they grant it. You promise it has no malware, no hidden data flows, and asks only for what its code uses.
- A **check**, then **harnsy's security reviewer** (an AI agent), then the **catalog curator** (an AI agent) look at every pull request. A **person** (the owner) decides each publication, each revoke and each block, and answers appeals. The AI reviews are posted as comments under the owner’s GitHub account, and the first line of each says an AI wrote it.
- Our signature and «Verified by harnsy» mean «reviewed against our rules», not «safe», «approved» or «free of bugs».
- We keep no e-mail of yours. Your GitHub name, your shown name and your description are public.

## 1. Who we are and what this is

1.1 «We» is Pavel Buchnev, individual entrepreneur registered in Georgia, identification number 345621981. The catalog is the public repository https://github.com/harnesy/plugins and its copy on plugins.harnsy.dev: a list of plugins (small Lua programs with pages, actions, hooks and checks) that people can install into harnsy.

1.2 A «plugin» is what you send in one pull request: one version folder (`plugin.yaml`, Lua scripts, page scripts, `LICENSE`, and its card: the `listing` block of `plugin.yaml` and the files under `listing/`, that is the summary, description, icon, docs pages and screenshots). «You» is the person or company that sends it. If you send for a company, you confirm you may bind it. You must be at least 18. You need a GitHub account and a **valid harnsy Max licence** (section 2.4). You must not be a person or company that we are forbidden by law to deal with.

1.3 GitHub is a separate service with its own terms and privacy statement. They apply to what you put on GitHub. We are not GitHub and these terms do not replace theirs. The licence for harnsy Max itself is a separate contract (the Terms of Sale and the licence agreement); these terms do not change it.

## 2. Sending a plugin

2.1 Your publisher id and details are in `publishers/<id>.json`: a shown name, a short description, and the GitHub account(s) allowed to send under this id. Only those accounts may send under it. Shown name and description are public. The ids `harnsy`, `local` and `own` are ours or reserved.

2.2 You send a plugin by opening a pull request that follows `CONTRIBUTING.md`. Opening it, with the confirmation boxes of the template ticked, means you accept these terms in the version named in the template. If you change a box later, nothing already granted is undone (3.4).

2.3 A published version never changes: a change is a new version folder.

2.4 **Licence of the author.** Only an author with a valid harnsy Max licence publishes. For each version your harnsy sends a proof from the licence service to plugins.harnsy.dev («Prepare for catalog»). The first proven version ties your publisher id to your licence; later versions must be proven with the same licence. When your licence ends (it expires, is cancelled or refunded), versions already published stay; new versions are refused until a valid licence proves them. A merged version that has no proof is not signed.

## 3. The licences you give

3.1 **Code, to everyone.** The files of each plugin version (`plugin.yaml`, Lua, page scripts, the card with its icon, docs pages and screenshots, every other file of the folder) are licensed to everyone, worldwide, under the licence in that version's `LICENSE` file: **MIT, Apache-2.0 or CC0-1.0, and no other**. No other licence applies to any file of the folder, the card included. Each file you did not write keeps its own copyright line and licence notice.

3.1a **The card goes out with the code.** The card is part of the version folder; its texts, icon and screenshots are under the same licence as the code, and you give credit as that licence asks. The names and logos «harnsy» and «plugins.harnsy.dev» and our signature are not part of any licence.

3.2 **To us.** You also give us a non-exclusive, worldwide, free licence to host, copy, check, pack, sign, list and display the plugin on plugins.harnsy.dev, put it in the feed and make it available for download and install in harnsy, show your shown name, description, GitHub link and the plugin's card, keep it in archives and backups, and take it off the list. It does not transfer ownership. The catalog packs your version folder as it is; we do not change your files. If a change is needed we ask you in the pull request.

3.3 **You keep your rights.** You may publish the same plugin elsewhere and under other terms for your own copies.

3.4 **Nothing can be called back from history.** The repository is public and keeps its history. Forks, clones, GitHub's copies, web archives, packs people already downloaded and copies others took under your licence stay where they are. Withdrawing or revoking a version (sections 6 and 8) stops us offering it and makes harnsy nodes that read the list turn it off; it does not remove it from history or from copies already taken, and the licence cannot be taken back for them.

3.5 If no copyright exists in some file (for example because a tool wrote it), you grant what rights you have, to that extent.

3.6 GitHub's terms also give GitHub and other GitHub users rights to what you post publicly (for example to view and fork it). Those rights come from GitHub's terms and are not limited by ours.

## 4. What you promise

You promise that, for each plugin version:

(a) **It is yours to license.** You wrote it or have the right to license it as in section 3, including your employer's or client's permission, and the terms of any tool (including AI tools) you used allow it.

(b) **Sources are lawful and credited.** Everything you adapted or copied in (libraries, snippets, text, images) is licensed for this use. Code you did not write is under MIT, Apache-2.0 or CC0 (or another licence that lets you relicense it as 3.1 requires) and its notice is kept. The `Sources` part of the pull request is complete and true.

(c) **It runs safely on other people's computers.** Your plugin runs on the computers of people you do not know, with the rights they grant it. It contains **no malware**: nothing that deletes, locks, damages or changes data the person did not expect, mines coins, scans networks, steals or collects credentials, keys, tokens or browser data, or spreads itself. It does what its card says and nothing else.

(d) **Rights, hosts and data flows are declared and true.** It asks **only** for the rights (`needs`), hosts, keys, MCP servers and chats its code uses, and the card (its description or a docs page) names each in plain words and says why. No hidden data flows: project data, messages or files leave the computer only to a host the plugin declares and the card names, and only for what the card says. If it reads the project and sends data out, the card says so in one sentence.

(e) **Code is plain.** Source code only: no compiled files, binaries, packed archives, obfuscated or minified-only code; no code that downloads and runs other code; nothing that runs by itself that the person did not switch on; no reading of credentials, keys or browser data. A page script runs in the person's dashboard with the person's rights: it does nothing the card does not say.

(f) **It does not get round harnsy.** It does not bypass harnsy's licence checks, editions or limits, the rights the person grants, or the checks of this catalog, and it does not hide from the review.

(g) **Other people's data.** If your plugin reads other people's messages, mail or forms, or sends personal data to a service of yours or of a third party, you are responsible for that under the law that applies to you, the card says so, and you tell the people concerned. If a third-party service's terms (for example a messenger or an API) forbid what the plugin does, it does not do it.

(h) **The card is honest.** It matches the code. No superlatives or promises such as «secure», «safe», «guaranteed»; no claim that harnsy or its signature guarantees anything; no promise of results. The title, name and icon do not use another's trademark or imply that the plugin is made or endorsed by harnsy or anyone else.

(i) **No secrets or personal data** in any file or screenshot: no passwords, keys, licence keys, proofs, tokens. If you committed a secret by mistake, change it at once: it stays in the history.

(j) **Your name is true to you.** Your shown name, description and GitHub link do not impersonate another person or company. The GitHub account is yours.

(k) **What you tell us is accurate:** the licence field, the credits, the description, the card.

## 5. How we review

5.1 Every pull request goes through four steps, in this order:
(1) an **automatic check** that only reads the files and never runs the plugin;
(2) a **code review** by harnsy's security reviewer, an **AI agent** run by us: it reads the Lua and page scripts of exactly this version;
(3) a **card review** by the catalog curator, also an **AI agent** run by us: it matches every line of the card to what the check reports;
(4) a **person**, the owner of the catalog, who reads the reviews and merges, asks for changes or closes.
The AI reviews are posted as **comments under the GitHub account of the owner of the catalog** (Pavel Buchnev), because the AI agents use that account. The first line of every AI review says it was written by an AI; **a comment from that account without this line is the owner’s own**. The AI reviewers recommend: they write comments only, they never use GitHub’s «Approve» or «Request changes» and they never merge. A review by an AI agent never stands in for the person's decision: only the owner’s own merge, or the owner’s own written decision to ask for changes or to close, is the decision.

5.2 Merging is the decision to accept. After a merge, our server packs the version from your folder, checks the proof, signs the pack with the catalog's key and publishes it on plugins.harnsy.dev, normally within an hour. We do not promise a review time and we are not obliged to accept or publish anything.

5.3 «Changes requested» means fix it in the same pull request, as often as you like. Any change to the code changes the version's identity: prepare it again on your node. A rejection is final for that plugin when it is a near-copy or has a harmful purpose.

5.4 **What «Verified by harnsy» and our signature mean.** The code review means: harnsy's security reviewer read this version's Lua and page scripts and found that they match what the plugin declares, and that project data leaves only to the declared hosts. It is a review of one version at one time. It is **not** a guarantee, certification or endorsement of safety, security, lawfulness or fitness, it is not a guarantee against bugs, and you must not say it is. It lifts only one thing for the person who installs: the extra consent tick for a plugin that reads the project and sends data out. Keys, hosts, chats and switching the plugin on in each project stay the person's.

5.5 Everyone can read the plugin's files, including its scripts, in the repository.

5.6 **Reasons.** Whenever we reject, close, revoke, restrict or block, we say in a comment on the pull request: what we did, the facts, the rule in these terms it rests on, whether an automated tool (an AI reviewer) took part, and how to appeal.

## 6. Taking down, revoking and blocking

6.1 We may revoke a published version (an entry in `revoked.json`), close your pull requests, or block you from sending, when: you break section 4; someone makes a well-founded claim about it; the law requires it; a harm is found after publication; or you send a stream of junk.

6.2 **A person decides every revoke and every block**: the owner of the catalog. The AI reviewers can only recommend. This also holds when the plugin looks malicious: we prepare the revoke and tell the owner at once; the owner decides. We do not promise how fast a decision is made. For a harmful plugin (for example one that pulls keys or data, deceives its users, or hides data flows) and for impersonating a real person or company, the block follows the same decision. For a false licence or credits, or for copying, we block after the second case. A block covers you and accounts you control.

6.3 A revoked version is added to `revoked.json`, and its file may be taken off the server. harnsy nodes read the list from time to time and turn a revoked version off and refuse to enable it again. A node that is offline, has the catalog switched off or cannot read the list does not know of the revoke until it can. The reason, as a short category (for example «code rules», «rights claim», «withdrawn by publisher») with no further details about you, is given in the comment on the pull request or Issue that revoked it; the list carries no reason. A block is shown in your publisher file as a date, without reasons. For a rights claim, a legal order, a secret or someone else's personal data we also remove the files from the server.

## 7. Appeals and reports

7.1 **Appeal.** Comment on the pull request, open an Issue in the repository, or write to hello@harnsy.dev. A person, not an AI reviewer, looks at it again and answers with reasons. You can also use any rights you have under law. Issues are public: do not put personal data in them; use e-mail for that.

7.2 **Report.** Anyone can report a plugin. Please write to hello@harnsy.dev (a report contains the reporter's name and e-mail, so an Issue is not the right place) with: why it is unlawful, harmful or breaks these terms; its exact address; the reporter's name and e-mail; and a statement that they believe in good faith that the report is accurate and complete. We may hide the plugin while we look, tell the publisher in the pull request, and tell the reporter what is already public: your shown name and GitHub account. A rights holder may also notify GitHub under GitHub's own rules.

## 8. When you stop

8.1 You may withdraw any version at any time by a pull request that adds it to `revoked.json`. Nodes that read the list turn it off (section 6.3); section 3.4 applies.

8.2 To leave completely, withdraw your plugins and ask at hello@harnsy.dev to remove your publisher file and the records in section 10.3. We remove them from the current files and from plugins.harnsy.dev; history and copies stay (section 10.6). Your publisher id stays reserved while any of your versions is published.

8.3 We may stop running the catalog or delist plugins. Copies already taken stay licensed. Where we can, we say so in the repository first.

## 9. Money, no partnership

9.1 Publishing is free. We pay you nothing and you pay us nothing; the harnsy Max licence is bought separately under its own terms. There are no prices: a plugin with a price is not listed. If we ever offer paid plugins, it will be under new terms you accept first.

9.2 You are not our employee, partner or agent. We choose how to list, order and remove plugins.

## 10. Your data

10.1 **Who is responsible.** Pavel Buchnev, individual entrepreneur, Georgia, for our part. Contact: hello@harnsy.dev. GitHub is responsible for its own service and your account there.

10.2 **What is public, and where.**
- On GitHub: your GitHub account name, your pull requests and comments, the files you send (the card with its icon, docs pages and screenshots included), your shown name and description, and the author name and e-mail written in your commits. GitHub offers a private «noreply» e-mail address for commits; we advise you to use it.
- On plugins.harnsy.dev: your shown name, description, a link to your GitHub account, your plugins with their cards (summary, description, icon, docs pages, screenshots), and the dates of your versions.
- Inside the files and in copies others take: the «© shown name» line and the credit.

10.3 **What we keep ourselves.** A copy of the repository on our server (Hetzner, Germany); review notes (public comments); in the catalog's database: which publisher id belongs to which coded licence ID, and, for each version, its identity (checksum), your GitHub login and the date; a block record, if any; letters you send us. The coded licence ID is a keyed hash of your licence ID that the catalog cannot reverse. We collect **no e-mail** from you. The proof from the licence service is valid for 10 minutes and is not kept; your licence key goes only to the licence service and never into a pull request. The web server of plugins.harnsy.dev logs each request (IP address, time, address requested without its query string, status, program or browser name; no referring page) and keeps the log for 7 days. plugins.harnsy.dev sets no cookies and has no analytics. The automatic check of a pull request runs on GitHub (its logs are GitHub's).

10.4 **Why.** To run the catalog and your contribution (contract), to check that only a paid licence publishes, to prevent abuse, keep block records and defend claims (our legitimate interest), and where the law requires.

10.5 **Who else.** The AI reviews are done by agents run by us on Anthropic's Claude (through Claude Code); they post their comments through the owner’s GitHub account (5.1). They read the pull request and your publisher file, which are already public, and they do not receive e-mail from you. Hetzner Online GmbH runs our servers for us. Reporters as in 7.2 and authorities where the law requires. We do not sell data.

10.6 **How long and your rights.** The publisher ↔ licence link and the version records: as long as a version of yours is published, and afterwards until you ask us to delete them; then only a note «publisher id + coded licence ID» stays, which keeps your id reserved while any of your versions is published. Letters you send us: 6 months after the last one. A block record: 1 year after the decision; then the block mark is removed from the current files and you may send again. The pull request, its comments, the repository history and forks stay on GitHub and with others, and we cannot erase them. You can ask at hello@harnsy.dev to see, correct, delete or restrict what we hold, or object. A published version stays when you ask us to delete the records (its checksum and your login are also in the public history), unless you also ask us to withdraw it. In a serious case (a secret, someone else's personal data, a court order) we may also rewrite the repository history or ask GitHub to remove content. You may complain to the data protection authority of Georgia or, if you live in the EU or UK, to your own.

10.7 **Decisions.** The AI reviewers can recommend changes or rejection in a comment. They cannot approve, merge, revoke or block, whatever account they post from. A person decides.

10.8 Data of the people who install your plugin is not ours and not yours by these terms: harnsy reads the catalog list and downloads packs without reporting what is installed, and data your plugin sends goes to the hosts the person allowed, under your responsibility (4(g)).

## 11. Liability

11.1 The catalog is free and provided «as is». We do not promise it will be available, that a plugin will be accepted or published, or how many people will install it.

11.2 You are responsible for what you send and for what your plugin does on other people's computers. We are not a party to the licence between you and the people who use your plugin. You will cover our reasonable costs from a third party's claim caused by your breach of section 4, as far as the law allows.

11.3 Our review and signature are not a promise about a plugin (5.4). Nothing in these terms limits liability that cannot be limited by law. Otherwise we are not liable for lost profit or indirect loss, and our liability for any other loss is limited to 100 euros.

## 12. Changes to these terms

The terms live in the repository as `TERMS.md`; every version stays in its history and the old versions stay readable. The pull request template names the version. The version current when you open or update a pull request applies to it. A change never reduces a licence already granted for copies already taken.

## 13. Law

These terms are governed by the laws of Georgia and the courts of Tbilisi have jurisdiction. The mandatory consumer rights of the country where you live stay. If a translation differs, the English text applies. Contact: hello@harnsy.dev.
