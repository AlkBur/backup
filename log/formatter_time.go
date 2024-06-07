package log

import (
	"log/slog"
	"time"
)

func TimezoneConverter(location *time.Location) Formatter {
	return FormatByKind(slog.KindTime, func(value slog.Value) slog.Value {
		t := value.Time()

		if location == nil {
			location = time.UTC
		}

		return slog.TimeValue(t.In(location))
	})
}

func TimeFormatter(timeFormat string, location *time.Location) Formatter {
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	return FormatByKind(slog.KindTime, func(value slog.Value) slog.Value {
		t := value.Time()

		if location != nil {
			t = t.In(location)
		}

		return slog.StringValue(t.Format(timeFormat))
	})
}

func FormatByKind(kind slog.Kind, formatter func(slog.Value) slog.Value) Formatter {
	return func(_ []string, attr slog.Attr) (slog.Value, bool) {
		value := attr.Value

		if value.Kind() == kind {
			return formatter(value), true
		}

		return value, false
	}
}
