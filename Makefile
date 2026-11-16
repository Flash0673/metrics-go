.PHONY: infra ## поднимает инфрастуктуру для проекта
infra:
	docker-compose -f ./infra/build/docker-compose.yaml up -d --force-recreate --wait

.PHONY: infra-stop ## останавливает контейнеры
infra-stop:
		docker-compose -f ./infra/build/docker-compose.yaml down