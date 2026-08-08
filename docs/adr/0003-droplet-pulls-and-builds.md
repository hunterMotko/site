# The droplet pulls and builds; it does not consume a registry image

Deployment is `git pull` followed by `docker compose up -d --build` on the
DigitalOcean droplet. The sibling `prebuilt` project publishes an image to GHCR
from CI and pulls that image instead, so a reader familiar with it will expect
the same here and wonder why it is missing.

The difference is deliberate. This site is a single low-traffic service on a
droplet that is already paid for and already has the Go toolchain available
through Docker. Building on the box costs a minute per deploy and removes the
registry, the publish step, and the credentials that go with them. There is no
second environment to keep in sync, which is the problem a registry actually
solves.

## Consequences

- CI builds and tests, and verifies the image builds, but does not publish or
  deploy. Deployment stays a manual step, captured as `make deploy` so the
  sequence is written down rather than remembered.
- Configuration cannot ride along in the image. It comes from
  `.env.production`, which exists only on the droplet — see
  [ADR-0004](0004-metrics-are-deliberately-minimal.md) for what the metrics
  half of that config does.
- Revisit if a second environment appears, or if build time on the droplet
  becomes disruptive.
