# Migrate Bug Tracking into System Section as Generic l8ui Component

## Context

L8Bugs Phase 1 is complete — the Bug, Feature, and BugsProject services work end-to-end. Currently, the app has two sidebar sections: "Bug Tracking" (standalone module) and "System". The user wants to:

1. **Cancel the standalone Bug Tracking module** — remove it as a separate section
2. **Add Tracking as a sub-module of the System section** — alongside Health, Security, Modules, Logs
3. **Make it a generic l8ui component** — files live under `l8ui/sys/tracking/`, reusable by any Layer8 app

This follows the pattern of `l8ui/sys/security/` which has enums, columns, forms, and a verification entry point under the shared l8ui library.

---

## Step 1: Create Desktop Tracking Component (`l8ui/sys/tracking/`)

Create 5 files, all following the `l8ui/sys/security/` pattern:

### 1.1 `l8ui/sys/tracking/l8tracking-enums.js`
- Copy from `bugs/tracking/tracking-enums.js`
- Replace `window.BugsTracking` → `window.L8Tracking` throughout
- Keep all enum definitions and renderers identical

### 1.2 `l8ui/sys/tracking/l8tracking-columns.js`
- Copy from `bugs/tracking/tracking-columns.js`
- Replace `BugsTracking` → `L8Tracking` throughout

### 1.3 `l8ui/sys/tracking/l8tracking-forms.js`
- Copy from `bugs/tracking/tracking-forms.js`
- Replace `BugsTracking` → `L8Tracking` throughout

### 1.4 `l8ui/sys/tracking/l8tracking.js`
- Verification entry point (same pattern as `l8security.js`)
- Checks `window.L8Tracking` has `columns`, `forms`, `primaryKeys`, `enums`

### 1.5 `l8ui/sys/tracking/l8tracking-reference.js`
- Move from `js/reference-registry-bugs.js`
- Registers BugsProject, Bug, Feature with `Layer8DReferenceRegistry`

---

## Step 2: Modify System Config and Init (Desktop)

### 2.1 `l8ui/sys/l8sys-config.js`
- Add `'tracking'` as **first** entry in `L8Sys.modules` (makes it the default tab):
  ```
  tracking: { label:'Tracking', icon:'🐛', services: [
      { key:'bugs', endpoint:'/20/Bug', model:'Bug', supportedViews:['table','kanban'] },
      { key:'features', endpoint:'/20/Feature', model:'Feature', supportedViews:['table','kanban'] },
      { key:'projects', endpoint:'/20/Project', model:'BugsProject' }
  ]}
  ```
- Add `'L8Tracking'` to `L8Sys.submodules`: `['L8Security', 'L8Tracking']`

### 2.2 `l8ui/sys/l8sys-init.js`
- Change `defaultModule: 'health'` → `'tracking'`
- Change `defaultService: 'users'` → `'bugs'`
- Change `sectionSelector: 'health'` → `'tracking'`
- Add `'L8Tracking'` to `requiredNamespaces`
- **Fix existing bug**: CRUD fallthrough uses `SYS` (undefined) — change to `L8Sys` on 3 lines

---

## Step 3: Update `sections/system.html` (Desktop)

- Add **Tracking tab** as first tab (with `active` class), remove `active` from Health
- Add **Tracking module content** as first content div (with `active` class):
  - Subnav: Bugs (active), Features, Projects
  - Service views with container IDs: `tracking-bugs-table-container`, `tracking-features-table-container`, `tracking-projects-table-container`
  - View switcher slots: `tracking-bugs-view-switcher`, etc.
- Remove `active` from Health module content
- Update subtitle to mention Tracking

---

## Step 4: Update `app.html` (Desktop)

### Remove:
- CSS: `bugs/bugs.css` (line 34)
- Script: `js/reference-registry-bugs.js` (line 146)
- Scripts: entire BUGS Module block (lines 223-229)
- Sidebar: "Bug Tracking" nav item (lines 109-112)

### Add:
- Tracking scripts in SYS Module block (after security, before modules):
  ```
  l8ui/sys/tracking/l8tracking-enums.js
  l8ui/sys/tracking/l8tracking-columns.js
  l8ui/sys/tracking/l8tracking-forms.js
  l8ui/sys/tracking/l8tracking.js
  l8ui/sys/tracking/l8tracking-reference.js
  ```
- Make System nav item `active` (the only sidebar item)

---

## Step 5: Update `js/sections.js` and `js/app.js` (Desktop)

### `js/sections.js`
- Remove `bugs` from `sections` object and `sectionInitializers`
- Only `system` remains

### `js/app.js`
- Change `loadSection('bugs')` → `loadSection('system')`

---

## Step 6: Create Mobile Tracking Files (`m/js/tracking/`)

### 6.1 `m/js/tracking/l8tracking-enums.js`
- Copy from `m/js/bugs/tracking-enums.js`
- Replace `MobileBugsTracking` → `MobileL8Tracking`

### 6.2 `m/js/tracking/l8tracking-columns.js`
- Copy from `m/js/bugs/tracking-columns.js`
- Replace `MobileBugsTracking` → `MobileL8Tracking`

### 6.3 `m/js/tracking/l8tracking-forms.js`
- Copy from `m/js/bugs/tracking-forms.js`
- Replace `MobileBugsTracking` → `MobileL8Tracking`

### 6.4 `m/js/tracking/layer8m-reference-registry-tracking.js`
- Copy from `m/js/bugs/layer8m-reference-registry-bugs.js`
- Same reference entries (BugsProject, Bug, Feature)

---

## Step 7: Update Mobile System Section (`m/sections/system.html`)

