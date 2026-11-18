# KSUForge desktop design QA

final result: passed

## Target and accepted adaptation

The user selected web direction 2 (dark studio) and requested an application that looks and behaves like a native desktop GUI. The selected graphite surfaces, subdued purple accent, device/task separation and minimal tool rail are retained. Removing the promotional hero, reducing control sizes, using edge-to-edge split panes and docking the actions/logs are intentional adaptations to that request.

## Evidence

Screenshots are stored under `/Users/wsy182/.codex/visualizations/2026/11/16/01a1e268-5357-7880-ba91-0cd6ab66552b/ksuforge-directions/`.

- Reference: `source-b-viewport.png`, 1440 × 1000, connected Xiaomi 15, LKM, idle.
- Implementation: `native-reference-size.png`, same viewport, density and state.
- Side-by-side review: `native-design-comparison.png`, source left, implementation right.
- Default content size: `native-connected.png`, 1180 × 800.
- Minimum content size: `native-minimum.png`, 720 × 532, allowing for the OS title bar within the 720 × 560 minimum window.
- Additional reviewed states: `native-empty.png`, `native-settings.png`, `native-error.png`, `native-english.png`, `native-success.png`.
- All final visible image assets completed decoding; no missing images or browser JavaScript errors.

## Findings and fixes

- [P2, fixed] The decorative switch track intercepted clicks addressed to the underlying checkbox. The transparent input now covers the full switch hit area. Browser interaction tests passed after the fix.
- [P2, fixed] Switching languages reconstructed the stepper and lost the failed-step color. Failed state is now restored in the step renderer. The final browser run covered error state followed by English settings.
- No remaining actionable P0/P1/P2 findings in the reviewed views.

## Required fidelity surfaces

- Typography: system UI stack and monospaced technical values; compact desktop hierarchy replaces the source's promotional headline. Labels, state text and paths were inspected at the default size. Technical values support text selection.
- Layout: 60 px utility rail, 48 px toolbar, independently scrolling device and configuration panes, persistent primary actions, docked logs and status bar. Browser checks confirmed persistent controls stay inside the default and minimum viewports with no document overflow. Preferences scroll separately with fixed Save/Cancel actions.
- Colors: reference graphite surfaces and lavender action/selection states are preserved. Success and failure use separate semantic colors. Focus indicators remain visible against the dark surfaces.
- Assets: local Lucide SVG icons and the existing KSUForge application icon are used; image loading was checked. No runtime image/CDN dependency was introduced.
- Copy: promotional headings were removed. Task actions reflect prepare/flash/running states. The native flash confirmation includes the model, build, target partition and slot. Chinese/English states were reviewed.

## Interaction validation

Passed: disconnected/connected states, LKM/GKI switching and GKI field disclosure, local ROM selection/clearing, preferences focus containment and return, keyboard shortcuts, confirmation cancellation/approval, configuration locking during a run, success path display, error disclosure, log expansion and language changes.

The browser exercised the real GUI frontend against an isolated stub of the Wails bridge. No device download or flash was performed. `make test` (all Go tests and `go vet`) and the macOS Wails production build passed. macOS native application/edit menus and dark window appearance are configured; Windows uses its dark native frame. Actual Windows rendering and physical-device flashing were not part of this UI verification.

## Implementation checklist

- [x] Adapt selected direction to a desktop shell.
- [x] Preserve real Wails bindings, file dialogs and flash confirmation.
- [x] Check fixed controls at default/minimum sizes and error/success states.
- [x] Inspect reference and implementation together and review focused states.
- [x] Run application tests and build the local macOS application.

## Follow-up polish

No blocking polish remains. Platform-specific checks outside macOS can follow the normal release validation.

# Independent log workspace follow-up

final result: passed

The selected log direction is the independent log page. The user explicitly removed the Return button and requested switching through the left toolbar only.

## Evidence and accepted changes

Screenshots are under the prior screenshot directory's `logs/` subdirectory:

- `selected-log-reference.png`: selected design at 1180 × 800, connected device and failed ROM acquisition.
- `final-log-page.png`: implementation at the same size, state, event sequence and theme.
- `final-log-comparison.png`: reference left and implementation right, inspected together.
- `final-workspace.png`: main workspace with no bottom console.
- `implemented-log-en.png` and `implemented-log-minimum.png`: English and 720 × 532 content viewport.

