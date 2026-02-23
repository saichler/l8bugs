(function() {
    'use strict';
    Layer8SectionConfigs.register('bugs', {
        title: 'Bug & Feature Tracking',
        subtitle: 'AI-First Issue Tracking for Software Teams',
        icon: '\uD83D\uDC1B',
        svgContent: Layer8SvgFactory.generate('bugs'),
        initFn: 'initializeBugs',
        modules: [
            {
                key: 'tracking', label: 'Tracking', icon: '\uD83D\uDC1B', isDefault: true,
                services: [
                    { key: 'bugs', label: 'Bugs', icon: '\uD83D\uDC1B', isDefault: true },
                    { key: 'features', label: 'Features', icon: '\u2728' },
                    { key: 'projects', label: 'Projects', icon: '\uD83D\uDCC1' }
                ]
            }
        ]
    });
})();
