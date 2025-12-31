package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("62")).
			Bold(true)
	vaultDir string
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

func init() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("error getting home directory", err)
	}
	vaultDir = fmt.Sprintf("%s/.ShellNotes", homeDir)
}

type model struct {
	newFileInput           textinput.Model
	createFileInputVisible bool
	currentFile            *os.File
	noteTextArea           textarea.Model
	list                   list.Model
	showingList            bool
	isFileUnsaved          bool
}
type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-5)

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+q":
			fmt.Println("User Clicked", msg.String())
			return m, tea.Quit
		case "esc":
			if m.createFileInputVisible {
				m.createFileInputVisible = false
				m.newFileInput.SetValue("")
				return m, nil
			}

			if m.currentFile != nil {
				// Save the full filename BEFORE any operations
				filename := m.currentFile.Name()
				shouldRemove := m.isFileUnsaved

				// Close the file first
				_ = m.currentFile.Close()

				// Clean up state BEFORE removal
				m.noteTextArea.SetValue("")
				m.currentFile = nil
				m.isFileUnsaved = false

				// Remove if it was unsaved (using saved filename)
				if shouldRemove {
					_ = os.Remove(filename)
				}

				return m, nil
			}

			if m.showingList {
				if m.list.FilterState() == list.Filtering {
					break
				}
				m.showingList = false
				return m, nil
			}
		case "ctrl+n":
			m.createFileInputVisible = true
			return m, cmd
		case "ctrl+l":
			noteList := listFiles()
			m.list.SetItems(noteList)
			m.showingList = true
			//todo show list
			return m, nil
		case "enter":

			if m.showingList {
				item, ok := m.list.SelectedItem().(item)
				if ok {
					filepath := fmt.Sprintf("%s/%s", vaultDir, item.title)
					content, err := os.ReadFile(filepath)
					if err != nil {
						log.Printf("Error reading file: %v", err)
						return m, nil
					}
					m.noteTextArea.SetValue(string(content))
					f, err := os.OpenFile(filepath, os.O_RDWR, 0644)
					if err != nil {
						log.Printf("Error reading file: %v", err)
						return m, nil
					}
					m.currentFile = f
					m.showingList = false
				}
				return m, nil
			}
			if m.currentFile != nil {
				break
			}
			// todo: create file
			filename := m.newFileInput.Value()
			if filename != "" {
				filePath := fmt.Sprintf("%s/%s.md", vaultDir, filename)

				if _, err := os.Stat(filePath); err == nil {
					return m, nil
				}

				f, err := os.Create(filePath)
				if err != nil {
					log.Fatalf("%v", err)
				}
				m.currentFile = f
				m.isFileUnsaved = true
				m.createFileInputVisible = false
				m.newFileInput.SetValue("")

			}
			return m, nil
		case "ctrl+s":
			// textarea value write it in file descriptor and close it
			if m.currentFile == nil {
				break
			}
			if err := m.currentFile.Truncate(0); err != nil {
				fmt.Println("cannot save the file")
				return m, nil
			}
			if _, err := m.currentFile.Seek(0, 0); err != nil {
				fmt.Println("cannot save the file")
			}

			if _, err := m.currentFile.WriteString(m.noteTextArea.Value()); err != nil {
				fmt.Println()
			}
			m.isFileUnsaved = false
			if err := m.currentFile.Close(); err != nil {
				fmt.Println("cannot close the file")
			}
			m.currentFile = nil
			m.noteTextArea.SetValue("")
			return m, nil

		case "ctrl+d":
			// Only allow deletion from list view
			if m.showingList {
				selectedItem, ok := m.list.SelectedItem().(item)
				if ok {
					filepath := fmt.Sprintf("%s/%s", vaultDir, selectedItem.title)

					// Delete the file
					if err := os.Remove(filepath); err != nil {
						log.Printf("Error deleting file: %v", err)
						return m, nil
					}

					// Refresh the list
					noteList := listFiles()
					m.list.SetItems(noteList)
				}
			}
			return m, nil
		}

	}

	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}

	if m.currentFile != nil {

		m.noteTextArea, cmd = m.noteTextArea.Update(msg)

	}
	if m.showingList {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {

	var style_heading = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("205")).
		PaddingLeft(2).
		PaddingRight(2)

	var style_help = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A6ADC8")). // soft gray text
		Background(lipgloss.Color("#1E1E2E")). // dark background
		Padding(0, 1).
		Bold(true)
	welcome := style_heading.Render("Welcome to ShellNotes📝")
	help := style_help.Render("Ctrl+N: new file | Ctrl+L: list | Esc: back | Ctrl+S: save | Ctrl+D: delete note in list | Ctrl+Q: quit")
	view := ""
	if m.createFileInputVisible {
		view = m.newFileInput.View()
	}
	if m.currentFile != nil {
		view = m.noteTextArea.View()
	}
	if m.showingList {
		view = m.list.View()
	}

	return fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, view, help)
}
func initializeModel() model {

	err := os.MkdirAll(vaultDir, 0750)
	if err != nil {
		log.Fatal(err)
	}
	//initialize new file input
	ti := textinput.New()
	ti.Placeholder = "What would you like to call it?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 220
	ti.Cursor.Style = cursorStyle
	ti.Prompt = "✏️ "

	//textarea

	ta := textarea.New()
	ta.Placeholder = "Write your note here..."
	ta.ShowLineNumbers = false
	ta.Focus()

	//list

	noteList := listFiles()
	finalList := list.New(noteList, list.NewDefaultDelegate(), 0, 0)
	finalList.Title = "All notes 📙"
	finalList.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("16")).Background(lipgloss.Color("254")).Padding(0, 1)
	return model{
		newFileInput:           ti,
		createFileInputVisible: false,
		noteTextArea:           ta,
		list:                   finalList,
	}
}
func main() {

	p := tea.NewProgram((initializeModel()))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error : %v", err)
		os.Exit((1))
	}

}

func listFiles() []list.Item {

	items := make([]list.Item, 0)

	entries, err := os.ReadDir(vaultDir)
	if err != nil {
		log.Fatal("Error reading notes list")
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			modTime := info.ModTime().Format("2006-01-02 15:04")

			items = append(items, item{title: entry.Name(),
				desc: fmt.Sprintf("Modified: %s", modTime)})
		}
	}

	return items
}
