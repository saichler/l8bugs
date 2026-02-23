// Main application initialization for L8Bugs

function getAuthHeaders() {
    const bearerToken = sessionStorage.getItem('bearerToken');
    return {
        'Authorization': bearerToken ? `Bearer ${bearerToken}` : '',
        'Content-Type': 'application/json'
    };
}

async function makeAuthenticatedRequest(url, options = {}) {
    const bearerToken = sessionStorage.getItem('bearerToken');

    if (!bearerToken) {
        console.error('No bearer token found');
        window.location.href = 'l8ui/login/index.html';
        return;
    }

    const headers = {
        'Authorization': `Bearer ${bearerToken}`,
        'Content-Type': 'application/json',
        ...options.headers
    };

    try {
        const response = await fetch(url, {
            ...options,
            headers: headers
        });

        if (response.status === 401) {
            sessionStorage.removeItem('bearerToken');
            window.location.href = 'l8ui/login/index.html';
            return;
        }

        return response;
    } catch (error) {
        console.error('API request failed:', error);
        throw error;
    }
}

function logout() {
    sessionStorage.removeItem('bearerToken');
    localStorage.removeItem('bearerToken');
    localStorage.removeItem('rememberedUser');
    window.location.href = 'l8ui/login/index.html';
}

document.addEventListener('DOMContentLoaded', async function() {
    if (typeof Layer8DConfig !== 'undefined') {
        await Layer8DConfig.load();
    }

    const bearerToken = sessionStorage.getItem('bearerToken');
    if (!bearerToken) {
        window.location.href = 'l8ui/login/index.html';
        return;
    }

    localStorage.setItem('bearerToken', bearerToken);
    window.bearerToken = bearerToken;

    const username = sessionStorage.getItem('currentUser') || 'Admin';
    document.querySelector('.username').textContent = username;

    // Module filter skipped — L8Bugs has a single module, no server-side ModConfig service
    // If Layer8DModuleFilter is present, mark it as loaded so downstream code doesn't block
    if (typeof Layer8DModuleFilter !== 'undefined') {
        Layer8DModuleFilter._loaded = true;
    }

    loadSection('bugs');

    const navLinks = document.querySelectorAll('.nav-link');
    navLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            navLinks.forEach(l => l.classList.remove('active'));
            this.classList.add('active');
            const section = this.getAttribute('data-section');
            loadSection(section);
        });
    });
});
