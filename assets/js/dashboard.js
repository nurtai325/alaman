document.addEventListener('DOMContentLoaded', () => {
	const elements = document.querySelectorAll('.number');
	elements.forEach((el) => {
		const targetValue = parseInt(el.dataset.target, 10);
		let currentValue = parseInt(targetValue / 1.2);
		const increment = () => {
			if (currentValue < targetValue) {
				currentValue += Math.ceil(targetValue / 100); // Increment in small steps
				if (currentValue > targetValue) currentValue = targetValue;
				let measure = el.classList.contains('scr') ? '%' : 'тг.'
				if (el.classList.contains('new-lead') || el.classList.contains('sale-count') || el.classList.contains('new-lead-today')) {
					measure = ''
				}
				el.textContent = `${currentValue} ${measure}`;
				setTimeout(increment, 20); // Refresh every 20ms
			}
		};
		el.classList.add('animate-count-up');
		increment();
	});
});
