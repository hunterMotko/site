# Metrics store no visitor identifiers and are not backed up

The site records its own pageviews rather than using a third-party analytics
service. Two restrictions on that are deliberate, and both look like omissions
to anyone who did not decide them.

**No IP addresses and no user agents.** The schema has no column for either.
There is no question they would answer that justifies holding visitor
identifiers on a personal site, and not collecting them is a stronger claim than
any retention policy. Referrer *is* stored, because it answers the one question
the metrics exist for — whether a link shared into Slack or LinkedIn brought
anyone — and it identifies the source, not the person.

**No backups.** The SQLite file lives on a named Docker volume and is not backed
up anywhere. Losing pageview history costs nothing that would change a decision,
and a backup job is a standing thing to maintain. `prebuilt` has offsite backups
because it holds a customer's inventory; this holds counters.

## Consequences

- Do not add an `ip`, `user_agent`, or `session_id` column. A test asserts their
  absence and will fail.
- Metrics degrade to off rather than failing startup: an unset `DB_PATH`, or a
  database that will not open, logs the reason and the site keeps serving pages.
- The volume must survive `docker compose up --build`, or every deploy silently
  resets the history. That is the failure this ADR exists to prevent.
