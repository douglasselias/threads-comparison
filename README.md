# Threads comparison

This repo is a small experiment to see how each language does multithreading.

The selected languages were C, Go and JavaScript.

All of them do the same thing, Create a thread to fetch a webpage and write into a file. Then wait all to finish.

One interesting thing to notice is that only C has real threads, managed by the OS, Go and JavaScript doesn't.

## Conclusion

JavaScript has by far the simplest code, followed by Go and C.

JavaScript and Go simplify a lot of the logic behind fetching a webpage and both have built-in functions to do that. C uses a library to handle the fetch and it still quite verbose. Doing without a library would require using sockets. Also, C is the only that it needs a struct to pass the id and the url to the thread function.

Despite being very different languages they all share the same structure: a way to create a thread, a collection of threads and a way to join all the threads. None of them is really any harder or any easier than the other (at least in this small code base).