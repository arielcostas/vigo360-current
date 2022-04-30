let artId = window.location.pathname.split("/")[3]
let form = document.getElementById("extra-image-upload")
let imageStatus = document.getElementById("image-status")

async function listarImagenes() {
	let resp = await fetch('/admin/async/fotosExtra?articulo='+artId)
	let body = await resp.json()

	let imagenes = []
	if (resp.status == 200) {
		imagenes = body
	}
	return imagenes
}

function crearFigura(i) {
	let figura = document.createElement("div")

	let img = document.createElement("img")
	img.setAttribute("src", "/static/extra/" + i)
	img.setAttribute("title", i)
	img.setAttribute("alt", i)
	img.addEventListener("click", _ => {
		navigator.clipboard.writeText(`![Texto Alternativo](/static/extra/${i})`)
		imageStatus.innerText = `Imagen copiada al cortapapeles`
		setTimeout(() => {
			imageStatus.innerText = ""
		}, 2000)
	})

	let deleteButton = document.createElement("button")
	deleteButton.setAttribute("class", "btn-error btn-large")
	deleteButton.innerText = "Eliminar"
	deleteButton.addEventListener("click", e => {
		e.preventDefault()
		if (window.confirm(`¿Seguro que desea eliminar ${i}? Esta acción es irreversible`)) {
			eliminarFotografia(i)
		}
	})

	figura.appendChild(img)
	figura.appendChild(deleteButton)
	return figura
}

async function eliminarFotografia(n) {
	let resp = await fetch("/admin/async/fotosExtra?foto=" + n, {
		method: "DELETE"
	})
	let body = await resp.json()

	if (resp.status == 200) {
		cargarImagenes()
		imageStatus.innerText = `Imagen borrada`
		setTimeout(() => {
			imageStatus.innerText = ""
		}, 2000)
	} else {
		imageStatus.innerText = body.errorMsg
	}
}

function cargarImagenes() {
	let imageList = document.getElementById("editor-image-grid")
	imageStatus.innerHTML = "<b>Cargando datos...</b>"
	let imagenes = listarImagenes().then((img) => {
		if (img.length < 1) {
			imageStatus.innerHTML = "<b>Ninguna imágen adicional</b>"
		} else {
			imageStatus.innerHTML = ""
			imageList.innerHTML = ""
			img.forEach(i => {
				let figure = crearFigura(i)
				imageList.appendChild(figure)
			})
		}
	})
}

cargarImagenes()

form.addEventListener("submit", async (e) => {
	e.preventDefault()
	archivo = document.getElementById("archivo")
	fd = new FormData()
	fd.append("foto", archivo.files[0])
	fd.append("articulo", artId)

	if (archivo.files[0].size > 20971520) {
		imageStatus.innerText = `La imagen seleccionada es demasiado pesada`
		return
	}
	
	resp = await fetch('/admin/async/fotosExtra', {
		method: "POST",
		body: fd
	})
	if (resp.status != 200) {
		body = await resp.json()
		imageStatus.innerText = body.errorMsg
	} else {
		imageStatus.innerText = `Imagen guardada exitosamente`
		setTimeout(() => {
			imageStatus.innerText = ""
		}, 2000)
	}

	cargarImagenes()
})