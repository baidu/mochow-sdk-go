/*
 * Copyright 2017 Baidu, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the
 * License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions
 * and limitations under the License.
 */

// logger.go - defines the logger structure and methods

// Package log implements the log facilities for BCE. It supports log to stderr, stdout as well as
// log to file with rotating. It is safe to be called by multiple goroutines.
// By using the package level function to use the default logger:
//
//	log.SetLogHandler(log.stdout | log.file) // default is log to stdout
//	log.SetLogDir("/tmp")                    // default is /tmp
//	log.SetRotateType(log.rotateDay)        // default is log.HOUR
//	log.SetRotateSize(1 << 30)               // default is 1GB
//	log.SetLogLevel(log.INFO)                // default is log.DEBUG
//	log.Debug(1, 1.2, "a")
//	log.Debugln(1, 1.2, "a")
//	log.Debugf(1, 1.2, "a")
//
// User can also create new logger without using the default logger:
//
//	customLogger := log.NewLogger()
//	customLogger.SetLogHandler(log.file)
//	customLogger.Debug(1, 1.2, "a")
//
// The log format can also support custom setting by using the following interface:
//
//	log.SetLogFormat([]string{log.fmtLevel, log.fmtTime, log.fmtMsg})
//
// Most of the cases just use the default format is enough:
//
//	[]string{fmtLevel, fmtLTime, fmtLocation, fmtMsg}
package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Handler uint8

// Constants for log handler flags, default is stdout
const (
	None   Handler = 0
	Stdout Handler = 1
	Stderr Handler = 1 << 1
	File   Handler = 1 << 2
)

type RotateStrategy uint8

// Constants for log rotating strategy when logging to file, default is by hour
const (
	RotateNone RotateStrategy = iota
	RotateDay
	RotateHour
	RotateMinute
	RotateSize

	DefaultRotateType          = RotateHour
	DefaultRotateSize    int64 = 1 << 30
	DefaultLogDir              = "/tmp"
	RotateSizeFilePrefix       = "rotating"
)

type Level uint8

// Constants for log levels, default is DEBUG
const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
	PANIC
)

var gLevelString = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL", "PANIC"}

// Constants of the log format components to support user custom specification
const (
	fmtLevel    = "level"
	fmtLTime    = "ltime"    // long time with microsecond
	fmtTime     = "time"     // just with second
	fmtLocation = "location" // caller's location with file, line, function
	fmtMsg      = "msg"
)

var (
	logFmtStr = map[string]string{
		fmtLevel:    "[%s]",
		fmtLTime:    "2006-01-02 15:04:05.000000",
		fmtTime:     "2006-01-02 15:04:05",
		fmtLocation: "%s:%d:%s:",
		fmtMsg:      "%s",
	}
	gDefaultLogFormat = []string{fmtLevel, fmtLTime, fmtLocation, fmtMsg}
)

type writerArgs struct {
	record     string
	rotateArgs interface{} // used for rotating: the size of the record or the logging time
}

// Logger defines the internal implementation of the log facility
type logger struct {
	writers        map[Handler]io.WriteCloser // the destination writer to log message
	writerChan     chan *writerArgs           // the writer channal to pass each record and time or size
	logFormat      []string
	levelThreshold Level
	handler        Handler

	// Fields that used when logging to file
	logDir     string
	logFile    string
	rotateType RotateStrategy
	rotateSize int64
	done       chan bool
}

