package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hjian0821/ksuforge/internal/cleanup"
	"github.com/hjian0821/ksuforge/internal/i18n"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) ScanStorage() (cleanup.Report, error) {
	a.storageMu.Lock()
	defer a.storageMu.Unlock()
	if a.State().Status == "running" {
		return cleanup.Report{}, i18n.NewError("error.storage_busy")
	}
	if a.cleanupPlan != nil {
		a.cleanupPlan.Close()
		a.cleanupPlan = nil
	}
	plan, err := cleanup.Scan(a.DefaultOutputDirectory())
	if err != nil {
		return cleanup.Report{}, i18n.WrapError("error.storage_scan", err)
	}
	a.cleanupPlan = plan
	return plan.Report, nil
}

func (a *App) CleanStorage(token string, ids []string) (cleanup.Result, error) {
	a.storageMu.Lock()
	defer a.storageMu.Unlock()
	if a.State().Status == "running" {
		return cleanup.Result{}, i18n.NewError("error.storage_busy")
	}
	p := a.cleanupPlan
	root, err := filepath.Abs(a.DefaultOutputDirectory())
	if err != nil || p == nil || token != p.Report.Token || root != p.Report.Root {
		return cleanup.Result{}, i18n.NewError("error.storage_rescan")
	}
	items, err := p.Selection(ids)
	if err != nil {
		return cleanup.Result{}, i18n.WrapError("error.storage_rescan", err)
	}
	if err := p.Validate(ids); err != nil {
		return cleanup.Result{}, i18n.WrapError("error.storage_rescan", err)
	}
	language := a.currentLanguage()
	var bytes int64
	count, protected := 0, false
	names := []string{}
	for _, item := range items {
		bytes += item.Bytes
		count += len(item.Files)
		protected = protected || item.Kind == "stock" || item.Kind == "patched"
		if len(names) < 8 {
			names = append(names, item.Files[0].Path)
		}
	}
	message := i18n.Text("storage.confirm.message", language, count, storageBytes(bytes), p.Report.Root)
	if protected {
		message += "\n\n" + i18n.Text("storage.confirm.images", language)
	}
	message += "\n\n" + strings.Join(names, "\n")
	if len(items) > len(names) {
		message += fmt.Sprintf("\n+ %d", len(items)-len(names))
	}
	cancel := i18n.Text("storage.cancel", language)
	clean := i18n.Text("storage.confirm.action", language)
	result, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type: runtime.WarningDialog, Title: i18n.Text("storage.confirm.title", language), Message: message,
		Buttons: []string{cancel, clean}, DefaultButton: cancel, CancelButton: cancel,
	})
	if err != nil {
		return cleanup.Result{}, i18n.WrapError("error.storage_clean", err)
	}
	if result != clean {
		return cleanup.Result{Cancelled: true}, nil
	}
	// The pinned plan checks identities again after the native confirmation.
	cleaned, err := p.Clean(ids)
	p.Close()
	a.cleanupPlan = nil
	if err != nil {
		return cleaned, i18n.WrapError("error.storage_rescan", err)
	}
	return cleaned, nil
}

func storageBytes(bytes int64) string {
	if bytes >= 1<<30 {
		return fmt.Sprintf("%.2f GiB", float64(bytes)/(1<<30))
	}
	return fmt.Sprintf("%.1f MiB", float64(bytes)/(1<<20))
}
