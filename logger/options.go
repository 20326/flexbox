package logger

import (
	"os"
	"time"

	"github.com/phuslu/log"
	"github.com/robfig/cron/v3"
)

type (
	OptionsFunc func(logger *log.Logger)
	CleanerFunc func(filename string, maxBackups int, matches []os.FileInfo)
)

func WithLevel(level string) OptionsFunc {
	return func(o *log.Logger) {
		o.Level = log.ParseLevel(level)
		writer := o.Writer
		switch writer.(type) {
		case *log.MultiWriter:
			writer.(*log.MultiWriter).ConsoleLevel = o.Level
			break
		default:
			break
		}
	}
}

func WithTimeFormat(timeFmt string) OptionsFunc {
	return func(o *log.Logger) {
		if timeFmt == "" {
			o.TimeFormat = "15:04:05"
			return
		}
		o.TimeFormat = timeFmt
	}
}

func WithMaxSize(maxSize int64) OptionsFunc {
	return func(o *log.Logger) {
		writer := o.Writer
		switch writer.(type) {
		case *log.MultiWriter:
			mw := writer.(*log.MultiWriter)
			mw.InfoWriter.(*log.FileWriter).MaxSize = maxSize
			mw.WarnWriter.(*log.FileWriter).MaxSize = maxSize
			mw.ErrorWriter.(*log.FileWriter).MaxSize = maxSize
			break
		case *log.FileWriter:
			fw := writer.(*log.FileWriter)
			fw.MaxSize = maxSize
			break
		default:
			break
		}
	}
}

func WithMaxBackups(maxBackups int) OptionsFunc {
	return func(o *log.Logger) {
		writer := o.Writer
		switch writer.(type) {
		case *log.MultiWriter:
			mw := writer.(*log.MultiWriter)
			mw.InfoWriter.(*log.FileWriter).MaxBackups = maxBackups
			mw.WarnWriter.(*log.FileWriter).MaxBackups = maxBackups
			mw.ErrorWriter.(*log.FileWriter).MaxBackups = maxBackups
			break
		case *log.FileWriter:
			fw := writer.(*log.FileWriter)
			fw.MaxBackups = maxBackups
			break
		default:
			break
		}
	}
}

func WithFileName(fileName string) OptionsFunc {
	return func(o *log.Logger) {
		writer := o.Writer
		switch writer.(type) {
		case *log.MultiWriter:
			mw := writer.(*log.MultiWriter)
			mw.InfoWriter.(*log.FileWriter).Filename = fileName + ".INFO"
			mw.WarnWriter.(*log.FileWriter).Filename = fileName + ".WARNING"
			mw.ErrorWriter.(*log.FileWriter).Filename = fileName + ".ERROR"
			break
		case *log.FileWriter:
			fw := writer.(*log.FileWriter)
			fw.Filename = fileName
			break
		default:
			break
		}
	}
}

func WithCronRunner(spec string) OptionsFunc {
	return func(o *log.Logger) {
		runner := cron.New(cron.WithSeconds(), cron.WithLocation(time.Local))
		_, _ = runner.AddFunc(spec, func() {
			writer := o.Writer
			switch writer.(type) {
			case *log.MultiWriter:
				mw := writer.(*log.MultiWriter)
				_ = mw.InfoWriter.(*log.FileWriter).Rotate()
				_ = mw.WarnWriter.(*log.FileWriter).Rotate()
				_ = mw.ErrorWriter.(*log.FileWriter).Rotate()
				break
			case *log.FileWriter:
				fw := writer.(*log.FileWriter)
				_ = fw.Rotate()
				break
			case *log.AsyncWriter:
				// TODO
				break
			case *log.ConsoleWriter:
				return
			default:
				log.Error().Msg("unknown logger writer")
				return
			}
		})
		go runner.Run()
	}
}
