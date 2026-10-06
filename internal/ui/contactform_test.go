package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestContactFocusCycleAndEscape(t *testing.T) {
	f := newContactForm("", "")
	f, _ = f.update(tea.KeyMsg{Type: tea.KeyEnter})
	for i := 0; i < fieldCount*10; i++ {
		want := i % fieldCount
		if f.focus != want {
			t.Fatalf("step %d: focus %d, want %d", i, f.focus, want)
		}
		f, _ = f.update(tea.KeyMsg{Type: tea.KeyTab})
	}
	f, _ = f.update(tea.KeyMsg{Type: tea.KeyShiftTab})
	if f.focus != fieldSubmit {
		t.Fatalf("reverse wrap: focus %d, want Submit", f.focus)
	}
	f, _ = f.update(tea.KeyMsg{Type: tea.KeyEsc})
	if f.editing {
		t.Fatal("Esc must leave editing mode")
	}
	// The full UI must route ordinary letters, including q, to the form.
	m := Model{keys: DefaultKeyMap(), activeTab: 4, scrollOffset: 50, contact: newContactForm("", "")}
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m = updated.(Model)
	if m.contact.name.Value() != "q" || m.scrollOffset != 0 {
		t.Fatal("editing must accept q and reset the stale scroll offset")
	}
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	if updated.(Model).activeTab != 5 {
		t.Fatal("Tab navigation must resume after Esc")
	}
	m.activeTab = 4
	for i := 0; i < 1000; i++ {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = updated.(Model)
	}
	if m.scrollOffset >= 1000 {
		t.Fatal("scrolling must stop at the end instead of accumulating invisible steps")
	}
	before := m.scrollOffset
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if updated.(Model).scrollOffset != before-1 {
		t.Fatal("scrolling back must respond immediately")
	}
}
