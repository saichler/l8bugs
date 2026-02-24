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
 * L8Dashboard - KPI Widgets and Charts for Bug Tracking
 * Fetches data from Bug/Feature/Sprint endpoints and computes metrics client-side.
 * Lives as a sub-nav item within the Tracking module.
 */
(function() {
    'use strict';

    const ENDPOINTS = {
        Bug: '/20/Bug',
        Feature: '/20/Feature',
        Sprint: '/20/Sprint'
    };

    let _charts = [];
    let _widgetContainerId = 'dashboard-widgets';
    let _chartContainerId = 'dashboard-charts';

    function _resolveEndpoint(path) {
        return typeof Layer8DConfig !== 'undefined'
            ? Layer8DConfig.resolveEndpoint(path) : '/bugs' + path;
    }

    function _getHeaders() {
        return typeof getAuthHeaders === 'function'
            ? getAuthHeaders() : { 'Content-Type': 'application/json' };
    }

    async function _fetchAll(model, endpoint) {
        const query = 'select * from ' + model + ' limit 500 page 0';
        const body = encodeURIComponent(JSON.stringify({ text: query }));
        try {
            const resp = await fetch(_resolveEndpoint(endpoint) + '?body=' + body, {
                method: 'GET', headers: _getHeaders()
            });
            if (!resp.ok) return [];
            const data = await resp.json();
            return data.list || [];
        } catch (e) {
            console.error('Dashboard fetch error (' + model + '):', e);
            return [];
        }
    }

    function _computeKPIs(bugs, features, sprints) {
        const now = Date.now() / 1000;
        const weekAgo = now - 7 * 86400;

        const openBugs = bugs.filter(b => b.status >= 1 && b.status <= 4).length;
        const openFeatures = features.filter(f => f.status >= 1 && f.status <= 5).length;

        const resolvedThisWeek = bugs.filter(b => {
            const rd = b.resolvedDate;
            return rd && rd > weekAgo;
        }).length;

        const triaged = bugs.filter(b => b.triageStatus === 3);
        const highConf = triaged.filter(b => b.aiConfidence > 70).length;
        const triageAccuracy = triaged.length > 0
            ? Math.round((highConf / triaged.length) * 100) : 0;

        const activeSprints = sprints.filter(s => s.status === 2).length;

        const overdue = bugs.filter(b => b.dueDate && b.dueDate < now && b.status < 5).length
            + features.filter(f => f.dueDate && f.dueDate < now && f.status < 6).length;

        return [
            { label: 'Open Bugs', value: openBugs, icon: 'bug' },
            { label: 'Open Features', value: openFeatures, icon: 'feature' },
            { label: 'Resolved This Week', value: resolvedThisWeek, icon: 'resolved' },
            { label: 'AI Triage Accuracy', value: triageAccuracy, icon: 'ai', suffix: '%' },
            { label: 'Active Sprints', value: activeSprints, icon: 'sprint' },
            { label: 'Overdue Items', value: overdue, icon: 'overdue' }
        ];
    }

    function _renderWidgets(kpis, containerId) {
        const container = document.getElementById(containerId);
        if (!container) return;

        const ICONS = {
            bug: '<svg viewBox="0 0 20 20" width="24" height="24"><circle cx="10" cy="10" r="8" fill="none" stroke="currentColor" stroke-width="1.5"/><path d="M7 8h6M7 12h6M10 5v10" stroke="currentColor" stroke-width="1.2" fill="none"/></svg>',
            feature: '<svg viewBox="0 0 20 20" width="24" height="24"><polygon points="10,2 12.5,7.5 18,8 14,12 15,18 10,15 5,18 6,12 2,8 7.5,7.5" fill="none" stroke="currentColor" stroke-width="1.2"/></svg>',
            resolved: '<svg viewBox="0 0 20 20" width="24" height="24"><circle cx="10" cy="10" r="8" fill="none" stroke="currentColor" stroke-width="1.5"/><path d="M6 10l3 3 5-6" fill="none" stroke="currentColor" stroke-width="1.5"/></svg>',
            ai: '<svg viewBox="0 0 20 20" width="24" height="24"><rect x="4" y="4" width="12" height="12" rx="3" fill="none" stroke="currentColor" stroke-width="1.2"/><circle cx="8" cy="9" r="1.5" fill="currentColor"/><circle cx="12" cy="9" r="1.5" fill="currentColor"/><path d="M8 13h4" stroke="currentColor" stroke-width="1.2" fill="none"/></svg>',
            sprint: '<svg viewBox="0 0 20 20" width="24" height="24"><path d="M4 16 L8 4 L12 12 L16 6" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/></svg>',
            overdue: '<svg viewBox="0 0 20 20" width="24" height="24"><circle cx="10" cy="10" r="8" fill="none" stroke="currentColor" stroke-width="1.5"/><path d="M10 6v5l3 2" fill="none" stroke="currentColor" stroke-width="1.5"/></svg>'
        };

        container.innerHTML = kpis.map(function(kpi) {
            return Layer8DWidget.render(
                { label: kpi.label, iconSvg: ICONS[kpi.icon] || '' },
                kpi.value,
                {}
            );
        }).join('');
    }

    function _destroyCharts() {
        _charts.forEach(function(c) { c.destroy(); });
        _charts = [];
    }

    function _renderCharts(bugs, features, containerId) {
        const container = document.getElementById(containerId);
        if (!container) return;

        _destroyCharts();
        container.innerHTML = '';

        _addChart(container, 'chart-bugs-priority', 'Bugs by Priority', bugs, {
            chartType: 'pie', categoryField: 'priority', aggregation: 'count'
        });

        _addChart(container, 'chart-bugs-status', 'Bugs by Status', bugs, {
            chartType: 'bar', categoryField: 'status', aggregation: 'count'
        });

        _addChart(container, 'chart-features-status', 'Features by Status', features, {
            chartType: 'bar', categoryField: 'status', aggregation: 'count'
        });

        var bugsWithComponent = bugs.filter(function(b) { return b.component; });
        _addChart(container, 'chart-top-components', 'Top Components', bugsWithComponent, {
            chartType: 'bar', categoryField: 'component', aggregation: 'count'
        });
    }

    function _addChart(parent, id, title, data, viewConfig) {
        var card = document.createElement('div');
        card.className = 'l8-dashboard-chart-card';

        var titleEl = document.createElement('div');
        titleEl.className = 'l8-dashboard-chart-title';
        titleEl.textContent = title;
        card.appendChild(titleEl);

        var chartDiv = document.createElement('div');
        chartDiv.id = id;
        chartDiv.style.height = '260px';
        chartDiv.style.overflow = 'hidden';
        card.appendChild(chartDiv);

        parent.appendChild(card);

        var chart = new Layer8DChart({
            containerId: id,
            columns: [],
            viewConfig: viewConfig
        });
        chart.init();
        chart.setData(data, data.length);

        // Disconnect resize observer to prevent render loops between charts
        if (chart._resizeObserver) {
            chart._resizeObserver.disconnect();
            chart._resizeObserver = null;
        }

        // Hide chart type controls — dashboard charts are fixed type
        var controls = chartDiv.querySelector('.layer8d-chart-controls');
        if (controls) controls.style.display = 'none';

        _charts.push(chart);
    }

    window.L8Dashboard = {
        initialize: function(widgetId, chartId) {
            _widgetContainerId = widgetId || 'dashboard-widgets';
            _chartContainerId = chartId || 'dashboard-charts';

            var wc = document.getElementById(_widgetContainerId);
            if (wc) wc.innerHTML = '<div class="l8-dashboard-loading">Loading dashboard...</div>';

            this.refresh();
        },

        refresh: async function() {
            var bugs = await _fetchAll('Bug', ENDPOINTS.Bug);
            var features = await _fetchAll('Feature', ENDPOINTS.Feature);
            var sprints = await _fetchAll('BugsSprint', ENDPOINTS.Sprint);

            var kpis = _computeKPIs(bugs, features, sprints);
            _renderWidgets(kpis, _widgetContainerId);
            _renderCharts(bugs, features, _chartContainerId);
        }
    };

})();
