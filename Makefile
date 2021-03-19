DYNAMODB_PORT=9000
DYNAMODB_ENDPOINT=http://localhost:$(DYNAMODB_PORT)
DYNAMODB_IMAGE=$(shell docker ps | grep amazon/dynamodb-local | cut -d' ' -f1)
DYNAMODB_ALL_IMAGES=$(shell docker ps -a | grep amazon/dynamodb-local | cut -d' ' -f1)
.PHONY: dynamodb-start
dynamodb-start:
	docker pull amazon/dynamodb-local
	docker run -d -p $(DYNAMODB_PORT):$(DYNAMODB_PORT) amazon/dynamodb-local -jar DynamoDBLocal.jar -sharedDb -port $(DYNAMODB_PORT)
dynamodb-stop:
	docker stop $(DYNAMODB_IMAGE)
dynamodb-remove:
	docker rm -f $(DYNAMODB_ALL_IMAGES)
dynamodb-list-table:
	aws dynamodb list-tables --endpoint-url $(DYNAMODB_ENDPOINT)

