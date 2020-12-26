package logger

import (
	"os"
	"time"

	"github.com/dchest/safefile"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(level zapcore.Level, outputfile string) (*zap.Logger, error) {
	err := mkFile(outputfile)
	if err != nil {
		return nil, errors.New("Error file creating")
	}
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = syslogTimeEncoder
	cfg.EncoderConfig.EncodeLevel = customLevelEncoder
	cfg.OutputPaths = []string{outputfile, "stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.Level = zap.NewAtomicLevelAt(level)

	logger, err := cfg.Build()
	if err != nil {
		return nil, errors.Wrap(err, "Logger build failed")
	}
	return logger, nil
}

func syslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("Jan 2 15:04:05"))
}

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.String() + "]")
}

func mkFile(path string) error {
	existFile := Exists(path)
	if !existFile {
		err := os.Mkdir("/tmp/tmpdir", 0755)
		if err != nil {
			return errors.Wrap(err, "Mkdir failed")
		}
		tmpfile, err := safefile.Create(path, 0755)
		if err != nil {
			return errors.Wrap(err, "Create tmpfile failed")
		}
		defer tmpfile.Close()
		return nil
	}
	return nil
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
