# smtp-service

One small Go service that sends email over SMTP for the other apps (gr-printing). One HTTP endpoint. It runs in Docker on prod-sp-01, port 9081, LAN only. The owner deploys it by hand: there is no auto-deploy.

## Git workflow

- Never push directly to the default branch. Work on a branch (`task/<n>-slug`, `bug/<n>-slug` or `chore/<n>-slug`) and open a PR with `Refs #<n>` in the body.
- Agents never merge. The team lead (the owner's lead-engineer session) reviews and merges.
- Commits and PRs carry **no** `Co-Authored-By: Claude` line and **no** "Generated with Claude Code" footer. This is the owner's rule.

## Privacy

Never log a recipient address, a subject or a body. Logs may say how many recipients and the outcome.

## Writing style

No em or en dashes. Short, plain sentences. This applies to code comments, commit messages, issues and PRs.

## Go style

Standard library first. Declare `err` once per function, then assign it with `=`, including `if err = f(); err != nil`. Never shadow or redeclare `err` with `:=`. Run `gofmt`, `go vet ./...` and `go test ./...` before every push.

## Production

Never touch prod-sp-01 from an agent session.
