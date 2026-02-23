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
// SYS Module - Configuration
// Module definitions and service mappings for System Administration

(function() {
    'use strict';

    // Create SYS namespace
    window.L8Sys = window.L8Sys || {};

    // SYS Module Configuration
    L8Sys.modules = {
        'tracking': {
            label: 'Tracking',
            icon: '\uD83D\uDC1B',
            services: [
                { key: 'bugs', label: 'Bugs', icon: '\uD83D\uDC1B',
                  endpoint: '/20/Bug', model: 'Bug',
                  supportedViews: ['table', 'kanban'],
                  viewConfig: {
                      laneField: 'status',
                      lanes: {
                          1: { label: 'Open', color: '#f59e0b' },
                          2: { label: 'Triaged', color: '#0ea5e9' },
                          3: { label: 'In Progress', color: '#8b5cf6' },
                          4: { label: 'In Review', color: '#6366f1' },
                          5: { label: 'Resolved', color: '#22c55e' },
                          6: { label: 'Closed', color: '#64748b' }
                      },
                      cardTitle: 'title',
                      cardSubtitle: 'bugNumber',
                      cardFields: ['priority', 'severity', 'assigneeId']
                  }
                },
                { key: 'features', label: 'Features', icon: '\u2728',
                  endpoint: '/20/Feature', model: 'Feature',
                  supportedViews: ['table', 'kanban'],
                  viewConfig: {
                      laneField: 'status',
                      lanes: {
                          1: { label: 'Proposed', color: '#f59e0b' },
                          2: { label: 'Triaged', color: '#0ea5e9' },
                          3: { label: 'Approved', color: '#22c55e' },
                          4: { label: 'In Progress', color: '#8b5cf6' },
                          5: { label: 'In Review', color: '#6366f1' },
                          6: { label: 'Done', color: '#10b981' },
                          7: { label: 'Closed', color: '#64748b' }
                      },
                      cardTitle: 'title',
                      cardSubtitle: 'featureNumber',
                      cardFields: ['priority', 'assigneeId']
                  }
                },
                { key: 'projects', label: 'Projects', icon: '\uD83D\uDCC1',
                  endpoint: '/20/Project', model: 'BugsProject' },
                { key: 'assignees', label: 'Assignees', icon: '\uD83D\uDC64',
                  endpoint: '/20/Assignee', model: 'BugsAssignee' },
                { key: 'sprints', label: 'Sprints', icon: '\uD83C\uDFC3',
                  endpoint: '/20/Sprint', model: 'BugsSprint' }
            ]
        },
        'health': {
            label: 'Health',
            icon: '💚',
            services: []
        },
        'security': {
            label: 'Security',
            icon: '🔒',
            services: [
                { key: 'users', label: 'Users', icon: '👤', endpoint: '/73/users', model: 'L8User' },
                { key: 'roles', label: 'Roles', icon: '🛡️', endpoint: '/74/roles', model: 'L8Role' },
                { key: 'credentials', label: 'Credentials', icon: '🔑', endpoint: '/75/Creds', model: 'L8Credentials' }
            ]
        },
        'modules': {
            label: 'Modules',
            icon: '🧩',
            services: []
        },
        'logs': {
            label: 'Logs',
            icon: '📋',
            services: []
        }
    };

    // Sub-module namespaces for service registry
    L8Sys.submodules = ['L8Security', 'L8Tracking'];

})();
