package about

type Education struct {
	School      string
	Title       string
	To          string
	From        string
	Description string
}

type Experience struct {
	Title        string
	Employment   string
	Company      string
	Location     string
	LocationType string
	Start        string
	End          string
	Description  []string
}

// Proof is a single piece of evidence, stated as a claim with the context
// that backs it. These lead the page — the reader should hit concrete
// results before job titles or a technology list.
type Proof struct {
	Claim   string
	Context string
}

// Org is a program worth naming rather than leaving as an unlabeled logo.
//
// DISCLOSURE NOTE: The Last Mile and The Next Chapter are named here, with
// blurbs describing education inside correctional facilities and reentry
// apprenticeships, and The Last Mile also appears in EducationList. Together
// these disclose an incarceration history independently of the Story
// paragraph — so trimming Story alone changes the wording without changing
// what the page reveals. Whether to disclose is a single decision that has to
// be applied to Story, Orgs, and EducationList at once, and it is Hunter's
// call, not a copy edit.
type Org struct {
	Name  string
	Href  string
	Logo  string
	Blurb string
}

type About struct {
	// Positioning is the opening line: what the work is and who it is for.
	// Deliberately not a job title — "Full Stack Software Engineer" is the
	// most contested phrase in the market and invites comparison on years
	// of experience, the weakest available axis.
	Positioning string
	Focus       string
	Proof       []Proof
	// Building is the only field on the site allowed to state a direction that
	// no artifact yet backs — everything else asserts something already
	// shipped. One entry per paragraph: HTML collapses newlines, so a
	// multi-paragraph string would render as one undifferentiated block.
	Building []string
	Story    string
	Orgs     []Org
	// Stack sits last, as a footnote. A technology list is a keyword filter
	// that screens people out more often than it draws them in.
	Stack          []string
	EducationList  []Education
	ExperienceList []Experience
}

func getEducation() []Education {
	return []Education{
		{
			School:      "Hack Reactor",
			Title:       "Advanced Full Stack Engineering Immersive",
			To:          "2022",
			From:        "2022",
			Description: "Completed Hack Reactor's highly rigorous 19-week full-stack engineering program, where I mastered core technologies including JavaScript, Node.js, React.js, Express.js, and MongoDB. The curriculum heavily emphasized professional software development workflows, including daily practice of Test-Driven Development (TDD), Agile methodologies, continuous integration, and disciplined code review practices. Furthermore, the program provided a solid foundation in advanced algorithms and data structures, applying these computer science fundamentals to collaborative, real-world software design and deployment challenges.",
		},
		{
			School:      "Qwasar Silicon Valley",
			Title:       "Full Stack Development, Software Engineering",
			To:          "2022",
			From:        "2021",
			Description: "Completed an intensive, project-based program focused on Full-Stack Development and Computer Software Engineering. The curriculum emphasized mastery of diverse technologies including C, Java, Ruby, and advanced SQL, alongside a rigorous application of advanced algorithms and data structures to solve complex problems. Gained deep proficiency in professional development workflows such as Test-Driven Development (TDD), Agile methodologies, continuous integration, and collaborative peer programming, cultivating robust interpersonal and communication skills essential for high-performance engineering teams.",
		},
		{
			School:      "The Last Mile",
			Title:       "Full Stack Development",
			To:          "2021",
			From:        "2018",
			Description: "Completed comprehensive full-stack web development training with a focus on the high-demand MERN stack (MongoDB, Express.js, React.js, Node.js). The program delivered deep, hands-on experience in engineering robust front-end and back-end solutions, encompassing database management, software optimization, and advanced troubleshooting techniques. The curriculum culminated in collaborative, real-world projects, honing crucial skills in technical problem-solving and delivering high-quality, maintainable code.",
		},
	}
}

func getExperience() []Experience {
	return []Experience{
		{
			Title:        "Full Stack Software Engineer",
			Employment:   "Freelance",
			Company:      "Remote",
			Location:     "US",
			LocationType: "Remote",
			Start:        "February 2023",
			End:          "April 2024",
			// Reordered so the storage work leads. The 50% number was earned in
			// Postgres, and the original ordering credited it to the dashboard.
			Description: []string{
				"Cut page load times by 50% in Postgres rather than in the client — indexing and query tuning against high-volume transactional data. Case handling came out 40% faster off the same work.",
				"Owned the delivery path end to end: GitHub Actions CI, Dockerised deploys, and the DigitalOcean infrastructure underneath them.",
				"Built the data-visualisation dashboards the API fed — React, TypeScript, Next.js — against the same schema I was tuning.",
			},
		},
		{
			Title:        "Software Engineer",
			Employment:   "Zoom",
			Company:      "Remote",
			Location:     "US",
			LocationType: "Remote",
			Start:        "June 2022",
			End:          "December 2022",
			Description: []string{
				"Worked on subscription and renewal workflow services in Java on Apache Tomcat. The renewal changes tracked to a 20% improvement in retention.",
				"Shipped billing-adjacent changes across Marketing, Operations, and Billing without breaking the systems on either side of ours.",
			},
		},
		{
			Title:        "Frontend Developer",
			Employment:   "Greenstand",
			Company:      "Remote",
			Location:     "US",
			LocationType: "Remote",
			Start:        "September 2021",
			End:          "March 2022",
			Description: []string{
				"Wrote the ingestion path for bulk packet data arriving with polymorphic structure and a long tail of edge cases, without dropping records.",
				"Built an offline-first messaging interface for the admin panel: writes queued through network loss and reconciled on reconnect.",
				"Worked with the backend team to pin down API contracts the ingestion code could actually rely on.",
			},
		},
	}
}

