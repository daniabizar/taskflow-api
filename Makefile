.PHONY: help build up down logs restart clean test

help:
	@echo "TaskFlow API - Available Commands:"
	@echo "  make build    - Build Docker images"
	@echo "  make up       - Start all services"
	@echo "  make down     - Stop all services"
	@echo "  make logs     - View logs"
	@echo "  make restart  - Restart all services"
	@echo "  make clean    - Remove all containers and volumes"
	@echo "  make test     - Run tests"

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

restart:
	docker-compose restart

clean:
	docker-compose down -v
	docker system prune -f

test:
	go test -v ./...
```

**SAVE file ini!**

---

### **STEP 11.4: Update .gitignore**

Tambahkan beberapa file Docker yang tidak perlu di-commit.

**Di VS Code:**
1. Buka file `.gitignore` yang sudah ada
2. Tambahkan di akhir file:
```
# Docker
*.log

# Binary
main