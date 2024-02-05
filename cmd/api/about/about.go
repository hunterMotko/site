package about

import ()

type Education struct {
	School string
	Title  string
	To     string
	From   string
}

type Experience struct {
	Title        string
	Employment   string
	Company      string
	Location     string
	LocationType string
	Start        string
	End          string
	Description  string
}

type About struct {
	Intro            string
	Body             string
	Conclusion       string
	EducationList    []Education
	ExperienceList   []Experience
	Skills           SkillList
	OrganizationList []Organization
}

type Skill struct {
	Name  string
	Level int
}

type SkillList []Skill

type Organization struct {
	Name  string
	URL   string
	Image string
}

func getEducation() []Education {
	return []Education{
		{
			School: "Hack Reactor",
			Title:  "Full Stack Engineering",
			To:     "2021",
			From:   "2020",
		},
		{
			School: "Qwasar Silicon Valley",
			Title:  "Full Stack Development, Computer Software Engineering",
			To:     "2020",
			From:   "2019",
		},
		{
			School: "The Last Mile",
			Title:  "Full Stack Development",
			To:     "2019",
			From:   "2018",
		},
	}
}

func getExperience() []Experience {
	return []Experience{
		{
			Title:        "Full Stack Developer",
			Employment:   "Freelance",
			Company:      "Remote",
			Location:     "US",
			LocationType: "Remote",
			Start:        "Feb 2023",
			End:          "Jan 2024",
			Description:  "Full stack development of http://ilpublicdefenderstats.org",
		},
		{
			Title:        "Software Engineering Apprentice",
			Employment:   "Zoom",
			Company:      "Remote",
			Location:     "US",
			LocationType: "Remote",
			Start:        "June 2022",
			End:          "Dec 2022",
			Description:  "Collaborated with the development team to implement subscription and unsubscribe pages and systems. Utilizing Java, Vue.js, Apache, FTL (FreeMarker Template Language), and other relevant technologies to develop robust and user-friendly features.",
		},
		{
			Title:        "Frontend Developer",
			Employment:   "Greenstand",
			Company:      "Remote",
			Location:     "US",
			LocationType: "Remote",
			Start:        "Sep 2021",
			End:          "Mar 2022",
			Description:  "Utilized proficiency in frontend technologies like HTML, CSS, JavaScript, and React to collaborate with design teams, ensuring the development of a user-friendly and efficient messaging interface.",
		},
	}
}

func getSkills() SkillList {
	return SkillList{
		{"JavaScript", 10},
		{"React", 10},
		{"Golang", 7},
		{"POSIX", 7},
		{"Ruby", 8},
		{"Ruby on Rails", 8},
		{"PostgreSQL", 9},
		{"MySql", 9},
		{"MongoDB", 8},
		{"DynamoDB", 7},
		{"Vue", 6},
		{"HTML", 10},
		{"CSS", 9},
		{"CI/CD", 5},
		{"Docker", 6},
		{"Java", 5},
		{"C", 4},
		{"Git", 8},
	}
}

func getOrganizations() []Organization {
	return []Organization{
		{"Hack Reactor", "https://www.hackreactor.com/", "hr.svg"},
		{"Qwasar Silicon Valley", "https://www.qwasar.io/", "qwasar.svg"},
		{"Next Chapter", "https://www.nextchapterbk.com/", "nch.svg"},
		{"The Last Mile", "https://thelastmile.org/", "tlm.svg"},
	}
}

// Intro, Body, Conclusion, EducationList, ExperienceList, Skills, OrganizationList
func GetAbout() About {
	return About{
		Intro: `
Hello, my name is Hunter and I am a Software Engineer with a passion for creating innovative web applications, 
command line interfaces, VIM, and solving real world problems. 
With a strong foundation in technologies such as 
React, Next.js, Node.js, Ruby, Ruby on Rails, Golang, Postgresql, MongoDB, and DynamoDB, 
I possess a diverse skill set that enables me to tackle complex projects.`,
		Body: `
I have three years of professional experience in the field, 
complemented by five years of continuous learning and self-improvement. 
Throughout my career, I have had the opportunity to work with diverse teams, 
collaborating on the development of admin panels, subscribe and unsubscribe platforms, 
and complete websites from design to production deployment.`,
		Conclusion: `
My educational background includes the successful completion of esteemed programs such 
as Hack Reactor, Qwasar Silicon Valley, and The Last Mile. 
These programs have equipped me with a solid understanding of full-stack development, 
backend technologies, databases, CI/CD, and integrations.
I am proud to be associated with organizations such as Next Chapter and The Last Mile, 
which have provided me with valuable industry exposure and networking opportunities. 
Additionally, my commitment to excellence has been recognized through the completion 
of various certifications and qualifications.`,
		EducationList:    getEducation(),
		ExperienceList:   getExperience(),
		Skills:           getSkills(),
		OrganizationList: getOrganizations(),
	}
}
