document.body.addEventListener('htmx:beforeOnLoad', function(evt) {
	if (evt.detail.xhr.status === 422) {
		evt.detail.shouldSwap = true;
		evt.detail.isError = false;
	}
});
document.body.addEventListener('htmx:afterRequest', function(event) {
	if (event.detail.xhr.status === 302) {
		window.location.href = '/';
	}
});
