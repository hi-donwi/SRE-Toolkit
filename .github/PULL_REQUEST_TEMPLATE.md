## What this changes

<!-- And why. The diff shows what; explain the reason. If it fixes a bug, say
     what the broken behaviour was. -->

## How it was verified

<!-- Not "tests pass" — what did you actually check? If you added a rule, where
     did the fixture come from? If you fixed a bug, how did you reproduce it
     first? -->

## Checklist

- [ ] `gofmt -l .` prints nothing
- [ ] `go vet ./...` is clean
- [ ] `go test -race ./...` passes
- [ ] New or changed rules have an evaluator test, covering the healthy case too
- [ ] New rule IDs are added to `SPEC.md` section 4
- [ ] `make licenses` re-run and committed, if dependencies changed
- [ ] Behaviour changes that would surprise an existing user are called out above
