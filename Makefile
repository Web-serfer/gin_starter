.PHONY: check test fmt lint security

# ОДНА команда для проверки всего
check: fmt lint test security
	@echo "✅ Все проверки пройдены!"

# Форматирование
fmt:
	@echo "🎨 Форматирование кода..."
	gofmt -w .
	goimports -w .

# Линтинг
lint:
	@echo "🔍 Проверка стиля и ошибок..."
	golangci-lint run

# Тесты
test:
	@echo "🧪 Запуск тестов..."
	go test -race -cover ./...

# Безопасность
security:
	@echo "🔐 Проверка безопасности..."
	govulncheck ./...

# Быстрая проверка перед коммитом
quick:
	gofmt -d . | grep -v "^$$" || true
	golangci-lint run --fast
	go test -short ./...