clean:
  rm file_1.html file_2.html file_3.html file_4.html

c:
  gcc main.c -o main -lcurl && ./main

go:
  go run main.go

js:
  node main.js