func (l *logger) logging(level Level, format string, args ...interface{}) {
	// Only log message that set the handler and is greater than or equal to the threshold
	if l.handler == None || level < l.levelThreshold {
		return
	}

	// Generate the log record string and pass it to the writer channel
	now := time.Now()
	pc, file, line, ok, funcname := uintptr(0), "???", 0, true, "???"
	pc, file, line, ok = runtime.Caller(2)
	if ok {
		funcname = runtime.FuncForPC(pc).Name()
		funcname = filepath.Ext(funcname)
		funcname = strings.TrimPrefix(funcname, ".")
		file = filepath.Base(file)
	}
	buf := make([]string, 0, len(l.logFormat))
	msg := fmt.Sprintf(format, args...)
	for _, f := range l.logFormat {
		if _, exists := logFmtStr[f]; !exists { // skip not supported part
			continue
		}
		fmtStr := logFmtStr[f]
		switch f {
		case fmtLevel:
			buf = append(buf, fmt.Sprintf(fmtStr, gLevelString[level]))
		case fmtLTime:
			buf = append(buf, now.Format(fmtStr))
		case fmtTime:
			buf = append(buf, now.Format(fmtStr))
		case fmtLocation:
			buf = append(buf, fmt.Sprintf(fmtStr, file, line, funcname))
		case fmtMsg:
			buf = append(buf, fmt.Sprintf(fmtStr, msg))
		}
	}
	record := strings.Join(buf, " ")
	if l.rotateType == RotateSize {
		l.writerChan <- &writerArgs{record, int64(len(record))}
	} else {
		l.writerChan <- &writerArgs{record, now}
	}

	// wait for current record done logging
}

func (l *logger) buildWriter(args interface{}) {
	if l.handler&Stdout == Stdout {
		l.writers[Stdout] = os.Stdout
	} else {
		delete(l.writers, Stdout)
	}
	if l.handler&Stderr == Stderr {
		l.writers[Stderr] = os.Stderr
	} else {
		delete(l.writers, Stderr)
	}
	if l.handler&File == File {
		l.writers[File] = l.buildFileWriter(args)
	} else {
		delete(l.writers, File)
	}
}

func (l *logger) buildFileWriter(args interface{}) io.WriteCloser {
	if l.handler&File != File {
		return os.Stderr
	}

	if len(l.logDir) == 0 {
		l.logDir = DefaultLogDir
	}
	if l.rotateType < RotateNone || l.rotateType > RotateSize {
		l.rotateType = DefaultRotateType
	}
	if l.rotateType == RotateSize && l.rotateSize == 0 {
		l.rotateSize = DefaultRotateSize
	}

	logFile, needCreateFile := "", false
	if l.rotateType == RotateSize {
		recordSize, _ := args.(int64)
		logFile, needCreateFile = l.buildFileWriterBySize(recordSize)
	} else {
		recordTime, _ := args.(time.Time)
		switch l.rotateType {
		case RotateNone:
			logFile = "default.log"
		case RotateDay:
			logFile = recordTime.Format("2006-01-02.log")
		case RotateHour:
			logFile = recordTime.Format("2006-01-02_15.log")
		case RotateMinute:
			logFile = recordTime.Format("2006-01-02_15-04.log")
		}
		if _, exist := getFileInfo(filepath.Join(l.logDir, logFile)); !exist {
			needCreateFile = true
		}
	}
	l.logFile = logFile
	logFile = filepath.Join(l.logDir, l.logFile)

	// Should create new file
	if needCreateFile {
		if w, ok := l.writers[File]; ok {
			w.Close()
		}
		if writer, err := os.Create(logFile); err == nil {
			return writer
		}
		return os.Stderr
	}

	// Already open the file
	if w, ok := l.writers[File]; ok {
		return w
	}

	// Newly open the file
	if writer, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND, 0666); err == nil {
		return writer
	}
	return os.Stderr
}

func (l *logger) buildFileWriterBySize(recordSize int64) (string, bool) {
	logFile, needCreateFile := "", false
	// First running the program and need to get filename by checking the existed files
	if len(l.logFile) == 0 {
		fname := fmt.Sprintf("%s-%s.0.log", RotateSizeFilePrefix, getSizeString(l.rotateSize))
		for {
			size, exist := getFileInfo(filepath.Join(l.logDir, fname))
			if !exist {
				logFile, needCreateFile = fname, true
				break
			}
			if exist && size+recordSize <= l.rotateSize {
				logFile, needCreateFile = fname, false
				break
			}
			fname = getNextFileName(fname)
		}
	} else { // check the file size to append to the existed file or create a new file
		currentFile := filepath.Join(l.logDir, l.logFile)
		size, exist := getFileInfo(currentFile)
		if !exist {
			logFile, needCreateFile = l.logFile, true
		} else {
			if size+recordSize > l.rotateSize { // size exceeded
				logFile, needCreateFile = getNextFileName(l.logFile), true
			} else {
				logFile, needCreateFile = l.logFile, false
			}
		}
	}
	return logFile, needCreateFile
}

