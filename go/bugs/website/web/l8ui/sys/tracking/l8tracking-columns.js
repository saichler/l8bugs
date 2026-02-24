(function() {
    'use strict';
    window.L8Tracking = window.L8Tracking || {};

    const col = window.Layer8ColumnFactory;
    const render = L8Tracking.render;

    L8Tracking.columns = {
        Bug: [
            ...col.id('bugId'),
            ...col.col('bugNumber', 'Bug #'),
            ...col.col('title', 'Title'),
            ...col.enum('status', 'Status', null, render.bugStatus),
            ...col.enum('priority', 'Priority', null, render.priority),
            ...col.enum('severity', 'Severity', null, render.severity),
            ...col.col('assigneeId', 'Assignee'),
            ...col.col('component', 'Component'),
            ...col.enum('resolution', 'Resolution', null, render.resolution),
            ...col.enum('triageStatus', 'AI Triage', null, render.triageStatus),
            ...col.col('aiConfidence', 'AI Conf.'),
            ...col.number('estimatedEffort', 'Est. Effort'),
            ...col.number('actualEffort', 'Actual Effort'),
            ...col.number('aiEstimatedEffort', 'AI Est. Effort'),
            ...col.number('aiEffortConfidence', 'AI Effort Conf.'),
            ...col.date('createdDate', 'Created'),
            ...col.date('resolvedDate', 'Resolved')
        ],

        Feature: [
            ...col.id('featureId'),
            ...col.col('featureNumber', 'Feature #'),
            ...col.col('title', 'Title'),
            ...col.enum('status', 'Status', null, render.featureStatus),
            ...col.enum('priority', 'Priority', null, render.priority),
            ...col.col('assigneeId', 'Assignee'),
            ...col.col('component', 'Component'),
            ...col.col('targetVersion', 'Target Version'),
            ...col.enum('triageStatus', 'AI Triage', null, render.triageStatus),
            ...col.col('aiConfidence', 'AI Conf.'),
            ...col.number('estimatedEffort', 'Est. Effort'),
            ...col.number('actualEffort', 'Actual Effort'),
            ...col.number('aiEstimatedEffort', 'AI Est. Effort'),
            ...col.date('createdDate', 'Created')
        ],

        BugsProject: [
            ...col.id('projectId'),
            ...col.col('name', 'Name'),
            ...col.col('key', 'Key'),
            ...col.enum('status', 'Status', null, render.projectStatus),
            ...col.enum('visibility', 'Visibility', null, render.projectVisibility),
            ...col.col('ownerId', 'Owner'),
            ...col.col('repositoryUrl', 'Repository')
        ],

        BugsAssignee: [
            ...col.id('assigneeId'),
            ...col.col('name', 'Name'),
            ...col.col('email', 'Email'),
            ...col.enum('assigneeType', 'Type', null, render.assigneeType),
            ...col.col('projectId', 'Project'),
            ...col.col('active', 'Active')
        ],

        BugsDigest: [
            ...col.id('digestId'),
            ...col.col('projectId', 'Project'),
            ...col.enum('period', 'Period', null, render.digestPeriod),
            ...col.date('startDate', 'Start Date'),
            ...col.date('endDate', 'End Date'),
            ...col.col('summary', 'Summary'),
            ...col.date('generatedDate', 'Generated')
        ],

        BugsSprint: [
            ...col.id('sprintId'),
            ...col.col('name', 'Name'),
            ...col.col('projectId', 'Project'),
            ...col.enum('status', 'Status', null, render.sprintStatus),
            ...col.col('goal', 'Goal'),
            ...col.date('startDate', 'Start Date'),
            ...col.date('endDate', 'End Date'),
            ...col.col('capacity', 'Capacity'),
            ...col.col('completedPoints', 'Completed')
        ]
    };
})();
