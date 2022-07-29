## 專看簡介

你們看到的這個庫名叫excel，當初的想法只是覺得這個名稱可能比較多人知道，另外想達到類似的功能(但僅侷限於檢視部分)

excel的功能其實已經很完善，但用COM控制很麻煩，而且不能很便捷的瀏覽到本機端的資料，需要強制把圖片篩進去，

而本專案的特色可以靠純餵入CSV檔案和指定static資料夾就能以CSV的方式瀏覽到圖片。

## Build

```yaml
go build -o greenViewer.exe -ldflags "-s -w"
```

## USAGE

根據UI介面，填入

1. CSV檔案
2. static資料夾路徑: [csv的內容可以以此路徑來當作相對位置](https://github.com/CarsonSlovoka/excel/blob/102eac62be07d4bf716d6b52a284ba6d827c41d4/app/urls/static/js/src/file/file.js#L214)

   此資料夾內還可以包含:
   - *.{png, ...}:
   - 📂 `fonts` 您[可以在static資料夾內，新增fonts的檔案](https://github.com/CarsonSlovoka/excel/blob/102eac62be07d4bf716d6b52a284ba6d827c41d4/app/urls/static/js/src/file/file.js#L212-L214)，那麼您就可以透過UI來選擇該字形
     - *.{ttf, ...}

## 第三方依賴

| Name | Desc |
| ---- | ---- |
[![bootstrap-table](https://github-readme-stats.vercel.app/api/pin?username=wenzhixin&repo=bootstrap-table)](https://github.com/wenzhixin/bootstrap-table) | 表格的製作主要依靠此庫，當然連同很多相關的plugin也包含在內