func (l *logger) SetHandler(h Handler) { l.handler = h }

func (l *logger) SetLogDir(dir string) { l.logDir = dir }

func (l *logger) SetLogLevel(level Level) { l.levelThreshold = level }

func (l *logger) SetLogFormat(format []string) { l.logFormat = format }

func (l *logger) SetRotateType(rotate RotateStrategy) { l.rotateType = rotate }

func (l *logger) SetRotateSize(size int64) { l.rotateSize = size }

func (l *logger) Debug(msg ...interface{}) { l.logging(DEBUG, "%s\n", concat(msg...)) }

func (l *logger) Debugln(msg ...interface{}) { l.logging(DEBUG, "%s\n", concat(msg...)) }

func (l *logger) Debugf(f string, msg ...interface{}) { l.logging(DEBUG, f+"\n", msg...) }

func (l *logger) Info(msg ...interface{}) { l.logging(INFO, "%s\n", concat(msg...)) }

func (l *logger) Infoln(msg ...interface{}) { l.logging(INFO, "%s\n", concat(msg...)) }

func (l *logger) Infof(f string, msg ...interface{}) { l.logging(INFO, f+"\n", msg...) }

func (l *logger) Warn(msg ...interface{}) { l.logging(WARN, "%s\n", concat(msg...)) }

func (l *logger) Warnln(msg ...interface{}) { l.logging(WARN, "%s\n", concat(msg...)) }

func (l *logger) Warnf(f string, msg ...interface{}) { l.logging(WARN, f+"\n", msg...) }

func (l *logger) Error(msg ...interface{}) { l.logging(ERROR, "%s\n", concat(msg...)) }

func (l *logger) Errorln(msg ...interface{}) { l.logging(ERROR, "%s\n", concat(msg...)) }

func (l *logger) Errorf(f string, msg ...interface{}) { l.logging(ERROR, f+"\n", msg...) }

func (l *logger) Fatal(msg ...interface{}) { l.logging(FATAL, "%s\n", concat(msg...)) }

func (l *logger) Fatalln(msg ...interface{}) { l.logging(FATAL, "%s\n", concat(msg...)) }

func (l *logger) Fatalf(f string, msg ...interface{}) { l.logging(FATAL, f+"\n", msg...) }

func (l *logger) Panic(msg ...interface{}) {
	record := concat(msg...)
	l.logging(PANIC, "%s\n", record)
	panic(record)
}

func (l *logger) Panicln(msg ...interface{}) {
	record := concat(msg...)
	l.logging(PANIC, "%s\n", record)
	panic(record)
}

func (l *logger) Panicf(format string, msg ...interface{}) {
	record := fmt.Sprintf(format, msg...)
	l.logging(PANIC, format+"\n", msg...)
	panic(record)
}

func (l *logger) Close() {
	select {
	case <-l.done:
		return
	default:
	}
	l.writerChan <- nil
}

func NewLogger() *logger {
	obj := &logger{
		writers:        make(map[Handler]io.WriteCloser, 3), // now only support 3 kinds of handler
		writerChan:     make(chan *writerArgs, 100),
		logFormat:      gDefaultLogFormat,
		levelThreshold: DEBUG,
		handler:        None,
		done:           make(chan bool),
	}
	// The backend writer goroutine to write each log record
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println(e)
			}
		}()
		for {
			select {
			case <-obj.done:
				return
			case args := <-obj.writerChan: // wait until a record comes to log
				if args == nil {
					close(obj.done)
					close(obj.writerChan)
					return
				}
				obj.buildWriter(args.rotateArgs)
				for _, w := range obj.writers {
					fmt.Fprint(w, args.record)
				}
			}
		}
	}()

	return obj
}