- Add **Tracking tab** as first tab (active), shift Health/Security/Modules to non-active
- Add `system-tracking-content` div with:
  - Service sub-tabs: Bugs (active), Features, Projects
  - `system-tracking-table` container
- Add `TRACKING_SERVICES` config using `MobileL8Tracking.columns`
- Add `loadTrackingService()` function (creates `Layer8MEditTable`)
- Update `switchSystemModule()` to handle `'tracking'` case
- Add tracking service tab event listeners in `initMobileSystem()`
- Default to tracking on init: call `loadTrackingService('bugs')` first

---

## Step 8: Update Mobile `m/app.html`

### Remove:
- Sidebar "Bug Tracking" item (lines 94-108)
- Script: `js/bugs/layer8m-reference-registry-bugs.js`
- Scripts: `js/bugs/tracking-enums.js`, `tracking-columns.js`, `tracking-forms.js`, `bugs-index.js`

### Add:
- Tracking scripts: `js/tracking/l8tracking-enums.js`, `l8tracking-columns.js`, `l8tracking-forms.js`
- Reference: `js/tracking/layer8m-reference-registry-tracking.js`
- Make System sidebar item `active`

---

## Step 9: Update Mobile `m/js/app-core.js`

- Remove `'bugs'` from `SECTIONS` — only `'system'` remains
- Change `currentSection = 'bugs'` → `'system'`
- Change default hash fallback to `'system'`
- Remove `if (sectionKey === 'bugs')` branch from `executeSectionScripts()`

---

## Step 10: Update Mobile Config Files

### `m/js/mobile-config-bugs.js`
- Remove `Layer8MConfig.registerModules()` call (tracking is now inline in system.html)
- Keep `Layer8MConfig.registerReferences()` call (still needed for reference picker)

### `m/js/layer8m-nav-config-bugs.js`
- Remove `bugs` module — keep only `system` in `modules` array
- Remove `bugs` submodule/service config

---

## Step 11: Delete Obsolete Files

### Desktop (9 files + 2 directories):
- `web/bugs/bugs-section-config.js`
- `web/bugs/bugs-config.js`
- `web/bugs/bugs-init.js`
- `web/bugs/bugs.css`
- `web/bugs/tracking/tracking-enums.js`
- `web/bugs/tracking/tracking-columns.js`
- `web/bugs/tracking/tracking-forms.js`
- `web/sections/bugs.html`
- `web/js/reference-registry-bugs.js`

### Mobile (6 files + 1 directory):
- `m/js/bugs/tracking-enums.js`
- `m/js/bugs/tracking-columns.js`
- `m/js/bugs/tracking-forms.js`
- `m/js/bugs/bugs-index.js`
- `m/js/bugs/layer8m-reference-registry-bugs.js`
- `m/sections/bugs.html`

---

## Step 12: Update `l8ui/GUIDE.md`

- Remove `bugs/bugs.css` from CSS loading order
- Remove all `bugs/*` from JS loading order
- Add `l8ui/sys/tracking/l8tracking-*.js` entries after security scripts
- Add "Tracking Sub-module" documentation section

---

## Critical Implementation Notes

1. **Container ID convention**: Desktop `Layer8DServiceRegistry.initializeServiceTable()` generates IDs as `${moduleKey}-${serviceKey}-table-container`. HTML must use `tracking-bugs-table-container` (not `bugs-table-container`).

2. **`SYS` bug in l8sys-init.js**: Lines 51, 64, 77 use `origOpenAdd.call(SYS, service)` where `SYS` is undefined. Previously harmless because the else branches were never reached (only Security models existed). Now tracking models WILL hit these branches. Must fix to `L8Sys`.

3. **sectionSelector == defaultModule**: Both must be `'tracking'`, and system.html must have `<div class="l8-module-content active" data-module="tracking">`.

---

## Verification

```bash
# Desktop: no bugs references in app.html
grep -c "bugs/" go/bugs/website/web/app.html  # expect 0

# Desktop: tracking scripts present
grep "l8tracking" go/bugs/website/web/app.html  # expect 5 hits

# Desktop: tracking tab in system.html
grep 'data-module="tracking"' go/bugs/website/web/sections/system.html  # expect 2+

# Desktop: app.js defaults to system
grep "loadSection" go/bugs/website/web/js/app.js  # expect 'system'

# Mobile: no bugs references
grep -c "js/bugs/" go/bugs/website/web/m/app.html  # expect 0

# Mobile: tracking tab in system.html
grep 'data-module="tracking"' go/bugs/website/web/m/sections/system.html  # expect 1+

# Mobile: defaults to system
grep "currentSection" go/bugs/website/web/m/js/app-core.js  # expect 'system'

# JS syntax check all new files
for f in go/bugs/website/web/l8ui/sys/tracking/*.js; do node -c "$f"; done

# Old files deleted
ls go/bugs/website/web/bugs/ 2>/dev/null  # expect not found
ls go/bugs/website/web/m/js/bugs/ 2>/dev/null  # expect not found
```

## File Summary

| Action | Count | Files |
|--------|-------|-------|
| Create | 9 | 5 desktop (`l8ui/sys/tracking/`) + 4 mobile (`m/js/tracking/`) |
| Modify | 12 | l8sys-config, l8sys-init, system.html x2, app.html x2, sections.js, app.js, app-core.js, mobile-config, nav-config, GUIDE.md |
| Delete | 15 | 9 desktop (`bugs/`, `sections/bugs.html`, `reference-registry-bugs.js`) + 6 mobile (`m/js/bugs/`, `m/sections/bugs.html`) |
