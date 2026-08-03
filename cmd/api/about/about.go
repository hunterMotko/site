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
	Building    string
	Story       string
	Orgs        []Org
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
			Description: []string{
				"Engineered and deployed high-performance data visualization dashboards using React, TypeScript, and Next.js, resulting in a 50% reduction in page load times and a 40% increase in user case handling speed.",
				"Spearheaded the full software delivery lifecycle, implementing robust GitHub CI/CD pipelines and Dockerized deployments to establish a scalable, efficient infrastructure on Digital Ocean.",
				"Optimized PostgreSQL database performance through advanced indexing and query tuning, directly contributing to the frontend speed improvements and ensuring data integrity across high-volume transactions.",
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
				"Engineered and optimized dynamic subscription workflow services using Java and Apache Tomcat, directly contributing to a 20% increase in user retention rates by enhancing renewal processes.",
				"Collaborated in an Agile/SCRUM environment with Marketing, Operations, and Billing partners to align technical execution with strategic business goals, ensuring successful, enterprise-grade solution delivery.",
				"Translated complex technical roadmaps into clear, tangible business value for non-technical stakeholders, facilitating consensus and accelerating the deployment of critical user-facing features.",
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
				"Developed a resilient, real-time messaging interface for the admin panel using React, implementing an offline-first architecture to maintain user functionality and queue data for synchronization upon network reconnection.",
				"Tackled complex data ingestion challenges by engineering processing logic for bulk packet data with polymorphic structures and significant edge cases to ensure data integrity and stability.",
				"Collaborated closely with backend engineers to establish and enforce strict API standards, resulting in the optimization of system-wide data flow and a significant increase in transactional efficiency.",
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
			Context: "GitHub Actions CI/CD and Dockerized deploys on DigitalOcean. Not handed " +
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
		Positioning: "I build the data infrastructure other systems run on.",
		Focus: `Ingestion pipelines, the APIs that serve them, and the deployment and
observability around both. Go and Python, Postgres, Docker. Working remote with teams
anywhere in the US.`,
		Proof: getProof(),
		Building: `Currently building a self-hosted market data platform — a Go service
ingesting federal economic data into Postgres and serving it over a documented API, with
migrations, rate limiting, and source-licensing constraints enforced at the query layer.
It runs at effectively zero cost on a single droplet.`,
		// NOTE: this paragraph is a starting draft. It is the most personal
		// thing on the site and should be rewritten in Hunter's own words
		// before this page goes in front of anyone.
		Story: `I came into this field through The Last Mile and The Next Chapter, then
Qwasar and Hack Reactor. That route selects hard for the part of engineering that is
mostly persistence — reading unfamiliar systems until they make sense, and staying with a
problem after the interesting part is over. Most of what I do now is that.`,
		Orgs: getOrgs(),
		Stack: []string{
			"Go", "Python", "Java", "TypeScript",
			"Postgres", "SQLite", "MongoDB",
			"Docker", "GitHub Actions", "DigitalOcean",
			"React", "Vue", "Next.js", "Node",
		},
		EducationList:  getEducation(),
		ExperienceList: getExperience(),
	}
}
