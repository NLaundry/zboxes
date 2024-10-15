package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

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

// Model to hold application state

type model struct {
	cursor       [4]int // Cursor positions for each level
	selected     [4]int // Selected indices at each level
	currentLevel int    // Current navigation level (0: ZBox, 1: ZPool, 2: Dataset, 3: Snapshot)

	width  int // Terminal width
	height int // Terminal height

	zbox ZBox // The local zbox
}

func initialModel() model {
	// Initialize the model
	zbox, err := getLocalZBox()
	if err != nil {
		fmt.Printf("Error initializing ZBox: %v\n", err)
		os.Exit(1)
	}

	return model{
		cursor:       [4]int{0, 0, 0, 0},
		selected:     [4]int{0, 0, 0, 0},
		currentLevel: 0, // Start at ZBox level
		zbox:         zbox,
	}
}

func getLocalZBox() (ZBox, error) {
	var zbox ZBox

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return zbox, err
	}
	zbox.Hostname = hostname

	// For Name and User, we can set default values
	zbox.Name = "Local ZBox"
	zbox.User = os.Getenv("USER")

	// Get ZPools
	zpools, err := getZPools()
	if err != nil {
		return zbox, err
	}
	zbox.ZPools = zpools

	return zbox, nil
}

func getZPools() ([]ZPool, error) {
	var zpools []ZPool

	cmd := exec.Command("zpool", "list", "-Hp", "-o", "name,health")
	output, err := cmd.Output()
	if err != nil {
		return zpools, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		fields := strings.Split(line, "\t")
		if len(fields) >= 2 {
			zpoolName := fields[0]
			zpoolHealth := fields[1]
			zpool := ZPool{
				Name:   zpoolName,
				Health: zpoolHealth,
			}

			// Get datasets for this zpool
			datasets, err := getDatasets(zpool.Name)
			if err != nil {
				return zpools, err
			}
			zpool.Datasets = datasets
			zpool.NumDatasets = len(datasets)

			// Count snapshots
			numSnapshots := 0
			for _, ds := range datasets {
				numSnapshots += len(ds.Snapshots)
			}
			zpool.NumSnapshots = numSnapshots

			zpools = append(zpools, zpool)
		}
	}

	return zpools, nil
}

func getDatasets(zpoolName string) ([]Dataset, error) {
	var datasets []Dataset

	cmd := exec.Command("zfs", "list", "-r", "-Hp", "-o", "name,used,avail,mountpoint", zpoolName)
	output, err := cmd.Output()
	if err != nil {
		return datasets, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		fields := strings.Split(line, "\t")
		if len(fields) >= 4 {
			datasetName := fields[0]
			used := fields[1]
			avail := fields[2]
			mountpoint := fields[3]

			dataset := Dataset{
				Name:       datasetName,
				Used:       used,
				Available:  avail,
				Mountpoint: mountpoint,
			}

			// Get snapshots for this dataset
			snapshots, err := getSnapshots(dataset.Name)
			if err != nil {
				return datasets, err
			}
			dataset.Snapshots = snapshots

			datasets = append(datasets, dataset)
		}
	}

	return datasets, nil
}

func getSnapshots(datasetName string) ([]Snapshot, error) {
	var snapshots []Snapshot

	cmd := exec.Command("zfs", "list", "-t", "snapshot", "-r", "-Hp", "-o", "name,used,creation", datasetName)
	output, err := cmd.Output()
	if err != nil {
		return snapshots, nil // It's okay if there are no snapshots
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		fields := strings.Split(line, "\t")
		if len(fields) >= 3 {
			snapshotName := fields[0]
			size := fields[1]
			creation := fields[2]
			// Convert creation time to desired format
			timestamp, err := parseUnixTimestamp(creation)
			var formattedDate string
			if err == nil {
				formattedDate = timestamp.Format("02-01-2006")
			} else {
				formattedDate = creation
			}

			snapshot := Snapshot{
				Name: snapshotName,
				Size: size,
				Date: formattedDate,
			}

			snapshots = append(snapshots, snapshot)
		}
	}

	return snapshots, nil
}

