const drop = (event) => {
	event.preventDefault();
	let taskID = event.dataTransfer.getData("text");
	let leadId = event.dataTransfer.getData("lead_id");
	if (event.target.id.includes("task")) {
		return;
	}
	if (!confirm("Тапсырыстың жеткізілгеніне сенімдісіз бе?")) {
		return
	}
	const url = `/leads/${leadId}/complete`;
	fetch(url, {
		method: 'POST',
	})
		.then(response => {
			if (response.ok) {
				let task = document.getElementById(taskID)
				task.classList.remove('bg-orange-500');
				task.classList.add('bg-green-500');
				event.target.appendChild(task);
			}
		})
		.catch(error => {
			console.error('Error:', error);
		});
}
const allowdrop = (event) => {
	event.preventDefault();
}
const drag = (event) => {
	event.dataTransfer.setData("text", event.target.id);
	event.dataTransfer.setData("lead_id", event.target.getAttribute('lead-id'));
}
document.addEventListener('DOMContentLoaded', () => {
	const leadCells = document.querySelectorAll('.lead-cell');
	leadCells.forEach(cell => {
		document.addEventListener(cell.id, () => {
			cell.remove();
		});
	});
});
let element;
const openLeadModal = (id) => {
	document.getElementById("leads-modal-id").value = id
	openModal()
}
let now = new Date();
let year = now.getFullYear();
let month = String(now.getMonth() + 1).padStart(2, '0'); // Months are 0-indexed
let day = String(now.getDate()).padStart(2, '0');
let hours = String(now.getHours()).padStart(2, '0');
let minutes = String(now.getMinutes()).padStart(2, '0');

let datetimeNow = `${year}-${month}-${day}T${hours}:${minutes}`;
document.getElementById('payment-at').value = datetimeNow;
let saletype = document.getElementById('saletype');
let deliverycostInput = document.getElementById('deliverycost');
let loancostInput = document.getElementById("loancost")
let items = document.getElementById('items');
let itemsumInput = document.getElementById("itemsum")
let fullsumInput = document.getElementById("fullsum")
const itemsElement = document.getElementById('items');
function calcSaleSum() {
	let fullcost = parseFloat(deliverycostInput.value)
	let itemscost = 0
	for (let child of itemsElement.children) {
		let selectElement = child.querySelector('select');
		let itemPrice = parseFloat(selectElement.getAttribute('price'));
		let inputElement = child.querySelector('input');
		let inputValue = parseFloat(inputElement.value);
		itemscost += itemPrice * inputValue;
	}
	itemsumInput.value = itemscost
	fullcost += itemscost
	let loancost = 0
	if (saletype.value === 'kaspi-loan') {
		loancost = fullcost * 0.15
	}
	console.log(itemscost)
	console.log(fullcost)
	console.log(loancost)
	loancostInput.value = loancost
	fullsumInput.value = fullcost + loancost
}
let loancostdiv = document.getElementById("loancostdiv")
saletype.addEventListener('change', function() {
	const selectedValue = this.value;
	if (selectedValue === 'kaspi-loan') {
		loancostdiv.style.display = 'block';
	} else {
		loancostdiv.style.display = 'none';
		loancost.value = '0';
	}
});
let deliveryTypeSelect = document.getElementById("delivery-type")
let deliveryCostDiv = document.getElementById("delivery-cost-div")
deliveryTypeSelect.addEventListener('change', function() {
	const selectedValue = this.value;
	if (selectedValue === 'no') {
		deliveryCostDiv.style.display = 'none';
		deliverycostInput.value = '0';
	} else {
		deliveryCostDiv.style.display = 'block';
	}
});
document.addEventListener('htmx:afterSwap', function(event) {
	if (event.target.id === 'items') {
		calcSaleSum();
	}
});
let itemsInput = document.getElementById("items-input")
function parseCartItems() {
	let parsedItems = ""
	for (let child of itemsElement.children) {
		let selectElement = child.querySelector('select');
		let inputElement = child.querySelector('input');
		let inputValue = parseFloat(inputElement.value);
		let productId = selectElement.value
		let quantity = inputValue
		parsedItems += `${productId},${quantity};`
	}
	itemsInput.value = parsedItems
}
const taskElements = document.querySelectorAll('.task');
taskElements.forEach((taskElement) => {
	document.addEventListener(taskElement.id, () => {
		taskElement.remove();
	});
});
let newLeadsSearch = document.getElementById("new-leads-search")
let assignedLeadsSearch = document.getElementById("assigned-leads-search")
let indeliveryLeadsSearch = document.getElementById("indelivery-leads-search")
let completedLeadsSearch = document.getElementById("completed-leads-search")
function filterLeads(inputName, leadsClass) {
	let searchQuery = ""
	if (inputName === "new") {
		searchQuery = newLeadsSearch.value.toLowerCase();
	} else if (inputName === "assigned") {
		searchQuery = assignedLeadsSearch.value.toLowerCase();
	} else if (inputName === "indelivery") {
		searchQuery = indeliveryLeadsSearch.value.toLowerCase();
	} else if (inputName === "completed") {
		searchQuery = completedLeadsSearch.value.toLowerCase();
	}
	const leads = document.querySelectorAll(leadsClass);
	leads.forEach(function(lead) {
		const text = lead.getAttribute('phone').toLowerCase();
		if (text.includes(searchQuery)) {
			lead.style.display = ''; // Show the lead
		} else {
			lead.style.display = 'none'; // Hide the lead
		}
	});
};
