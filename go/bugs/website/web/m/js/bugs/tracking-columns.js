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
 * Mobile Bugs Tracking Module - Column Definitions
 * Desktop Equivalent: bugs/tracking/tracking-columns.js
 */
(function() {
    'use strict';
    window.MobileBugsTracking = window.MobileBugsTracking || {};

    const col = window.Layer8ColumnFactory;
    const enums = MobileBugsTracking.enums;
    const render = MobileBugsTracking.render;

    MobileBugsTracking.columns = {
        Bug: [
            ...col.id('bugId'),
            ...col.col('bugNumber', 'Bug #'),
            ...col.col('title', 'Title'),
            ...col.enum('status', 'Status', enums.BUG_STATUS_VALUES, render.bugStatus),
            ...col.enum('priority', 'Priority', enums.PRIORITY_VALUES, render.priority),
            ...col.enum('severity', 'Severity', enums.SEVERITY_VALUES, render.severity),
            ...col.col('assigneeId', 'Assignee'),
            ...col.col('component', 'Component'),
            ...col.enum('resolution', 'Resolution', enums.RESOLUTION_VALUES, render.resolution),
            ...col.date('createdDate', 'Created'),
            ...col.date('resolvedDate', 'Resolved')
        ],

        Feature: [
            ...col.id('featureId'),
            ...col.col('featureNumber', 'Feature #'),
            ...col.col('title', 'Title'),
            ...col.enum('status', 'Status', enums.FEATURE_STATUS_VALUES, render.featureStatus),
            ...col.enum('priority', 'Priority', enums.PRIORITY_VALUES, render.priority),
            ...col.col('assigneeId', 'Assignee'),
            ...col.col('component', 'Component'),
            ...col.col('targetVersion', 'Target Version'),
            ...col.date('createdDate', 'Created')
        ],

        BugsProject: [
            ...col.id('projectId'),
            ...col.col('name', 'Name'),
            ...col.col('key', 'Key'),
            ...col.enum('status', 'Status', enums.PROJECT_STATUS_VALUES, render.projectStatus),
            ...col.enum('visibility', 'Visibility', null, render.projectVisibility),
            ...col.col('ownerId', 'Owner')
        ]
    };

    MobileBugsTracking.primaryKeys = {
        Bug: 'bugId',
        Feature: 'featureId',
        BugsProject: 'projectId'
    };
})();