The Return button is removed as requested. Export moves to the page header; auto-scroll remains a task-level viewer control. Actual timestamps replace the illustrative elapsed times. A task failure without an original timestamp shows a dash rather than an invented time. The illustration's absolute claim that nothing was written is not displayed for arbitrary real failures, which can occur after flashing has started.

## Fidelity review

- Typography and spacing retain the selected compact desktop table, 164 px stage rail, native-sized toolbar and independently scrolling records.
- Graphite, lavender selection, warning and error tokens match the approved direction.
- Existing local library SVGs are reused. No generated or external image dependencies were added.
- Copy is bilingual; timestamps, severity, message and expandable raw events come from the task transcript.
- There is no page-level Return button, bottom console or automatic navigation on failure.

## Validation

Passed browser checks: empty state, sidebar-only switching, live task polling while the log page is open, preservation of KMI and local ROM selections, level/stage filtering, search including raw event fields, HTML escaping, event-detail expansion, raw transcript copying, native export bridge dispatch, failure without forced navigation, English settings, keyboard search and minimum viewport controls. No browser JavaScript errors were recorded. The native export uses the Wails save dialog; its OS dialog interaction was not automated.

`make test` and the local macOS Wails production build passed. The stream logging tests now verify real timestamps and INFO/WARN/ERROR levels, so the viewer does not infer severity from localized message text.

No remaining P0/P1/P2 findings. Existing physical-device and other-platform validation limits remain unchanged.

# Icon palette and storage cleanup follow-up

final result: passed

The user authorized the storage cleanup proposal and requested the app icon's palette. This supersedes the prior lavender tokens. The native desktop layout, independent log page, saved automatic-flashing preference and arrow-free primary button remain in place.

## Evidence

Screenshots are in `/Users/wsy182/.codex/visualizations/2026/11/16/01a1e268-5357-7880-ba91-0cd6ab66552b/storage/`:

- `icon-palette-comparison.png`: source icon alongside the implemented 1180 × 800 workspace, inspected together.
- `storage-default.png`: grouped artifacts, recommended ROM selection and allocation estimates.
- `storage-protected.png`: explicit stock-image selection and warning.
- `storage-result.png`: deletion result with retained images and partial download.
- `storage-minimum-en.png`: English at 720 × 532, scrolling content with fixed, visible footer.

Graphite dominates the interface; icon green identifies actions and selection, icon blue supports source/device information. Amber and red retain warning/error meanings. Core text contrast checks passed: main text 14.19:1, secondary 7.36:1, muted 5.02:1, primary button text 8.98:1 and information blue 7.94:1 against their corresponding surfaces. Status remains represented by text and icons as well as color.

## Behavior and verification

Storage is a separate settings tab. Only recognized files in the saved output directory are scanned. The default selection includes complete ROMs with a SHA-256-verified exact-build stock backup and matching recorded ROM checksum prefix. Partial downloads, stock and patched images remain unchecked. Image sidecars travel with their images. Explicit image selection displays a warning. One native confirmation defaults to Cancel. Tasks block scanning/deletion; scan tokens and file/root identity checks reject stale selections. Allocation estimates are approximate, not guarantees of physical disk space reclaimed.

Backend tests use isolated temporary directories: recommended defaults, explicit image-and-sidecar deletion, preserved unknown/external files, symbolic links, replaced roots, unknown IDs, changed files, missing directories, tampered backups and strict stock reuse without network access. A regression test also proves a tampered extraction cache is re-extracted instead of reused.

Browser checks passed for defaults, cancellation, ID-only dispatch, protected selection, success, partial errors, active task restrictions, English, minimum window, scan failure and empty state. Independent log and persisted automatic-flashing preference regressions passed. No JavaScript errors. Intermittent preview-server connection resets were isolated by fulfilling browser requests with the exact repository assets; final screenshots use those assets.

`make test` (including `go vet`), `make build`, macOS signature verification and `git diff --check` passed. Wails bindings were generated. OS-native confirmation interaction was not automated; the bridge dispatch and deletion core were tested separately. No real user artifacts were deleted during QA. Physical-device and non-macOS runtime checks remain outside this validation.

No remaining P0/P1/P2 findings for the requested changes.
