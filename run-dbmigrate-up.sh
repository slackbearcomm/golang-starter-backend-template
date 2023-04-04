# Load Config File
if [ "$#" -eq 1 ]; then
    if [ -e "config/$1.env" ]; then
        cp "config/$1.env" .env.local
    fi
elif [ ! -e ".env.local" ]; then
    cp config/local.env .env.local
fi

# Run Application
go run cmd/main.go -migrateup