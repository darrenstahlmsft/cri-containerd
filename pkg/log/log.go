/*
Copyright 2018 The Containerd Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package log

import "github.com/sirupsen/logrus"

// TODO(random-liu): Add trace support in containerd.

// TraceLevel is the log level for trace.
const TraceLevel = logrus.Level(uint32(logrus.DebugLevel + 1))

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (logrus.Level, error) {
	if lvl == "trace" {
		return TraceLevel, nil
	}
	return logrus.ParseLevel(lvl)
}

// Trace logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	if logrus.GetLevel() >= TraceLevel {
		logrus.Debug(args...)
	}
}

// Tracef logs a message at level Trace on the standard logger.
func Tracef(format string, args ...interface{}) {
	if logrus.GetLevel() >= TraceLevel {
		logrus.Debugf(format, args...)
	}
}