func parseUnixTimestamp(ts string) (time.Time, error) {
	unixTime, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(unixTime, 0), nil
}

// Catppuccin color palette (Mocha)
var (
	rosewater = lipgloss.Color("#f5e0dc")
	flamingo  = lipgloss.Color("#f2cdcd")
	pink      = lipgloss.Color("#f5c2e7")
	mauve     = lipgloss.Color("#cba6f7")
	red       = lipgloss.Color("#f38ba8")
	maroon    = lipgloss.Color("#eba0ac")
	peach     = lipgloss.Color("#fab387")
	yellow    = lipgloss.Color("#f9e2af")
	green     = lipgloss.Color("#a6e3a1")
	teal      = lipgloss.Color("#94e2d5")
	sky       = lipgloss.Color("#89dceb")
	sapphire  = lipgloss.Color("#74c7ec")
	blue      = lipgloss.Color("#89b4fa")
	lavender  = lipgloss.Color("#b4befe")
	text      = lipgloss.Color("#cdd6f4")
	subtext1  = lipgloss.Color("#bac2de")
	subtext0  = lipgloss.Color("#a6adc8")
	overlay2  = lipgloss.Color("#9399b2")
	overlay1  = lipgloss.Color("#7f849c")
	overlay0  = lipgloss.Color("#6c7086")
	surface2  = lipgloss.Color("#585b70")
	surface1  = lipgloss.Color("#45475a")
	surface0  = lipgloss.Color("#313244")
	base      = lipgloss.Color("#1e1e2e")
	mantle    = lipgloss.Color("#181825")
	crust     = lipgloss.Color("#11111b")
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(peach).
		Background(crust).
		Padding(0, 1).
		BorderStyle(lipgloss.Border{Bottom: "─"}).
		BorderBottom(true).
		BorderForeground(surface2)

	normalTextStyle = lipgloss.NewStyle().
		Foreground(text).
		Background(base)

	cursorStyle = lipgloss.NewStyle().
		Foreground(crust).
		Background(sky).
		Bold(true)

	selectedStyle = lipgloss.NewStyle().
		Foreground(crust).
		Background(mauve).
		Bold(true)

	borderStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(surface2).
		Background(base)

	instructionStyle = lipgloss.NewStyle().
		Foreground(subtext0).
		Background(base)

	activeColumnStyle = lipgloss.NewStyle().
		Background(surface0)
)

// Init function

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
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

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View function to render the UI in a ranger-style interface

