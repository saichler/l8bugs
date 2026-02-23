(function() {
    'use strict';
    const factory = window.Layer8EnumFactory;
    const { createStatusRenderer, renderEnum, renderDate } = Layer8DRenderers;

    window.L8Tracking = window.L8Tracking || {};

    const BUG_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Open', 'open', 'layer8d-status-pending'],
        ['Triaged', 'triaged', 'layer8d-status-active'],
        ['In Progress', 'in-progress', 'layer8d-status-active'],
        ['In Review', 'in-review', 'layer8d-status-active'],
        ['Resolved', 'resolved', 'layer8d-status-active'],
        ['Closed', 'closed', 'layer8d-status-inactive'],
        ['Reopened', 'reopened', 'layer8d-status-pending'],
        ['Won\'t Fix', 'wont-fix', 'layer8d-status-terminated'],
        ['Duplicate', 'duplicate', 'layer8d-status-terminated'],
        ['Cannot Reproduce', 'cannot-reproduce', 'layer8d-status-terminated']
    ]);

    const FEATURE_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Proposed', 'proposed', 'layer8d-status-pending'],
        ['Triaged', 'triaged', 'layer8d-status-active'],
        ['Approved', 'approved', 'layer8d-status-active'],
        ['In Progress', 'in-progress', 'layer8d-status-active'],
        ['In Review', 'in-review', 'layer8d-status-active'],
        ['Done', 'done', 'layer8d-status-active'],
        ['Closed', 'closed', 'layer8d-status-inactive'],
        ['Rejected', 'rejected', 'layer8d-status-terminated'],
        ['Deferred', 'deferred', 'layer8d-status-inactive']
    ]);

    const PRIORITY = factory.create([
        ['Unspecified', null, ''],
        ['Critical', 'critical', 'layer8d-status-terminated'],
        ['High', 'high', 'layer8d-status-pending'],
        ['Medium', 'medium', 'layer8d-status-active'],
        ['Low', 'low', 'layer8d-status-inactive']
    ]);

    const SEVERITY = factory.create([
        ['Unspecified', null, ''],
        ['Blocker', 'blocker', 'layer8d-status-terminated'],
        ['Major', 'major', 'layer8d-status-pending'],
        ['Minor', 'minor', 'layer8d-status-active'],
        ['Trivial', 'trivial', 'layer8d-status-inactive']
    ]);

    const RESOLUTION = factory.create([
        ['Unspecified', null, ''],
        ['Fixed', 'fixed', 'layer8d-status-active'],
        ['Won\'t Fix', 'wont-fix', 'layer8d-status-terminated'],
        ['Duplicate', 'duplicate', 'layer8d-status-inactive'],
        ['Cannot Reproduce', 'cannot-reproduce', 'layer8d-status-inactive'],
        ['By Design', 'by-design', 'layer8d-status-inactive'],
        ['Obsolete', 'obsolete', 'layer8d-status-terminated']
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
        ['Active', 'active', 'layer8d-status-active'],
        ['Archived', 'archived', 'layer8d-status-inactive']
    ]);

    const PROJECT_VISIBILITY = factory.create([
        ['Unspecified', null, ''],
        ['Public', 'public', ''],
        ['Private', 'private', ''],
        ['Internal', 'internal', '']
    ]);

    const MILESTONE_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Open', 'open', 'layer8d-status-active'],
        ['Closed', 'closed', 'layer8d-status-inactive']
    ]);

    const TRIAGE_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Pending', 'pending', 'layer8d-status-pending'],
        ['In Progress', 'in-progress', 'layer8d-status-active'],
        ['Completed', 'completed', 'layer8d-status-active'],
        ['Failed', 'failed', 'layer8d-status-terminated'],
        ['Skipped', 'skipped', 'layer8d-status-inactive']
    ]);

    const SPRINT_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Planning', 'planning', 'layer8d-status-pending'],
        ['Active', 'active', 'layer8d-status-active'],
        ['Completed', 'completed', 'layer8d-status-inactive']
    ]);

    L8Tracking.enums = {
        BUG_STATUS: BUG_STATUS.enum,
        BUG_STATUS_CLASSES: BUG_STATUS.classes,
        FEATURE_STATUS: FEATURE_STATUS.enum,
        FEATURE_STATUS_CLASSES: FEATURE_STATUS.classes,
        PRIORITY: PRIORITY.enum,
        PRIORITY_CLASSES: PRIORITY.classes,
        SEVERITY: SEVERITY.enum,
        SEVERITY_CLASSES: SEVERITY.classes,
        RESOLUTION: RESOLUTION.enum,
        RESOLUTION_CLASSES: RESOLUTION.classes,
        ASSIGNEE_TYPE: ASSIGNEE_TYPE.enum,
        AUTHOR_TYPE: AUTHOR_TYPE.enum,
        PROJECT_STATUS: PROJECT_STATUS.enum,
        PROJECT_STATUS_CLASSES: PROJECT_STATUS.classes,
        PROJECT_VISIBILITY: PROJECT_VISIBILITY.enum,
        MILESTONE_STATUS: MILESTONE_STATUS.enum,
        MILESTONE_STATUS_CLASSES: MILESTONE_STATUS.classes,
        SPRINT_STATUS: SPRINT_STATUS.enum,
        SPRINT_STATUS_CLASSES: SPRINT_STATUS.classes,
        TRIAGE_STATUS: TRIAGE_STATUS.enum,
        TRIAGE_STATUS_CLASSES: TRIAGE_STATUS.classes
    };

    L8Tracking.render = {
        bugStatus: createStatusRenderer(BUG_STATUS.enum, BUG_STATUS.classes),
        featureStatus: createStatusRenderer(FEATURE_STATUS.enum, FEATURE_STATUS.classes),
        priority: createStatusRenderer(PRIORITY.enum, PRIORITY.classes),
        severity: createStatusRenderer(SEVERITY.enum, SEVERITY.classes),
        resolution: createStatusRenderer(RESOLUTION.enum, RESOLUTION.classes),
        assigneeType: renderEnum(ASSIGNEE_TYPE.enum),
        projectStatus: createStatusRenderer(PROJECT_STATUS.enum, PROJECT_STATUS.classes),
        projectVisibility: renderEnum(PROJECT_VISIBILITY.enum),
        milestoneStatus: createStatusRenderer(MILESTONE_STATUS.enum, MILESTONE_STATUS.classes),
        sprintStatus: createStatusRenderer(SPRINT_STATUS.enum, SPRINT_STATUS.classes),
        triageStatus: createStatusRenderer(TRIAGE_STATUS.enum, TRIAGE_STATUS.classes),
        date: renderDate
    };
})();
