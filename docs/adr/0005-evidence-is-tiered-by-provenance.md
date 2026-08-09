# Evidence on /work is tiered by provenance, not by quality

`/work` sorts projects into three tiers — case study, systems project, bench
project — and the line between them is *where the work came from*, never how
good it is. A case study ran in production for a client. A systems project is
written up at exactly the same depth, decisions and caveat included, but had no
client. A bench project is small and self-directed, listed in a group rather
than written up.

The rule this encodes: **a project cannot be promoted for being good.** Only a
client promotes work into a case study. That produces results that look wrong
from outside — an unpaid daemon with 2,370 lines, tests beside every module, and
five documents of its own is not a case study, while a shed builder's inventory
system is.

That is the intended behaviour. A stranger reading this site cannot verify craft
from a repository; they have no time and no reason to try. What they *can*
verify in seconds is that somebody paid for a thing and it is still running.
"Had a customer" is the scarce, checkable fact, and it is the only one the top
tier asserts. The moment that word also means "was built seriously," it stops
being checkable and the tier stops doing any work at all.

The obvious alternative was one tier with a looser definition — drop the client
requirement from "case study" and let anything sufficiently substantial in. It
was rejected because it spends the credibility of the two projects that have
clients in order to decorate the ones that don't. The three-tier shape keeps
"had a customer" and "was built seriously" as separate facts, which is what
lets unpaid work be presented at full depth without either inflating it or
burying it.

The split also has to survive a career change. Work built while moving toward
embedded and low-level engineering will have no client for some time by
definition. A vocabulary in which unpaid means lesser would have quietly
filed that entire direction under "toy projects."

## Consequences

- A systems project has no `Client` field, deliberately. A `Client: none` row
  would draw the eye to an absence instead of to the work.
- The systems section is guarded by `{{ if .Systems }}` and the tier is
  currently empty. A visible but unpopulated section advertises the gap rather
  than the work; it appears when there is something to put in it.
- Bench projects carry `Year` and `Kind` (`from spec` / `original`). Both are
  load-bearing. Without `Year` the group is an undated pile and cannot show a
  trajectory, which is the only reason it earns space. Without `Kind` a
  reimplementation can be read as original work — and a reviewer who recognises
  `my_ls` and works that out unaided will discount everything else on the page.
  A test asserts the rendered meta line, not the data, so the section cannot
  silently lose the thing that justifies it.
- Provenance is a struct field rather than a sentence of prose, so every future
  entry is forced to declare what it is. Prose would have needed rewriting each
  time and would eventually drift out of date. It also avoids a wrong blanket
  label: the C reimplementations were assumed to be 2022 coursework and are
  mostly 2026 work, and labelling the group by era would have buried the most
  recent thing on the page.
- Do not promote a project between tiers because it improved. Promotion out of
  the bench requires a full writeup; promotion into a case study requires a
  client.
