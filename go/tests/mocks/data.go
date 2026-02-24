/*
© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package mocks

// Project data
var projectNames = []string{
	"Backend API", "Frontend App", "Mobile Client", "Data Pipeline",
	"Auth Service",
}

var projectKeys = []string{
	"BAPI", "FAPP", "MOBC", "DPIP", "AUTH",
}

var projectDescriptions = []string{
	"Core REST API service powering all clients",
	"React-based web application for end users",
	"Cross-platform mobile application (iOS/Android)",
	"ETL and data processing pipeline",
	"Authentication and authorization microservice",
}

// Assignee data
var assigneeNames = []string{
	"Alice Chen", "Bob Martinez", "Carol Kim", "David Okafor",
	"Eve Johansson", "Frank Patel", "Grace Wu", "CodeBot AI",
}

var assigneeEmails = []string{
	"alice@example.com", "bob@example.com", "carol@example.com",
	"david@example.com", "eve@example.com", "frank@example.com",
	"grace@example.com", "codebot@example.com",
}

// Bug data
var bugTitles = []string{
	"Login page crashes on Safari",
	"API returns 500 on large payload",
	"Memory leak in WebSocket handler",
	"Pagination breaks with special characters in search",
	"Dark mode toggle not persisted across sessions",
	"Race condition in concurrent writes",
	"Incorrect timezone conversion for UTC+13",
	"File upload silently fails over 10MB",
	"OAuth callback URL mismatch in production",
	"Missing CSRF token validation on POST endpoints",
	"Database connection pool exhaustion under load",
	"Broken link in password reset email",
	"Infinite scroll stops loading after page 5",
	"Export CSV generates corrupt UTF-8 for CJK characters",
	"Push notification not delivered on Android 14",
	"Session expires during active form editing",
	"Sort by date shows future dates first",
	"Search results include deleted records",
	"Duplicate email sent on retry",
	"Chart tooltip overlaps axis labels",
}

var bugDescriptions = []string{
	"When opening the login page in Safari 17, the page crashes with a blank white screen. No errors in console.",
	"Sending a POST request with a payload larger than 1MB causes an HTTP 500 response. Works fine under 1MB.",
	"The WebSocket connection handler leaks goroutines when clients disconnect abruptly.",
	"Searching for terms containing '&' or '%' breaks the pagination offset calculation.",
	"User preference for dark mode resets to light mode after browser restart.",
}

// Feature data
var featureTitles = []string{
	"Dark mode support",
	"Export to CSV",
	"Two-factor authentication",
	"Real-time collaboration",
	"API rate limiting",
	"Bulk import from CSV",
	"Custom dashboard widgets",
	"Email notification preferences",
	"Audit trail viewer",
	"Webhook configuration UI",
}

var featureDescriptions = []string{
	"Add a system-wide dark mode theme with automatic OS detection.",
	"Allow users to export any table view to CSV format with column selection.",
	"Implement TOTP-based two-factor authentication for all user accounts.",
	"Enable real-time collaborative editing of issues and features.",
	"Add configurable API rate limiting per user/API key with dashboard.",
}

// Component data
var componentNames = []string{
	"Authentication", "API Gateway", "Database", "Frontend",
	"Notifications", "Search", "File Storage", "Analytics",
}

// Sprint data
var sprintNames = []string{
	"Sprint 1 - Foundation", "Sprint 2 - Core Features",
	"Sprint 3 - Integration", "Sprint 4 - Polish",
	"Sprint 5 - Performance", "Sprint 6 - Release",
}

var sprintGoals = []string{
	"Establish project foundations and core data models",
	"Implement primary user-facing features",
	"Integrate third-party services and APIs",
	"UI polish, accessibility improvements, and bug fixes",
	"Performance optimization and load testing",
	"Final testing, documentation, and release preparation",
}

// Digest data
var digestSummaries = []string{
	"Strong progress on authentication module. 8 bugs resolved, 3 new features completed.",
	"Focus on API stability. Resolved critical memory leak. 2 features in review.",
	"Sprint planning complete. Backlog groomed. 5 high-priority bugs triaged.",
	"Release preparation underway. All blockers resolved. Final QA in progress.",
}
