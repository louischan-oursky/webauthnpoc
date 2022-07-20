document.querySelector("#create-form").addEventListener("submit", (e) => {
  e.preventDefault();
  e.stopPropagation();
  postForJSON(
    "/create-options",
    new URLSearchParams(new FormData(e.currentTarget))
  ).then(createOptions => {
    const base64URLChallenge = createOptions.publicKey.challenge;
    const challenge = base64DecToArr(base64URLToBase64(base64URLChallenge));
    createOptions.publicKey.challenge = challenge;

    const base64URLUserID = createOptions.publicKey.user.id;
    const userID = base64DecToArr(base64URLToBase64(base64URLUserID));
    createOptions.publicKey.user.id = userID;

    if (createOptions.publicKey.excludeCredentials != null) {
      for (const c of createOptions.publicKey.excludeCredentials) {
        c.id = base64DecToArr(base64URLToBase64(c.id));
      }
    }

    return navigator.credentials.create(createOptions);
  }).then(credential => {
    const credentialJSON = serializePublicKeyCredentialAttestation(credential);
    postJSONForText("/register", credentialJSON).then((text) => {
      alert(text);
    }, (err) => {
      handleError(err);
    });
  }).catch(err => {
    handleError(err);
  });
});
