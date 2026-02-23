(function() {
    'use strict';
    Layer8ModuleConfigFactory.create({
        namespace: 'Bugs',
        modules: {
            'tracking': {
                label: 'Tracking',
                icon: '\uD83D\uDC1B',
                services: [
                    { key: 'bugs', label: 'Bugs', icon: '\uD83D\uDC1B',
                      endpoint: '/20/Bug', model: 'Bug',
                      supportedViews: ['table', 'kanban'] },
                    { key: 'features', label: 'Features', icon: '\u2728',
                      endpoint: '/20/Feature', model: 'Feature',
                      supportedViews: ['table', 'kanban'] },
                    { key: 'projects', label: 'Projects', icon: '\uD83D\uDCC1',
                      endpoint: '/20/Project', model: 'BugsProject' }
                ]
            }
        },
        submodules: ['Tracking']
    });
})();
