# meetings-pdf

[Click here to download the executable](https://github.com/douglasselias/meetings-pdf/releases/download/initial_release/meetings_server_windows.exe)

Then open [http://localhost:8080](http://localhost:8080)

or build manually:

`go build -o main.exe ./main.go && ./main.exe`

if you are using linux you can cross compile to windows:

`GOOS=windows GOARCH=amd64 go build -o meetings_server_windows.exe ./main.go`
