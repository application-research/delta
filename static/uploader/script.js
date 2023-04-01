const form = document.querySelector('form');
const loader = document.getElementById('loader');
const fileInput = document.getElementById('file-input');
const minerInput = document.getElementById('miner-input');
const apiKeyInput = document.getElementById('api-input');
const messageArea = document.getElementById('message-area');
const contentInput = document.getElementById('content-input');
const downloadBtn = document.getElementById('download-btn');
const protocol = location.protocol;
const requestButton = document.getElementById('request-button');
const responseArea = document.getElementById('response-area');
const urlSelect = document.getElementById('url-select');

let defaultUrl = `${protocol}//${location.host}`;
let urlInput = document.getElementById('url-input');
defaultUrl = "http://localhost:1414";
urlInput.value = defaultUrl;
contentInput.value = defaultUrl + "/open/stats/content/";


// Download the response as a file
downloadBtn.addEventListener('click', () => {
    const text = responseArea.value;
    const blob = new Blob([text], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'status-check-for'+contentInput+'.txt';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
});

// Update the URL input field when the dropdown value changes
urlSelect.addEventListener('change', () => {
    const selectedUrl = urlSelect.value;
    urlInput.value = selectedUrl;
    contentInput.value = selectedUrl + "/open/stats/content/";
});

//  Submit the form
form.addEventListener('submit', (event) => {
    event.preventDefault();
    loader.classList.add('loading');
    const file = fileInput.files[0];
    const miner = minerInput.value;
    const apiKey = apiKeyInput.value
    const formData = new FormData();
    formData.append('data', file);
    formData.append("metadata", "{\"deal_verify_state\":\"verified\", \"duration_in_days\":521,\"start_epoch_in_days\": 1,\"auto_retry\":true}");
    messageArea.value = 'Uploading...';
    fetch(urlInput.value + '/api/v1/deal/end-to-end', {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${apiKey}`,
            'X-MinerID': miner
        },
        body: formData

    })
        .then(response => {
            if (response.ok) {
                return response.json();
            } else {
                throw new Error('Upload failed');
            }
        })
        .then(data => {
            messageArea.value = `${JSON.stringify(data, undefined,4)}</div>`;
        })
        .catch(error => {
            messageArea.value = `<div class="message error">${error}</div>`;
        });
});


//  Submit the form
requestButton.addEventListener('click', () => {
    let contentId = document.getElementById('content-id-input').value;
    const url = contentInput.value + contentId;
    fetch(url)
        .then(response => {
            if (response.ok) {
                return response.json();
            } else {
                throw new Error('Request failed');
            }
        })
        .then(data => {
            responseArea.value = JSON.stringify(data, undefined,4);
        })
        .catch(error => {
            responseArea.value = JSON.stringify(error.message, undefined,4);
        });
});