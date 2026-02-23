(function() {
    'use strict';
    const ref = window.Layer8RefFactory;

    Layer8DReferenceRegistry.register({
        ...ref.simple('BugsProject', 'projectId', 'name', 'Project'),
        ...ref.simple('Bug', 'bugId', 'title', 'Bug'),
        ...ref.simple('Feature', 'featureId', 'title', 'Feature')
    });
})();
