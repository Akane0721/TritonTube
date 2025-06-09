package web

import (
	"encoding/json"
	"html/template"
	"time"
)

var funcMap = template.FuncMap{
	"formatTime": func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	},
	"toJSON": func(v interface{}) (template.JS, error) {
		b, err := json.Marshal(v)
		return template.JS(b), err
	},
}

const indexHTML = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8" />
  <title>TritonTube</title>
  <link
    rel="stylesheet"
    href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css"
  />
  <script src="https://unpkg.com/vue@3/dist/vue.global.prod.js"></script>
</head>
<body class="bg-light">
  <div id="app" class="container py-5">
    <h1 class="mb-4 text-primary">🎥 TritonTube</h1>

    <h3>Upload an MP4 Video</h3>
    <div class="mb-3">
      <input
        type="file"
        class="form-control"
        accept="video/mp4"
        ref="fileInput"
        @change="onFileChange"
      />
    </div>
    <button
      class="btn btn-success mb-3"
      :disabled="!file"
      @click="upload"
    >
      Upload
    </button>

    <div v-if="progress >= 0" class="mb-4">
      <div class="progress">
        <div
          class="progress-bar"
          role="progressbar"
          :style="{ width: progress + '%' }"
          v-text="progress + '%'"
        ></div>
      </div>
    </div>

    <h3>Watchlist</h3>
    <ul class="list-group">
      <li
        v-for="v in videos"
        :key="v.Id"
        class="list-group-item d-flex justify-content-between align-items-center"
      >
        <a :href="'/videos/' + v.Id" class="link">{{ v.Id }}</a>
        <small class="text-muted">{{ formatTime(v.UploadedAt) }}</small>
      </li>
      <li v-if="videos.length === 0" class="list-group-item">
        No videos uploaded yet.
      </li>
    </ul>
  </div>

  <script>
    const { createApp } = Vue;
    createApp({
      data() {
        return {
          file: null,
          progress: -1,
          videos: {{ . | toJSON }},
        };
      },
      methods: {
        onFileChange(evt) {
          this.file = evt.target.files[0];
        },
        upload() {
          if (!this.file) return;
          const form = new FormData();
          form.append("file", this.file);
          const xhr = new XMLHttpRequest();
          xhr.open("POST", "/upload");
          xhr.upload.addEventListener("progress", (e) => {
            if (e.lengthComputable) {
              this.progress = Math.floor((e.loaded / e.total) * 100);
            }
          });
          xhr.onload = () => {
            if (xhr.status === 200) location.reload();
            else alert("Upload error: " + xhr.statusText);
          };
          xhr.onerror = () => alert("Upload failed");
          xhr.send(form);
        }
      }
    }).mount('#app');
  </script>

  <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>
`

const videoHTML = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8" />
  <title>{{.Id}} - TritonTube</title>
  <link
    rel="stylesheet"
    href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css"
  />
  <script src="https://cdn.dashjs.org/latest/dash.all.min.js"></script>
</head>
<body class="bg-dark text-light">
  <div class="container py-5">
    <h1 class="text-warning">{{.Id}}</h1>
    <p>Uploaded at: {{ formatTime .UploadedAt }}</p>

    <div class="ratio ratio-16x9 mb-3">
      <video id="dashPlayer" controls class="rounded bg-black"></video>
    </div>

    <script>
      var url = "/content/{{.Id}}/manifest.mpd";
      var player = dashjs.MediaPlayer().create();
      player.initialize(document.querySelector("#dashPlayer"), url, false);
    </script>

    <a href="/" class="btn btn-outline-light mt-3">← Back to Home</a>
  </div>
</body>
</html>
`
