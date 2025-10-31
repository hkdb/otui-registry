package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Plugin struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags,omitempty"`
	License     string    `json:"license,omitempty"`
	Repository  string    `json:"repository"`
	Author      string    `json:"author"`
	Stars       int       `json:"stars,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Language    string    `json:"language,omitempty"`
	InstallType string    `json:"install_type"`
	Package     string    `json:"package,omitempty"`
	Verified    bool      `json:"verified"`
	Official    bool      `json:"official"`
}

type GitHubRepo struct {
	StargazersCount int    `json:"stargazers_count"`
	Language        string `json:"language"`
	License         *struct {
		SPDXID string `json:"spdx_id"`
	} `json:"license"`
	PushedAt time.Time `json:"pushed_at"`
}

var categoryMap = map[string]string{
	"productivity":  "productivity",
	"development":   "development",
	"ai":            "ai",
	"data":          "data",
	"finance":       "finance",
	"communication": "communication",
	"security":      "security",
	"utility":       "utility",
	"integration":   "integration",

	"aggregators":                   "utility",
	"aerospace-and-astrodynamics":   "productivity",
	"art-and-culture":               "productivity",
	"bio":                           "data",
	"browser-automation":            "development",
	"cloud-platforms":               "integration",
	"code-execution":                "development",
	"coding-agents":                 "development",
	"command-line":                  "development",
	"customer-data-platforms":       "data",
	"databases":                     "data",
	"data-platforms":                "data",
	"developer-tools":               "development",
	"delivery":                      "integration",
	"data-science-tools":            "data",
	"embedded-system":               "development",
	"file-systems":                  "utility",
	"finance--fintech":              "finance",
	"gaming":                        "productivity",
	"knowledge--memory":             "ai",
	"location-services":             "integration",
	"marketing":                     "productivity",
	"monitoring":                    "development",
	"multimedia-process":            "productivity",
	"news--information":             "productivity",
	"productivity-and-organization": "productivity",
	"project-management":            "productivity",
	"search":                        "utility",
	"smart-home":                    "integration",
	"security-and-privacy":          "security",
	"social-media":                  "communication",
	"sports":                        "productivity",
	"testing":                       "development",
	"travel":                        "productivity",
	"version-control":               "development",
	"web-automation":                "development",
	"web3--blockchain":              "integration",
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: parse-registry <input.md> <output.json> [custom-plugins.json]")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	var customPluginsFile string
	if len(os.Args) >= 4 {
		customPluginsFile = os.Args[3]
	}

	plugins, err := parseMarkdown(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing markdown: %v\n", err)
		os.Exit(1)
	}

	enrichPlugins(plugins)

	// Load and merge custom plugins if provided
	if customPluginsFile != "" {
		customPlugins, err := loadCustomPlugins(customPluginsFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading custom plugins: %v\n", err)
			os.Exit(1)
		}
		enrichPlugins(customPlugins)
		plugins = mergePlugins(plugins, customPlugins)
		fmt.Printf("Merged %d custom plugins\n", len(customPlugins))
	}

	// Generate JSON file (all plugins, not just top 50)
	if err := generateJSONFile(plugins, outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating JSON file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s with %d plugins\n", outputFile, len(plugins))
}

func parseMarkdown(filename string) ([]Plugin, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var plugins []Plugin
	scanner := bufio.NewScanner(file)

	linkRegex := regexp.MustCompile(`^\s*-\s*\[([^\]]+)\]\(([^)]+)\).*?-\s*(.+)`)
	categoryRegex := regexp.MustCompile(`<a name="([^"]+)"></a>`)
	currentCategory := "utility"

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "### ") {
			if matches := categoryRegex.FindStringSubmatch(line); matches != nil {
				anchorName := matches[1]
				if mapped, ok := categoryMap[anchorName]; ok {
					currentCategory = mapped
				}
			}
			continue
		}

		matches := linkRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		name := strings.TrimSpace(matches[1])
		repo := strings.TrimSpace(matches[2])
		desc := strings.TrimSpace(matches[3])

		if !strings.HasPrefix(repo, "https://github.com/") {
			continue
		}

		id := generateID(name, repo)
		author := extractAuthor(repo)

		plugin := Plugin{
			ID:          id,
			Name:        name,
			Description: desc,
			Category:    currentCategory,
			Repository:  repo,
			Author:      author,
			InstallType: "",
			Package:     "",
			Verified:    false,
			Official:    strings.Contains(strings.ToLower(author), "anthropic") || strings.Contains(strings.ToLower(author), "modelcontextprotocol"),
		}

		plugins = append(plugins, plugin)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return plugins, nil
}

