# huntermotko.dev

A personal site whose job is to convert a stranger — a hiring manager, a
recruiter, or a prospective client — into a conversation. Its content is
modelled as Go data rather than written into templates, so the wording and the
structure can be argued about separately.

## Language

### Page content

**Positioning line**:
The single sentence stating what Hunter does and for whom. It is not a job
title, and it opens both the home page and `/about`.
_Avoid_: tagline, headline, title, elevator pitch

**Focus**:
The paragraph directly under the positioning line, naming the technologies and
the working arrangement. Supports the positioning line; never replaces it.
_Avoid_: summary, bio, blurb

**Proof block**:
A claim about work done, paired with the context that makes it checkable. Lives
on `/about`. Distinct from a case study: a proof block asserts an outcome in one
or two sentences, where a case study shows the work.
_Avoid_: achievement, bullet, highlight

**Case study**:
A full write-up of one project on `/work` — client, role, period, links,
summary, decisions, stack. Reserved for work that ran in production for someone
other than Hunter.
_Avoid_: project, portfolio piece, showcase

**Decision**:
One engineering choice inside a case study, stated with the reasoning behind it.
The unit a senior reviewer actually reads. Not to be confused with an ADR, which
records a decision about *this* repo.
_Avoid_: highlight, technical detail, challenge

**Caveat**:
The limitation a case study volunteers about itself — the thing a careful
reviewer would otherwise find and wonder why it went unmentioned.
_Avoid_: weakness, disclaimer, caveat emptor

**Systems project**:
A smaller piece of evidence for work below the framework layer. Rendered as a
group rather than individually, because the point they make is collective.
_Avoid_: side project, experiment, toy

**Story**:
The account of how Hunter entered engineering. Deliberately separate from proof
blocks: it explains the route, it does not argue competence.
_Avoid_: bio, background, journey

### Instrumentation

**Pageview**:
One recorded GET of a page route. Excludes static assets, HEAD requests, and any
route not in the page set — so the count means "a person loaded a page," not
"the server handled a request."
_Avoid_: hit, visit, impression, request

**Page route**:
A URL that renders a page and is therefore eligible to be counted and listed in
the sitemap. `/`, `/work`, `/about`. Distinct from `/stats`, `/resume.pdf`, and
static assets, which are none of those things.
_Avoid_: endpoint, path
