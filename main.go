package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Data structures for ZBoxes, ZPools, Datasets, and Snapshots

type Snapshot struct {
	Name   string
	Size   string
	Date   string
	Status string
}

type Dataset struct {
	Name       string
	Used       string
	Available  string
	Mountpoint string
	Snapshots  []Snapshot
}

type ZPool struct {
	Name         string
	Health       string
	LastScrub    string
	NumDatasets  int
	NumSnapshots int
	Datasets     []Dataset
}

type ZBox struct {
	Name     string
	IP       string
	Hostname string
	User     string
	ZPools   []ZPool
}

// Mock data representing ZBoxes, ZPools, Datasets, and Snapshots

var zboxes = []ZBox{
	{
		Name:     "ZBox1",
		IP:       "192.168.1.10",
		Hostname: "zbox1.local",
		User:     "admin",
		ZPools: []ZPool{
			{
				Name:         "ZPool1",
				Health:       "ONLINE",
				LastScrub:    "2023-10-01",
				NumDatasets:  2,
				NumSnapshots: 5,
				Datasets: []Dataset{
					{
						Name:       "Dataset1",
						Used:       "10G",
						Available:  "90G",
						Mountpoint: "/mnt/dataset1",
						Snapshots: []Snapshot{
							{Name: "snap1", Size: "1G", Date: "2023-09-30", Status: "OK"},
							{Name: "snap2", Size: "1G", Date: "2023-09-25", Status: "OK"},
						},
					},
					{
						Name:       "Dataset2",
						Used:       "20G",
						Available:  "80G",
						Mountpoint: "/mnt/dataset2",
						Snapshots: []Snapshot{
							{Name: "snap3", Size: "2G", Date: "2023-09-28", Status: "OK"},
						},
					},
				},
			},
			{
				Name:         "ZPool2",
				Health:       "DEGRADED",
				LastScrub:    "2023-09-25",
				NumDatasets:  1,
				NumSnapshots: 2,
				Datasets: []Dataset{
					{
						Name:       "Dataset3",
						Used:       "5G",
						Available:  "95G",
						Mountpoint: "/mnt/dataset3",
						Snapshots: []Snapshot{
							{Name: "snap4", Size: "500M", Date: "2023-09-20", Status: "OK"},
						},
					},
				},
			},
		},
	},
	{
		Name:     "ZBox2",
		IP:       "192.168.1.11",
		Hostname: "zbox2.local",
		User:     "admin",
		ZPools: []ZPool{
			{
				Name:         "ZPoolA",
				Health:       "ONLINE",
				LastScrub:    "2023-09-28",
				NumDatasets:  2,
				NumSnapshots: 3,
				Datasets: []Dataset{
					{
						Name:       "DatasetX",
						Used:       "15G",
						Available:  "85G",
						Mountpoint: "/mnt/datasetX",
						Snapshots: []Snapshot{
							{Name: "snap5", Size: "1.5G", Date: "2023-09-27", Status: "OK"},
						},
					},
					{
						Name:       "DatasetY",
						Used:       "25G",
						Available:  "75G",
						Mountpoint: "/mnt/datasetY",
						Snapshots: []Snapshot{
							{Name: "snap6", Size: "2.5G", Date: "2023-09-25", Status: "OK"},
							{Name: "snap7", Size: "2.5G", Date: "2023-09-20", Status: "OK"},
						},
					},
				},
			},
			{
				Name:         "ZPoolB",
				Health:       "FAULTED",
				LastScrub:    "2023-09-20",
				NumDatasets:  1,
				NumSnapshots: 1,
				Datasets: []Dataset{
					{
						Name:       "DatasetZ",
						Used:       "30G",
						Available:  "70G",
						Mountpoint: "/mnt/datasetZ",
						Snapshots: []Snapshot{
							{Name: "snap8", Size: "3G", Date: "2023-09-18", Status: "OK"},
						},
					},
				},
			},
		},
	},
}

// Model to hold application state

type model struct {
	cursor       [4]int // Cursor positions for each level
	selected     [4]int // Selected indices at each level
	currentLevel int    // Current navigation level (0: ZBox, 1: ZPool, 2: Dataset, 3: Snapshot)
}

