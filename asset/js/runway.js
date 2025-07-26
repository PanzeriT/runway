document.body.addEventListener('htmx:responseError', function (evt) {
  if (evt.detail.xhr.status === 401) {
    window.location.href = '/login'
  }
});
