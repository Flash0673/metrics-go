package logger

import "go.uber.org/zap"

// Logger будет доступен всему коду как синглтон.
// // Никакой код навыка, кроме функции InitLogger, не должен модифицировать эту переменную.
// // По умолчанию установлен no-op-логер, который не выводит никаких сообщений.
var Logger *zap.Logger = zap.NewNop()

// New инициализирует синглтон логера с необходимым уровнем логирования.
func New(level string) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	Logger = zl
	return nil
}
