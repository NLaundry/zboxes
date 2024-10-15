// model.go
package main

import (
	"fmt"
	"strings"
    "os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	cursor       [MaxLevel]int // Cursor positions for each level
	selected     [MaxLevel]int // Selected indices at each level
	currentLevel int           // Current navigation level

	width  int // Terminal width
	height int // Terminal height

	zbox ZBox // The local ZBox
}

func initialModel() model {
	zbox, err := NewLocalZBox()
	if err != nil {
		fmt.Printf("Error initializing ZBox: %v\n", err)
		os.Exit(1)
	}

	return model{
		currentLevel: LevelZBox,
		zbox:         zbox,
	}
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

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
			max := len(m.getItems(m.currentLevel)) - 1
			if m.cursor[m.currentLevel] < max {
				m.cursor[m.currentLevel]++
			}

		// Move right (increase level)
		case "l":
			if m.currentLevel < LevelSnapshot && len(m.getItems(m.currentLevel)) > 0 {
				// Update selected index at current level
				m.selected[m.currentLevel] = m.cursor[m.currentLevel]
				m.currentLevel++
				// Reset cursor for new level
				m.cursor[m.currentLevel] = 0
			}

		// Move left (decrease level)
		case "h":
			if m.currentLevel > LevelZBox {
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
		parentLevel := m.currentLevel - 1
		if parentLevel < LevelZBox {
			parentLevel = LevelZBox
		}
		parentTitle := levelTitles[parentLevel]
		leftColumn.WriteString(titleStyle.Render(parentTitle) + "\n")
		parentItems := m.getItems(parentLevel)
		if len(parentItems) == 0 {
			leftColumn.WriteString(normalTextStyle.Render("No items\n"))
		} else {
			for i, item := range parentItems {
				lineStyle := normalTextStyle
				if m.selected[parentLevel] == i {
					lineStyle = selectedStyle
					item = lineStyle.Render(item)
				} else {
					item = lineStyle.Render(item)
				}
				prefix := "  "
				if m.selected[parentLevel] == i {
					prefix = "➤ "
				}
				leftColumn.WriteString(fmt.Sprintf("%s%s\n", prefix, item))
			}
		}

		columnStyle := borderStyle.Width(columnWidth)
		if parentLevel == m.currentLevel {
			columnStyle = columnStyle.Inherit(activeColumnStyle)
		}
		leftColumnStr := columnStyle.Render(leftColumn.String())
		columns = append(columns, leftColumnStr)
	}

	// Middle Column: Current Items
	{
		var middleColumn strings.Builder
		currentTitle := levelTitles[m.currentLevel]
		middleColumn.WriteString(titleStyle.Render(currentTitle) + "\n")
		currentItems := m.getItems(m.currentLevel)
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
		if nextLevel > LevelSnapshot {
			nextLevel = LevelSnapshot
		}
		nextTitle := levelTitles[nextLevel]
		nextColumn.WriteString(titleStyle.Render(nextTitle) + "\n")
		nextItems := m.getItems(nextLevel)
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

func (m model) getItems(level int) []string {
	switch level {
	case LevelZBox:
		return []string{m.zbox.Name}
	case LevelZPool:
		items := make([]string, len(m.zbox.ZPools))
		for i, zpool := range m.zbox.ZPools {
			items[i] = zpool.Name
		}
		return items
	case LevelDataset:
		if len(m.zbox.ZPools) == 0 {
			return nil
		}
		selectedZPool := m.zbox.ZPools[m.selected[LevelZPool]]
		items := make([]string, len(selectedZPool.Datasets))
		for i, dataset := range selectedZPool.Datasets {
			items[i] = dataset.Name
		}
		return items
	case LevelSnapshot:
		if len(m.zbox.ZPools) == 0 {
			return nil
		}
		selectedZPool := m.zbox.ZPools[m.selected[LevelZPool]]
		if len(selectedZPool.Datasets) == 0 {
			return nil
		}
		selectedDataset := selectedZPool.Datasets[m.selected[LevelDataset]]
		items := make([]string, len(selectedDataset.Snapshots))
		for i, snapshot := range selectedDataset.Snapshots {
			items[i] = snapshot.Name
		}
		return items
	default:
		return nil
	}
}

func (m model) getPreview() string {
	var previewText string
	switch m.currentLevel {
	case LevelZBox:
		// ZBox details
		zbox := m.zbox
		previewText += fmt.Sprintf("Name: %s\n", zbox.Name)
		previewText += fmt.Sprintf("Hostname: %s\n", zbox.Hostname)
		previewText += fmt.Sprintf("User: %s\n", zbox.User)
	case LevelZPool:
		// ZPool details
		if len(m.zbox.ZPools) > 0 {
			zpool := m.zbox.ZPools[m.cursor[LevelZPool]]
			previewText += fmt.Sprintf("Name: %s\n", zpool.Name)
			previewText += fmt.Sprintf("Health: %s\n", zpool.Health)
			previewText += fmt.Sprintf("Datasets: %d\n", zpool.NumDatasets)
			previewText += fmt.Sprintf("Snapshots: %d\n", zpool.NumSnapshots)
		}
	case LevelDataset:
		// Dataset details
		if len(m.zbox.ZPools) > 0 {
			zpool := m.zbox.ZPools[m.selected[LevelZPool]]
			if len(zpool.Datasets) > 0 {
				dataset := zpool.Datasets[m.cursor[LevelDataset]]
				previewText += fmt.Sprintf("Name: %s\n", dataset.Name)
				previewText += fmt.Sprintf("Used: %s\n", dataset.Used)
				previewText += fmt.Sprintf("Available: %s\n", dataset.Available)
				previewText += fmt.Sprintf("Mountpoint: %s\n", dataset.Mountpoint)
			}
		}
	case LevelSnapshot:
		// Snapshot details
		if len(m.zbox.ZPools) > 0 {
			zpool := m.zbox.ZPools[m.selected[LevelZPool]]
			if len(zpool.Datasets) > 0 {
				dataset := zpool.Datasets[m.selected[LevelDataset]]
				if len(dataset.Snapshots) > 0 {
					snapshot := dataset.Snapshots[m.cursor[LevelSnapshot]]
					previewText += fmt.Sprintf("Name: %s\n", snapshot.Name)
					previewText += fmt.Sprintf("Size: %s\n", snapshot.Size)
					previewText += fmt.Sprintf("Date: %s\n", snapshot.Date)
				}
			}
		}
	}
	return previewText
}