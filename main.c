#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>

#include <curl/curl.h>

size_t write_callback(void *contents, size_t size, size_t nmemb, FILE *file) {
  return fwrite(contents, size, nmemb, file);
}

#define total_threads 4

typedef struct {
  char *url;
  int id;
} Task;

void *download_webpage(void *arg) {
  Task *task = (Task *)arg;
  printf("Worker %d started\n", task->id);
  char filename[100];
  snprintf(filename, sizeof(filename), "file_%d.html", task->id);

  curl_global_init(CURL_GLOBAL_DEFAULT);
  CURL *curl = curl_easy_init();
  curl_easy_setopt(curl, CURLOPT_URL, task->url);
  // Set the callback function to write data to a file
  FILE *file = fopen(filename, "wb");
  curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
  curl_easy_setopt(curl, CURLOPT_WRITEDATA, file);

  // Perform the request
  curl_easy_perform(curl);

  // Cleanup
  fclose(file);
  curl_easy_cleanup(curl);

  // Cleanup libcurl
  curl_global_cleanup();

  printf("Worker %d finished\n", task->id);
  return NULL;
}

int main() {
  pthread_t threads[total_threads] = {};
  char *urls[total_threads] = {
      "https://example.com",
      "https://www.google.com/",
      "https://store.steampowered.com/",
      "https://pomodorotimer.online/",
  };
  Task *tasks[total_threads] = {};

  for (int i = 0; i < total_threads; i++) {
    Task *task = malloc(sizeof(Task));
    task->id = i;
    task->url = urls[i];
    tasks[i] = task;
    pthread_create(&threads[i], NULL, download_webpage, tasks[i]);
  }

  for (int i = 0; i < total_threads; i++) {
    pthread_join(threads[i], NULL);
    free(tasks[i]);
  }

  printf("Finished\n");
}