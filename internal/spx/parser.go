package spx

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ParseProfile parse spx profile from file
func ParseProfile(filename string) (*Profile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	var scanner *bufio.Scanner

	if filepath.Ext(filename) == ".gz" {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return nil, fmt.Errorf("cannot create gzip reader: %w", err)
		}
		defer gzReader.Close()
		scanner = bufio.NewScanner(gzReader)
	} else {
		scanner = bufio.NewScanner(file)
	}

	profile := &Profile{
		Events:    make([]Event, 0),
		Functions: make(map[int]string),
	}

	inEvents := false
	inFunctions := false
	functionIndex := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		// switch sections
		if line == "[events]" {
			inEvents = true
			inFunctions = false
			continue
		}

		if line == "[functions]" {
			inEvents = false
			inFunctions = true
			continue
		}

		// parse events
		if inEvents {
			event, err := parseEvent(line)
			if err != nil {
				continue // пропускаем неверные строки
			}
			profile.Events = append(profile.Events, event)
		}

		// parse functions
		if inFunctions {
			profile.Functions[functionIndex] = line
			functionIndex++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return profile, nil
}

// parseEvent parse event
func parseEvent(line string) (Event, error) {
	parts := strings.Fields(line)
	if len(parts) != 4 {
		return Event{}, fmt.Errorf("invalid event format: %s", line)
	}

	functionID, err := strconv.Atoi(parts[0])
	if err != nil {
		return Event{}, fmt.Errorf("invalid function ID: %s", parts[0])
	}

	eventType, err := strconv.Atoi(parts[1])
	if err != nil {
		return Event{}, fmt.Errorf("invalid event type: %s", parts[1])
	}

	time, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return Event{}, fmt.Errorf("invalid time: %s", parts[2])
	}

	memory, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return Event{}, fmt.Errorf("invalid memory: %s", parts[3])
	}

	return Event{
		FunctionID: functionID,
		EventType:  eventType,
		Time:       time,
		Memory:     memory,
	}, nil
}
