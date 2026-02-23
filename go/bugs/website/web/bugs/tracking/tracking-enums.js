(function() {
    'use strict';
    const factory = window.Layer8EnumFactory;
    const { createStatusRenderer, renderEnum, renderDate } = Layer8DRenderers;

    window.BugsTracking = window.BugsTracking || {};

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

    BugsTracking.enums = {
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
        MILESTONE_STATUS_CLASSES: MILESTONE_STATUS.classes
    };

    BugsTracking.render = {
        bugStatus: createStatusRenderer(BUG_STATUS.enum, BUG_STATUS.classes),
        featureStatus: createStatusRenderer(FEATURE_STATUS.enum, FEATURE_STATUS.classes),
        priority: createStatusRenderer(PRIORITY.enum, PRIORITY.classes),
        severity: createStatusRenderer(SEVERITY.enum, SEVERITY.classes),
        resolution: createStatusRenderer(RESOLUTION.enum, RESOLUTION.classes),
        assigneeType: renderEnum(ASSIGNEE_TYPE.enum),
        projectStatus: createStatusRenderer(PROJECT_STATUS.enum, PROJECT_STATUS.classes),
        projectVisibility: renderEnum(PROJECT_VISIBILITY.enum),
        milestoneStatus: createStatusRenderer(MILESTONE_STATUS.enum, MILESTONE_STATUS.classes),
        date: renderDate
    };
})();
