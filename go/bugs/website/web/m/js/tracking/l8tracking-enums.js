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
 * Mobile L8Tracking Module - Enum Definitions using Layer8EnumFactory
 * Desktop Equivalent: l8ui/sys/tracking/l8tracking-enums.js
 */
(function() {
    'use strict';

    const factory = window.Layer8EnumFactory;
    const { createStatusRenderer, createEnumRenderer, renderDate } = Layer8MRenderers;

    window.MobileL8Tracking = window.MobileL8Tracking || {};

    const BUG_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Open', 'open', 'status-pending'],
        ['Triaged', 'triaged', 'status-active'],
        ['In Progress', 'in-progress', 'status-active'],
        ['In Review', 'in-review', 'status-active'],
        ['Resolved', 'resolved', 'status-active'],
        ['Closed', 'closed', 'status-inactive'],
        ['Reopened', 'reopened', 'status-pending'],
        ['Won\'t Fix', 'wont-fix', 'status-terminated'],
        ['Duplicate', 'duplicate', 'status-terminated'],
        ['Cannot Reproduce', 'cannot-reproduce', 'status-terminated']
    ]);

    const FEATURE_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Proposed', 'proposed', 'status-pending'],
        ['Triaged', 'triaged', 'status-active'],
        ['Approved', 'approved', 'status-active'],
        ['In Progress', 'in-progress', 'status-active'],
        ['In Review', 'in-review', 'status-active'],
        ['Done', 'done', 'status-active'],
        ['Closed', 'closed', 'status-inactive'],
        ['Rejected', 'rejected', 'status-terminated'],
        ['Deferred', 'deferred', 'status-inactive']
    ]);

    const PRIORITY = factory.create([
        ['Unspecified', null, ''],
        ['Critical', 'critical', 'status-terminated'],
        ['High', 'high', 'status-pending'],
        ['Medium', 'medium', 'status-active'],
        ['Low', 'low', 'status-inactive']
    ]);

    const SEVERITY = factory.create([
        ['Unspecified', null, ''],
        ['Blocker', 'blocker', 'status-terminated'],
        ['Major', 'major', 'status-pending'],
        ['Minor', 'minor', 'status-active'],
        ['Trivial', 'trivial', 'status-inactive']
    ]);

    const RESOLUTION = factory.create([
        ['Unspecified', null, ''],
        ['Fixed', 'fixed', 'status-active'],
        ['Won\'t Fix', 'wont-fix', 'status-terminated'],
        ['Duplicate', 'duplicate', 'status-inactive'],
        ['Cannot Reproduce', 'cannot-reproduce', 'status-inactive'],
        ['By Design', 'by-design', 'status-inactive'],
        ['Obsolete', 'obsolete', 'status-terminated']
    ]);

    const ASSIGNEE_TYPE = factory.create([
        ['Unspecified', null, ''],
        ['Human', 'human', ''],
        ['AI Agent', 'ai-agent', '']
    ]);

    const AUTHOR_TYPE = factory.create([
        ['Unspecified', null, ''],
        ['Human', 'human', ''],
        ['AI Agent', 'ai-agent', ''],
        ['System', 'system', '']
    ]);

    const PROJECT_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Active', 'active', 'status-active'],
        ['Archived', 'archived', 'status-inactive']
    ]);

    const PROJECT_VISIBILITY = factory.create([
        ['Unspecified', null, ''],
        ['Public', 'public', ''],
        ['Private', 'private', ''],
        ['Internal', 'internal', '']
    ]);

    const MILESTONE_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Open', 'open', 'status-active'],
        ['Closed', 'closed', 'status-inactive']
    ]);

    const SPRINT_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Planning', 'planning', 'status-pending'],
        ['Active', 'active', 'status-active'],
        ['Completed', 'completed', 'status-inactive']
    ]);

    const TRIAGE_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Pending', 'pending', 'status-pending'],
        ['In Progress', 'in-progress', 'status-active'],
        ['Completed', 'completed', 'status-active'],
        ['Failed', 'failed', 'status-terminated'],
        ['Skipped', 'skipped', 'status-inactive']
    ]);

    const DIGEST_PERIOD = factory.create([
        ['Unspecified', null, ''],
        ['Daily', 'daily', 'status-active'],
        ['Weekly', 'weekly', 'status-active'],
        ['Custom', 'custom', 'status-pending']
    ]);

    const WEBHOOK_EVENT_TYPE = factory.create([
        ['Unspecified', null, ''],
        ['Issue Created', 'issue-created', ''],
        ['Issue Updated', 'issue-updated', ''],
        ['Status Changed', 'status-changed', ''],
        ['Assigned', 'assigned', ''],
        ['Commented', 'commented', '']
    ]);

    MobileL8Tracking.enums = {
        BUG_STATUS: BUG_STATUS.enum,
        BUG_STATUS_VALUES: BUG_STATUS.values,
        BUG_STATUS_CLASSES: BUG_STATUS.classes,
        FEATURE_STATUS: FEATURE_STATUS.enum,
        FEATURE_STATUS_VALUES: FEATURE_STATUS.values,
        FEATURE_STATUS_CLASSES: FEATURE_STATUS.classes,
        PRIORITY: PRIORITY.enum,
        PRIORITY_VALUES: PRIORITY.values,
        PRIORITY_CLASSES: PRIORITY.classes,
        SEVERITY: SEVERITY.enum,
        SEVERITY_VALUES: SEVERITY.values,
        SEVERITY_CLASSES: SEVERITY.classes,
        RESOLUTION: RESOLUTION.enum,
        RESOLUTION_VALUES: RESOLUTION.values,
        RESOLUTION_CLASSES: RESOLUTION.classes,
        ASSIGNEE_TYPE: ASSIGNEE_TYPE.enum,
        AUTHOR_TYPE: AUTHOR_TYPE.enum,
        PROJECT_STATUS: PROJECT_STATUS.enum,
        PROJECT_STATUS_VALUES: PROJECT_STATUS.values,
        PROJECT_STATUS_CLASSES: PROJECT_STATUS.classes,
        PROJECT_VISIBILITY: PROJECT_VISIBILITY.enum,
        MILESTONE_STATUS: MILESTONE_STATUS.enum,
        MILESTONE_STATUS_VALUES: MILESTONE_STATUS.values,
        MILESTONE_STATUS_CLASSES: MILESTONE_STATUS.classes,
        SPRINT_STATUS: SPRINT_STATUS.enum,
        SPRINT_STATUS_VALUES: SPRINT_STATUS.values,
        SPRINT_STATUS_CLASSES: SPRINT_STATUS.classes,
        TRIAGE_STATUS: TRIAGE_STATUS.enum,
        TRIAGE_STATUS_VALUES: TRIAGE_STATUS.values,
        TRIAGE_STATUS_CLASSES: TRIAGE_STATUS.classes,
        DIGEST_PERIOD: DIGEST_PERIOD.enum,
        DIGEST_PERIOD_VALUES: DIGEST_PERIOD.values,
        DIGEST_PERIOD_CLASSES: DIGEST_PERIOD.classes,
        WEBHOOK_EVENT_TYPE: WEBHOOK_EVENT_TYPE.enum
    };

    MobileL8Tracking.render = {
        bugStatus: createStatusRenderer(BUG_STATUS.enum, BUG_STATUS.classes),
        featureStatus: createStatusRenderer(FEATURE_STATUS.enum, FEATURE_STATUS.classes),
        priority: createStatusRenderer(PRIORITY.enum, PRIORITY.classes),
        severity: createStatusRenderer(SEVERITY.enum, SEVERITY.classes),
        resolution: createStatusRenderer(RESOLUTION.enum, RESOLUTION.classes),
        assigneeType: createEnumRenderer(ASSIGNEE_TYPE.enum),
        projectStatus: createStatusRenderer(PROJECT_STATUS.enum, PROJECT_STATUS.classes),
        projectVisibility: createEnumRenderer(PROJECT_VISIBILITY.enum),
        milestoneStatus: createStatusRenderer(MILESTONE_STATUS.enum, MILESTONE_STATUS.classes),
        sprintStatus: createStatusRenderer(SPRINT_STATUS.enum, SPRINT_STATUS.classes),
        triageStatus: createStatusRenderer(TRIAGE_STATUS.enum, TRIAGE_STATUS.classes),
        digestPeriod: createStatusRenderer(DIGEST_PERIOD.enum, DIGEST_PERIOD.classes),
        webhookEventType: createEnumRenderer(WEBHOOK_EVENT_TYPE.enum),
        date: renderDate
    };
})();