func initialModel() model {
	return model{
		cursor:       [4]int{0, 0, 0, 0},
		selected:     [4]int{0, 0, 0, 0},
		currentLevel: 0, // Start at ZBox level
	}
}

// Init function

func (m model) Init() tea.Cmd {
	// No initialization needed
	return nil
}

// Update function

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		// Navigate up
		case "k":
			if m.cursor[m.currentLevel] > 0 {
				m.cursor[m.currentLevel]--
			}

		// Navigate down
		case "j":
			max := len(m.getCurrentItems()) - 1
			if m.cursor[m.currentLevel] < max {
				m.cursor[m.currentLevel]++
			}

		// Move right (increase level)
		case "l":
			if m.currentLevel < 3 && len(m.getCurrentItems()) > 0 {
				// Update selected index at current level
				m.selected[m.currentLevel] = m.cursor[m.currentLevel]
				m.currentLevel++
				// Reset cursor for new level
				m.cursor[m.currentLevel] = 0
			}

		// Move left (decrease level)
		case "h":
			if m.currentLevel > 0 {
				m.currentLevel--
				// Optionally reset cursor to selected index
				m.cursor[m.currentLevel] = m.selected[m.currentLevel]
			}

		// Quit the application
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

// View function to render the UI in a ranger-style interface

func (m model) View() string {
	var b strings.Builder

	// Column widths
	columnWidth := 30
	previewWidth := 50

	// Styles
	titleStyle := lipgloss.NewStyle().Bold(true).Underline(true)
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Bold(true)
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(true)
	invertedStyle := func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("229")).Render(s)
	}

	// Level titles
	levelTitles := []string{"ZBoxes", "ZPools", "Datasets", "Snapshots"}

	// Left Column: Parent Items
	var leftColumn strings.Builder
	parentTitle := "Parent Level"
	if m.currentLevel > 0 {
		parentTitle = levelTitles[m.currentLevel-1]
	}
	leftColumn.WriteString(titleStyle.Render(parentTitle) + "\n")
	parentItems := m.getParentItems()
	for i, item := range parentItems {
		lineStyle := lipgloss.NewStyle()
		prefix := "  "
		if m.selected[m.currentLevel-1] == i {
			lineStyle = selectedStyle
			prefix = "➤ "
			item = invertedStyle(item)
		}
		leftColumn.WriteString(fmt.Sprintf("%s%s\n", prefix, lineStyle.Render(item)))
	}

	// Middle Column: Current Items
	var middleColumn strings.Builder
	middleColumn.WriteString(titleStyle.Render(levelTitles[m.currentLevel]) + "\n")
	currentItems := m.getCurrentItems()
	for i, item := range currentItems {
		cursor := "  "
		lineStyle := lipgloss.NewStyle()
		if m.cursor[m.currentLevel] == i {
			cursor = "➤ "
			lineStyle = cursorStyle
			item = invertedStyle(item)
		}
		middleColumn.WriteString(fmt.Sprintf("%s%s\n", cursor, lineStyle.Render(item)))
	}

	// Right Column: Preview
	var previewPane strings.Builder
	previewPane.WriteString(titleStyle.Render("Details") + "\n")
	previewText := m.getPreview()
	previewPane.WriteString(previewText)

	// Combine columns
	columns := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(columnWidth).Render(leftColumn.String()),
		lipgloss.NewStyle().Width(columnWidth).Render(middleColumn.String()),
		lipgloss.NewStyle().Width(previewWidth).Border(lipgloss.NormalBorder()).Render(previewPane.String()),
	)

	// Instructions at the bottom
	instructions := "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(
		"Navigate with h (left), j (down), k (up), l (right). Press 'q' to quit.",
	)

	// Final output
	b.WriteString(columns)
	b.WriteString(instructions)

	return b.String()
}

// Functions to get items for the current and parent levels

