let seriesSelect = document.getElementById("series-select")
let seriesNum = document.getElementById("serie-num")

seriesNum.disabled = seriesSelect.value == ""

seriesSelect.addEventListener("change", () => {
	seriesNum.disabled = seriesSelect.value == ""
})
