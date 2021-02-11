sleep 100
echo "integration tests running"
go test ./integration-test/...
echo "docker-compose down"
docker-compose down