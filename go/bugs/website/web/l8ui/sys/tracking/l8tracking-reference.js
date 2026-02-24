(function() {
    'use strict';
    const ref = window.Layer8RefFactory;

    Layer8DReferenceRegistry.register({
        ...ref.simple('BugsProject', 'projectId', 'name', 'Project'),
        ...ref.simple('BugsAssignee', 'assigneeId', 'name', 'Assignee'),
        ...ref.simple('Bug', 'bugId', 'title', 'Bug'),
        ...ref.simple('Feature', 'featureId', 'title', 'Feature'),
        ...ref.simple('BugsSprint', 'sprintId', 'name', 'Sprint'),
        ...ref.simple('BugsDigest', 'digestId', 'summary', 'Digest')
    });
})();
