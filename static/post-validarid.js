var artIDfield = document.getElementById("art-id")
artIDfield.addEventListener("keyup", e => {
	e.target.value = e.target.value.replaceAll(" ", "-").toLowerCase()
	console.log(e.target.value)
})