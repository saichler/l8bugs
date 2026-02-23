/*
© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
/**
 * Mobile L8Bugs App Core - Navigation and initialization
 */
(function() {
    'use strict';

    const SECTIONS = {
        'system': 'sections/system.html'
    };

    let currentSection = 'system';
    let sectionCache = {};

    window.MobileApp = {
        async init() {
            if (!Layer8MAuth.requireAuth()) return;

            await Layer8MConfig.load();
            this.updateUserInfo();

            const token = Layer8MAuth.getBearerToken();

            // Module filter skipped — L8Bugs has a single module, no server-side ModConfig service
            if (typeof Layer8DModuleFilter !== 'undefined') {
                Layer8DModuleFilter._loaded = true;
            }

            this.initSidebar();

            document.getElementById('refresh-btn')?.addEventListener('click', () => {
                this.loadSection(currentSection, true);
            });

            const hash = window.location.hash.slice(1);
            const section = SECTIONS[hash] ? hash : 'system';
            await this.loadSection(section);

            window.addEventListener('hashchange', () => {
                const newSection = window.location.hash.slice(1);
                if (SECTIONS[newSection] && newSection !== currentSection) {
                    this.loadSection(newSection);
                }
            });
        },

        updateUserInfo() {
            const username = Layer8MAuth.getUsername();
            const initial = username.charAt(0).toUpperCase();

            document.getElementById('user-name').textContent = username;
            document.getElementById('user-avatar').textContent = initial;
        },

        initSidebar() {
            const menuToggle = document.getElementById('menu-toggle');
            const overlay = document.getElementById('sidebar-overlay');

            menuToggle?.addEventListener('click', () => this.openSidebar());
            overlay?.addEventListener('click', () => this.closeSidebar());

            document.querySelectorAll('.sidebar-item[data-section]').forEach(item => {
                item.addEventListener('click', async (e) => {
                    e.preventDefault();
                    const section = item.dataset.section;
                    this.closeSidebar();
                    await this.loadSection(section);
                });
            });
        },

        openSidebar() {
            document.getElementById('sidebar')?.classList.add('open');
            document.getElementById('sidebar-overlay')?.classList.add('visible');
            document.body.style.overflow = 'hidden';
        },

        closeSidebar() {
            document.getElementById('sidebar')?.classList.remove('open');
            document.getElementById('sidebar-overlay')?.classList.remove('visible');
            document.body.style.overflow = '';
        },

        async loadSection(sectionKey, forceReload) {
            const url = SECTIONS[sectionKey];
            if (!url) return;

            currentSection = sectionKey;
            window.location.hash = sectionKey;

            document.querySelectorAll('.sidebar-item').forEach(item => {
                item.classList.toggle('active', item.dataset.section === sectionKey);
            });

            const contentArea = document.getElementById('content-area');

            if (!forceReload && sectionCache[sectionKey]) {
                contentArea.innerHTML = sectionCache[sectionKey];
                this.executeSectionScripts(contentArea, sectionKey);
                return;
            }

            contentArea.innerHTML = '<div class="loading">Loading</div>';

            try {
                const token = Layer8MAuth.getBearerToken();
                const response = await fetch(url, {
                    headers: { 'Authorization': 'Bearer ' + token }
                });
                if (!response.ok) throw new Error('Failed to load section');
                const html = await response.text();

                sectionCache[sectionKey] = html;
                contentArea.innerHTML = html;
                this.executeSectionScripts(contentArea, sectionKey);
            } catch (error) {
                console.error('Failed to load section:', error);
                contentArea.innerHTML = '<div class="empty-state"><div class="empty-state-icon">&#x26A0;&#xFE0F;</div><div class="empty-state-title">Failed to load</div></div>';
            }
        },

        executeSectionScripts(container, sectionKey) {
            const scripts = container.querySelectorAll('script');
            scripts.forEach(script => {
                const newScript = document.createElement('script');
                newScript.textContent = script.textContent;
                script.parentNode.replaceChild(newScript, script);
            });

            if (sectionKey === 'system') {
                if (typeof initMobileSystem === 'function') {
                    initMobileSystem();
                }
            }
        },

        applyModuleFilter() {
            if (typeof Layer8DModuleFilter === 'undefined') return;
            document.querySelectorAll('.sidebar-item[data-section]').forEach(item => {
                const section = item.dataset.section;
                if (section === 'system') return;
                if (!Layer8DModuleFilter.isModuleEnabled(section)) {
                    item.style.display = 'none';
                }
            });
        },

        logout() {
            sessionStorage.removeItem('bearerToken');
            localStorage.removeItem('bearerToken');
            localStorage.removeItem('rememberedUser');
            window.location.href = '../l8ui/login/index.html';
        }
    };

    document.addEventListener('DOMContentLoaded', () => MobileApp.init());
})();
