const sections = {
    bugs: 'sections/bugs.html',
    system: 'sections/system.html'
};

const sectionInitializers = {
    bugs: () => {
        if (typeof initializeBugs === 'function') {
            initializeBugs();
        }
    },
    system: () => {
        if (typeof initializeL8Sys === 'function') {
            initializeL8Sys();
        }
    }
};

function loadSection(sectionName) {
    const contentArea = document.getElementById('content-area');
    if (!contentArea) return;

    const sectionPath = sections[sectionName];
    if (!sectionPath) {
        contentArea.innerHTML = '<div class="section-placeholder"><h2>Section not found</h2></div>';
        return;
    }

    fetch(sectionPath)
        .then(response => {
            if (!response.ok) throw new Error('Failed to load section');
            return response.text();
        })
        .then(html => {
            contentArea.innerHTML = html;

            // Execute any inline scripts in the loaded HTML
            const scripts = contentArea.querySelectorAll('script');
            scripts.forEach(script => {
                const newScript = document.createElement('script');
                if (script.src) {
                    newScript.src = script.src;
                } else {
                    newScript.textContent = script.textContent;
                }
                document.body.appendChild(newScript);
                document.body.removeChild(newScript);
            });

            // Call the section initializer if it exists
            if (sectionInitializers[sectionName]) {
                sectionInitializers[sectionName]();
            }
        })
        .catch(error => {
            console.error('Error loading section:', error);
            contentArea.innerHTML = '<div class="section-placeholder"><h2>Error loading section</h2></div>';
        });
}
