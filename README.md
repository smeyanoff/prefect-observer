# Prefect Observer

[![Go](https://img.shields.io/badge/lang-Go-blue.svg)]()
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)]()
[![Docker Pulls](https://img.shields.io/docker/pulls/smeyanoff/prefect-observer.svg)]()
[![Helm Chart](https://img.shields.io/badge/helm-chart-blue.svg)]()

![image](https://github.com/user-attachments/assets/7a9119f5-4056-4a2b-a786-d6bb266001ea)


<!-- README на русском -->

## Содержание

- [Описание](#описание)
- [Основные возможности](#основные-возможности)
- [Архитектура](#архитектура)
- [Требования](#требования)
- [Начало работы](#начало-работы)
  - [1. Клонирование репозитория](#1-клонирование-репозитория)
  - [2. Конфигурация переменных окружения](#2-конфигурация-переменных-окружения)
  - [3. Локальный запуск с Docker Compose](#3-локальный-запуск-с-docker-compose)
  - [4. Развёртывание в Kubernetes с Helm](#4-развёртывание-в-kubernetes-с-helm)
- [API Endpoints](#api-endpoints)
- [Разработка](#разработка)
- [Лицензия](#лицензия)
- [Contributing](#contributing)

---

## Описание

**Prefect Observer** — это небольшой pet-проект, цель которого — строить более высокие абстракции поверх Prefect тасков, объединяя их в бизнес-логические рабочие процессы ("sendposts") и обеспечивая мониторинг их выполнения в реальном времени.

## Основные возможности

- **Observer Runner**: таск-наблюдатель, который мониторит выполнение указанного Prefect deployment по его ID в течение последних 24 часов.
- **Sequential Runner**: последовательно запускает один Prefect таск по ID деплоймента.
- **Parallel Runner**: одновременно запускает несколько Prefect тасков по списку ID деплойментов.
- **Отчет об ошибках**: детально показывает, на каком этапе (какой конкретной Prefect таске) в sendpost произошла ошибка.
- **WebSocket-уведомления**: потоковые обновления статусов через WebSocket.
- **API-документация**: встроенная Swagger UI для интерактивного изучения API.
- **DDD-архитектура бэкенда**: все компоненты бэкенда реализованы по принципам Domain-Driven Design для четкого разделения доменной логики.
- **Автоматизация развёртывания**: Docker Compose для локального запуска и Helm-чарт для Kubernetes.

## Архитектура

1. **Бэкенд (Go/Gin, DDD)**
   - REST API и WebSocket-сервер на базе [Gin](https://github.com/gin-gonic/gin).
   - Слой доменной логики (domain) четко отделен от приложения (application) и инфраструктуры (infrastructure) согласно DDD.
   - PostgreSQL через GORM.
   - Интеграция с Prefect V2 для оркестрации и опроса статусов.
   - Сервисы для sendpost и stage, фабрика раннеров.
   - Сервис уведомлений по WebSocket.

2. **Фронтенд**
   - SPA для создания, настройки и запуска sendpost-рабочих процессов.
   - Отображение статусов тасков и ошибок в реальном времени.

3. **Helm-чарт**
   - Шаблоны для Deployment, Service, ConfigMap и WebSocket Ingress.
   - Файлы `values.yaml` и `values-prod.yaml` для настройки.

## Требования

- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- (Опционально) [Helm](https://helm.sh/) для Kubernetes
- Git

## Начало работы

### 1. Клонирование репозитория

```bash
git clone https://github.com/smeyanoff/prefect-observer.git
cd prefect-observer
```

### 2. Конфигурация переменных окружения

Скопируйте пример и настройте переменные:

```bash
cp .env.example .env
# Измените в .env:
# OBSERVER_DB_HOST, OBSERVER_DB_PORT, OBSERVER_DB_DATABASE, OBSERVER_DB_USER, OBSERVER_DB_PWD
# OBSERVER_APP_PORT, OBSERVER_APP_PREFECTAPIURL, OBSERVER_APP_STAGESTATUSQUERYTIMEOUT, OBSERVER_APP_NUMWORKERS, OBSERVER_APP_HOST
```

### 3. Локальный запуск с Docker Compose

```bash
docker-compose up --build
```

- **Backend**: `http://localhost:${OBSERVER_APP_PORT}/v1`
- **Swagger UI**: `http://localhost:${OBSERVER_APP_PORT}/swagger/index.html`
- **Frontend**: `http://localhost/`

### 4. Развёртывание в Kubernetes с Helm (опционально)

```bash
helm install observer ./observer --values observer/values.yaml
# Для production override:
# helm install observer ./observer --values observer/values-prod.yaml
```

---

## API Endpoints

### Sendpost Workflows

| Метод | Путь | Описание |
| ------ | ---- | -------- |
| POST | `/v1/sendposts` | Создать sendpost workflow |
| GET | `/v1/sendposts` | Получить список всех sendposts |
| GET | `/v1/sendposts/:sendpost_id` | Получить details sendpost |
| DELETE | `/v1/sendposts/:sendpost_id` | Удалить sendpost |
| POST | `/v1/sendposts/:sendpost_id/parameters` | Добавить/обновить параметр |
| DELETE | `/v1/sendposts/:sendpost_id/parameters/:key` | Удалить параметр |

### Stage Management

| Метод | Путь | Описание |
| ------ | ---- | -------- |
| POST | `/v1/sendposts/:sendpost_id/stages` | Добавить этап (Prefect task) |
| GET | `/v1/sendposts/:sendpost_id/stages` | Список этапов |
| GET | `/v1/sendposts/:sendpost_id/stages/:stage_id` | Информация об этапе |
| PATCH | `/v1/sendposts/:sendpost_id/stages/:stage_id` | Блок/разблок этапа |
| PUT | `/v1/sendposts/:sendpost_id/stages/:stage_id` | Обновить параметры этапа |
| DELETE | `/v1/sendposts/:sendpost_id/stages/:stage_id` | Удалить этап |
| GET | `/v1/prefectV2/:deployment_id/parameters` | Параметры Prefect deployment |

### Workflow Execution & Notifications

| Метод | Путь | Описание |
| ------ | ---- | -------- |
| POST | `/v1/sendposts/:sendpost_id/run` | Запустить runner |
| GET | `/v1/sendposts/:sendpost_id/run/ws` | WebSocket для live updates |

---

## Разработка

- **Backend**: `/backend` (Go модули, `main.go`, сервисы, контроллеры)
- **Frontend**: `/frontend` (SPA фреймворк)
- **Helm**: `/observer` (чарт)
- `Makefile` содержит команды для lint, тестов и сборки.

## Лицензия

Этот проект лицензирован под Apache 2.0 License.

## Contributing

PR, Issues и feature requests приветствуются! Форкайте репозиторий и присылайте PR.

---

<!-- English README -->

[![Go](https://img.shields.io/badge/lang-Go-blue.svg)]()
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)]()
[![Docker Pulls](https://img.shields.io/docker/pulls/smeyanoff/prefect-observer.svg)]()
[![Helm Chart](https://img.shields.io/badge/helm-chart-blue.svg)]()

# Prefect Observer

## Table of Contents

- [Description](#description)
- [Features](#features)
- [Architecture Overview](#architecture-overview)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints-1)
- [Development](#development)
- [License](#license)
- [Contributing](#contributing-1)

---

## Description

**Prefect Observer** is a lightweight pet project designed to build higher-level abstractions on top of Prefect tasks, uniting them into business-logic-driven workflows ("sendposts") and providing real-time monitoring of their execution.

## Features

- **Observer Runner**: Monitors a specified Prefect deployment by its ID and checks whether the associated task has run within the last 24 hours.
- **Sequential Runner**: Triggers a single Prefect deployment task by its deployment ID in a sequential fashion.
- **Parallel Runner**: Executes multiple Prefect deployment tasks concurrently by their deployment IDs.
- **Error Reporting**: Identifies and displays which specific stage (Prefect task) in a sendpost workflow failed.
- **WebSocket Notifications**: Real-time updates on workflow progression and task statuses via WebSocket endpoints.
- **API Documentation**: Built-in Swagger UI for interactive API exploration.
- **Domain-Driven Design (DDD) Backend**: The backend is implemented following DDD principles, separating domain logic, application services, and infrastructure.
- **Deployment Automation**: Includes Docker Compose for local setup and a Helm chart for Kubernetes deployment.

## Architecture Overview

The project consists of three main components:

1. **Backend (Go/Gin, DDD)**
   - REST API and WebSocket server built with [Gin](https://github.com/gin-gonic/gin).
   - Domain layer, application services, and infrastructure layer separated according to DDD.
   - PostgreSQL persistence via GORM.
   - Prefect V2 client integration for workflow orchestration and status checks.
   - Modular services for sendpost and stage management, along with runner factories.
   - Notification service using WebSocket for live updates.

2. **Frontend**
   - Single-page application (SPA) enabling users to create, configure, and run sendpost workflows.
   - Real-time display of task statuses and detailed error messages.

3. **Observer (Helm Chart)**
   - Helm chart for Kubernetes deployment (`/observer` directory).
   - Supports customization via `values.yaml` and `values-prod.yaml`.
   - Templates for Deployments, Services, ConfigMaps, and WebSocket ingress.

## Prerequisites

- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- (Optional) [Helm](https://helm.sh/) for Kubernetes deployments
- Git

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/smeyanoff/prefect-observer.git
cd prefect-observer
```

### 2. Configure environment variables

```bash
cp .env.example .env
# OBSERVER_DB_HOST, OBSERVER_DB_PORT, OBSERVER_DB_DATABASE, OBSERVER_DB_USER, OBSERVER_DB_PWD
# OBSERVER_APP_PORT, OBSERVER_APP_PREFECTAPIURL, OBSERVER_APP_STAGESTATUSQUERYTIMEOUT, OBSERVER_APP_NUMWORKERS, OBSERVER_APP_HOST
```

### 3. Run locally with Docker Compose

```bash
docker-compose up --build
```

- **Backend**: `http://localhost:${OBSERVER_APP_PORT}/v1`
- **Swagger UI**: `http://localhost:${OBSERVER_APP_PORT}/swagger/index.html`
- **Frontend**: `http://localhost/`

### 4. Deploy to Kubernetes with Helm (optional)

```bash
helm install observer ./observer --values observer/values.yaml
# For production overrides:
# helm install observer ./observer --values observer/values-prod.yaml
```

## API Endpoints

Refer to the sections below for full API details.

## Development

- **Backend**: `/backend` (Go modules, `main.go`, services, controllers)
- **Frontend**: `/frontend` (SPA framework)
- **Observer**: `/observer` (Helm chart)
- **Makefile**: commands for linting, testing, and builds.

## License

Apache 2.0 License

## Contributing

Contributions, issues, and feature requests are welcome!

