function setupAutofill() {
  if (typeof PublicKeyCredential.isConditionalMediationAvailable === "function") {
    PublicKeyCredential.isConditionalMediationAvailable().then((ok) => {
      if (ok) {
        postForJSON(
          "/get-options-conditional",
          undefined
        ).then(getOptions => {
          deserializeGetOptions(getOptions);
          return navigator.credentials.get(getOptions);
        }).then(credential => {
          signIn(credential);
          setupAutofill();
        }).catch((_err) => {
          // We always need to maintain a conditional mediation promise so that when
          // the user autofill in any way, we can receive the credential.
          setupAutofill();
        });
      }
    });
  }
}
setupAutofill();