func (m model) getParentItems() []string {
	var items []string
	if m.currentLevel == 0 {
		// No parent items at the top level
		return items
	}
	switch m.currentLevel {
	case 1:
		{
			// Parent level is ZBoxes
			for _, zbox := range zboxes {
				items = append(items, zbox.Name)
			}
		}
	case 2:
		{
			// Parent level is ZPools
			zbox := zboxes[m.selected[0]]
			for _, zpool := range zbox.ZPools {
				items = append(items, zpool.Name)
			}
		}
	case 3:
		{
			// Parent level is Datasets
			zbox := zboxes[m.selected[0]]
			zpool := zbox.ZPools[m.selected[1]]
			for _, dataset := range zpool.Datasets {
				items = append(items, dataset.Name)
			}
		}
	}
	return items
}

func (m model) getCurrentItems() []string {
	var items []string
	switch m.currentLevel {
	case 0:
		{
			// Current level is ZBoxes
			for _, zbox := range zboxes {
				items = append(items, zbox.Name)
			}
		}
	case 1:
		{
			// Current level is ZPools
			zbox := zboxes[m.selected[0]]
			for _, zpool := range zbox.ZPools {
				items = append(items, zpool.Name)
			}
		}
	case 2:
		{
			// Current level is Datasets
			zbox := zboxes[m.selected[0]]
			zpool := zbox.ZPools[m.selected[1]]
			for _, dataset := range zpool.Datasets {
				items = append(items, dataset.Name)
			}
		}
	case 3:
		{
			// Current level is Snapshots
			zbox := zboxes[m.selected[0]]
			zpool := zbox.ZPools[m.selected[1]]
			dataset := zpool.Datasets[m.selected[2]]
			for _, snapshot := range dataset.Snapshots {
				items = append(items, snapshot.Name)
			}
		}
	}
	return items
}

// Function to get preview details of the selected item

func (m model) getPreview() string {
	var previewText string
	switch m.currentLevel {
	case 0:
		// ZBox details
		if len(zboxes) > 0 {
			zbox := zboxes[m.cursor[0]]
			previewText += fmt.Sprintf("Name: %s\n", zbox.Name)
			previewText += fmt.Sprintf("IP: %s\n", zbox.IP)
			previewText += fmt.Sprintf("Hostname: %s\n", zbox.Hostname)
			previewText += fmt.Sprintf("User: %s\n", zbox.User)
		}
	case 1:
		// ZPool details
		zbox := zboxes[m.selected[0]]
		if len(zbox.ZPools) > 0 {
			zpool := zbox.ZPools[m.cursor[1]]
			previewText += fmt.Sprintf("Name: %s\n", zpool.Name)
			previewText += fmt.Sprintf("Health: %s\n", zpool.Health)
			previewText += fmt.Sprintf("Last Scrub: %s\n", zpool.LastScrub)
			previewText += fmt.Sprintf("Datasets: %d\n", zpool.NumDatasets)
			previewText += fmt.Sprintf("Snapshots: %d\n", zpool.NumSnapshots)
		}
	case 2:
		// Dataset details
		zbox := zboxes[m.selected[0]]
		zpool := zbox.ZPools[m.selected[1]]
		if len(zpool.Datasets) > 0 {
			dataset := zpool.Datasets[m.cursor[2]]
			previewText += fmt.Sprintf("Name: %s\n", dataset.Name)
			previewText += fmt.Sprintf("Used: %s\n", dataset.Used)
			previewText += fmt.Sprintf("Available: %s\n", dataset.Available)
			previewText += fmt.Sprintf("Mountpoint: %s\n", dataset.Mountpoint)
		}
	case 3:
		// Snapshot details
		zbox := zboxes[m.selected[0]]
		zpool := zbox.ZPools[m.selected[1]]
		dataset := zpool.Datasets[m.selected[2]]
		if len(dataset.Snapshots) > 0 {
			snapshot := dataset.Snapshots[m.cursor[3]]
			previewText += fmt.Sprintf("Name: %s\n", snapshot.Name)
			previewText += fmt.Sprintf("Size: %s\n", snapshot.Size)
			previewText += fmt.Sprintf("Date: %s\n", snapshot.Date)
			previewText += fmt.Sprintf("Status: %s\n", snapshot.Status)
		}
	}
	return previewText
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
