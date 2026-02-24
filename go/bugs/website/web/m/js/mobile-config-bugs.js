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
 * Mobile Config Data - Bugs Modules & Reference Registry
 * Registers Bugs module navigation and reference picker configs into Layer8MConfig.
 */
(function() {
    'use strict';

    // Bugs Reference Registry for mobile picker
    Layer8MConfig.registerReferences({
        BugsProject: {
            idColumn: 'projectId', displayColumn: 'name', endpoint: '/20/Project',
            displayField: 'name', idField: 'projectId', searchFields: ['name', 'key']
        },
        Bug: {
            idColumn: 'bugId', displayColumn: 'title', endpoint: '/20/Bug',
            displayField: 'title', idField: 'bugId', searchFields: ['title', 'bugNumber']
        },
        Feature: {
            idColumn: 'featureId', displayColumn: 'title', endpoint: '/20/Feature',
            displayField: 'title', idField: 'featureId', searchFields: ['title', 'featureNumber']
        },
        BugsDigest: {
            idColumn: 'digestId', displayColumn: 'summary', endpoint: '/20/Digest',
            displayField: 'summary', idField: 'digestId', searchFields: ['summary']
        }
    });

})();
