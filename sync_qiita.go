package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	qiitaAPIBase = "https://qiita.com/api/v2"
	qiitaUser    = "ShinguAkira"
)

type Tag struct {
	Name string `json:"name"`
}

type Article struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Tags    []Tag  `json:"tags"`
	URL     string `json:"url"`
	Private bool   `json:"private"`
}

func main() {
	token  := flag.String("token", os.Getenv("QIITA_TOKEN"), "Qiita API token (or $QIITA_TOKEN)")
	outDir := flag.String("out", filepath.Join(repoRoot(), "article"), "output directory")
	dryRun := flag.Bool("dry-run", false, "show what would change without writing")
	flag.Parse()

	articles, err := fetchAll(*token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("fetched %d articles\n\n", len(articles))

	if !*dryRun {
		if err := os.MkdirAll(*outDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "mkdir error: %v\n", err)
			os.Exit(1)
		}
	}

	updated, skipped := 0, 0
	for _, a := range articles {
		name    := sanitize(a.Title) + ".md"
		path    := filepath.Join(*outDir, name)
		content := render(a)

		existing, err := os.ReadFile(path)
		if err == nil && string(existing) == content {
			skipped++
			continue
		}

		action := "update"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			action = "create"
		}

		if *dryRun {
			fmt.Printf("[dry-run] %-6s %s\n", action, path)
		} else {
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "write error %s: %v\n", path, err)
				continue
			}
			fmt.Printf("%-6s %s\n", action, path)
		}
		updated++
	}

	fmt.Printf("\ntotal %d — %d changed, %d unchanged\n", len(articles), updated, skipped)
}

func fetchAll(token string) ([]Article, error) {
	var all []Article
	for page := 1; ; page++ {
		url := fmt.Sprintf("%s/users/%s/items?per_page=100&page=%d", qiitaAPIBase, qiitaUser, page)
		batch, err := apiGet[[]Article](url, token)
		if err != nil {
			return nil, err
		}
		all = append(all, batch...)
		if len(batch) < 100 {
			break
		}
	}
	return all, nil
}

func apiGet[T any](url, token string) (T, error) {
	var zero T
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return zero, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return zero, fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)
	}
	var v T
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return zero, err
	}
	return v, nil
}

func render(a Article) string {
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("title: " + a.Title + "\n")
	b.WriteString("tags:\n")
	for _, t := range a.Tags {
		b.WriteString("  - " + t.Name + "\n")
	}
	fmt.Fprintf(&b, "private: %v\n", a.Private)
	b.WriteString("qiita_url: " + a.URL + "\n")
	b.WriteString("---\n\n")
	body := strings.TrimRight(a.Body, "\n")
	b.WriteString(body + "\n")
	return b.String()
}

func repoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}

var reInvalid = regexp.MustCompile(`[<>:"/\\|?*` + "\x00-\x1f]")

func sanitize(title string) string {
	s := reInvalid.ReplaceAllString(title, "")
	s = strings.TrimSpace(s)
	if s == "" {
		return "untitled"
	}
	return s
}
