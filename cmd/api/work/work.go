// Package work holds the contents of /work: the projects that are evidence
// rather than assertion.
//
// The ordering rule for this page is that a reader should be able to verify a
// claim before they are asked to believe a summary. So each case study leads
// with a link — a live URL first where one exists, because a stranger can check
// a running site in five seconds and cannot check a repo in five minutes — and
// the engineering decisions come before the stack list.
//
// Evidence is tiered, and the tiers are about provenance rather than quality.
// A CaseStudy ran in production for a client. A SystemsProject is written up at
// the same depth but had no client. A BenchProject is small and listed in a
// group. Keeping "somebody paid for this" separate from "this was built
// seriously" is the whole point: collapsing them would either inflate the
// unpaid work or bury it.
//
// Deliberately absent: guided coursework. Boot.dev projects (gator, pokedex,
// chirpy, notely, bookbot, static_site_gen, page-crawler, and the rest) exist in
// thousands of byte-identical copies on GitHub, and reviewers who recognise
// them read them as curriculum, not work. Note this is not a rule against
// exercises — several bench projects are implementations written from a
// published specification, which is a different thing and is labelled as such.
// The distinction is whether the shape of the solution was handed to you.
package work

// Kind records how a bench project came about. Rendered on the page rather than
// left for the reader to infer: someone who recognises my_ls as a classic
// exercise and finds it presented as original work discounts everything else on
// the page, and the label costs nothing that inference was not going to take
// anyway.
type Kind string

const (
	FromSpec Kind = "from spec"
	Original Kind = "original"
)

// Decision is one engineering choice and the reasoning behind it. These are the
// payload of a case study — a senior reviewer skims the stack and reads these,
// because the stack says what you have touched and a decision says how you
// think.
type Decision struct {
	Title string
	Body  string
}

// CaseStudy is a project written up in full, for a client. Live and Repo are
// optional; the template omits any link whose URL is empty, so an unpublished
// project never renders a dead anchor.
type CaseStudy struct {
	Name    string
	Tagline string
	Client  string
	Role    string
	Period  string
	Live    string
	Repo    string
	// Summary is one entry per paragraph. HTML collapses newlines, so a single
	// multi-paragraph string renders as one undifferentiated block — the breaks
	// have to survive as structure, not as whitespace.
	Summary []string
	// Caveat states, in the project's own writeup, whatever a careful reviewer
	// would otherwise discover and wonder why you omitted. Volunteering it
	// converts a weakness into the kind of tradeoff conversation that separates
	// people who have run things from people who have read about running them.
	Caveat    string
	Decisions []Decision
	Stack     []string
}

// SystemsProject is a full writeup — the same depth as a case study, decisions
// and caveat included — for work that had no client. It carries no Client field
// because that is precisely what distinguishes it, and a "Client: none" row
// would draw the eye to an absence rather than to the work.
//
// The tier is defined and rendered but currently empty: the section does not
// appear until there is something in it, because a visibly waiting slot
// advertises the gap rather than the work.
type SystemsProject struct {
	Name    string
	Tagline string
	Role    string
	Period  string
	Live    string
	Repo    string
	Summary []string
	Caveat  string

	Decisions []Decision
	Stack     []string
}

// BenchProject is small, self-directed work: listed as a group rather than
// written up, because individually they are modest and collectively they show
// something a case study cannot — that the building never stopped, and which
// direction it has been moving.
//
// Year is what makes the group work. Without it the list reads as one undated
// pile, and a pile cannot show a trajectory; with it the same entries show C in
// 2026 sitting above Go in 2025, which is the actual claim.
type BenchProject struct {
	Name  string
	Repo  string
	Lang  string
	Year  string
	Kind  Kind
	Blurb string
}

type Work struct {
	Intro        string
	CaseStudies  []CaseStudy
	SystemsIntro string
	Systems      []SystemsProject
	BenchIntro   string
	Bench        []BenchProject
	Closing      string
}

