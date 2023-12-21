# meetings-pdf

[Click here to download the html files](https://github.com/douglasselias/meetings-pdf/archive/refs/tags/second_release.zip)

[Click here to download the executable](https://github.com/douglasselias/meetings-pdf/releases/download/second_release/meetings_server_windows.exe)

Unzip the html files then put the executable inside the folder then start the executable.

Then open [http://localhost:8080](http://localhost:8080)

or build manually:

`go build -o main.exe ./main.go`

if you are using linux you can cross compile to windows:

`GOOS=windows GOARCH=amd64 go build -o meetings_server_windows.exe ./main.go`
