(function() {
    'use strict';

    if (typeof window.L8Tracking === 'undefined') {
        console.error('Tracking module not properly initialized. Ensure all l8tracking-*.js files are loaded.');
        return;
    }

    var requiredProps = ['columns', 'forms', 'primaryKeys', 'enums'];
    for (var i = 0; i < requiredProps.length; i++) {
        if (!L8Tracking[requiredProps[i]]) {
            console.error('L8Tracking.' + requiredProps[i] + ' not found. Ensure all l8tracking-*.js files are loaded.');
            return;
        }
    }

    console.log('Tracking module initialized');
})();
