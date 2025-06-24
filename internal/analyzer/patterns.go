package analyzer

import "regexp"

// FrameworkPatterns holds regex patterns for different frameworks
type FrameworkPatterns struct {
	Name         string
	FilePatterns []string // File extensions to look for
	Patterns     []Pattern
}

// Pattern represents a regex pattern for finding endpoints
type Pattern struct {
	Regex          *regexp.Regexp
	MethodIndex    int // Capture group index for HTTP method
	PathIndex      int // Capture group index for path
	FunctionIndex  int // Capture group index for function name (optional)
	IsMethodFirst  bool // If true, method comes before path in regex
}

// GetAllPatterns returns patterns for all supported frameworks
func GetAllPatterns() []FrameworkPatterns {
	return []FrameworkPatterns{
		// Express.js / Node.js
		{
			Name:         "Express",
			FilePatterns: []string{".js", ".ts", ".mjs"},
			Patterns: []Pattern{
				{
					// app.get('/path', handler)
					Regex:         regexp.MustCompile(`app\.(get|post|put|delete|patch|options|head)\s*\(\s*['"\` + "`" + `]([^'"\` + "`" + `]+)['"\` + "`" + `]`),
					MethodIndex:   1,
					PathIndex:     2,
					IsMethodFirst: true,
				},
				{
					// router.get('/path', handler)
					Regex:         regexp.MustCompile(`router\.(get|post|put|delete|patch|options|head)\s*\(\s*['"\` + "`" + `]([^'"\` + "`" + `]+)['"\` + "`" + `]`),
					MethodIndex:   1,
					PathIndex:     2,
					IsMethodFirst: true,
				},
				{
					// app.route('/path').get(handler)
					Regex:         regexp.MustCompile(`app\.route\s*\(\s*['"\` + "`" + `]([^'"\` + "`" + `]+)['"\` + "`" + `]\s*\)\s*\.(get|post|put|delete|patch)`),
					MethodIndex:   2,
					PathIndex:     1,
					IsMethodFirst: false,
				},
			},
		},
		// Flask / Python
		{
			Name:         "Flask",
			FilePatterns: []string{".py"},
			Patterns: []Pattern{
				{
					// @app.route('/path', methods=['GET'])
					Regex:         regexp.MustCompile(`@app\.route\s*\(\s*['"]([^'"]+)['"]\s*(?:,\s*methods\s*=\s*\[['"](\w+)['"]\])?`),
					MethodIndex:   2,
					PathIndex:     1,
					IsMethodFirst: false,
				},
				{
					// @blueprint.route('/path', methods=['POST'])
					Regex:         regexp.MustCompile(`@\w+\.route\s*\(\s*['"]([^'"]+)['"]\s*(?:,\s*methods\s*=\s*\[['"](\w+)['"]\])?`),
					MethodIndex:   2,
					PathIndex:     1,
					IsMethodFirst: false,
				},
			},
		},
		// FastAPI / Python
		{
			Name:         "FastAPI",
			FilePatterns: []string{".py"},
			Patterns: []Pattern{
				{
					// @app.get("/path")
					Regex:         regexp.MustCompile(`@app\.(get|post|put|delete|patch)\s*\(\s*["']([^"']+)["']`),
					MethodIndex:   1,
					PathIndex:     2,
					IsMethodFirst: true,
				},
				{
					// @router.get("/path")
					Regex:         regexp.MustCompile(`@router\.(get|post|put|delete|patch)\s*\(\s*["']([^"']+)["']`),
					MethodIndex:   1,
					PathIndex:     2,
					IsMethodFirst: true,
				},
			},
		},
		// Spring Boot / Java
		{
			Name:         "Spring",
			FilePatterns: []string{".java"},
			Patterns: []Pattern{
				{
					// @GetMapping("/path")
					Regex:         regexp.MustCompile(`@(Get|Post|Put|Delete|Patch)Mapping\s*\(\s*["']([^"']+)["']`),
					MethodIndex:   1,
					PathIndex:     2,
					IsMethodFirst: true,
				},
				{
					// @RequestMapping(value = "/path", method = RequestMethod.GET)
					Regex:         regexp.MustCompile(`@RequestMapping\s*\([^)]*value\s*=\s*["']([^"']+)["'][^)]*method\s*=\s*RequestMethod\.(\w+)`),
					MethodIndex:   2,
					PathIndex:     1,
					IsMethodFirst: false,
				},
			},
		},
		// Gin / Go
		{
			Name:         "Gin",
			FilePatterns: []string{".go"},
			Patterns: []Pattern{
				{
					// router.GET("/path", handler)
					Regex:         regexp.MustCompile(`router\.(GET|POST|PUT|DELETE|PATCH)\s*\(\s*["` + "`" + `]([^"` + "`" + `]+)["` + "`" + `]`),
					MethodIndex:   1,
					PathIndex:     2,
					IsMethodFirst: true,
				},
				{
					// r.GET("/path", handler)
					Regex:         regexp.MustCompile(`\br\.(GET|POST|PUT|DELETE|PATCH)\s*\(\s*["` + "`" + `]([^"` + "`" + `]+)["` + "`" + `]`),
					MethodIndex:   1,
					PathIndex:     2,
					IsMethodFirst: true,
				},
			},
		},
		// Echo / Go
		{
			Name:         "Echo",
			FilePatterns: []string{".go"},
			Patterns: []Pattern{
				{
					// e.GET("/path", handler)
					Regex:         regexp.MustCompile(`e\.(GET|POST|PUT|DELETE|PATCH)\s*\(\s*["` + "`" + `]([^"` + "`" + `]+)["` + "`" + `]`),
					MethodIndex:   1,
					PathIndex:     2,
					IsMethodFirst: true,
				},
			},
		},
		// Ruby on Rails
		{
			Name:         "Rails",
			FilePatterns: []string{".rb"},
			Patterns: []Pattern{
				{
					// get '/path', to: 'controller#action'
					Regex:         regexp.MustCompile(`^\s*(get|post|put|patch|delete)\s+['"]([^'"]+)['"]`),
					MethodIndex:   1,
					PathIndex:     2,
					IsMethodFirst: true,
				},
				{
					// resources :users (generates multiple routes)
					Regex:         regexp.MustCompile(`^\s*resources\s+:(\w+)`),
					MethodIndex:   0, // Special case - generates multiple methods
					PathIndex:     1,
					IsMethodFirst: false,
				},
			},
		},
		// ASP.NET Core / C#
		{
			Name:         "ASP.NET",
			FilePatterns: []string{".cs"},
			Patterns: []Pattern{
				{
					// [HttpGet("/path")]
					Regex:         regexp.MustCompile(`\[Http(Get|Post|Put|Delete|Patch)\s*\(\s*["']([^"']+)["']\s*\)`),
					MethodIndex:   1,
					PathIndex:     2,
					IsMethodFirst: true,
				},
				{
					// [Route("api/[controller]")]
					Regex:         regexp.MustCompile(`\[Route\s*\(\s*["']([^"']+)["']\s*\)`),
					MethodIndex:   0, // No method in route attribute
					PathIndex:     1,
					IsMethodFirst: false,
				},
			},
		},
	}
} 