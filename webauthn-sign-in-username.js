document.querySelector("#get-form").addEventListener("submit", (e) => {
  e.preventDefault();
  e.stopPropagation();
  postForJSON(
    "/get-options-modal",
    new URLSearchParams(window.location.search)
  ).then(getOptions => {
    deserializeGetOptions(getOptions);
    return navigator.credentials.get(getOptions);
  }).then(credential => {
    signIn(credential);
  }).catch(err => {
    handleError(err);
  });
});