func getCaseStudies() []CaseStudy {
	return []CaseStudy{
		{
			Name:    "Prebuilt Sheds LLC",
			Tagline: "Marketing site and inventory system for a shed builder.",
			Client:  "Prebuilt Sheds LLC",
			Role:    "Sole engineer — build, deploy, and operations",
			Period:  "2026 — ongoing",
			Repo:    "https://github.com/hunterMotko/prebuilt",
			Summary: []string{
				`Server-rendered Go, SQLite, and htmx: one static binary, four direct
dependencies, no JavaScript toolchain and no CSS framework. Public pages plus a
management area behind auth, CSRF, and a separate body limit.`,
				`It runs on a DigitalOcean droplet I administer — nginx and certbot in
front, fail2ban, scheduled and offsite backups — with CI gating format, vet,
tests, a smoke run, and govulncheck before publishing an image to GHCR.`,
				`It has a real customer, which is the part that makes it engineering
rather than a demo. Nobody files a bug against a portfolio project.`,
			},
			Decisions: []Decision{
				{
					Title: "Pure-Go SQLite, so the binary is actually static",
					Body: `modernc.org/sqlite compiles without cgo, so CGO_ENABLED=0 produces a
static binary and the runtime image needs no libc matching. It trades some raw
speed for a materially simpler build and a smaller image — at this traffic the
speed is irrelevant and the build simplicity is not.`,
				},
				{
					Title: "Connection pool pinned to one connection",
					Body: `SQLite's PRAGMA foreign_keys is per-connection, and database/sql pools
connections transparently. Without SetMaxOpenConns(1), the ON DELETE CASCADE on
inventory photos would have applied only on whichever connection happened to
have run the pragma — a corruption bug that appears under concurrency and never
in testing. SQLite serializes writes anyway, so the pool was buying nothing.`,
				},
				{
					Title: "Money is int64 cents, rounded rather than truncated",
					Body: `int64(dollars * 100) silently loses a cent whenever the float
multiplication lands just below the integer. Prices are stored and computed in
cents, and the conversion rounds. Storing money in a float is the bug that shows
up in an invoice six months later and cannot be reconstructed.`,
				},
				{
					Title: "WAL and a busy timeout, after a real lock race",
					Body: `A SQLite lock race turned CI red intermittently on main. The fix was
enabling WAL and giving every connection — including the smoke script's sqlite3
— a busy timeout. Intermittent CI failures get muted rather than fixed more
often than any other class of bug; this one was reproduced and closed.`,
				},
				{
					Title: "Every CI job body is a make target",
					Body: `The pipeline holds no commands of its own. When CI and the developer
run different things they drift, and "passes locally, red in CI" becomes
routine — so make ci reproduces the entire workflow on a laptop.`,
				},
				{
					Title: "The Dockerfile is the source of truth for the Go version",
					Body: `setup-go treats go.mod's go directive as an exact pin rather than the
minimum it is, so a directive of 1.25.0 installed exactly 1.25.0 and govulncheck
reported 26 stdlib findings already fixed in later patches — none of which were
present in the shipped image. CI now parses the Go version out of the Dockerfile
and installs the newest patch of that minor, mirroring how golang:<minor>-alpine
resolves.`,
				},
			},
			Stack: []string{"Go", "Echo", "SQLite", "htmx", "Docker", "GitHub Actions", "nginx", "certbot", "DigitalOcean"},
		},
		{
			Name:    "Illinois Public Defender Statistics",
			Tagline: "Statewide public-defense data, mapped, for all 102 Illinois counties.",
			Client:  "Illinois Supreme Court & the Administrative Office of Illinois Courts",
			Role:    "Sole engineer on the application",
			Period:  "2024 — still in production",
			Live:    "https://ilpublicdefenderstats.org",
			Repo:    "https://github.com/hunterMotko/il_pd_st",
			Summary: []string{
				`Commissioned by the Illinois Supreme Court and the Administrative
Office of Illinois Courts, after the Sixth Amendment Center's evaluation of the
state's public defense system found structural deficiencies in oversight and
independence. The application gives public defenders, researchers, and the
courts a way to explore resource distribution across all 102 counties.`,
				`The interesting problem was serving one dataset through three different
geographic hierarchies — counties, judicial circuits, and appellate districts —
each needing its own identifiers, colour scales, and legend breaks to drive
choropleth rendering. That shape drove the API: twelve route handlers, organised
by hierarchy rather than by table, returning exactly what the map layer needs
and nothing else. Query time came out of PostgreSQL aggregates, functions, and
indexes rather than out of caching in front of slow queries.`,
				`It has been serving continuously since delivery.`,
			},
			Caveat: `Worth stating plainly: I was the sole code contributor, but not the
whole project. A Northwestern University team collected and analysed the
underlying data, and the court commissioned the work. Worth stating too that
this ran on managed Postgres and App Platform build-on-push — real deployment,
but a platform I consumed rather than administered. Prebuilt is the other half
of that answer, and having shipped both is what makes the comparison worth
having.`,
			Decisions: []Decision{
				{
					Title: "The API is shaped by the map, not by the tables",
					Body: `Separate colors, legend, and ids endpoints per hierarchy, plus a
detail route per feature. A generic CRUD surface would have made the client
assemble a choropleth from four round trips and reimplement the classification
breaks in JavaScript; putting the breaks server-side kept the legend and the
fills provably consistent.`,
				},
				{
					Title: "Query optimisation in Postgres, not in a cache",
					Body: `Aggregates, functions, and indexes brought response times down at the
source. A cache in front of a slow query would have hidden the problem and added
an invalidation bug — and for data that updates rarely and is read constantly,
the index is the correct answer.`,
				},
				{
					Title: "Performance treated as a measured requirement",
					Body: `First Contentful Paint, Largest Contentful Paint, Total Blocking
Time, and layout shift were tracked as numbers, not impressions. Interactive
maps are heavy by default, and "feels fine on my laptop" is not a finding.`,
				},
			},
			Stack: []string{"TypeScript", "Next.js", "React", "PostgreSQL", "Tailwind", "DigitalOcean App Platform"},
		},
	}
}

