---
title: "Quick start"
description: "Run your first v2ex command."
weight: 30
---

Once `v2ex` is on your `PATH`:

```bash
v2ex --help       # see the command tree
v2ex version      # build info
```

This is a fresh scaffold, so the command tree is just `version` for now. Add
your first real command in `cli/`, build on the `v2ex-cli` library package,
and document it here.

A good first command usually fetches one thing and prints it as JSON, so the
output pipes straight into `jq` and the rest of your tools.
