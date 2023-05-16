gen:
	mockgen -source=storage/interface.go -destination=storage/mock_storage/mock.go
	mockgen -source=shortener/interface.go -destination=shortener/mock_shortener/mock.go