func (m model) View() string {
	var b strings.Builder

	// Level titles
	levelTitles := []string{"ZBox", "ZPools", "Datasets", "Snapshots"}

	// Calculate column widths
	totalWidth := m.width
	if totalWidth <= 0 {
		totalWidth = 80 // default width
	}

	// Half the screen for columns, half for preview
	columnAreaWidth := totalWidth / 2
	previewWidth := totalWidth - columnAreaWidth

	// Number of columns is always 3 (parent, current, next)
	numColumns := 3

	// Divide the column area width among the columns
	columnWidth := columnAreaWidth / numColumns

	// Ensure columns have a minimum width
	if columnWidth < 20 {
		columnWidth = 20
	}

	// Collect columns
	columns := []string{}

	// Left Column: Parent Items
	{
		var leftColumn strings.Builder
		parentTitle := levelTitles[m.currentLevel]
		leftColumn.WriteString(titleStyle.Render(parentTitle) + "\n")
		parentItems := m.getParentItems()
		if len(parentItems) == 0 {
			leftColumn.WriteString(normalTextStyle.Render("No items\n"))
		} else {
			for i, item := range parentItems {
				lineStyle := normalTextStyle
				if m.selected[m.currentLevel-1] == i {
					lineStyle = selectedStyle
					item = lineStyle.Render(item)
				} else {
					item = lineStyle.Render(item)
				}
				prefix := "  "
				if m.selected[m.currentLevel-1] == i {
					prefix = "➤ "
				}
				leftColumn.WriteString(fmt.Sprintf("%s%s\n", prefix, item))
			}
		}

		columnStyle := borderStyle.Width(columnWidth)
		if m.currentLevel-1 == m.currentLevel-1 {
			columnStyle = columnStyle.Inherit(activeColumnStyle)
		}
		leftColumnStr := columnStyle.Render(leftColumn.String())
		columns = append(columns, leftColumnStr)
	}

	// Middle Column: Current Items
	{
		var middleColumn strings.Builder
		middleColumn.WriteString(titleStyle.Render(levelTitles[m.currentLevel]) + "\n")
		currentItems := m.getCurrentItems()
		if len(currentItems) == 0 {
			middleColumn.WriteString(normalTextStyle.Render("No items\n"))
		} else {
			for i, item := range currentItems {
				lineStyle := normalTextStyle
				prefix := "  "
				if m.cursor[m.currentLevel] == i {
					lineStyle = cursorStyle
					item = lineStyle.Render(item)
					prefix = "➤ "
				} else {
					item = lineStyle.Render(item)
				}
				middleColumn.WriteString(fmt.Sprintf("%s%s\n", prefix, item))
			}
		}

		columnStyle := borderStyle.Width(columnWidth)
		if m.currentLevel == m.currentLevel {
			columnStyle = columnStyle.Inherit(activeColumnStyle)
		}
		middleColumnStr := columnStyle.Render(middleColumn.String())
		columns = append(columns, middleColumnStr)
	}

	// Next Level Column
	{
		var nextColumn strings.Builder
		nextLevel := m.currentLevel + 1
		if nextLevel > 3 {
			nextLevel = 3
		}
		nextTitle := levelTitles[nextLevel]
		nextColumn.WriteString(titleStyle.Render(nextTitle) + "\n")
		nextItems := m.getNextItems()
		if len(nextItems) == 0 {
			nextColumn.WriteString(normalTextStyle.Render("No items\n"))
		} else {
			for _, item := range nextItems {
				item = normalTextStyle.Render(item)
				nextColumn.WriteString(fmt.Sprintf("  %s\n", item))
			}
		}

		columnStyle := borderStyle.Width(columnWidth)
		if nextLevel == m.currentLevel {
			columnStyle = columnStyle.Inherit(activeColumnStyle)
		}
		nextColumnStr := columnStyle.Render(nextColumn.String())
		columns = append(columns, nextColumnStr)
	}

	// Right Column: Preview
	var previewPane strings.Builder
	previewPane.WriteString(titleStyle.Width(previewWidth).Render("Details") + "\n")
	previewText := m.getPreview()
	previewPane.WriteString(normalTextStyle.Render(previewText))

	// Combine columns
	columns = append(columns, borderStyle.Width(previewWidth).Render(previewPane.String()))

	// Join all columns
	view := lipgloss.JoinHorizontal(lipgloss.Top, columns...)

	// Instructions at the bottom
	instructions := "\n" + instructionStyle.Render(
		"Navigate with h (left), j (down), k (up), l (right). Press 'q' to quit.",
	)

	// Final output
	b.WriteString(view)
	b.WriteString(instructions)

	return b.String()
}

// Functions to get items for the current, parent, and next levels

func (m model) getParentItems() []string {
	var items []string
	if m.currentLevel == 0 {
		// No parent items at the top level
		return items
	}
	switch m.currentLevel {
	case 1:
		// Parent level is ZBox
		items = append(items, m.zbox.Name)
	case 2:
		// Parent level is ZPools
		for _, zpool := range m.zbox.ZPools {
			items = append(items, zpool.Name)
		}
	case 3:
		// Parent level is Datasets
		zpool := m.zbox.ZPools[m.selected[1]]
		for _, dataset := range zpool.Datasets {
			items = append(items, dataset.Name)
		}
	}
	return items
}

