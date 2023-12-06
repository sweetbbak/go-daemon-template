default:
    go build
    ./dae
    sleep 0.5
    ./dae -e "notify-send hello client"
    sleep 0.5
    ./dae -s stop
