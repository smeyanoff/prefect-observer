

.PHONY: submit swag swagFmt runApp .validateMod .generateMock buildApp generateFrontend goTests helmDryRun

## Подготавливает commit message
submit: swag swagFmt
	git commit -m "Изменены файлы: $(shell git diff --name-only --cached | sed ':a;N;$!ba;s/\n/, /g')"

## Обновляет документацию swagger
swag: swagFmt
	cd ./backend && swag init

## Форматирует документацию swagger
swagFmt:
	cd ./backend && swag fmt

## Запускает приложение локально
runApp: .env swag
	cd ./backend && go run main.go

## Подготавливает приложение для сборки в docker контейнер
.validateMod:
	cd ./backend && \
	go mod tidy && \
	go mod verify && \
	go mod vendor

## Обновляет Mock объекты
.generateMock:
	cd ./backend && \
	go generate ./...

## Собирает приложение в docker контейнер
buildApp: swag goTests
	docker compose up --build -d

## Генерирует фронтенд из swagger
generateFrontend:
	cd ./frontend && \
	npx @openapitools/openapi-generator-cli generate \
	-i ../backend/docs/swagger.json -g typescript-axios -o ./src/api

## Запускает тесты go
goTests: .validateMod .generateMock
	cd ./backend && go test -v ./...

## Запускает helm dry run
helmDryRun:
	helm install observer --dry-run --debug ./observer


#################################################################################
# Self Documenting Commands                                                     #
#################################################################################

.DEFAULT_GOAL := help

# Inspired by <http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html>
# sed script explained:
# /^##/:
# 	* save line in hold space
# 	* purge line
# 	* Loop:
# 		* append newline + line to hold space
# 		* go to next line
# 		* if line starts with doc comment, strip comment character off and loop
# 	* remove target prerequisites
# 	* append hold space (+ newline) to line
# 	* replace newline plus comments by `---`
# 	* print line
# Separate expressions are necessary because labels cannot be delimited by
# semicolon; see <http://stackoverflow.com/a/11799865/1968>
.PHONY: help
help:
	@echo "$$(tput bold)Available rules:$$(tput sgr0)"
	@echo
	@sed -n -e "/^## / { \
		h; \
		s/.*//; \
		:doc" \
		-e "H; \
		n; \
		s/^## //; \
		t doc" \
		-e "s/:.*//; \
		G; \
		s/\\n## /---/; \
		s/\\n/ /g; \
		p; \
	}" ${MAKEFILE_LIST} \
	| LC_ALL='C' sort --ignore-case \
	| awk -F '---' \
		-v ncol=$$(tput cols) \
		-v indent=19 \
		-v col_on="$$(tput setaf 6)" \
		-v col_off="$$(tput sgr0)" \
	'{ \
		printf "%s%*s%s ", col_on, -indent, $$1, col_off; \
		n = split($$2, words, " "); \
		line_length = ncol - indent; \
		for (i = 1; i <= n; i++) { \
			line_length -= length(words[i]) + 1; \
			if (line_length <= 0) { \
				line_length = ncol - indent - length(words[i]) - 1; \
				printf "\n%*s ", -indent, " "; \
			} \
			printf "%s ", words[i]; \
		} \
		printf "\n"; \
	}' \
	| more $(shell test $(shell uname) = Darwin && echo '--no-init --raw-control-chars')