build:
	docker build -t localhost:5001/dapr-go-samples:latest .
	docker push localhost:5001/dapr-go-samples:latest

start-workflow:
	curl -XPOST localhost:8080/workflow

event-workflow:
	curl -XPOST localhost:8080/workflow/event