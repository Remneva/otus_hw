sleep 40
echo "integration tests running"
go test ./integration-test/...
docker-compose down