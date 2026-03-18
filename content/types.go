package content

type Portfolio struct {
	Name     string       `yaml:"name"`
	Title    string       `yaml:"title"`
	Location string       `yaml:"location"`
	About    string       `yaml:"about"`
	Resume   string       `yaml:"resume"`

	Experience []Experience `yaml:"experience"`
	Projects   []Project    `yaml:"projects"`
	Skills     SkillSet     `yaml:"skills"`
	Contact    Contact      `yaml:"contact"`
}

type Experience struct {
	Role       string   `yaml:"role"`
	Company    string   `yaml:"company"`
	Location   string   `yaml:"location"`
	Period     string   `yaml:"period"`
	Highlights []string `yaml:"highlights"`
}

type Project struct {
	Name        string `yaml:"name"`
	Tech        string `yaml:"tech"`
	Description string `yaml:"description"`
}

type SkillSet struct {
	Categories []SkillCategory `yaml:"categories"`
}

type SkillCategory struct {
	Name  string   `yaml:"name"`
	Items []string `yaml:"items"`
}

type Contact struct {
	GitHub   string `yaml:"github"`
	LinkedIn string `yaml:"linkedin"`
	Email    string `yaml:"email"`
	Blog     string `yaml:"blog"`
}
