sleep 30
echo "integration tests running"
go test ./integration-test/...
docker-compose down