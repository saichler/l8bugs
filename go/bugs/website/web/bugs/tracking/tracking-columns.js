(function() {
    'use strict';
    window.BugsTracking = window.BugsTracking || {};

    const col = window.Layer8ColumnFactory;
    const render = BugsTracking.render;

    BugsTracking.columns = {
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
            ...col.date('createdDate', 'Created')
        ],

        BugsProject: [
            ...col.id('projectId'),
            ...col.col('name', 'Name'),
            ...col.col('key', 'Key'),
            ...col.enum('status', 'Status', null, render.projectStatus),
            ...col.enum('visibility', 'Visibility', null, render.projectVisibility),
            ...col.col('ownerId', 'Owner')
        ]
    };
})();
