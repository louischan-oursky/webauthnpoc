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

function deserializeGetOptions(getOptions) {
  const base64URLChallenge = getOptions.publicKey.challenge;
  const challenge = base64DecToArr(base64URLToBase64(base64URLChallenge));
  getOptions.publicKey.challenge = challenge;
  if (getOptions.publicKey.allowCredentials) {
    for (const c of getOptions.publicKey.allowCredentials) {
      c.id = base64DecToArr(base64URLToBase64(c.id));
    }
  }
}

function signIn(credential) {
  const credentialJSON = serializePublicKeyCredentialAssertion(credential);
  postJSONForText("/sign-in", credentialJSON).then((text) => {
    alert(text);
    console.log(text);
  }, (err) => {
    handleError(err);
  });
}

function handleError(err) {
  // Cancel
  if (err instanceof DOMException && err.name === "NotAllowedError") {
    console.log(err.message);
    return;
  }

  alert(err);
  console.error(err);
}
