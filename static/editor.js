const editor = CodeMirror.fromTextArea(document.getElementById("editor"), {
    lineNumbers: true,
    keyMap: "vim",
    mode: "text/x-go"
});

const output = document.getElementById("output");

function save() {
    const content = editor.getValue();
    fetch('/save', {
        method: 'POST',
        body: content,
        headers: { 'Content-Type': 'text/plain' }
    })
        .then(res => res.text())
        .then(() => alert('Saved to server'))
        .catch(err => console.error(err));
}

function load() {
    fetch('/load')
        .then(res => res.text())
        .then(text => editor.setValue(text))
        .catch(() => alert('File not found'));
}

function run() {
    const code = editor.getValue();
    fetch('/run', {
        method: 'POST',
        body: code,
        headers: { 'Content-Type': 'text/plain' }
    })
        .then(res => res.text())
        .then(data => {
            output.textContent = data;
        })
        .catch(err => {
            output.textContent = 'Error running code:\n' + err;
        });
}

function download() {
    const content = editor.getValue();
    const blob = new Blob([content], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'main.go';
    a.click();
    URL.revokeObjectURL(url);
}

function openFile() {
    const fileInput = document.getElementById('fileInput');
    fileInput.click();
}

document.getElementById('fileInput').addEventListener('change', function(event) {
    const file = event.target.files[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = function(e) {
        editor.setValue(e.target.result);
    };
    reader.readAsText(file);
});
