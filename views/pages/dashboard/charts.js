(async function() {
	let weekData = {{ .WeekSalesData }};
	new Chart(
		document.getElementById('week-sales'),
		{
			type: 'bar',
			data: {
				labels: weekData.map(row => row.Day),
				datasets: [
					{
						label: 'күннің кассасы',
						data: weekData.map(row => row.Amount)
					}
				]
			}
		}
	);
	let monthData = {{ .MonthSalesData }};
	new Chart(
		document.getElementById('month-sales'),
		{
			type: 'bar',
			data: {
				labels: monthData.map(row => row.Day),
				datasets: [
					{
						label: 'күннің кассасы',
						data: monthData.map(row => row.Amount)
					}
				]
			}
		}
	);
})();
