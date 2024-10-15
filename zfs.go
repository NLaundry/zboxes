// zfs.go
package main

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Data Structures

type Snapshot struct {
	Name string
	Size string
	Date string
}

func NewSnapshot(name, size, date string) Snapshot {
	return Snapshot{
		Name: name,
		Size: size,
		Date: date,
	}
}

type Dataset struct {
	Name       string
	Used       string
	Available  string
	Mountpoint string
	Snapshots  []Snapshot
}

func NewDataset(name, used, available, mountpoint string) (Dataset, error) {
	dataset := Dataset{
		Name:       name,
		Used:       used,
		Available:  available,
		Mountpoint: mountpoint,
	}

	// Get snapshots for this dataset
	snapshots, err := getSnapshots(name)
	if err != nil {
		return dataset, err
	}
	dataset.Snapshots = snapshots

	return dataset, nil
}

type ZPool struct {
	Name         string
	Health       string
	NumDatasets  int
	NumSnapshots int
	Datasets     []Dataset
}

func NewZPool(name, health string) (ZPool, error) {
	zpool := ZPool{
		Name:   name,
		Health: health,
	}

	// Get datasets for this zpool
	datasets, err := getDatasets(name)
	if err != nil {
		return zpool, err
	}
	zpool.Datasets = datasets
	zpool.NumDatasets = len(datasets)

	// Count snapshots
	numSnapshots := 0
	for _, ds := range datasets {
		numSnapshots += len(ds.Snapshots)
	}
	zpool.NumSnapshots = numSnapshots

	return zpool, nil
}

type ZBox struct {
	Name     string
	Hostname string
	User     string
	ZPools   []ZPool
}

func NewLocalZBox() (ZBox, error) {
	var zbox ZBox

	hostname, err := os.Hostname()
	if err != nil {
		return zbox, err
	}
	zbox.Hostname = hostname
	zbox.Name = "Local ZBox"
	zbox.User = os.Getenv("USER")

	zpools, err := getZPools()
	if err != nil {
		return zbox, err
	}
	zbox.ZPools = zpools

	return zbox, nil
}

// Functions to get ZFS data

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

			zpool, err := NewZPool(zpoolName, zpoolHealth)
			if err != nil {
				return zpools, err
			}

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

			dataset, err := NewDataset(datasetName, used, avail, mountpoint)
			if err != nil {
				return datasets, err
			}

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

			snapshot := NewSnapshot(snapshotName, size, formattedDate)
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