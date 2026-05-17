package cmd

import (
	"fmt"
	"sort"
	"strings"
)

func parseStatus(raw string) (state, uptime, health string) {
	switch {
	case strings.HasPrefix(raw, "Up"):
		state = "up"
	case strings.HasPrefix(raw, "Exited"):
		state = "exited"
	case strings.HasPrefix(raw, "Restarting"):
		state = "restarting"
	default:
		state = "unknown"
	}

	rest := raw[strings.Index(raw, " ")+1:]

	openIdx := strings.Index(rest, " (")
	if openIdx == -1 {
		uptime = rest
		return
	}

	uptime = rest[:openIdx]

	closeIdx := strings.Index(rest, ")")
	if closeIdx == -1 {
		return
	}

	inner := rest[openIdx+2 : closeIdx]
	health = strings.TrimPrefix(inner, "health: ")
	return
}

func parseContainers(out string) []Container {
	containers := []Container{}
	lines := strings.Split(out, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) != 4 {
			continue
		}

		state, uptime, health := parseStatus(parts[1])

		c := Container{
			Name:   parts[0],
			Image:  parts[2],
			Stack:  parts[3],
			State:  state,
			Uptime: uptime,
			Health: health,
		}
		containers = append(containers, c)
	}
	return containers
}

func sortContainers(containers []Container, col string, asc bool) []Container {
	sort.Slice(containers, func(i, j int) bool {
		if col == "uptime" {
			ai := uptimeToMinutes(containers[i].Uptime)
			bi := uptimeToMinutes(containers[j].Uptime)
			if asc {
				return ai < bi
			}
			return ai > bi
		}

		var a, b string
		switch col {
		case "name":
			a, b = containers[i].Name, containers[j].Name
		case "stack":
			a, b = containers[i].Stack, containers[j].Stack
		case "health":
			a, b = containers[i].Health, containers[j].Health
		default:
			a, b = containers[i].Name, containers[j].Name
		}

		if asc {
			return a < b
		}
		return a > b
	})
	return containers
}

func shortImage(image string) string {
	parts := strings.Split(image, "/")
	return parts[len(parts)-1]
}

func truncate(s string, max int) string {
    if len(s) <= max {
        return s
    }
    return s[:max-3] + "..."
}

func uptimeToMinutes(uptime string) int {
	if strings.HasPrefix(uptime, "About an") || uptime == "an hour" {
		return 60
	}

	parts := strings.Fields(uptime)
	if len(parts) < 2 {
		return 0
	}

	val := 0
	fmt.Sscanf(parts[0], "%d", &val)

	switch parts[1] {
	case "seconds", "second":
		return val
	case "minutes", "minute":
		return val
	case "hours", "hour":
		return val * 60
	case "days", "day":
		return val * 60 * 24
	case "weeks", "week":
		return val * 60 * 24 * 7
	}
	return 0
}

func normalizeUptime(uptime string) string {
    if strings.HasPrefix(uptime, "About an") || uptime == "an hour" {
        return "~1 hour"
    }
    return uptime
}
