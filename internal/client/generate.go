package client

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i github.com/en7ka/chat-server/internal/client/db.TxManager -o ./mocks/ -s "_mock.go"