func generateID(name, repo string) string {
	id := strings.ToLower(name)
	id = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(id, "-")
	id = strings.Trim(id, "-")
	return id
}

func extractAuthor(repo string) string {
	parts := strings.Split(strings.TrimPrefix(repo, "https://github.com/"), "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func cleanRepoPath(repoURL string) string {
	path := strings.TrimPrefix(repoURL, "https://github.com/")
	path = strings.TrimSuffix(path, "/")

	if idx := strings.Index(path, "/tree/"); idx != -1 {
		path = path[:idx]
	}
	if idx := strings.Index(path, "/blob/"); idx != -1 {
		path = path[:idx]
	}

	return path
}

func detectInstallType(language string, repoPath string) (string, string) {
	lowerLang := strings.ToLower(language)

	parts := strings.Split(repoPath, "/")
	var packageName string
	if len(parts) >= 2 {
		packageName = parts[1]
	}

	switch lowerLang {
	case "python":
		return "pip", packageName
	case "typescript", "javascript":
		return "npm", repoPath
	case "go":
		return "go", "github.com/" + repoPath + "@latest"
	default:
		return "manual", ""
	}
}

func enrichPlugins(plugins []Plugin) {
	token := os.Getenv("GITHUB_TOKEN")

	for i := range plugins {
		repoPath := cleanRepoPath(plugins[i].Repository)

		ghData, err := fetchGitHubData(repoPath, token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to fetch GitHub data for %s: %v\n", repoPath, err)
			plugins[i].InstallType = "manual"
			plugins[i].Package = ""
			continue
		}

		plugins[i].Stars = ghData.StargazersCount
		plugins[i].Language = ghData.Language
		if ghData.License != nil {
			plugins[i].License = ghData.License.SPDXID
		}
		plugins[i].UpdatedAt = ghData.PushedAt

		installType, pkg := detectInstallType(ghData.Language, repoPath)
		plugins[i].InstallType = installType
		plugins[i].Package = pkg

		// Special case: if this is from modelcontextprotocol/servers, mark as manual install
    if strings.Contains(repoPath, "modelcontextprotocol/servers") {
        plugins[i].InstallType = "manual"
        // Keep the original package name or set a descriptive one
        if plugins[i].Package == "" {
            plugins[i].Package = "modelcontextprotocol/servers"
        }
    }

		time.Sleep(100 * time.Millisecond)
	}
}

func fetchGitHubData(repoPath, token string) (*GitHubRepo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s", repoPath)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, string(body))
	}

	var ghData GitHubRepo
	if err := json.NewDecoder(resp.Body).Decode(&ghData); err != nil {
		return nil, err
	}

	return &ghData, nil
}

func loadCustomPlugins(filename string) ([]Plugin, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var customPlugins []Plugin
	if err := json.Unmarshal(data, &customPlugins); err != nil {
		return nil, err
	}

	// Convert simple format to full Plugin format
	for i := range customPlugins {
		if customPlugins[i].ID == "" {
			customPlugins[i].ID = generateID(customPlugins[i].Name, customPlugins[i].Repository)
		}
		if customPlugins[i].Author == "" {
			customPlugins[i].Author = extractAuthor(customPlugins[i].Repository)
		}
	}

	return customPlugins, nil
}

func mergePlugins(awesomePlugins, customPlugins []Plugin) []Plugin {
	// Create map of awesome plugins by repository URL for deduplication
	pluginMap := make(map[string]Plugin)
	for _, p := range awesomePlugins {
		pluginMap[p.Repository] = p
	}

	// Custom plugins override awesome plugins (deduplicate by repo URL)
	for _, p := range customPlugins {
		pluginMap[p.Repository] = p
	}

	// Convert map back to slice
	result := make([]Plugin, 0, len(pluginMap))
	for _, p := range pluginMap {
		result = append(result, p)
	}

	return result
}

func generateJSONFile(plugins []Plugin, filename string) error {
	// Sort plugins for consistent output
	sort.Slice(plugins, func(i, j int) bool {
		if plugins[i].Official != plugins[j].Official {
			return plugins[i].Official
		}
		if plugins[i].Stars != plugins[j].Stars {
			return plugins[i].Stars > plugins[j].Stars
		}
		return plugins[i].Name < plugins[j].Name
	})

	data, err := json.MarshalIndent(plugins, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
