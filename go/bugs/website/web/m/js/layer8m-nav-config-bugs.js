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
/**
 * Mobile L8Bugs - Navigation Configuration
 * Defines modules, submodules, and services for mobile navigation
 */
(function() {
    'use strict';

    window.Layer8MNavConfig = {
        modules: [
            { key: 'bugs', label: 'Bug Tracking', icon: 'bugs' },
            { key: 'system', label: 'System', icon: 'system' }
        ],

        moduleIcons: {
            bugs: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M8 2l1.5 1.5"/><path d="M14.5 2L16 3.5"/><path d="M9.5 5.5L8 7"/><path d="M14.5 5.5L16 7"/><circle cx="12" cy="14" r="6"/><path d="M12 8v2"/><path d="M6 14H4"/><path d="M20 14h-2"/></svg>',
            system: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>'
        },

        bugs: {
            subModules: [
                { key: 'tracking', label: 'Tracking', icon: 'bugs' }
            ],

            services: {
                'tracking': [
                    {
                        key: 'bugs',
                        label: 'Bugs',
                        icon: 'bugs',
                        endpoint: '/20/Bug',
                        model: 'Bug',
                        idField: 'bugId',
                        supportedViews: ['table', 'kanban']
                    },
                    {
                        key: 'features',
                        label: 'Features',
                        icon: 'bugs',
                        endpoint: '/20/Feature',
                        model: 'Feature',
                        idField: 'featureId',
                        supportedViews: ['table', 'kanban']
                    },
                    {
                        key: 'projects',
                        label: 'Projects',
                        icon: 'bugs',
                        endpoint: '/20/Project',
                        model: 'BugsProject',
                        idField: 'projectId'
                    }
                ]
            }
        }
    };
})();
