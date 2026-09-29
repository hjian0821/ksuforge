package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/hjian0821/ksuforge/internal/i18n"
	"github.com/hjian0821/ksuforge/internal/paths"
)

type languageContextKey struct{}

func WithLanguage(ctx context.Context, language string) context.Context {
	return context.WithValue(ctx, languageContextKey{}, i18n.Normalize(language))
}

func Language(ctx context.Context) string {
	if ctx != nil {
		if language, ok := ctx.Value(languageContextKey{}).(string); ok {
			return i18n.Normalize(language)
		}
	}
	return i18n.English
}

// WithWriter returns a logger that writes structured records to the process
// logger and mirrors them to the supplied task stream. The stream handler omits
// timestamps and levels so GUI and CLI transcripts stay compact.
func WithWriter(out io.Writer, language ...string) *slog.Logger {
	if out == nil {
		return slog.Default()
	}
	lang := i18n.English
	if len(language) > 0 {
		lang = i18n.Normalize(language[0])
	}
	var stream slog.Handler = slog.NewTextHandler(out, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.TimeKey || attr.Key == slog.LevelKey {
				return slog.Attr{}
			}
			return attr
		},
	})
	localizedStream := localizingHandler{next: stream, language: lang}
	return slog.New(multiHandler{slog.Default().Handler(), localizedStream})
}

type localizingHandler struct {
	next     slog.Handler
	language string
}

func (h localizingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h localizingHandler) Handle(ctx context.Context, record slog.Record) error {
	localized := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	record.Attrs(func(attr slog.Attr) bool {
		localized.AddAttrs(localizeAttr(attr, h.language))
		return true
	})
	message, ok := i18n.Lookup(record.Message, h.language)
	if ok {
		event := record.Message
		localized.Message = message
		localized.AddAttrs(slog.String("event", event))
	}
	return h.next.Handle(ctx, localized)
}

func localizeAttr(attr slog.Attr, language string) slog.Attr {
	attr.Value = attr.Value.Resolve()
	if attr.Value.Kind() == slog.KindAny {
		if err, ok := attr.Value.Any().(error); ok {
			return slog.String(attr.Key, i18n.ErrorText(err, language))
		}
	}
	if attr.Value.Kind() == slog.KindGroup {
		group := attr.Value.Group()
		for i := range group {
			group[i] = localizeAttr(group[i], language)
		}
		return slog.Attr{Key: attr.Key, Value: slog.GroupValue(group...)}
	}
	return attr
}

func (h localizingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return localizingHandler{next: h.next.WithAttrs(attrs), language: h.language}
}

func (h localizingHandler) WithGroup(name string) slog.Handler {
	return localizingHandler{next: h.next.WithGroup(name), language: h.language}
}

type multiHandler []slog.Handler

func (h multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h multiHandler) Handle(ctx context.Context, record slog.Record) error {
	var firstErr error
	for _, handler := range h {
		if handler.Enabled(ctx, record.Level) {
			if err := handler.Handle(ctx, record.Clone()); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (h multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := make(multiHandler, len(h))
	for i, handler := range h {
		next[i] = handler.WithAttrs(attrs)
	}
	return next
}

func (h multiHandler) WithGroup(name string) slog.Handler {
	next := make(multiHandler, len(h))
	for i, handler := range h {
		next[i] = handler.WithGroup(name)
	}
	return next
}

// Options configures the process logger.
type Options struct {
	Level  slog.Level
	Stderr bool
	Dir    string
	Prefix string
}

// Setup builds a slog logger that appends to a daily file under the OS log
// directory. When Stderr is set, records also go to stderr. If the log file
// cannot be opened the logger falls back to stderr so logging never blocks the
// application. The returned function flushes and closes the file.
func Setup(opts Options) (*slog.Logger, func()) {
	var writers []io.Writer
	if opts.Stderr {
		writers = append(writers, os.Stderr)
	}
	file, err := openLogFile(opts.Dir, opts.Prefix)
	if err == nil {
		writers = append(writers, file)
	}
	if len(writers) == 0 {
		writers = append(writers, os.Stderr)
	}
	baseHandler := slog.NewTextHandler(io.MultiWriter(writers...), &slog.HandlerOptions{Level: opts.Level})
	handler := localizingHandler{next: baseHandler, language: i18n.English}
	logger := slog.New(handler)
	slog.SetDefault(logger)
	if err != nil && opts.Stderr {
		logger.Warn("logging.file_unavailable", "error", err)
	}
	closeFn := func() {}
	if file != nil {
		closeFn = func() { _ = file.Close() }
	}
	return logger, closeFn
}

func openLogFile(dir, prefix string) (*os.File, error) {
	if dir == "" {
		dir = paths.LogDir()
	}
	if prefix == "" {
		prefix = "ksuforge"
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	name := prefix + "-" + time.Now().Format("20060102") + ".log"
	return os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
}
