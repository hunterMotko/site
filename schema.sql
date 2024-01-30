CREATE TABLE education (
  id SERIAL PRIMARY KEY,
  school VARCHAR(255),
  title VARCHAR(255),
  to TIMESTAMP,
  from TIMESTAMP
);

CREATE TABLE experience (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255),
  employment VARCHAR(255),
  company VARCHAR(255),
  location VARCHAR(255),
  location_type VARCHAR(255),
  to TIMESTAMP,
  from TIMESTAMP
);

CREATE TABLE skills (
  id SERIAL PRIMARY KEY,
  skill VARCHAR(255)
);

CREATE TABLE projects (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255),
  description VARCHAR(255),
  link VARCHAR(255),
  image VARCHAR(255),
  to TIMESTAMP,
  from TIMESTAMP
);

CREATE TABLE education_skills (
  id SERIAL PRIMARY KEY,
  FOREIGN KEY (skills_id) REFERENCES education(id),
  FOREIGN KEY (education_id) REFERENCES education(id)
)

CREATE TABLE experience_skills (
  id SERIAL PRIMARY KEY,
  FOREIGN KEY (skills_id) REFERENCES experience(id),
  FOREIGN KEY (experience_id) REFERENCES experience(id)
)
