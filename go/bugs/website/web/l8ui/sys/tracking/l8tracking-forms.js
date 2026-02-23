(function() {
    'use strict';
    window.L8Tracking = window.L8Tracking || {};

    const f = window.Layer8FormFactory;
    const enums = L8Tracking.enums;

    L8Tracking.forms = {
        Bug: f.form('Bug', [
            f.section('Bug Details', [
                ...f.text('bugNumber', 'Bug #'),
                ...f.text('title', 'Title', true),
                ...f.reference('projectId', 'Project', 'BugsProject', true),
                ...f.select('status', 'Status', enums.BUG_STATUS),
                ...f.select('priority', 'Priority', enums.PRIORITY),
                ...f.select('severity', 'Severity', enums.SEVERITY)
            ]),
            f.section('Description', [
                ...f.textarea('description', 'Description'),
                ...f.textarea('stepsToReproduce', 'Steps to Reproduce'),
                ...f.textarea('expectedBehavior', 'Expected Behavior'),
                ...f.textarea('actualBehavior', 'Actual Behavior')
            ]),
            f.section('Assignment', [
                ...f.reference('assigneeId', 'Assignee', 'BugsAssignee'),
                ...f.select('assigneeType', 'Assignee Type', enums.ASSIGNEE_TYPE),
                ...f.text('reporterId', 'Reporter'),
                ...f.text('component', 'Component'),
                ...f.text('labels', 'Labels')
            ]),
            f.section('Technical', [
                ...f.text('environment', 'Environment'),
                ...f.textarea('stackTrace', 'Stack Trace'),
                ...f.text('affectedVersion', 'Affected Version'),
                ...f.text('fixVersion', 'Fix Version')
            ]),
            f.section('AI Analysis', [
                ...f.select('triageStatus', 'Triage Status', enums.TRIAGE_STATUS),
                ...f.number('aiConfidence', 'AI Confidence'),
                ...f.select('aiSuggestedPriority', 'AI Suggested Priority', enums.PRIORITY),
                ...f.select('aiSuggestedSeverity', 'AI Suggested Severity', enums.SEVERITY),
                ...f.text('aiSuggestedComponent', 'AI Suggested Component'),
                ...f.reference('aiSuggestedAssigneeId', 'AI Suggested Assignee', 'BugsAssignee'),
                ...f.textarea('aiRootCause', 'AI Root Cause'),
                ...f.text('triageError', 'Triage Error')
            ]),
            f.section('Resolution', [
                ...f.select('resolution', 'Resolution', enums.RESOLUTION),
                ...f.date('resolvedDate', 'Resolved Date'),
                ...f.text('linkedPrUrl', 'Linked PR URL'),
                ...f.text('linkedBranch', 'Linked Branch'),
                ...f.reference('duplicateOfId', 'Duplicate Of', 'Bug'),
                ...f.reference('parentBugId', 'Parent Bug', 'Bug')
            ]),
            f.section('Tracking', [
                ...f.date('dueDate', 'Due Date'),
                ...f.number('estimatedEffort', 'Estimated Effort'),
                ...f.number('voteCount', 'Votes'),
                ...f.number('watcherCount', 'Watchers')
            ]),
            f.section('Comments', [
                ...f.inlineTable('comments', 'Comments', [
                    { key: 'commentId', label: 'ID', hidden: true },
                    { key: 'authorId', label: 'Author', type: 'text' },
                    { key: 'authorType', label: 'Type', type: 'select', options: enums.AUTHOR_TYPE },
                    { key: 'body', label: 'Comment', type: 'textarea' },
                    { key: 'isInternal', label: 'Internal', type: 'checkbox' },
                    { key: 'createdDate', label: 'Created', type: 'date' }
                ])
            ]),
            f.section('Attachments', [
                ...f.inlineTable('attachments', 'Attachments', [
                    { key: 'attachmentId', label: 'ID', hidden: true },
                    { key: 'filename', label: 'Filename', type: 'text' },
                    { key: 'contentType', label: 'Type', type: 'text' },
                    { key: 'size', label: 'Size', type: 'number' },
                    { key: 'url', label: 'URL', type: 'text' }
                ])
            ]),
            f.section('Activity Log', [
                ...f.inlineTable('activity', 'Activity', [
                    { key: 'entryId', label: 'ID', hidden: true },
                    { key: 'actorId', label: 'Actor', type: 'text' },
                    { key: 'actorType', label: 'Type', type: 'select', options: enums.AUTHOR_TYPE },
                    { key: 'action', label: 'Action', type: 'text' },
                    { key: 'fieldName', label: 'Field', type: 'text' },
                    { key: 'oldValue', label: 'Old Value', type: 'text' },
                    { key: 'newValue', label: 'New Value', type: 'text' },
                    { key: 'timestamp', label: 'When', type: 'date' }
                ])
            ])
        ]),

        Feature: f.form('Feature', [
            f.section('Feature Details', [
                ...f.text('featureNumber', 'Feature #'),
                ...f.text('title', 'Title', true),
                ...f.reference('projectId', 'Project', 'BugsProject', true),
                ...f.select('status', 'Status', enums.FEATURE_STATUS),
                ...f.select('priority', 'Priority', enums.PRIORITY)
            ]),
            f.section('Description', [
                ...f.textarea('description', 'Description'),
                ...f.textarea('userStory', 'User Story'),
                ...f.textarea('acceptanceCriteria', 'Acceptance Criteria')
            ]),
            f.section('Assignment', [
                ...f.reference('assigneeId', 'Assignee', 'BugsAssignee'),
                ...f.select('assigneeType', 'Assignee Type', enums.ASSIGNEE_TYPE),
                ...f.text('reporterId', 'Reporter'),
                ...f.text('component', 'Component'),
                ...f.text('labels', 'Labels')
            ]),
            f.section('Planning', [
                ...f.text('targetVersion', 'Target Version'),
                ...f.date('dueDate', 'Due Date'),
                ...f.number('estimatedEffort', 'Estimated Effort'),
                ...f.number('voteCount', 'Votes'),
                ...f.number('watcherCount', 'Watchers')
            ]),
            f.section('AI Analysis', [
                ...f.select('triageStatus', 'Triage Status', enums.TRIAGE_STATUS),
                ...f.number('aiConfidence', 'AI Confidence'),
                ...f.select('aiSuggestedPriority', 'AI Suggested Priority', enums.PRIORITY),
                ...f.text('aiSuggestedComponent', 'AI Suggested Component'),
                ...f.reference('aiSuggestedAssigneeId', 'AI Suggested Assignee', 'BugsAssignee'),
                ...f.textarea('aiBreakdown', 'AI Breakdown'),
                ...f.text('triageError', 'Triage Error')
            ]),
            f.section('Links', [
                ...f.text('linkedPrUrl', 'Linked PR URL'),
                ...f.text('linkedBranch', 'Linked Branch'),
                ...f.reference('parentFeatureId', 'Parent Feature', 'Feature')
            ]),
            f.section('Comments', [
                ...f.inlineTable('comments', 'Comments', [
                    { key: 'commentId', label: 'ID', hidden: true },
                    { key: 'authorId', label: 'Author', type: 'text' },
                    { key: 'authorType', label: 'Type', type: 'select', options: enums.AUTHOR_TYPE },
                    { key: 'body', label: 'Comment', type: 'textarea' },
                    { key: 'isInternal', label: 'Internal', type: 'checkbox' },
                    { key: 'createdDate', label: 'Created', type: 'date' }
                ])
            ]),
            f.section('Attachments', [
                ...f.inlineTable('attachments', 'Attachments', [
                    { key: 'attachmentId', label: 'ID', hidden: true },
                    { key: 'filename', label: 'Filename', type: 'text' },
                    { key: 'contentType', label: 'Type', type: 'text' },
                    { key: 'size', label: 'Size', type: 'number' },
                    { key: 'url', label: 'URL', type: 'text' }
                ])
            ]),
            f.section('Activity Log', [
                ...f.inlineTable('activity', 'Activity', [
                    { key: 'entryId', label: 'ID', hidden: true },
                    { key: 'actorId', label: 'Actor', type: 'text' },
                    { key: 'actorType', label: 'Type', type: 'select', options: enums.AUTHOR_TYPE },
                    { key: 'action', label: 'Action', type: 'text' },
                    { key: 'fieldName', label: 'Field', type: 'text' },
                    { key: 'oldValue', label: 'Old Value', type: 'text' },
                    { key: 'newValue', label: 'New Value', type: 'text' },
                    { key: 'timestamp', label: 'When', type: 'date' }
                ])
            ])
        ]),

        BugsAssignee: f.form('Assignee', [
            f.section('Assignee Details', [
                ...f.text('name', 'Name', true),
                ...f.text('email', 'Email'),
                ...f.select('assigneeType', 'Type', enums.ASSIGNEE_TYPE),
                ...f.reference('projectId', 'Project', 'BugsProject'),
                ...f.checkbox('active', 'Active')
            ])
        ]),

        BugsProject: f.form('Project', [
            f.section('Project Details', [
                ...f.text('name', 'Name', true),
                ...f.text('key', 'Key', true),
                ...f.textarea('description', 'Description'),
                ...f.text('ownerId', 'Owner'),
                ...f.select('status', 'Status', enums.PROJECT_STATUS),
                ...f.select('visibility', 'Visibility', enums.PROJECT_VISIBILITY),
                ...f.text('defaultAssigneeId', 'Default Assignee')
            ]),
            f.section('Labels', [
                ...f.inlineTable('labels', 'Labels', [
                    { key: 'labelId', label: 'ID', hidden: true },
                    { key: 'name', label: 'Name', type: 'text', required: true },
                    { key: 'color', label: 'Color', type: 'text' },
                    { key: 'description', label: 'Description', type: 'text' }
                ])
            ]),
            f.section('Components', [
                ...f.inlineTable('components', 'Components', [
                    { key: 'componentId', label: 'ID', hidden: true },
                    { key: 'name', label: 'Name', type: 'text', required: true },
                    { key: 'description', label: 'Description', type: 'text' },
                    { key: 'leadId', label: 'Lead', type: 'text' },
                    { key: 'defaultAssigneeId', label: 'Default Assignee', type: 'text' }
                ])
            ]),
            f.section('Milestones', [
                ...f.inlineTable('milestones', 'Milestones', [
                    { key: 'milestoneId', label: 'ID', hidden: true },
                    { key: 'name', label: 'Name', type: 'text', required: true },
                    { key: 'description', label: 'Description', type: 'text' },
                    { key: 'status', label: 'Status', type: 'select', options: enums.MILESTONE_STATUS },
                    { key: 'dueDate', label: 'Due Date', type: 'date' },
                    { key: 'completionPercentage', label: 'Completion %', type: 'number' }
                ])
            ])
        ]),

        BugsSprint: f.form('Sprint', [
            f.section('Sprint Details', [
                ...f.text('name', 'Name', true),
                ...f.reference('projectId', 'Project', 'BugsProject', true),
                ...f.select('status', 'Status', enums.SPRINT_STATUS),
                ...f.textarea('goal', 'Goal')
            ]),
            f.section('Schedule', [
                ...f.date('startDate', 'Start Date'),
                ...f.date('endDate', 'End Date'),
                ...f.number('capacity', 'Capacity (Story Points)'),
                ...f.number('completedPoints', 'Completed Points')
            ])
        ])
    };

    L8Tracking.primaryKeys = {
        Bug: 'bugId',
        Feature: 'featureId',
        BugsProject: 'projectId',
        BugsAssignee: 'assigneeId',
        BugsSprint: 'sprintId'
    };
})();
