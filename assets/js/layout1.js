document.body.addEventListener('htmx:beforeOnLoad', function(evt) {
	if (evt.detail.xhr.status === 422) {
		evt.detail.shouldSwap = true;
		evt.detail.isError = false;
	} else if (evt.detail.xhr.status === 302) {
		window.location.href = '/';
	}
});
