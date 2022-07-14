function base64URLToBase64(base64url) {
	let base64 = base64url.replace(/-/g, "+").replace(/_/g, "/");
	if (base64.length % 4 !== 0) {
		const count = 4 - base64.length % 4;
		base64 += "=".repeat(count);
	}
	return base64;
}

function base64ToBase64URL(base64) {
	return base64.replace(/\+/g, "-").replace(/\//g, "_").replace(/=/g, "");
}

function trimNewline(str) {
	return str.replace(/\r/g, "").replace(/\n/g, "");
}
