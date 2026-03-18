package content

import (
	"testing"
)

func TestLoadPortfolio(t *testing.T) {
	p, err := LoadPortfolio()
	if err != nil {
		t.Fatalf("LoadPortfolio: %v", err)
	}

	if p.Name == "" {
		t.Fatal("expected name to be set")
	}
	if p.Title == "" {
		t.Fatal("expected title to be set")
	}
	if p.Location == "" {
		t.Fatal("expected location to be set")
	}
	if p.About == "" {
		t.Fatal("expected about to be set")
	}
}

func TestLoadPortfolio_Experience(t *testing.T) {
	p, err := LoadPortfolio()
	if err != nil {
		t.Fatalf("LoadPortfolio: %v", err)
	}

	if len(p.Experience) == 0 {
		t.Fatal("expected at least one experience entry")
	}

	for i, exp := range p.Experience {
		if exp.Role == "" {
			t.Fatalf("experience[%d]: expected role", i)
		}
		if exp.Company == "" {
			t.Fatalf("experience[%d]: expected company", i)
		}
		if exp.Period == "" {
			t.Fatalf("experience[%d]: expected period", i)
		}
		if len(exp.Highlights) == 0 {
			t.Fatalf("experience[%d]: expected highlights", i)
		}
	}
}

func TestLoadPortfolio_Projects(t *testing.T) {
	p, err := LoadPortfolio()
	if err != nil {
		t.Fatalf("LoadPortfolio: %v", err)
	}

	if len(p.Projects) == 0 {
		t.Fatal("expected at least one project")
	}

	for i, proj := range p.Projects {
		if proj.Name == "" {
			t.Fatalf("project[%d]: expected name", i)
		}
		if proj.Tech == "" {
			t.Fatalf("project[%d]: expected tech", i)
		}
		if proj.Description == "" {
			t.Fatalf("project[%d]: expected description", i)
		}
	}
}

func TestLoadPortfolio_Skills(t *testing.T) {
	p, err := LoadPortfolio()
	if err != nil {
		t.Fatalf("LoadPortfolio: %v", err)
	}

	if len(p.Skills.Categories) == 0 {
		t.Fatal("expected at least one skill category")
	}

	for i, cat := range p.Skills.Categories {
		if cat.Name == "" {
			t.Fatalf("skill category[%d]: expected name", i)
		}
		if len(cat.Items) == 0 {
			t.Fatalf("skill category[%d]: expected items", i)
		}
	}
}

func TestLoadPortfolio_Contact(t *testing.T) {
	p, err := LoadPortfolio()
	if err != nil {
		t.Fatalf("LoadPortfolio: %v", err)
	}

	if p.Contact.GitHub == "" {
		t.Fatal("expected GitHub contact")
	}
	if p.Contact.LinkedIn == "" {
		t.Fatal("expected LinkedIn contact")
	}
	if p.Contact.Email == "" {
		t.Fatal("expected email contact")
	}
}
