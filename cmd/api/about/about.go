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

type About struct {
	Summary        string
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

func GetAbout() About {
	return About{
		Summary: `
I am a Full Stack Software Engineer with three years of professional experience and a strong foundation in full-stack development. My expertise spans React.js, Vue.js, Next.js, Node.js, Golang, and Java, supported by robust database skills in Postgresql and MongoDB. A graduate of Hack Reactor, Qwasar Silicon Valley, and The Last Mile, I have successfully delivered complex projects—from admin panels to production-ready websites—focusing on scalable architecture and seamless CI/CD integration.`,
		EducationList:  getEducation(),
		ExperienceList: getExperience(),
	}
}
