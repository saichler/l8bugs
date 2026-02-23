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
 * Mobile Reference Registry - Tracking Module
 * Reference configurations for Bug Tracking models
 */
const refTrackingM = window.Layer8RefFactory;

window.Layer8MReferenceRegistryTracking = {
    ...refTrackingM.simple('BugsProject', 'projectId', 'name', 'Project'),
    ...refTrackingM.simple('BugsAssignee', 'assigneeId', 'name', 'Assignee'),
    ...refTrackingM.simple('Bug', 'bugId', 'title', 'Bug'),
    ...refTrackingM.simple('Feature', 'featureId', 'title', 'Feature'),
    ...refTrackingM.simple('BugsSprint', 'sprintId', 'name', 'Sprint')
};

// Register with the central registry
Layer8MReferenceRegistry.register(window.Layer8MReferenceRegistryTracking);
