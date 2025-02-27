package utils

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

func Modulo(a, b int) int {
	result := a % b
	if result < 0 {
		result += b
	}

	return result
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Bounded(a, lower, upper int) int {
	return MaxInt(lower, MinInt(a, upper))
}

type TreeNode interface {
	Children() []TreeNode
	Traversable() bool
	SetTraversable(traversable bool, recursive bool)
}

func DepthFirstTraversal(node TreeNode, traverseAll bool) []TreeNode {
	explored := make([]TreeNode, 0)
	toBeExplored := []TreeNode{node}

	for len(toBeExplored) > 0 {
		node = toBeExplored[len(toBeExplored)-1]
		toBeExplored = toBeExplored[:len(toBeExplored)-1]
		if traverseAll || node.Traversable() {
			children := node.Children()
			for i := len(children) - 1; i >= 0; i-- {
				toBeExplored = append(toBeExplored, children[i])
			}
		}

		explored = append(explored, node)
	}

	return explored
}

func RepoHostOwnerAndName(repositoryURL string) (string, string, string, error) {
	// Turn "git@host:path.git" into "host/path" so that it is compatible with url.Parse()
	if strings.HasPrefix(repositoryURL, "git@") {
		repositoryURL = strings.TrimPrefix(repositoryURL, "git@")
		repositoryURL = strings.Replace(repositoryURL, ":", "/", 1)
	}
	repositoryURL = strings.TrimSuffix(repositoryURL, ".git")

	u, err := url.Parse(repositoryURL)
	if err != nil {
		return "", "", "", err
	}
	if u.Host == "" && !strings.Contains(repositoryURL, "://") {
		// example.com/aaa/bbb is parsed as url.URL{Host: "", Path:"example.com/aaa/bbb"}
		// but we expect url.URL{Host: "example.com", Path:"/aaa/bbb"}. Adding a scheme fixes this.
		//
		u, err = url.Parse("https://" + repositoryURL)
		if err != nil {
			return "", "", "", err
		}
	}

	components := strings.FieldsFunc(u.Path, func(c rune) bool { return c == '/' })
	if len(components) < 2 {
		err := fmt.Errorf("invalid repository path: %q (expected at least two components)",
			u.String())
		return "", "", "", err
	}

	return u.Hostname(), components[0], components[1], nil
}

func Prefix(s string, prefix string) string {
	builder := strings.Builder{}
	for _, line := range strings.Split(s, "\n") {
		builder.WriteString(fmt.Sprintf("%s%s\n", prefix, line))
	}

	return builder.String()
}

type NullString struct {
	Valid  bool
	String string
}

type NullTime struct {
	Valid bool
	Time  time.Time
}

func NullTimeFromTime(t *time.Time) NullTime {
	if t == nil {
		return NullTime{}
	}
	return NullTime{
		Time:  *t,
		Valid: true,
	}
}

func NullTimeFromString(s string) (t NullTime, err error) {
	if s != "" {
		t.Time, err = time.Parse(time.RFC3339, s)
		t.Valid = err == nil
	}

	return
}

func MinNullTime(times ...NullTime) NullTime {
	result := NullTime{}
	for _, t := range times {
		if result.Valid {
			if t.Valid && t.Time.Before(result.Time) {
				result = t
			}
		} else {
			result = t
		}
	}
	return result
}

func MaxNullTime(times ...NullTime) NullTime {
	result := NullTime{}
	for _, t := range times {
		if result.Valid {
			if t.Valid && t.Time.After(result.Time) {
				result = t
			}
		} else {
			result = t
		}
	}
	return result
}

type NullDuration struct {
	Valid    bool
	Duration time.Duration
}

func (d NullDuration) String() string {
	if !d.Valid {
		return "-"
	}

	minutes := d.Duration / time.Minute
	seconds := (d.Duration - minutes*time.Minute) / time.Second

	if minutes == 0 {
		if seconds == 0 {
			return "<1s"
		}
		return fmt.Sprintf("%ds", seconds)
	}
	return fmt.Sprintf("%dm%02ds", minutes, seconds)
}

func NullSub(after NullTime, before NullTime) NullDuration {
	return NullDuration{
		Valid:    after.Valid && before.Valid,
		Duration: after.Time.Sub(before.Time),
	}
}

func getEnvWithDefault(key string, d string) string {
	value := os.Getenv(key)
	if value == "" {
		value = d
	}
	return value
}

// Return possible locations of configuration files based on
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
func XDGConfigLocations(filename string) []string {
	confHome := getEnvWithDefault("XDG_CONFIG_HOME", path.Join(os.Getenv("HOME"), ".config"))
	locations := []string{
		path.Join(confHome, filename),
	}

	dirs := getEnvWithDefault("XDG_CONFIG_DIRS", "/etc/xdg")
	for _, dir := range strings.Split(dirs, ":") {
		locations = append(locations, path.Join(dir, filename))
	}

	return locations
}
