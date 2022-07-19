function serializePublicKeyCredentialAttestation(credential) {
  const attestationObject = trimNewline(base64ToBase64URL(base64EncArr(new Uint8Array(credential.response.attestationObject))));
  const clientDataJSON = trimNewline(base64ToBase64URL(base64EncArr(new Uint8Array(credential.response.clientDataJSON))));
  let transports = [];
  if (credential.response.getTransports) {
    transports = credential.response.getTransports();
  }
  const clientExtensionResults = credential.getClientExtensionResults();
  return {
    id: credential.id,
    rawId: credential.id,
    type: credential.type,
    response: {
      attestationObject,
      clientDataJSON,
      transports,
    },
    clientExtensionResults,
  };
}

function serializePublicKeyCredentialAssertion(credential) {
  const authenticatorData = trimNewline(base64ToBase64URL(base64EncArr(new Uint8Array(credential.response.authenticatorData))));
  const clientDataJSON = trimNewline(base64ToBase64URL(base64EncArr(new Uint8Array(credential.response.clientDataJSON))));
  const signature = trimNewline(base64ToBase64URL(base64EncArr(new Uint8Array(credential.response.signature))));
  const userHandle = trimNewline(base64ToBase64URL(base64EncArr(new Uint8Array(credential.response.userHandle))));
  const clientExtensionResults = credential.getClientExtensionResults();
  return {
    id: credential.id,
    rawId: credential.id,
    type: credential.type,
    response: {
      authenticatorData,
      clientDataJSON,
      signature,
      userHandle,
    },
    clientExtensionResults,
  }
}

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
      alert(err);
      console.error(err);
    });
  }).catch(err => {
    alert(err);
    console.error(err);
  });
});

document.querySelector("#get-form").addEventListener("submit", (e) => {
  e.preventDefault();
  e.stopPropagation();
  postForJSON(
    "/get-options",
    undefined
  ).then(getOptions => {
    const base64URLChallenge = getOptions.publicKey.challenge;
    const challenge = base64DecToArr(base64URLToBase64(base64URLChallenge));
    getOptions.publicKey.challenge = challenge;

    return navigator.credentials.get(getOptions);
  }).then(credential => {
    const credentialJSON = serializePublicKeyCredentialAssertion(credential);
    postJSONForText("/sign-in", credentialJSON).then((text) => {
      alert(text);
      console.log(text);
    }, (err) => {
      alert(err);
      console.error(err);
    });
  }).catch(err => {
    alert(err);
    console.error(err);
  });
});