func (m model) getCurrentItems() []string {
	var items []string
	switch m.currentLevel {
	case 0:
		// Current level is ZBox
		items = append(items, m.zbox.Name)
	case 1:
		// Current level is ZPools
		for _, zpool := range m.zbox.ZPools {
			items = append(items, zpool.Name)
		}
	case 2:
		// Current level is Datasets
		if len(m.zbox.ZPools) > 0 {
			zpool := m.zbox.ZPools[m.selected[1]]
			for _, dataset := range zpool.Datasets {
				items = append(items, dataset.Name)
			}
		}
	case 3:
		// Current level is Snapshots
		if len(m.zbox.ZPools) > 0 {
			zpool := m.zbox.ZPools[m.selected[1]]
			if len(zpool.Datasets) > 0 {
				dataset := zpool.Datasets[m.selected[2]]
				for _, snapshot := range dataset.Snapshots {
					items = append(items, snapshot.Name)
				}
			}
		}
	}
	return items
}

func (m model) getNextItems() []string {
	var items []string
	nextLevel := m.currentLevel + 1
	if nextLevel > 3 {
		return items
	}
	switch nextLevel {
	case 1:
		// Next level is ZPools
		for _, zpool := range m.zbox.ZPools {
			items = append(items, zpool.Name)
		}
	case 2:
		// Next level is Datasets
		if len(m.zbox.ZPools) > 0 {
			zpool := m.zbox.ZPools[m.cursor[1]]
			for _, dataset := range zpool.Datasets {
				items = append(items, dataset.Name)
			}
		}
	case 3:
		// Next level is Snapshots
		if len(m.zbox.ZPools) > 0 {
			zpool := m.zbox.ZPools[m.selected[1]]
			if len(zpool.Datasets) > 0 {
				dataset := zpool.Datasets[m.cursor[2]]
				for _, snapshot := range dataset.Snapshots {
					items = append(items, snapshot.Name)
				}
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
		zbox := m.zbox
		previewText += fmt.Sprintf("Name: %s\n", zbox.Name)
		previewText += fmt.Sprintf("Hostname: %s\n", zbox.Hostname)
		previewText += fmt.Sprintf("User: %s\n", zbox.User)
	case 1:
		// ZPool details
		if len(m.zbox.ZPools) > 0 {
			zpool := m.zbox.ZPools[m.cursor[1]]
			previewText += fmt.Sprintf("Name: %s\n", zpool.Name)
			previewText += fmt.Sprintf("Health: %s\n", zpool.Health)
			previewText += fmt.Sprintf("Datasets: %d\n", zpool.NumDatasets)
			previewText += fmt.Sprintf("Snapshots: %d\n", zpool.NumSnapshots)
		}
	case 2:
		// Dataset details
		if len(m.zbox.ZPools) > 0 {
			zpool := m.zbox.ZPools[m.selected[1]]
			if len(zpool.Datasets) > 0 {
				dataset := zpool.Datasets[m.cursor[2]]
				previewText += fmt.Sprintf("Name: %s\n", dataset.Name)
				previewText += fmt.Sprintf("Used: %s\n", dataset.Used)
				previewText += fmt.Sprintf("Available: %s\n", dataset.Available)
				previewText += fmt.Sprintf("Mountpoint: %s\n", dataset.Mountpoint)
			}
		}
	case 3:
		// Snapshot details
		if len(m.zbox.ZPools) > 0 {
			zpool := m.zbox.ZPools[m.selected[1]]
			if len(zpool.Datasets) > 0 {
				dataset := zpool.Datasets[m.selected[2]]
				if len(dataset.Snapshots) > 0 {
					snapshot := dataset.Snapshots[m.cursor[3]]
					previewText += fmt.Sprintf("Name: %s\n", snapshot.Name)
					previewText += fmt.Sprintf("Size: %s\n", snapshot.Size)
					previewText += fmt.Sprintf("Date: %s\n", snapshot.Date)
				}
			}
		}
	}
	return previewText
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
