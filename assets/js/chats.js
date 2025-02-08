let chatsSearch = document.getElementById("chats-search")
function filterChats() {
	let searchQuery = chatsSearch.value.toLowerCase();
	const chats = document.querySelectorAll(".chat");
	chats.forEach(function(chat) {
		const text = chat.getAttribute('chat-desc').toLowerCase();
		if (text.includes(searchQuery)) {
			chat.style.display = '';
		} else {
			chat.style.display = 'none';
		}
	});
};
document.addEventListener('htmx:afterSwap', function(event) {
	if (event.target.id === 'messages-box') {
		setTimeout(() => {
			let container = document.getElementById('messages-container');
			container.scrollTop = container.scrollHeight;
		}, 50)
	}
});
