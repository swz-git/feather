let file;

function dropHandler(e) {
  let dataTransfer = e.dataTransfer;
  file = dataTransfer.files[0];
  if (!file) return;
  console.log("Loading dropped file...", file);
  document.querySelector("img").src = URL.createObjectURL(file);
  document.querySelector(".fileName").innerText = file.name;
  document.querySelector(".fileType").innerText = file.type;
  document.querySelector(".fileSize").innerText = file.size + " bytes";
  document.querySelectorAll(".page").forEach((page) => {
    page.classList.toggle("hidden");
  });

  e.preventDefault();
}

function downloadFile() {
  if (!file) return;
  console.log("Downloading file...", file);
  const url = URL.createObjectURL(file);
  const a = document.createElement("a");
  a.href = url;
  a.download = file.name;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
}

async function uploadFile() {
  if (!file) return;
  console.log("Uploading file...", file);
  let req = await fetch(`${location.origin}/chunksize`, {
    method: "GET",
  });

  const CHUNK_SIZE_KB = +(await req.text()); // in kb
  const CHUNK_SIZE = CHUNK_SIZE_KB * 1000; // in bytes
  const chunkCount = file.size / CHUNK_SIZE;

  console.log("Chunk count", chunkCount);
  const fileID = (Math.random() * 10 + "").replace(".", "");
  let fileURL = null;
  for (let chunkId = 0; chunkId < chunkCount + 1; chunkId++) {
    const chunk = file.slice(
      chunkId * CHUNK_SIZE,
      chunkId * CHUNK_SIZE + CHUNK_SIZE
    );
    console.log(
      "Uploading chunk",
      chunkId,
      "of",
      chunkCount,
      "\nSize is",
      chunk.size
    );
    let req = await fetch(`${location.origin}/upload?id=${fileID}`, {
      method: "POST",
      headers: {
        "content-type": "application/octet-stream",
        "content-length": chunk.length,
        "file-name": file.name,
      },
      body: chunk,
    });
    fileURL = await req.text();
    document.querySelector("progress").classList.remove("hidden");
    document.querySelector("progress").value = (chunkId * 100) / chunkCount;
  }
  document.querySelector("progress").classList.add("hidden");
  document.getElementById("resultingurl").classList.remove("hidden");
  document.getElementById("resultingurl").href = location.origin + fileURL;
  document.getElementById("resultingurl").innerText = location.origin + fileURL;
}

function onLoad() {
  const dropZone = document.getElementById("dropzone");

  dropZone.addEventListener("dragenter", function (e) {
    e.preventDefault();
    dropZone.classList.add("dropping");
  });

  dropZone.addEventListener("dragleave", function (e) {
    e.preventDefault();
    dropZone.classList.remove("dropping");
  });

  dropZone.addEventListener("drop", function (e) {
    e.preventDefault();
    dropZone.classList.remove("dropping");
    dropHandler(e);
  });

  dropZone.addEventListener("dragover", function (e) {
    e.preventDefault();
  });
}
