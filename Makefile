DYNAMODB_PORT=9000
DYNAMODB_ENDPOINT=http://localhost:$(DYNAMODB_PORT)
DYNAMODB_IMAGE=$(shell docker ps | grep amazon/dynamodb-local | cut -d' ' -f1)
DYNAMODB_ALL_IMAGES=$(shell docker ps -a | grep amazon/dynamodb-local | cut -d' ' -f1)
.PHONY: dynamodb-start
dynamodb-start:
	docker pull amazon/dynamodb-local
	docker run -d -p $(DYNAMODB_PORT):$(DYNAMODB_PORT) amazon/dynamodb-local -jar DynamoDBLocal.jar -sharedDb -port $(DYNAMODB_PORT)
	aws dynamodb create-table --endpoint-url $(DYNAMODB_ENDPOINT) \
		--billing-mode PAY_PER_REQUEST \
		--table-name articles \
		--attribute-definitions \
			AttributeName=id,AttributeType=S \
		--key-schema \
			AttributeName=id,KeyType=HASH
	aws dynamodb put-item --endpoint-url $(DYNAMODB_ENDPOINT) \
		--table-name articles \
		--item '{
			"id": { "S": "63BD3670-65E9-4D10-9FAA-85F0F0BD9F24" },
			"title": { "S": "Congrats!" },
			"content": { "S": "This is a first article in the blog." },
			"created_at": { "S": "2021-02-15T15:00:00+09:00" },
			"updated_at": { "S": "2021-02-15T15:00:00+09:00" },
		}'
dynamodb-stop:
	docker stop $(DYNAMODB_IMAGE)
dynamodb-remove:
	docker rm -f $(DYNAMODB_ALL_IMAGES)
dynamodb-list-table:
	aws dynamodb list-tables --endpoint-url $(DYNAMODB_ENDPOINT)

