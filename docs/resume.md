# Resume updates

`content/portfolio.yaml` points at the canonical website PDF:

https://mayurathavale.com/mayur_athavale_resume.pdf

The TUI embeds the URL, not the PDF. Publishing a new PDF at that URL updates what visitors receive without changing this YAML, rebuilding the binary or restarting the SSH service.

Use the publisher in the website repository:

```powershell
./scripts/publish-resume.ps1 -Pdf "$HOME/Downloads/mayur_athavale_resume.pdf" -Push
```

Run it from the `mayurathavale.com` checkout. It updates both website PDF copies and pushes the existing Pages deployment. A small private Google Apps Script can sync that public PDF into the existing Drive file after a one-time authorization; no Drive desktop client is required. See [the publishing guide](https://github.com/mayurathavale18/mayurathavale.com/blob/main/docs/resume-publishing.md) for setup and its permission scope.

Keep this canonical URL rather than a dated Drive upload link. Other portfolio content still requires the normal build/deploy process.
