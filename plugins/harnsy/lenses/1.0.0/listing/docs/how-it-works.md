# How it works

A **lens** is a name and an instruction: what to pull out of a recording. The
plugin starts every project with three you may change or delete: «Meeting
summary», «Decisions and tasks», «Risks and open questions».

1. Attach a sound or video file to an item of the project (a call, a meeting).
2. Open the plugin's page, pick the recording under «Analyze a recording»,
   tick the lenses and press **Analyze**.
3. harnsy transcribes the recording in the background. The page shows
   «Transcribing 3/8» while it works. The transcript is attached to the
   recording's item as `<name>.transcript.md`.
4. The plugin then creates one item for a role of the project: the setting
   `role`; empty, the default, is the project's lead, whatever its role is
   called. Its title is the task («Lenses: write lenses.md and submit —
   <recording>»); the item links the recording and
   the transcript and lists the lenses with their instructions; it never
   copies the transcript's text.
5. The role's agent reads the transcript, writes `lenses.md` with a section
   for each lens and submits the item. The page shows «Done» and a link.

A recording that already has a transcript is not transcribed again.

## Without ffmpeg

harnsy cuts long recordings and video with ffmpeg on its machine. Without it
only a short recording can be transcribed; a long one or a video is marked
«Needs ffmpeg». Install ffmpeg on the machine of harnsy; no restart is needed.

## What it asks for, and why

- `items.read`: to list the project's recordings and to read the state of the
  items it created.
- `items.write`: to create the item for the role.
- `transcribe`: to have harnsy turn the recording into text. harnsy sends the
  recording to the node's dictation connection: a speech service the person
  chose in harnsy (it may be outside this machine). It always counts as
  sending out, so the plugin is «linked» and asks your consent when you
  enable it.
- `sql`: its own file in the project for the lenses and the analyses.

## What it does not do

- It reads no other attachment than sound and video, and no file of the
  project's folder.
- It sends nothing anywhere itself: no web address, key or chat. The
  recording leaves the machine only through harnsy's dictation connection.
- It never puts the transcript's text into an item or a message.

## Limits

At most 10 transcriptions and 10 new items an hour in a project (harnsy's
limits for a plugin), and the node's monthly limit of minutes.
