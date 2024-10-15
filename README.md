# ZBoxes | A ZFS Manager TUI

A terminal user interface (TUI) application to manage ZFS file systems on multiple ZBoxes within a local network. This application provides a ranger-style interface using the Bubble Tea and Lip Gloss Go libraries.

Table of Contents

	•	Introduction
	•	Features
	•	Installation
	•	Usage
	•	Project Structure
	•	Main Flows
	•	Data Structures
	•	Configuration
	•	Future Work
	•	License

Introduction

ZFS Manager TUI is designed to help system administrators manage ZFS pools, datasets, and snapshots across multiple ZBoxes (servers running ZFS) on a local network. The application provides a user-friendly, command-line interface for navigating and managing ZFS resources without the need for complex commands.

Features

	•	Ranger-style Navigation: Navigate through ZBoxes, ZPools, Datasets, and Snapshots using keyboard shortcuts.
	•	Real-time Data: Fetches live data from the local ZFS installation.
	•	Customizable Interface: Easily change color schemes through a configuration file.
	•	Extensible Design: Prepared for future integration with remote ZBoxes via SSH.

Installation

Prerequisites

	•	Go 1.16 or higher
	•	ZFS installed on the local machine
	•	Git (to clone the repository)

Steps

	1.	Clone the Repository

git clone https://github.com/yourusername/zfs-manager-tui.git
cd zfs-manager-tui


	2.	Install Dependencies

go get


	3.	Build the Application

go build -o zfs-manager


	4.	Run the Application

./zfs-manager



Usage

Navigate through the ZFS resources using the following keyboard shortcuts:

	•	h: Move left (up a level)
	•	j: Move down (next item)
	•	k: Move up (previous item)
	•	l: Move right (down a level)
	•	q: Quit the application

Project Structure

The project is organized into several files, each responsible for different aspects of the application:

	•	main.go: The entry point of the application.
	•	model.go: Contains the Bubble Tea model and logic for updating and rendering the UI.
	•	zfs.go: Defines data structures and functions to interact with ZFS.
	•	styles.go: Defines the UI styles and reads color schemes from the configuration file.
	•	config.go: Handles loading configuration settings from a TOML file.
	•	constants.go: Defines constants used across the application.
	•	ssh.go: Stubbed for future SSH functionality to connect to remote ZBoxes.
	•	config.toml: Configuration file for customizing the color scheme.

File Breakdown

main.go

	•	Initializes the application.
	•	Loads the configuration.
	•	Starts the Bubble Tea program.

model.go

	•	Defines the model struct, which holds the application state.
	•	Contains the Init, Update, and View methods required by Bubble Tea.
	•	Handles user input and updates the UI accordingly.

zfs.go

	•	Contains data structures for ZBox, ZPool, Dataset, and Snapshot.
	•	Provides constructors for these data structures.
	•	Implements functions to fetch data from ZFS using system commands.

styles.go

	•	Defines UI styles using Lip Gloss.
	•	The styles are initialized based on the configuration file.

config.go

	•	Defines the Config struct, which represents the configuration settings.
	•	Implements the LoadConfig function to read the settings from config.toml.

constants.go

	•	Defines constants for navigation levels (LevelZBox, LevelZPool, etc.).
	•	These constants improve code readability and maintainability.

ssh.go

	•	Contains a stub for future SSH functionality.
	•	Will be used to connect and fetch data from remote ZBoxes.

config.toml

	•	A TOML file that allows customization of the UI color scheme.
	•	Users can modify this file to change the look and feel of the application.

Main Flows

The application operates in a hierarchical manner, allowing users to navigate through different levels of ZFS resources:

	1.	ZBox Level
	•	Displays the local machine (and eventually remote machines).
	•	Shows basic information like the hostname and user.
	2.	ZPool Level
	•	Lists all ZPools on the selected ZBox.
	•	Provides health status and counts of datasets and snapshots.
	3.	Dataset Level
	•	Lists all datasets within the selected ZPool.
	•	Shows usage statistics and mount points.
	4.	Snapshot Level
	•	Lists all snapshots of the selected dataset.
	•	Displays size and creation date.

Navigation

	•	Left Column: Parent items of the current level.
	•	Middle Column: Items at the current level.
	•	Right Column: Child items of the current selection.
	•	Details Pane: Provides detailed information about the selected item.

Data Structures

Snapshot

Represents a ZFS snapshot.

	•	Fields:
	•	Name: The name of the snapshot.
	•	Size: The size of the snapshot.
	•	Date: The creation date of the snapshot.

Dataset

Represents a ZFS dataset.

	•	Fields:
	•	Name: The name of the dataset.
	•	Used: Space used by the dataset.
	•	Available: Space available to the dataset.
	•	Mountpoint: The mount point of the dataset.
	•	Snapshots: A slice of Snapshot instances.

ZPool

Represents a ZFS storage pool.

	•	Fields:
	•	Name: The name of the pool.
	•	Health: Health status of the pool.
	•	NumDatasets: Number of datasets in the pool.
	•	NumSnapshots: Total number of snapshots in the pool.
	•	Datasets: A slice of Dataset instances.

ZBox

Represents a machine with ZFS installed.

	•	Fields:
	•	Name: A friendly name for the ZBox.
	•	Hostname: The hostname of the machine.
	•	User: The current user.
	•	ZPools: A slice of ZPool instances.

Configuration

The application uses a config.toml file for configuration, allowing users to customize the color scheme.

Example config.toml

[colors]
title = "#fab387"
normal_text = "#cdd6f4"
cursor = "#89dceb"
selected = "#cba6f7"
border = "#585b70"
instruction = "#a6adc8"
active_column_bg = "#313244"

Customizing Colors

	•	Title: Color of the column titles.
	•	Normal Text: Default text color.
	•	Cursor: Color of the highlighted item when navigating.
	•	Selected: Color of selected items.
	•	Border: Color of the column borders.
	•	Instruction: Color of the instruction text at the bottom.
	•	Active Column Background: Background color of the active column.

Future Work

	•	Remote ZBoxes via SSH: Implement functionality to connect to and manage ZBoxes over SSH.
	•	Action Commands: Add the ability to perform actions like creating or deleting datasets and snapshots.
	•	Improved Error Handling: Enhance error messages and handling for a better user experience.
	•	Unit Tests: Write unit tests for critical components to ensure reliability.

License

This project is licensed under the MIT License. See the LICENSE file for details.

Note: This application currently only supports managing ZFS on the local machine. SSH functionality for remote ZBoxes is planned for future releases.
