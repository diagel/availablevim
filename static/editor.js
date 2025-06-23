const editor = CodeMirror.fromTextArea(document.getElementById('editor'), {
    lineNumbers: true,
    keyMap: "vim",
    extensions: []
});

// Сохранение на сервер
function save() {
    const content = editor.getValue();
    fetch('/save', {
        method: 'POST',
        body: content
    }).then(() => alert("Файл сохранён на сервере"));
}

// Загрузка с сервера
function load() {
    fetch('/load')
        .then(res => res.text())
        .then(text => editor.setValue(text))
        .catch(() => alert("Файл не найден"));
}

// Скачивание файла
function download() {
    const content = editor.getValue();
    const blob = new Blob([content], {type: "text/plain"});
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "my-file.txt";
    a.click();
    URL.revokeObjectURL(url);
}

// Открытие файла с компьютера
document.getElementById("fileInput").addEventListener("change", function(event) {
    const file = event.target.files[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = function(e) {
        const contents = e.target.result;
        editor.setValue(contents); // Загружаем содержимое в редактор
    };
    reader.readAsText(file);
});