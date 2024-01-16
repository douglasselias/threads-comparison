// @ts-check
const fs = require("fs")

async function downloadWebpage(id, url) {
  console.log("Worker %d started", id)

  const response = await fetch(url)
  const text = await response.text()
  fs.writeFileSync(`file_${id}.html`, text)
  console.log("Worker %d finished", id)
}

async function main() {
  const urls = [
    "https://example.com",
    "https://www.google.com/",
    "https://store.steampowered.com/",
    "https://pomodorotimer.online/",
  ]

  await Promise.all(
    urls.map((url, i) => downloadWebpage(i, url))
  )

  console.log('Finished')
}

main()
