function postForJSON(url, body) {
	return new Promise((resolve, reject) => {
		const xhr = new XMLHttpRequest();
		xhr.addEventListener("load", () => {
			resolve(xhr.response);
		});
		xhr.addEventListener("error", (err) => reject(err));
		xhr.responseType = "json";
		xhr.open("POST", url);
		xhr.send(body);
	});
}

function postJSONForText(url, value) {
	return new Promise((resolve, reject) => {
		const xhr = new XMLHttpRequest();
		xhr.addEventListener("load", () => {
			resolve(xhr.response);
		});
		xhr.addEventListener("error", (err) => reject(err));
		xhr.responseType = "text";
		xhr.open("POST", url);
		xhr.setRequestHeader("content-type", "application/json");
		xhr.send(JSON.stringify(value));
	});
}
