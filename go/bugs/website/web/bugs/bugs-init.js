(function() {
    'use strict';
    Layer8DModuleFactory.create({
        namespace: 'Bugs',
        defaultModule: 'tracking',
        defaultService: 'bugs',
        sectionSelector: 'tracking',
        initializerName: 'initializeBugs',
        requiredNamespaces: ['Tracking']
    });
})();