// getBench returns bench projects newest first. The ordering is the argument:
// read top to bottom it shows C in 2026 above Go in 2025, which is the
// trajectory. Sorted any other way it is a list of repositories.
func getBench() []BenchProject {
	return []BenchProject{
		{
			Name:  "my_blockchain",
			Repo:  "https://github.com/hunterMotko/my_blockchain",
			Lang:  "C",
			Year:  "2026",
			Kind:  Original,
			Blurb: "A command REPL over a linked-list block store in C, with file-backed persistence and an explicit synced/unsynced state.",
		},
		{
			Name:  "Unix primitives in C",
			Repo:  "https://github.com/hunterMotko?tab=repositories&q=my_",
			Lang:  "C",
			Year:  "2024 — 2026",
			Kind:  FromSpec,
			Blurb: "ls, tar, and printf implemented against the specification rather than wrapped — argument parsing, flags, and edge cases included.",
		},
		{
			Name:  "csvq",
			Repo:  "https://github.com/hunterMotko/csvq",
			Lang:  "Go",
			Year:  "2025",
			Kind:  Original,
			Blurb: "A jq-equivalent for CSV: a small query language with its own parser, reading from stdin or a file.",
		},
		{
			Name:  "go_networking",
			Repo:  "https://github.com/hunterMotko/go_networking",
			Lang:  "Go",
			Year:  "2025",
			Kind:  Original,
			Blurb: "Hand-rolled TCP, UDP, TLS, HTTP, and unix-socket servers — the layer most web work sits on top of without looking at.",
		},
		{
			Name:  "bdg",
			Repo:  "https://github.com/hunterMotko/bdg",
			Lang:  "Go",
			Year:  "2025",
			Kind:  Original,
			Blurb: "A terminal budgeting tool, built because the alternatives all wanted a bank login.",
		},
		{
			Name:  "my_sqlite",
			Repo:  "https://github.com/hunterMotko/my_sqlite",
			Lang:  "Ruby",
			Year:  "2022",
			Kind:  FromSpec,
			Blurb: "A working subset of SQLite — parser, storage, and query execution — written to find out what the file format is actually doing.",
		},
	}
}

func GetWork() Work {
	return Work{
		Intro: `Two projects in production, and the smaller work underneath them. Both
of the first two are running right now and one of them you can open in a new
tab, which is the only portfolio claim that verifies itself.`,
		CaseStudies: getCaseStudies(),

		// Empty by design — see the note on SystemsProject. The section renders
		// only when there is something in it.
		SystemsIntro: `Built to the same standard, without a client.`,
		Systems:      nil,

		BenchIntro: `Small on purpose. They exist because the fastest way to stop
guessing about a layer is to implement it, and because I have never stopped
building things for the sake of building them. Newest first.`,
		Bench: getBench(),

		Closing: `Currently open to remote platform, backend, and systems work, or
hybrid in northern Michigan.`,
	}
}
