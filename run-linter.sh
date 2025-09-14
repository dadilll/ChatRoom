SERVICES=("service_auth" "room_service" "notification_service" "message_service")
ERRORS=0

for dir in "${SERVICES[@]}"; do
    echo "===================="
    echo "Linting $dir"
    echo "===================="

    # Запускаем линтер для сервиса
    if ! (cd "$dir" && golangci-lint run); then
        echo "Errors found in $dir"
        ERRORS=1
    fi

    echo ""
done

if [ $ERRORS -eq 1 ]; then
    echo "Linting finished: some errors found."
    exit 1
else
    echo "Linting finished: no errors."
fi