func getProof() []Proof {
	return []Proof{
		{
			Claim: "Built ingestion for messy, high-volume data without losing records",
			Context: "Bulk packet data arriving with polymorphic structures and a long tail of " +
				"edge cases. Wrote the processing logic that kept it consistent, and the " +
				"offline-first sync that queued writes through network loss and reconciled " +
				"them on reconnect.",
		},
		{
			Claim: "Cut page load times 50% by fixing the database, not the frontend",
			Context: "Indexing and query tuning against Postgres under high-volume transactional " +
				"load. The frontend numbers — 50% faster loads, 40% faster case handling — " +
				"came out of the storage layer.",
		},
		{
			Claim: "Own the path from commit to running service",
			Context: "GitHub Actions CI and Dockerised deploys, and the Linux box " +
				"underneath them — nginx, TLS renewal, fail2ban, backups. Not handed " +
				"off to someone else's platform team.",
		},
	}
}

func getOrgs() []Org {
	return []Org{
		{
			Name:  "The Last Mile",
			Href:  "https://www.thelastmile.org/",
			Logo:  "images/tlm.svg",
			Blurb: "Software engineering education inside correctional facilities.",
		},
		{
			Name:  "The Next Chapter",
			Href:  "https://www.nextchapterbk.com/",
			Logo:  "images/nch.svg",
			Blurb: "Apprenticeships for engineers reentering the workforce.",
		},
		{
			Name:  "Hack Reactor",
			Href:  "https://www.hackreactor.com/",
			Logo:  "images/hr.svg",
			Blurb: "19-week advanced full-stack immersive.",
		},
		{
			Name:  "Qwasar Silicon Valley",
			Href:  "https://www.qwasar.io/",
			Logo:  "images/qwasar.svg",
			Blurb: "Project-based systems and algorithms program.",
		},
	}
}

func GetAbout() About {
	return About{
		// REVIEW: this positioning was rewritten when the market-data platform
		// was abandoned. The previous line — "I build the data infrastructure
		// other systems run on" — was written to be carried by that project,
		// and with it gone nothing on the site supported the claim. This
		// version is backed by the two projects on /work: one where the whole
		// machine is mine, one commissioned by a state supreme court and still
		// serving. Check it still sounds like you.
		Positioning: "I build production systems and run the machines they live on.",
		Focus: `Go, C, and TypeScript; Postgres and SQLite; Docker and Linux. From the
schema through the service to the reverse proxy and the pipeline that ships it — I would
rather own the whole path than hand half of it to someone else. Remote, or hybrid in
northern Michigan.`,
		Proof: getProof(),
		Building: []string{
			`Building and operating production systems for local businesses — currently a
marketing site and inventory system for a shed builder. Server-rendered Go, pure-Go
SQLite, htmx; one static binary, four direct dependencies, no JavaScript toolchain. It
runs on a DigitalOcean droplet I administer, behind nginx and certbot, with CI gating
format, vet, tests, a smoke run, and govulncheck before an image is published. It has a
real customer, which is the part that makes it interesting.`,
			// The direction, stated here and nowhere else on the site. Deliberately
			// unspecific about the branch — firmware, embedded Linux, and controls are
			// different markets and which one is still open. Naming one now would put a
			// claim on the page that neither the artifacts nor the geography support,
			// and would need rewriting the moment it resolves.
			`Alongside that I have been working further down — C, Linux, and the layers
underneath the ones I have been shipping on. That is the direction; the projects on the
work page are what I can point at today, and the bench is where the rest of it is
happening in the open.`,
		},
		// DECISION REQUIRED — see the note in Orgs below. The writing here is
		// professional; what needs deciding is whether to disclose the route at
		// all, and that cannot be settled by editing this paragraph alone.
		Story: `I came into engineering through The Last Mile and The Next Chapter, then
Qwasar and Hack Reactor. It is not the usual route, and it selected hard for the part of
the job that turns out to be most of the job: reading unfamiliar systems until they make
sense, and staying with a problem long after the interesting part is over. That is what I
still spend most of my time doing.`,
		Orgs: getOrgs(),
		// Trimmed to what is worth being screened *for*. React, Vue, Next.js,
		// Node, and MongoDB came out: they are real experience, they are already
		// visible in the case study and the work history, and as keywords they
		// were the loudest frontend signal on a page arguing the opposite. A
		// stack list is a filter, and every entry that pulls the wrong way costs
		// more than it returns.
		Stack: []string{
			"Go", "C", "Python", "Java", "TypeScript",
			"Postgres", "SQLite",
			"Linux", "Docker", "GitHub Actions", "DigitalOcean",
		},
		EducationList:  getEducation(),
		ExperienceList: getExperience(),
	}
}
