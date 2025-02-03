var modal = document.getElementById("modal");
var span = document.getElementsByClassName("close")[0];
function openModal() {
	modal.style.display = "block";
}
function closeModal() {
	modal.style.display = "none";
}
window.onclick = function(event) {
	if (event.target == modal) {
		modal.style.display = "none";
	}
}
document.body.addEventListener("openModal", function(evt) {
	openModal()
})
document.body.addEventListener("closeModal", function(evt) {
	closeModal()
})
